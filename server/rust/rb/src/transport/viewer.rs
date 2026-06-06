use std::sync::Arc;
use std::time::Duration;

use prost::Message as ProstMessage;

use crate::client::ice_viewer::ViewerHandle;
use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::SignalEnvelope;
use crate::transport::signal::SignalClient;
use crate::ice_link;
use crate::rb_http_client::RbHttpClient;
use crate::transport::link::{pick_viewer_link, ViewerLink};

pub async fn run_viewer(signal_url: String, room: String, list_root: String, ice_timeout_ms: u32) -> Result<(), String> {
    tracing::info!(%signal_url, %room, %list_root, ice_timeout_ms, "viewer run");
    let sig = Arc::new(SignalClient::connect(&signal_url).await?);
    let ice = Arc::new(ViewerHandle::new(ice_timeout_ms));
    let sig_pump = sig.clone();
    let ice_pump = ice.clone();
    tokio::spawn(async move { signal_pump_viewer(sig_pump, ice_pump).await });
    sig.register(&room, "viewer").await?;
    let link = pick_viewer_link(&sig, &ice).await?;
    let resp: crate::remotebrowse_svc::proto::ListFilesResponse = match link {
        ViewerLink::Tcp(sock) => {
            let mut c = RbHttpClient::connect_io(sock).await?;
            c.post_proto(
                "/rb/v1/list",
                &crate::remotebrowse_svc::proto::ListFilesRequest {
                    root_path: list_root.clone(),
                    ..Default::default()
                },
            )
            .await?
        }
        ViewerLink::Ice(ice_w) => {
            let mut c = ice_link::connect(ice_w).await?;
            c.post_proto(
                "/rb/v1/list",
                &crate::remotebrowse_svc::proto::ListFilesRequest {
                    root_path: list_root.clone(),
                    ..Default::default()
                },
            )
            .await?
        }
    };
    tracing::info!(%list_root, count = resp.entries.len(), "viewer list_files ok");
    for e in &resp.entries {
        eprintln!("{} size={}", e.name, e.size);
    }
    Ok(())
}

async fn signal_pump_viewer(sig: Arc<SignalClient>, ice: Arc<ViewerHandle>) {
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
