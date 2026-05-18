use std::sync::Arc;
use std::time::Duration;

use prost::Message as ProstMessage;

use crate::client::ice_agent::AgentHandle;
use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::SignalEnvelope;
use crate::transport::link::{gather_endpoints, listen_direct, pick_agent_link, AgentLink};
use crate::transport::signal::SignalClient;
use crate::transport::wire_agent;

pub async fn run_agent(signal_url: String, room: String, root: String, ice_timeout_ms: u32) -> Result<(), String> {
    let sig = Arc::new(SignalClient::connect(&signal_url).await?);
    let ice = Arc::new(AgentHandle::new(ice_timeout_ms));
    let sig_pump = sig.clone();
    let ice_pump = ice.clone();
    tokio::spawn(async move { signal_pump(sig_pump, ice_pump).await });
    sig.register(&room, "agent").await?;
    let (ln, port) = listen_direct().await?;
    sig.send_peer_endpoints(gather_endpoints(port)).await?;
    let link = pick_agent_link(&sig, ln, &ice).await?;
    match link {
        AgentLink::Tcp(t) => wire_agent::serve_tcp(t, root).await,
        AgentLink::Ice(i) => wire_agent::serve_ice(i, root).await,
    }
}

async fn signal_pump(sig: Arc<SignalClient>, ice: Arc<AgentHandle>) {
    loop {
        tokio::select! {
            env = sig.recv_envelope() => {
                let Ok(env) = env else { break };
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
