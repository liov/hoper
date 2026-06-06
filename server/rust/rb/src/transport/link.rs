use std::time::Duration;

use tokio::net::{TcpListener, TcpStream};

use crate::client::ice_agent::AgentHandle;
use crate::p2p_http;
use crate::client::ice_stream::IceWire;
use crate::client::ice_viewer::ViewerHandle;
use crate::signal_proto::{PeerEndpoint, PeerEndpoints, RelayToken};
use crate::transport::signal::SignalClient;

const DIRECT_PORT: u16 = 19091;
const DIRECT_PICK_SECS: u64 = 20;
const ROLE_VIEWER: u8 = 0;
const ROLE_AGENT: u8 = 1;

/// `RB_DIRECT_PORT`：直连监听端口；`0` 表示由系统分配（多 Agent 同机时可用）。
pub fn direct_listen_port() -> u16 {
    std::env::var("RB_DIRECT_PORT")
        .ok()
        .and_then(|s| s.parse().ok())
        .unwrap_or(DIRECT_PORT)
}

pub async fn listen_direct() -> Result<(TcpListener, u16), String> {
    let req = direct_listen_port();
    let ln = TcpListener::bind(("0.0.0.0", req)).await.map_err(|e| e.to_string())?;
    let local = ln.local_addr().map_err(|e| e.to_string())?;
    let bound = local.port();
    tracing::info!(requested = req, bound, %local, "direct tcp bind");
    Ok((ln, bound))
}

pub fn gather_endpoints(port: u16) -> PeerEndpoints {
    let mut items = Vec::new();
    if let Ok(ifaces) = if_addrs::get_if_addrs() {
        for iface in ifaces {
            if iface.is_loopback() {
                continue;
            }
            match iface.ip() {
                std::net::IpAddr::V4(ip) => {
                    items.push(PeerEndpoint { host: ip.to_string(), port: port as u32 });
                }
                std::net::IpAddr::V6(ip) => {
                    if !ip.is_loopback() && !ip.is_unspecified() {
                        items.push(PeerEndpoint { host: format!("[{ip}]"), port: port as u32 });
                    }
                }
            }
        }
    }
    if !items.is_empty() {
        let addrs: Vec<String> = items.iter().map(|e| format!("{}:{}", e.host, e.port)).collect();
        tracing::debug!(addrs = %addrs.join(","), "agent advertised endpoints");
    }
    PeerEndpoints { items }
}

pub enum AgentLink {
    Tcp(TcpStream),
    Ice(IceWire),
}

pub enum ViewerLink {
    Tcp(TcpStream),
    Ice(IceWire),
}

pub async fn pick_viewer_link(sig: &SignalClient, ice: &ViewerHandle) -> Result<ViewerLink, String> {
    tracing::debug!("viewer pick link: try direct tcp");
    if let Some(t) = try_direct_viewer(sig).await? {
        tracing::info!("viewer link: direct tcp");
        return Ok(ViewerLink::Tcp(t));
    }
    tracing::debug!("viewer pick link: try ice");
    if let Some(i) = try_ice_viewer(ice).await? {
        tracing::info!("viewer link: ice");
        return Ok(ViewerLink::Ice(i));
    }
    let tok = sig.wait_relay_token().await?;
    tracing::info!(relay = %tok.relay_host, port = tok.relay_port, session = %tok.session_id, "viewer link: relay");
    Ok(ViewerLink::Tcp(dial_relay_viewer(tok).await?))
}

/// 后台持续 accept :19091，供 IP 直连（无需先走信令配对）。
pub fn spawn_direct_inbound_acceptor(ln: TcpListener, sandbox: Option<String>, room: String) {
    tokio::spawn(async move {
        loop {
            match ln.accept().await {
                Ok((sock, addr)) => {
                    let sb = sandbox.clone();
                    let room_c = room.clone();
                    tracing::info!(%addr, "agent direct inbound (http2)");
                    tokio::spawn(async move {
                        let _ = p2p_http::serve_h2(sock, sb, Some(room_c)).await;
                    });
                }
                Err(e) => {
                    tracing::warn!(%e, "agent direct accept loop exit");
                    break;
                }
            }
        }
    });
}

pub async fn pick_agent_link(
    sig: &SignalClient,
    ice: &AgentHandle,
    viewer_eps: Option<PeerEndpoints>,
    pending_relay: Option<RelayToken>,
) -> Result<AgentLink, String> {
    tracing::debug!("agent pick link: try direct tcp outbound");
    if let Some(t) = try_direct_agent_outbound(sig, viewer_eps).await? {
        tracing::info!("agent link: direct tcp outbound");
        return Ok(AgentLink::Tcp(t));
    }
    tracing::debug!("agent pick link: try ice");
    if let Some(i) = try_ice_agent(ice).await? {
        tracing::info!("agent link: ice");
        return Ok(AgentLink::Ice(i));
    }
    let tok = match pending_relay {
        Some(t) => t,
        None => sig.wait_relay_token().await?,
    };
    tracing::info!(relay = %tok.relay_host, port = tok.relay_port, session = %tok.session_id, "agent link: relay");
    Ok(AgentLink::Tcp(dial_relay(tok).await?))
}

async fn dial_peer_endpoints(eps: &PeerEndpoints) -> Option<TcpStream> {
    let mut tasks = Vec::new();
    for ep in &eps.items {
        if ep.host.is_empty() || ep.port == 0 {
            continue;
        }
        let host = ep.host.clone();
        let port = ep.port as u16;
        tracing::debug!(%host, port, "agent direct: dial viewer");
        tasks.push(tokio::spawn(async move {
            match TcpStream::connect((host.as_str(), port)).await {
                Ok(sock) => Some((host, port, sock)),
                Err(e) => {
                    tracing::trace!(%host, port, err = %e, "agent direct: dial failed");
                    None
                }
            }
        }));
    }
    while !tasks.is_empty() {
        let (res, _idx, remain) = futures_util::future::select_all(tasks).await;
        tasks = remain;
        if let Ok(Some((host, port, sock))) = res {
            tracing::info!(%host, port, "agent direct: connected");
            return Some(sock);
        }
    }
    None
}

async fn try_direct_agent_outbound(
    sig: &SignalClient,
    viewer_eps: Option<PeerEndpoints>,
) -> Result<Option<TcpStream>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(DIRECT_PICK_SECS);
    let mut outbound = viewer_eps;
    let mut tried_outbound = false;
    loop {
        if tokio::time::Instant::now() >= deadline {
            tracing::debug!("agent direct outbound: pick timeout");
            return Ok(None);
        }
        if outbound.is_none() {
            match tokio::time::timeout(Duration::from_millis(200), sig.wait_peer_endpoints()).await {
                Ok(Ok(eps)) => outbound = Some(eps),
                _ => {}
            }
        }
        if let Some(eps) = outbound.take() {
            if eps.items.is_empty() || tried_outbound {
                tokio::time::sleep(Duration::from_millis(50)).await;
                continue;
            }
            tried_outbound = true;
            if let Some(t) = dial_peer_endpoints(&eps).await {
                return Ok(Some(t));
            }
            tracing::debug!("agent direct outbound: dial failed, retry");
        }
        tokio::time::sleep(Duration::from_millis(50)).await;
    }
}

async fn try_ice_agent(ice: &AgentHandle) -> Result<Option<IceWire>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(12);
    while tokio::time::Instant::now() < deadline {
        if let Some(w) = ice.try_take_wire() {
            return Ok(Some(w));
        }
        if ice.state_code() < 0 {
            tracing::debug!("agent ice: failed");
            return Ok(None);
        }
        tokio::time::sleep(Duration::from_millis(20)).await;
    }
    tracing::debug!("agent ice: timeout");
    Ok(None)
}

pub async fn dial_relay(tok: RelayToken) -> Result<TcpStream, String> {
    dial_relay_role(tok, ROLE_AGENT).await
}

pub async fn dial_relay_viewer(tok: RelayToken) -> Result<TcpStream, String> {
    dial_relay_role(tok, ROLE_VIEWER).await
}

async fn dial_relay_role(tok: RelayToken, role: u8) -> Result<TcpStream, String> {
    let addr = format!("{}:{}", tok.relay_host, tok.relay_port);
    tracing::debug!(%addr, role, session = %tok.session_id, "relay dial");
    let mut sock = TcpStream::connect(&addr).await.map_err(|e| e.to_string())?;
    write_relay_join(&mut sock, &tok.session_id, role).await?;
    tracing::info!(%addr, role, "relay joined (raw tcp for h2)");
    Ok(sock)
}

async fn try_direct_viewer(sig: &SignalClient) -> Result<Option<TcpStream>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(DIRECT_PICK_SECS);
    loop {
        if tokio::time::Instant::now() >= deadline {
            tracing::debug!("viewer direct: pick timeout");
            return Ok(None);
        }
        match tokio::time::timeout(Duration::from_millis(200), sig.wait_peer_endpoints()).await {
            Ok(Ok(eps)) => {
                for ep in eps.items {
                    if ep.host.is_empty() || ep.port == 0 {
                        continue;
                    }
                    tracing::debug!(host = %ep.host, port = ep.port, "viewer direct: dial agent");
                    match TcpStream::connect((ep.host.as_str(), ep.port as u16)).await {
                        Ok(sock) => {
                            tracing::info!(host = %ep.host, port = ep.port, "viewer direct: connected");
                            return Ok(Some(sock));
                        }
                        Err(e) => tracing::debug!(host = %ep.host, port = ep.port, err = %e, "viewer direct: dial failed"),
                    }
                }
            }
            _ => {}
        }
        tokio::time::sleep(Duration::from_millis(50)).await;
    }
}

async fn try_ice_viewer(ice: &ViewerHandle) -> Result<Option<IceWire>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(12);
    while tokio::time::Instant::now() < deadline {
        if let Some(w) = ice.try_take_wire() {
            return Ok(Some(w));
        }
        if ice.state_code() < 0 {
            tracing::debug!("viewer ice: failed");
            return Ok(None);
        }
        tokio::time::sleep(Duration::from_millis(20)).await;
    }
    tracing::debug!("viewer ice: timeout");
    Ok(None)
}

async fn write_relay_join(sock: &mut TcpStream, session_id: &str, role: u8) -> Result<(), String> {
    use tokio::io::AsyncWriteExt;
    let id = uuid::Uuid::parse_str(session_id).map_err(|e| e.to_string())?;
    let mut buf = Vec::with_capacity(22);
    buf.extend_from_slice(b"RBRL");
    buf.push(1);
    buf.extend_from_slice(id.as_bytes());
    buf.push(role);
    sock.write_all(&buf).await.map_err(|e| e.to_string())?;
    Ok(())
}
