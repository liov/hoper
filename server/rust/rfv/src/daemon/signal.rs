use std::collections::HashMap;
use std::sync::{Arc, Mutex};

use axum::extract::ws::{Message, WebSocket, WebSocketUpgrade};
use axum::extract::State;
use axum::response::IntoResponse;
use axum::Json;
use futures_util::{SinkExt, StreamExt};
use prost::Message as ProstMessage;
use serde::Serialize;
use tokio::sync::mpsc;
use uuid::Uuid;

use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::{RegisterReq, RegisterResp, RelayToken, SignalEnvelope};

const ROLE_VIEWER: &str = "viewer";
const ROLE_AGENT: &str = "agent";
const MAX_PENDING: usize = 128;

struct Peer {
    id: String,
    tx: mpsc::UnboundedSender<Vec<u8>>,
}

struct Room {
    viewer: Option<Peer>,
    agent: Option<Peer>,
    pend_to_viewer: Vec<Vec<u8>>,
    pend_to_agent: Vec<Vec<u8>>,
}

#[derive(Clone)]
pub struct Hub {
    inner: Arc<HubInner>,
}

struct HubInner {
    rooms: Mutex<HashMap<String, Room>>,
    relay_tcp_addr: String,
}

impl Hub {
    pub fn new(relay: std::net::SocketAddr) -> Self {
        Self {
            inner: Arc::new(HubInner {
                rooms: Mutex::new(HashMap::new()),
                relay_tcp_addr: relay.to_string(),
            }),
        }
    }

    pub fn relay_tcp_addr(&self) -> &str {
        &self.inner.relay_tcp_addr
    }

    fn remove_peer_from_room(&self, code: &str, peer_id: &str) {
        let mut g = self.inner.rooms.lock().expect("rooms");
        let Some(room) = g.get_mut(code) else {
            return;
        };
        if room.viewer.as_ref().is_some_and(|p| p.id == peer_id) {
            room.viewer = None;
        }
        if room.agent.as_ref().is_some_and(|p| p.id == peer_id) {
            room.agent = None;
        }
        if room.viewer.is_none() && room.agent.is_none() {
            g.remove(code);
        }
    }

    fn on_register(&self, tx: &mpsc::UnboundedSender<Vec<u8>>, req: RegisterReq, room_code: &mut String, peer_id: &mut String, role: &mut String) {
        if req.room_code.is_empty() || (req.role != ROLE_VIEWER && req.role != ROLE_AGENT) {
            let _ = send_proto(tx, error_env("bad register"));
            return;
        }
        *room_code = req.room_code.clone();
        *role = req.role.clone();
        *peer_id = Uuid::new_v4().to_string();
        let mut g = self.inner.rooms.lock().expect("rooms");
        let room = g.entry(room_code.clone()).or_default();
        if *role == ROLE_VIEWER {
            if room.viewer.is_some() {
                let _ = send_proto(tx, error_env("viewer busy"));
                return;
            }
            room.viewer = Some(Peer { id: peer_id.clone(), tx: tx.clone() });
        } else if room.agent.is_some() {
            let _ = send_proto(tx, error_env("agent busy"));
            return;
        } else {
            room.agent = Some(Peer { id: peer_id.clone(), tx: tx.clone() });
        }
        let (viewer, agent) = (room.viewer.as_ref(), room.agent.as_ref());
        if *role == ROLE_VIEWER {
            flush_pending(tx, &mut room.pend_to_viewer);
        } else {
            flush_pending(tx, &mut room.pend_to_agent);
        }
        let _ = send_proto(
            tx,
            SignalEnvelope {
                payload: Some(Payload::RegisterAck(RegisterResp {
                    peer_id: peer_id.clone(),
                    room_id: room_code.clone(),
                })),
                ..Default::default()
            },
        );
        if viewer.is_some() && agent.is_some() && !self.inner.relay_tcp_addr.is_empty() {
            let (host, port) = parse_relay_addr(&self.inner.relay_tcp_addr);
            let tok = SignalEnvelope {
                payload: Some(Payload::RelayToken(RelayToken {
                    session_id: Uuid::new_v4().to_string(),
                    relay_host: host,
                    relay_port: port,
                    psk: Vec::new(),
                })),
                ..Default::default()
            };
            if let Some(v) = room.viewer.as_ref() {
                let _ = send_proto(&v.tx, tok.clone());
            }
            if let Some(a) = room.agent.as_ref() {
                let _ = send_proto(&a.tx, tok);
            }
        }
    }

    fn forward_ice(&self, room_code: &str, role: &str, env: SignalEnvelope) {
        let b = env.encode_to_vec();
        if room_code.is_empty() || role.is_empty() {
            return;
        }
        let other = if role == ROLE_AGENT { ROLE_VIEWER } else { ROLE_AGENT };
        let mut g = self.inner.rooms.lock().expect("rooms");
        let Some(room) = g.get_mut(room_code) else {
            return;
        };
        let target = if other == ROLE_VIEWER { room.viewer.as_ref() } else { room.agent.as_ref() };
        if let Some(peer) = target {
            let _ = peer.tx.send(b);
            return;
        }
        if other == ROLE_VIEWER {
            append_pending(&mut room.pend_to_viewer, b);
        } else {
            append_pending(&mut room.pend_to_agent, b);
        }
    }
}

impl Default for Room {
    fn default() -> Self {
        Self { viewer: None, agent: None, pend_to_viewer: Vec::new(), pend_to_agent: Vec::new() }
    }
}

pub async fn ws_handler(ws: WebSocketUpgrade, State(hub): State<Hub>) -> impl IntoResponse {
    ws.on_upgrade(move |socket| handle_socket(socket, hub))
}

#[derive(Serialize)]
pub struct Health {
    #[serde(rename = "signalWs")]
    signal_ws: String,
    #[serde(rename = "relayTcp")]
    relay_tcp: String,
    #[serde(rename = "rfvGrpc")]
    rfv_grpc: String,
    #[serde(rename = "thumbCache")]
    thumb_cache: String,
}

pub async fn health(State(hub): State<Hub>) -> Json<Health> {
    Json(Health {
        signal_ws: "/rb/signal".into(),
        relay_tcp: hub.relay_tcp_addr().into(),
        rfv_grpc: std::env::var("RFV_GRPC_ADDR").unwrap_or_else(|_| "127.0.0.1:50051".into()),
        thumb_cache: String::new(),
    })
}

async fn handle_socket(socket: WebSocket, hub: Hub) {
    let (mut sink, mut stream) = socket.split();
    let (tx, mut rx) = mpsc::unbounded_channel::<Vec<u8>>();
    tokio::spawn(async move {
        while let Some(buf) = rx.recv().await {
            if sink.send(Message::Binary(buf.into())).await.is_err() {
                break;
            }
        }
    });
    let mut room_code = String::new();
    let mut peer_id = String::new();
    let mut role = String::new();
    while let Some(Ok(msg)) = stream.next().await {
        let Message::Binary(data) = msg else { continue };
        let Ok(env) = SignalEnvelope::decode(data.as_ref()) else { continue };
        match env.payload {
            Some(Payload::Register(req)) => hub.on_register(&tx, req, &mut room_code, &mut peer_id, &mut role),
            Some(Payload::IceParameters(_)) | Some(Payload::IceCandidate(_)) | Some(Payload::IceComplete(_)) | Some(Payload::PeerEndpoints(_)) => {
                hub.forward_ice(&room_code, &role, env);
            }
            _ => {}
        }
    }
    if !room_code.is_empty() && !peer_id.is_empty() {
        hub.remove_peer_from_room(&room_code, &peer_id);
    }
}

fn append_pending(pending: &mut Vec<Vec<u8>>, b: Vec<u8>) {
    if pending.len() >= MAX_PENDING {
        pending.remove(0);
    }
    pending.push(b);
}

fn flush_pending(tx: &mpsc::UnboundedSender<Vec<u8>>, pending: &mut Vec<Vec<u8>>) {
    for b in pending.drain(..) {
        let _ = tx.send(b);
    }
}

fn send_proto(tx: &mpsc::UnboundedSender<Vec<u8>>, env: SignalEnvelope) -> Result<(), ()> {
    tx.send(env.encode_to_vec()).map_err(|_| ())
}

fn error_env(msg: &str) -> SignalEnvelope {
    SignalEnvelope {
        payload: Some(Payload::Error(msg.to_string())),
        ..Default::default()
    }
}

fn parse_relay_addr(addr: &str) -> (String, u32) {
    match addr.rsplit_once(':') {
        Some((host, port)) => match port.parse::<u32>() {
            Ok(p) => (host.to_string(), p),
            Err(_) => (host.to_string(), 19090),
        },
        None => (addr.to_string(), 19090),
    }
}
