use std::time::Duration;

use tokio::net::{TcpListener, TcpStream};

use crate::client::ice_agent::AgentHandle;
use crate::client::ice_stream::IceWire;
use crate::client::ice_viewer::ViewerHandle;
use crate::signal_proto::{PeerEndpoint, PeerEndpoints, RelayToken};
use crate::transport::signal::SignalClient;
use crate::transport::tcp_wire::TcpWire;

const DIRECT_PORT: u16 = 19091;
const ROLE_VIEWER: u8 = 0;
const ROLE_AGENT: u8 = 1;

pub async fn listen_direct() -> Result<(TcpListener, u16), String> {
    let ln = TcpListener::bind(("0.0.0.0", DIRECT_PORT)).await.map_err(|e| e.to_string())?;
    Ok((ln, DIRECT_PORT))
}

pub fn gather_endpoints(port: u16) -> PeerEndpoints {
    let mut items = Vec::new();
    if let Ok(ifaces) = if_addrs::get_if_addrs() {
        for iface in ifaces {
            if iface.is_loopback() {
                continue;
            }
            if let std::net::IpAddr::V4(ip) = iface.ip() {
                items.push(PeerEndpoint { host: ip.to_string(), port: port as u32 });
            }
        }
    }
    PeerEndpoints { items }
}

pub enum AgentLink {
    Tcp(TcpWire),
    Ice(IceWire),
}

pub enum ViewerLink {
    Tcp(TcpWire),
    Ice(IceWire),
}

pub async fn pick_viewer_link(sig: &SignalClient, ice: &ViewerHandle) -> Result<ViewerLink, String> {
    if let Some(t) = try_direct_viewer(sig).await? {
        return Ok(ViewerLink::Tcp(t));
    }
    if let Some(i) = try_ice_viewer(ice).await? {
        return Ok(ViewerLink::Ice(i));
    }
    let tok = sig.wait_relay_token().await?;
    Ok(ViewerLink::Tcp(dial_relay_viewer(tok).await?))
}

pub async fn pick_agent_link(sig: &SignalClient, ln: TcpListener, ice: &AgentHandle) -> Result<AgentLink, String> {
    if let Some(t) = try_direct_agent(sig, ln).await? {
        return Ok(AgentLink::Tcp(t));
    }
    if let Some(i) = try_ice_agent(ice).await? {
        return Ok(AgentLink::Ice(i));
    }
    let tok = sig.wait_relay_token().await?;
    Ok(AgentLink::Tcp(dial_relay(tok).await?))
}

async fn try_direct_agent(sig: &SignalClient, ln: TcpListener) -> Result<Option<TcpWire>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(5);
    loop {
        if tokio::time::Instant::now() >= deadline {
            return Ok(None);
        }
        tokio::select! {
            accept = ln.accept() => {
                if let Ok((sock, _)) = accept {
                    return Ok(Some(TcpWire::direct(sock)));
                }
            }
            eps = sig.wait_peer_endpoints() => {
                if let Ok(eps) = eps {
                    for ep in eps.items {
                        if ep.host.is_empty() || ep.port == 0 {
                            continue;
                        }
                        if let Ok(sock) = TcpStream::connect((ep.host.as_str(), ep.port as u16)).await {
                            return Ok(Some(TcpWire::direct(sock)));
                        }
                    }
                }
            }
            _ = tokio::time::sleep(Duration::from_millis(50)) => {}
        }
    }
}

async fn try_ice_agent(ice: &AgentHandle) -> Result<Option<IceWire>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(12);
    while tokio::time::Instant::now() < deadline {
        if let Some(w) = ice.try_take_wire() {
            return Ok(Some(w));
        }
        if ice.state_code() < 0 {
            return Ok(None);
        }
        tokio::time::sleep(Duration::from_millis(20)).await;
    }
    Ok(None)
}

pub async fn dial_relay(tok: RelayToken) -> Result<TcpWire, String> {
    dial_relay_role(tok, ROLE_AGENT).await
}

pub async fn dial_relay_viewer(tok: RelayToken) -> Result<TcpWire, String> {
    dial_relay_role(tok, ROLE_VIEWER).await
}

async fn dial_relay_role(tok: RelayToken, role: u8) -> Result<TcpWire, String> {
    let addr = format!("{}:{}", tok.relay_host, tok.relay_port);
    let mut sock = TcpStream::connect(&addr).await.map_err(|e| e.to_string())?;
    write_relay_join(&mut sock, &tok.session_id, role).await?;
    Ok(TcpWire::relay(sock))
}

async fn try_direct_viewer(sig: &SignalClient) -> Result<Option<TcpWire>, String> {
    let deadline = tokio::time::Instant::now() + Duration::from_secs(5);
    loop {
        if tokio::time::Instant::now() >= deadline {
            return Ok(None);
        }
        match tokio::time::timeout(Duration::from_millis(200), sig.wait_peer_endpoints()).await {
            Ok(Ok(eps)) => {
                for ep in eps.items {
                    if ep.host.is_empty() || ep.port == 0 {
                        continue;
                    }
                    if let Ok(sock) = TcpStream::connect((ep.host.as_str(), ep.port as u16)).await {
                        return Ok(Some(TcpWire::direct(sock)));
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
            return Ok(None);
        }
        tokio::time::sleep(Duration::from_millis(20)).await;
    }
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
