use std::sync::Arc;
use std::sync::Mutex;
use std::sync::atomic::{AtomicBool, Ordering};
use std::time::Duration;

use futures_util::{SinkExt, StreamExt};
use prost::Message as ProstMessage;
use tokio::sync::mpsc;
use tokio::task::JoinHandle;
use tokio_tungstenite::{connect_async, tungstenite::Message};

use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::{DeviceCapabilities, RegisterReq, RegisterResp, RelayToken, SignalEnvelope, PeerEndpoints};

pub struct SignalClient {
    out_tx: mpsc::UnboundedSender<Vec<u8>>,
    in_rx: Arc<tokio::sync::Mutex<mpsc::UnboundedReceiver<SignalEnvelope>>>,
    peer_rx: Arc<tokio::sync::Mutex<mpsc::UnboundedReceiver<PeerEndpoints>>>,
    relay_rx: Arc<tokio::sync::Mutex<mpsc::UnboundedReceiver<RelayToken>>>,
    alive: Arc<AtomicBool>,
    ws_task: Mutex<Option<JoinHandle<()>>>,
}

impl SignalClient {
    pub async fn connect(url: &str) -> Result<Self, String> {
        tracing::info!(%url, "signal ws connect");
        let (ws, resp) = connect_async(url).await.map_err(|e| e.to_string())?;
        tracing::debug!(status = ?resp.status(), "signal ws connected");
        let (mut sink, mut stream) = ws.split();
        let (out_tx, mut out_rx) = mpsc::unbounded_channel::<Vec<u8>>();
        let (in_tx, in_rx) = mpsc::unbounded_channel::<SignalEnvelope>();
        let (peer_tx, peer_rx) = mpsc::unbounded_channel::<PeerEndpoints>();
        let (relay_tx, relay_rx) = mpsc::unbounded_channel::<RelayToken>();
        let alive = Arc::new(AtomicBool::new(true));
        let alive_task = alive.clone();
        let ws_task = tokio::spawn(async move {
            let mut ping = tokio::time::interval(Duration::from_secs(25));
            ping.set_missed_tick_behavior(tokio::time::MissedTickBehavior::Delay);
            loop {
                tokio::select! {
                    Some(buf) = out_rx.recv() => {
                        if sink.send(Message::Binary(buf.into())).await.is_err() {
                            tracing::debug!("signal ws send failed");
                            break;
                        }
                    }
                    msg = stream.next() => {
                        match msg {
                            Some(Ok(Message::Binary(data))) => {
                                let Ok(env) = SignalEnvelope::decode(data.as_ref()) else { continue };
                                match env.payload {
                                    Some(Payload::RegisterAck(_)) | Some(Payload::Error(_)) => { let _ = in_tx.send(env); }
                                    Some(Payload::IceParameters(_)) | Some(Payload::IceCandidate(_)) | Some(Payload::IceComplete(_)) => { let _ = in_tx.send(env); }
                                    Some(Payload::PeerEndpoints(e)) => {
                                        tracing::debug!(n = e.items.len(), "signal recv peer_endpoints");
                                        let _ = peer_tx.send(e);
                                    }
                                    Some(Payload::RelayToken(t)) => {
                                        tracing::info!(relay = %t.relay_host, port = t.relay_port, session = %t.session_id, "signal recv relay_token");
                                        let _ = relay_tx.send(t);
                                    }
                                    _ => {}
                                }
                            }
                            Some(Ok(Message::Ping(p))) => {
                                let _ = sink.send(Message::Pong(p)).await;
                            }
                            Some(Ok(Message::Pong(_))) => {}
                            Some(Ok(Message::Close(frame))) => {
                                tracing::debug!(?frame, "signal ws peer close");
                                break;
                            }
                            None => {
                                tracing::debug!("signal ws stream ended");
                                break;
                            }
                            Some(Ok(_)) => {}
                            Some(Err(e)) => {
                                tracing::debug!(%e, "signal ws read error");
                                break;
                            }
                        }
                    }
                    _ = ping.tick() => {
                        if sink.send(Message::Ping(vec![].into())).await.is_err() {
                            tracing::debug!("signal ws ping failed");
                            break;
                        }
                    }
                }
            }
            alive_task.store(false, Ordering::Release);
        });
        Ok(Self {
            out_tx,
            in_rx: Arc::new(tokio::sync::Mutex::new(in_rx)),
            peer_rx: Arc::new(tokio::sync::Mutex::new(peer_rx)),
            relay_rx: Arc::new(tokio::sync::Mutex::new(relay_rx)),
            alive,
            ws_task: Mutex::new(Some(ws_task)),
        })
    }

    pub fn signal_alive(&self) -> bool {
        self.alive.load(Ordering::Acquire)
    }

    /// 结束信令 WS 任务（重连前必须调用，避免僵尸连接）。
    pub fn shutdown(&self) {
        self.alive.store(false, Ordering::Release);
        if let Ok(mut g) = self.ws_task.lock() {
            if let Some(h) = g.take() {
                h.abort();
            }
        }
    }

    pub async fn register(&self, room: &str, role: &str) -> Result<RegisterResp, String> {
        tracing::debug!(%room, %role, "signal register send");
        let env = SignalEnvelope {
            payload: Some(Payload::Register(RegisterReq {
                room_code: room.into(),
                role: role.into(),
                caps: Some(DeviceCapabilities {
                    platform: std::env::consts::OS.to_string(),
                    ..Default::default()
                }),
            })),
            ..Default::default()
        };
        self.send(env).await?;
        let resp = self.recv_register().await?;
        Ok(resp)
    }

    pub async fn send(&self, env: SignalEnvelope) -> Result<(), String> {
        self.out_tx.send(env.encode_to_vec()).map_err(|_| "signal closed".into())
    }

    pub async fn send_bytes(&self, data: Vec<u8>) -> Result<(), String> {
        self.out_tx.send(data).map_err(|_| "signal closed".into())
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

    /// Viewer 发来 peer_endpoints 后再建直连（勿与 relay_token 竞态；relay 在 direct/ice 失败后再取）。
    pub async fn wait_viewer_peer_endpoints(&self) -> Result<PeerEndpoints, String> {
        tracing::info!("agent waiting for viewer peer_endpoints");
        let alive = self.alive.clone();
        tokio::select! {
            eps = async {
                self.peer_rx.lock().await.recv().await
            } => {
                let eps = eps.ok_or_else(|| "signal closed".to_string())?;
                tracing::info!(n = eps.items.len(), "agent viewer peer_endpoints ready");
                Ok(eps)
            }
            () = async {
                Self::wait_until_dead(alive).await;
            } => Err("signal closed".into()),
        }
    }

    async fn wait_until_dead(alive: Arc<AtomicBool>) {
        while alive.load(Ordering::Acquire) {
            tokio::time::sleep(Duration::from_millis(200)).await;
        }
    }

    async fn recv_register(&self) -> Result<RegisterResp, String> {
        loop {
            let env = self.recv_envelope().await?;
            if let Some(Payload::Error(e)) = env.payload {
                tracing::warn!(%e, "signal register error");
                return Err(e);
            }
            if let Some(Payload::RegisterAck(a)) = env.payload {
                tracing::info!(peer_id = %a.peer_id, room_id = %a.room_id, "signal register ok");
                return Ok(a);
            }
        }
    }
}

