use std::collections::HashMap;
use std::net::SocketAddr;
use std::sync::{Arc, Mutex};

use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::{TcpListener, TcpStream};
use tokio::sync::{Mutex as AsyncMutex, Notify};
use uuid::Uuid;

const MAGIC: &[u8; 4] = b"RBRL";
const FRAME_VER: u8 = 1;
const ROLE_VIEWER: u8 = 0;
const ROLE_AGENT: u8 = 1;

struct SessionPair {
    conns: [Option<TcpStream>; 2],
    bridging: bool,
    notify: Notify,
}

struct RelayHub {
    sessions: Mutex<HashMap<String, Arc<AsyncMutex<SessionPair>>>>,
}

impl RelayHub {
    fn new() -> Self {
        Self { sessions: Mutex::new(HashMap::new()) }
    }

    fn session(&self, sid: &str) -> Arc<AsyncMutex<SessionPair>> {
        let mut g = self.sessions.lock().expect("relay sessions");
        if let Some(s) = g.get(sid) {
            return s.clone();
        }
        let s = Arc::new(AsyncMutex::new(SessionPair {
            conns: [None, None],
            bridging: false,
            notify: Notify::new(),
        }));
        g.insert(sid.to_string(), s.clone());
        s
    }

    fn delete_session(&self, sid: &str) {
        self.sessions.lock().expect("relay sessions").remove(sid);
    }
}

/// 绑定 `RB_RELAY_TCP`（默认 `0.0.0.0:0`）。下发给客户端的地址见 `relay_advertise_addr`。
pub async fn listen() -> Result<String, std::io::Error> {
    let addr = std::env::var("RB_RELAY_TCP").unwrap_or_else(|_| "0.0.0.0:0".to_string());
    let hub = Arc::new(RelayHub::new());
    let listener = bind_relay(&addr).await?;
    let local = listener.local_addr()?;
    let advertise = relay_advertise_addr(local);
    tracing::info!(%local, %advertise, bind = %addr, "relay tcp listening");
    tokio::spawn(async move {
        loop {
            match listener.accept().await {
                Ok((sock, peer)) => {
                    tracing::debug!(%peer, "relay tcp accept");
                    let hub2 = hub.clone();
                    tokio::spawn(async move {
                        hub2.handle_conn(sock).await;
                    });
                }
                Err(e) => {
                    tracing::warn!(%e, "relay accept error, continue");
                    tokio::time::sleep(std::time::Duration::from_millis(200)).await;
                }
            }
        }
    });
    Ok(advertise)
}

async fn bind_relay(addr: &str) -> Result<TcpListener, std::io::Error> {
    loop {
        match TcpListener::bind(addr).await {
            Ok(ln) => return Ok(ln),
            Err(e) => {
                tracing::warn!(%addr, %e, "relay bind failed, retry in 1s");
                tokio::time::sleep(std::time::Duration::from_secs(1)).await;
            }
        }
    }
}

fn relay_advertise_addr(local: SocketAddr) -> String {
    if let Ok(s) = std::env::var("RB_RELAY_ADVERTISE") {
        return s;
    }
    let port = local.port();
    if let Ok(host) = std::env::var("RB_RELAY_HOST") {
        let host = host.trim().trim_matches(|c| c == '[' || c == ']');
        if !host.is_empty() {
            return format!("{host}:{port}");
        }
    }
    let ip = local.ip();
    if !ip.is_unspecified() && !ip.is_loopback() {
        return format!("{ip}:{port}");
    }
    if let Ok(http) = std::env::var("RB_HTTP") {
        if let Ok(sa) = http.parse::<SocketAddr>() {
            let hip = sa.ip();
            if !hip.is_unspecified() && !hip.is_loopback() {
                return format!("{hip}:{port}");
            }
        }
    }
    if let Some(ip) = first_non_loopback_ipv4() {
        tracing::info!(%ip, "relay advertise from local interface");
        return format!("{ip}:{port}");
    }
    tracing::warn!(
        "relay has no public advertise address; set RB_RELAY_ADVERTISE or RB_RELAY_HOST for cross-device"
    );
    format!("{ip}:{port}")
}

fn first_non_loopback_ipv4() -> Option<std::net::Ipv4Addr> {
    let ifaces = if_addrs::get_if_addrs().ok()?;
    for iface in ifaces {
        if iface.is_loopback() {
            continue;
        }
        if let std::net::IpAddr::V4(ip) = iface.ip() {
            if !ip.is_unspecified() && !ip.is_loopback() {
                return Some(ip);
            }
        }
    }
    None
}

impl RelayHub {
    async fn handle_conn(&self, mut sock: TcpStream) {
        let (session_id, role) = match read_join(&mut sock).await {
            Ok(v) => v,
            Err(e) => {
                tracing::debug!(err = %e, "relay join rejected");
                let _ = sock.shutdown().await;
                return;
            }
        };
        let idx = role as usize;
        let sp = self.session(&session_id);
        let mut g = sp.lock().await;
        if let Some(mut old) = g.conns[idx].take() {
            tracing::warn!(%session_id, role, "relay role replaced");
            let _ = old.shutdown().await;
        }
        tracing::info!(%session_id, role, "relay peer joined");
        g.conns[idx] = Some(sock);
        g.notify.notify_waiters();
        let other = 1 - idx;
        let deadline = tokio::time::Instant::now() + std::time::Duration::from_secs(120);
        loop {
            if g.conns[other].is_some() {
                break;
            }
            if tokio::time::Instant::now() >= deadline {
                tracing::warn!(%session_id, role, "relay wait peer timeout");
                if let Some(mut sock) = g.conns[idx].take() {
                    let _ = sock.shutdown().await;
                }
                drop(g);
                self.maybe_delete_empty(&session_id, &sp).await;
                return;
            }
            drop(g);
            tokio::time::sleep(std::time::Duration::from_millis(50)).await;
            g = sp.lock().await;
        }
        if g.bridging {
            drop(g);
            return;
        }
        g.bridging = true;
        let mut viewer = g.conns[0].take().expect("viewer");
        let mut agent = g.conns[1].take().expect("agent");
        drop(g);
        tracing::info!(%session_id, "relay bridge start (raw tcp)");
        if let Err(e) = tokio::io::copy_bidirectional(&mut viewer, &mut agent).await {
            tracing::debug!(%session_id, err = %e, "relay bridge ended");
        }
        tracing::info!(%session_id, "relay bridge end");
    }

    async fn maybe_delete_empty(&self, sid: &str, sp: &Arc<AsyncMutex<SessionPair>>) {
        let g = sp.lock().await;
        if g.conns[0].is_none() && g.conns[1].is_none() {
            drop(g);
            self.delete_session(sid);
        }
    }

}

async fn read_join(sock: &mut TcpStream) -> Result<(String, u8), std::io::Error> {
    let mut magic = [0u8; 4];
    sock.read_exact(&mut magic).await?;
    if &magic != MAGIC {
        return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "bad magic"));
    }
    let mut hdr = [0u8; 18];
    sock.read_exact(&mut hdr).await?;
    if hdr[0] != FRAME_VER {
        return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "bad frame"));
    }
    let role = hdr[17];
    if role != ROLE_VIEWER && role != ROLE_AGENT {
        return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "bad role"));
    }
    let id = Uuid::from_bytes(hdr[1..17].try_into().expect("uuid"));
    Ok((id.to_string(), role))
}

