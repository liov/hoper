//! ICE（UDP）上可靠分块字节流 + HTTP/2 protobuf。UDP 每包 `[4B len][payload]`，多路复用由 HTTP/2 提供。
use std::sync::Arc;

use tokio::io::{AsyncReadExt, AsyncWriteExt};
use webrtc_util::conn::Conn;

use crate::client::ice_stream::IceWire;
use crate::p2p_http;
use crate::rb_http_client::RbHttpClient;

const MAX_CHUNK: usize = 256 * 1024;

#[cfg(feature = "agent-serve")]
pub async fn serve(ice: IceWire, sandbox: Option<String>, room: Option<String>) -> Result<(), String> {
    tracing::info!(sandbox = sandbox.as_deref().unwrap_or(""), "ice link http2 serve");
    let (grpc_end, bridge_end) = tokio::io::duplex(4 << 20);
    let conn = ice.into_conn();
    tokio::spawn(bridge(conn, bridge_end));
    p2p_http::serve_h2_io(grpc_end, sandbox, room).await
}

pub async fn connect(ice: IceWire) -> Result<RbHttpClient, String> {
    tracing::info!("ice link http2 client");
    let (grpc_end, bridge_end) = tokio::io::duplex(4 << 20);
    let conn = ice.into_conn();
    tokio::spawn(bridge(conn, bridge_end));
    p2p_http::connect_h2_io(grpc_end).await
}

async fn bridge(conn: Arc<dyn Conn + Send + Sync>, duplex: tokio::io::DuplexStream) {
    let (mut dup_rd, mut dup_wr) = tokio::io::split(duplex);
    let conn2 = conn.clone();
    let ice_to_dup = async move {
        let mut buf = vec![0u8; 2048];
        loop {
            let n = match conn.recv(&mut buf).await {
                Ok(n) if n > 0 => n,
                Ok(_) => continue,
                Err(_) => break,
            };
            if dup_wr.write_all(&(n as u32).to_be_bytes()).await.is_err() {
                break;
            }
            if dup_wr.write_all(&buf[..n]).await.is_err() {
                break;
            }
        }
    };
    let dup_to_ice = async move {
        loop {
            let mut sz = [0u8; 4];
            if dup_rd.read_exact(&mut sz).await.is_err() {
                break;
            }
            let n = u32::from_be_bytes(sz) as usize;
            if n > MAX_CHUNK {
                break;
            }
            let mut payload = vec![0u8; n];
            if n > 0 && dup_rd.read_exact(&mut payload).await.is_err() {
                break;
            }
            if conn2.send(&payload).await.is_err() {
                break;
            }
        }
    };
    let _ = tokio::join!(ice_to_dup, dup_to_ice);
}
