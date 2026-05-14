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

pub async fn listen() -> Result<SocketAddr, std::io::Error> {
    let addr = std::env::var("RB_RELAY_TCP").unwrap_or_else(|_| "127.0.0.1:0".to_string());
    let hub = Arc::new(RelayHub::new());
    let listener = TcpListener::bind(&addr).await?;
    let local = listener.local_addr()?;
    tokio::spawn(async move {
        loop {
            let Ok((sock, _)) = listener.accept().await else {
                break;
            };
            let hub2 = hub.clone();
            tokio::spawn(async move {
                hub2.handle_conn(sock).await;
            });
        }
    });
    Ok(local)
}

impl RelayHub {
    async fn handle_conn(&self, mut sock: TcpStream) {
        let (session_id, role) = match read_join(&mut sock).await {
            Ok(v) => v,
            Err(_) => {
                let _ = sock.shutdown().await;
                return;
            }
        };
        let idx = role as usize;
        let sp = self.session(&session_id);
        let mut g = sp.lock().await;
        if g.conns[idx].is_some() {
            drop(g);
            let _ = sock.shutdown().await;
            return;
        }
        g.conns[idx] = Some(sock);
        g.notify.notify_waiters();
        let other = 1 - idx;
        let deadline = tokio::time::Instant::now() + std::time::Duration::from_secs(120);
        loop {
            if g.conns[other].is_some() {
                break;
            }
            if tokio::time::Instant::now() >= deadline {
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
        self.bridge(&session_id, &mut viewer, &mut agent, &sp).await;
    }

    async fn maybe_delete_empty(&self, sid: &str, sp: &Arc<AsyncMutex<SessionPair>>) {
        let g = sp.lock().await;
        if g.conns[0].is_none() && g.conns[1].is_none() {
            drop(g);
            self.delete_session(sid);
        }
    }

    async fn bridge(&self, sid: &str, viewer: &mut TcpStream, agent: &mut TcpStream, sp: &Arc<AsyncMutex<SessionPair>>) {
        let (mut v_read, mut v_write) = viewer.split();
        let (mut a_read, mut a_write) = agent.split();
        let v2a = pipe(&mut a_write, &mut v_read);
        let a2v = pipe(&mut v_write, &mut a_read);
        let _ = tokio::join!(v2a, a2v);
        let mut g = sp.lock().await;
        g.conns = [None, None];
        g.bridging = false;
        drop(g);
        self.delete_session(sid);
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

async fn read_data_frame(sock: &mut (impl AsyncReadExt + Unpin)) -> Result<Vec<u8>, std::io::Error> {
    let mut sz = [0u8; 4];
    sock.read_exact(&mut sz).await?;
    let n = u32::from_be_bytes(sz) as usize;
    if n > 1 << 20 {
        return Err(std::io::Error::new(std::io::ErrorKind::InvalidData, "frame too large"));
    }
    if n == 0 {
        return Ok(Vec::new());
    }
    let mut buf = vec![0u8; n];
    sock.read_exact(&mut buf).await?;
    Ok(buf)
}

async fn write_data_frame(sock: &mut (impl AsyncWriteExt + Unpin), payload: &[u8]) -> Result<(), std::io::Error> {
    sock.write_all(&u32::to_be_bytes(payload.len() as u32)).await?;
    if payload.is_empty() {
        return Ok(());
    }
    sock.write_all(payload).await?;
    Ok(())
}

async fn pipe(dst: &mut (impl AsyncWriteExt + Unpin), src: &mut (impl AsyncReadExt + Unpin)) {
    loop {
        let payload = match read_data_frame(src).await {
            Ok(v) => v,
            Err(_) => return,
        };
        if write_data_frame(dst, &payload).await.is_err() {
            return;
        }
    }
}
