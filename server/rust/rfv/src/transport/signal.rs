use std::sync::Arc;

use futures_util::{SinkExt, StreamExt};
use prost::Message as ProstMessage;
use tokio::sync::mpsc;
use tokio_tungstenite::{connect_async, tungstenite::Message};

use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::{RegisterReq, RegisterResp, RelayToken, SignalEnvelope, PeerEndpoints};

pub struct SignalClient {
    out_tx: mpsc::UnboundedSender<Vec<u8>>,
    in_rx: Arc<tokio::sync::Mutex<mpsc::UnboundedReceiver<SignalEnvelope>>>,
    peer_rx: Arc<tokio::sync::Mutex<mpsc::UnboundedReceiver<PeerEndpoints>>>,
    relay_rx: Arc<tokio::sync::Mutex<mpsc::UnboundedReceiver<RelayToken>>>,
}

impl SignalClient {
    pub async fn connect(url: &str) -> Result<Self, String> {
        let (ws, _) = connect_async(url).await.map_err(|e| e.to_string())?;
        let (mut sink, mut stream) = ws.split();
        let (out_tx, mut out_rx) = mpsc::unbounded_channel::<Vec<u8>>();
        let (in_tx, in_rx) = mpsc::unbounded_channel::<SignalEnvelope>();
        let (peer_tx, peer_rx) = mpsc::unbounded_channel::<PeerEndpoints>();
        let (relay_tx, relay_rx) = mpsc::unbounded_channel::<RelayToken>();
        tokio::spawn(async move {
            loop {
                tokio::select! {
                    Some(buf) = out_rx.recv() => {
                        if sink.send(Message::Binary(buf.into())).await.is_err() {
                            break;
                        }
                    }
                    msg = stream.next() => {
                        let Some(Ok(Message::Binary(data))) = msg else { break };
                        let Ok(env) = SignalEnvelope::decode(data.as_ref()) else { continue };
                        match env.payload {
                            Some(Payload::RegisterAck(_)) | Some(Payload::Error(_)) => { let _ = in_tx.send(env); }
                            Some(Payload::IceParameters(_)) | Some(Payload::IceCandidate(_)) | Some(Payload::IceComplete(_)) => { let _ = in_tx.send(env); }
                            Some(Payload::PeerEndpoints(e)) => { let _ = peer_tx.send(e); }
                            Some(Payload::RelayToken(t)) => { let _ = relay_tx.send(t); }
                            _ => {}
                        }
                    }
                }
            }
        });
        Ok(Self {
            out_tx,
            in_rx: Arc::new(tokio::sync::Mutex::new(in_rx)),
            peer_rx: Arc::new(tokio::sync::Mutex::new(peer_rx)),
            relay_rx: Arc::new(tokio::sync::Mutex::new(relay_rx)),
        })
    }

    pub async fn register(&self, room: &str, role: &str) -> Result<RegisterResp, String> {
        let env = SignalEnvelope {
            payload: Some(Payload::Register(RegisterReq {
                room_code: room.into(),
                role: role.into(),
                ..Default::default()
            })),
            ..Default::default()
        };
        self.send(env).await?;
        let resp = self.recv_register().await?;
        Ok(resp)
    }

    pub async fn send(&self, env: SignalEnvelope) -> Result<(), String> {
        self.out_tx.send(env.encode_to_vec()).map_err(|e| e.to_string())
    }

    pub async fn send_bytes(&self, data: Vec<u8>) -> Result<(), String> {
        self.out_tx.send(data).map_err(|e| e.to_string())
    }

    pub async fn send_peer_endpoints(&self, eps: PeerEndpoints) -> Result<(), String> {
        self.send(SignalEnvelope { payload: Some(Payload::PeerEndpoints(eps)), ..Default::default() }).await
    }

    pub async fn recv_envelope(&self) -> Result<SignalEnvelope, String> {
        self.in_rx.lock().await.recv().await.ok_or_else(|| "signal closed".into())
    }

    pub async fn wait_peer_endpoints(&self) -> Result<PeerEndpoints, String> {
        self.peer_rx.lock().await.recv().await.ok_or_else(|| "signal closed".into())
    }

    pub async fn wait_relay_token(&self) -> Result<RelayToken, String> {
        self.relay_rx.lock().await.recv().await.ok_or_else(|| "signal closed".into())
    }

    async fn recv_register(&self) -> Result<RegisterResp, String> {
        loop {
            let env = self.recv_envelope().await?;
            if let Some(Payload::Error(e)) = env.payload {
                return Err(e);
            }
            if let Some(Payload::RegisterAck(a)) = env.payload {
                return Ok(a);
            }
        }
    }
}
