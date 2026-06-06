use std::sync::Arc;
use std::time::Duration;

use prost::Message as ProstMessage;

use crate::client::ice_agent::AgentHandle;
use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::SignalEnvelope;
use crate::transport::link::{gather_endpoints, listen_direct, pick_agent_link, spawn_direct_inbound_acceptor, AgentLink};
use crate::transport::signal::SignalClient;
use crate::ice_link;
use crate::p2p_http;

/// `sandbox`：可选 `RB_AGENT_SANDBOX`；未设置时不限制可访问路径。
/// 永不返回：信令断开或任意错误仅退避重连，始终保持可接受 Viewer 连接。
pub async fn run_agent(
    signal_url: String,
    room: String,
    sandbox: Option<String>,
    ice_timeout_ms: u32,
) {
    let sandbox = sandbox.filter(|s| !s.is_empty());
    tracing::info!(%signal_url, %room, ice_timeout_ms, sandbox = sandbox.as_deref().unwrap_or("(none)"), "agent run");
    let mut listen_backoff = Duration::from_secs(1);
    let direct_port = loop {
        match listen_direct().await {
            Ok((ln, port)) => {
                spawn_direct_inbound_acceptor(ln, sandbox.clone(), room.clone());
                break port;
            }
            Err(e) => {
                tracing::warn!(
                    err = %e,
                    wait_secs = listen_backoff.as_secs(),
                    hint = "端口可能被其它 rb Agent 占用，结束旧进程或设置 RB_DIRECT_PORT=0",
                    "agent direct listen failed"
                );
                tokio::time::sleep(listen_backoff).await;
                listen_backoff = (listen_backoff * 2).min(Duration::from_secs(30));
            }
        }
    };
    tracing::info!(direct_port, "agent direct acceptor ready (reused across signal reconnects)");
    let mut backoff = Duration::from_secs(1);
    loop {
        match run_agent_signal(&signal_url, &room, sandbox.clone(), ice_timeout_ms, direct_port).await {
            Err(e) => {
                tracing::warn!(err = %e, wait_secs = backoff.as_secs(), "agent signal reconnect");
                tokio::time::sleep(backoff).await;
                backoff = (backoff * 2).min(Duration::from_secs(30));
            }
            Ok(()) => {
                tracing::warn!("agent signal session ended unexpectedly, reconnect");
                backoff = Duration::from_secs(1);
            }
        }
    }
}

async fn run_agent_signal(
    signal_url: &str,
    room: &str,
    sandbox: Option<String>,
    ice_timeout_ms: u32,
    direct_port: u16,
) -> Result<(), String> {
    let sig = Arc::new(SignalClient::connect(signal_url).await?);
    let result = run_agent_session(&sig, room, sandbox, ice_timeout_ms, direct_port).await;
    sig.shutdown();
    result
}

async fn run_agent_session(
    sig: &Arc<SignalClient>,
    room: &str,
    sandbox: Option<String>,
    ice_timeout_ms: u32,
    port: u16,
) -> Result<(), String> {
    let ack = sig.register(room, "agent").await?;
    tracing::info!(peer_id = %ack.peer_id, room_id = %ack.room_id, "agent registered");
    let ice = Arc::new(AgentHandle::new(ice_timeout_ms));
    let sig_pump = sig.clone();
    let ice_pump = ice.clone();
    tokio::spawn(async move { signal_pump(sig_pump, ice_pump).await });
    loop {
        if !sig.signal_alive() {
            return Err("signal closed".into());
        }
        let eps = gather_endpoints(port);
        tracing::debug!(count = eps.items.len(), direct_port = port, "agent peer endpoints");
        if let Err(e) = sig.send_peer_endpoints(eps).await {
            if !sig.signal_alive() {
                return Err(e);
            }
            tracing::warn!(%e, "agent send peer_endpoints failed");
            continue;
        }
        let viewer_eps = sig.wait_viewer_peer_endpoints().await?;
        let eps2 = gather_endpoints(port);
        tracing::debug!(count = eps2.items.len(), "agent peer endpoints (after viewer ready)");
        if let Err(e) = sig.send_peer_endpoints(eps2).await {
            if !sig.signal_alive() {
                return Err(e);
            }
            tracing::warn!(%e, "agent send peer_endpoints failed (after viewer ready)");
            continue;
        }
        let link = match pick_agent_link(sig, &ice, Some(viewer_eps), None).await {
            Ok(l) => l,
            Err(e) => {
                tracing::warn!(%e, "agent pick link failed, wait next viewer");
                continue;
            }
        };
        let sig_watch = sig.clone();
        let sandbox_c = sandbox.clone();
        let room_c = room.to_string();
        let r = tokio::select! {
            r = async {
                match link {
                    AgentLink::Tcp(sock) => {
                        tracing::info!("agent link: tcp http+grpc");
                        p2p_http::serve_h2(sock, sandbox_c, Some(room_c)).await
                    }
                    AgentLink::Ice(ice) => {
                        tracing::info!("agent link: ice http+grpc");
                        ice_link::serve(ice, sandbox_c, Some(room_c)).await
                    }
                }
            } => r,
            () = async {
                while sig_watch.signal_alive() {
                    tokio::time::sleep(Duration::from_millis(200)).await;
                }
            } => {
                tracing::warn!("signal ws lost during grpc session");
                Err("signal closed".into())
            }
        };
        match r {
            Ok(()) => tracing::info!("agent grpc session ended, waiting next viewer"),
            Err(e) => {
                if !sig.signal_alive() {
                    return Err(e);
                }
                tracing::warn!(%e, "agent grpc session ended, keep signal and wait next viewer");
            }
        }
    }
}

async fn signal_pump(sig: Arc<SignalClient>, ice: Arc<AgentHandle>) {
    loop {
        tokio::select! {
            env = sig.recv_envelope() => {
                let Ok(env) = env else {
                    tracing::debug!("agent signal pump ended (ws closed, wire may continue)");
                    break;
                };
                if is_ice_payload(&env) {
                    ice.push(&env.encode_to_vec());
                }
            }
            _ = tokio::time::sleep(Duration::from_millis(15)) => {
                while let Some(b) = ice.poll_out() {
                    if sig.send_bytes(b).await.is_err() {
                        return;
                    }
                }
            }
        }
    }
}

fn is_ice_payload(env: &SignalEnvelope) -> bool {
    matches!(
        env.payload,
        Some(Payload::IceParameters(_)) | Some(Payload::IceCandidate(_)) | Some(Payload::IceComplete(_))
    )
}
