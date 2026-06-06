//! 单路 HTTP/2 客户端：POST + `application/protobuf`（ICE duplex / 本机调试）；载荷可选 zstd。
use http::{Method, Request};
use http_body_util::{BodyExt, Full};
use hyper::body::Incoming;
use hyper_util::rt::{TokioExecutor, TokioIo};
use prost::Message;
use tokio::io::{AsyncRead, AsyncWrite};

use crate::proto_zstd;

type BoxBody = Full<bytes::Bytes>;

const RB_PROTO: &str = "application/protobuf";

pub struct RbHttpClient {
    send: hyper::client::conn::http2::SendRequest<BoxBody>,
}

impl RbHttpClient {
    pub async fn connect_io<I>(io: I) -> Result<Self, String>
    where
        I: AsyncRead + AsyncWrite + Unpin + Send + 'static,
    {
        let io = TokioIo::new(io);
        let (send, conn) = hyper::client::conn::http2::Builder::new(TokioExecutor::new())
            .handshake(io)
            .await
            .map_err(|e| e.to_string())?;
        tokio::spawn(async move {
            if let Err(e) = conn.await {
                tracing::debug!("rb http2 client conn end: {e}");
            }
        });
        let mut send = send;
        send.ready().await.map_err(|e| e.to_string())?;
        Ok(Self { send })
    }

    pub fn fork(&self) -> Self {
        Self { send: self.send.clone() }
    }

    pub async fn post_proto<Req, Resp>(&mut self, path: &str, req: &Req) -> Result<Resp, String>
    where
        Req: Message,
        Resp: Message + Default,
    {
        let plain = req.encode_to_vec();
        let (body, zstd_req) = proto_zstd::encode_request(&plain);
        let mut builder = Request::builder()
            .method(Method::POST)
            .uri(path)
            .header(http::header::CONTENT_TYPE, RB_PROTO)
            .header(http::header::ACCEPT, RB_PROTO)
            .header(http::header::ACCEPT_ENCODING, proto_zstd::ZSTD);
        if zstd_req {
            builder = builder.header(http::header::CONTENT_ENCODING, proto_zstd::ZSTD);
        }
        let http_req = builder
            .body(Full::new(bytes::Bytes::from(body)))
            .map_err(|e| e.to_string())?;
        let resp = self.send.send_request(http_req).await.map_err(|e| e.to_string())?;
        let status = resp.status();
        let headers = resp.headers().clone();
        let body = read_body(resp.into_body()).await?;
        if !status.is_success() {
            return Err(format!("http {status}: {}", String::from_utf8_lossy(&body)));
        }
        let raw = proto_zstd::decode_http_body(&headers, &body)?;
        Resp::decode(raw.as_slice()).map_err(|e| e.to_string())
    }

    pub async fn get_proto<Resp>(&mut self, path: &str) -> Result<Resp, String>
    where
        Resp: Message + Default,
    {
        let http_req = Request::builder()
            .method(Method::GET)
            .uri(path)
            .header(http::header::ACCEPT, RB_PROTO)
            .header(http::header::ACCEPT_ENCODING, proto_zstd::ZSTD)
            .body(Full::new(bytes::Bytes::new()))
            .map_err(|e| e.to_string())?;
        let resp = self.send.send_request(http_req).await.map_err(|e| e.to_string())?;
        let status = resp.status();
        let headers = resp.headers().clone();
        let body = read_body(resp.into_body()).await?;
        if !status.is_success() {
            return Err(format!("http {status}: {}", String::from_utf8_lossy(&body)));
        }
        let raw = proto_zstd::decode_http_body(&headers, &body)?;
        Resp::decode(raw.as_slice()).map_err(|e| e.to_string())
    }
}

async fn read_body(body: Incoming) -> Result<Vec<u8>, String> {
    let collected = body.collect().await.map_err(|e| e.to_string())?;
    Ok(collected.to_bytes().to_vec())
}
