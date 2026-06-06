//! P2P 数据面：TCP / ICE 上 HTTP/2 + protobuf body（同端口 `/rb/v1/media` Range）。
use tokio::io::{AsyncRead, AsyncWrite};
use tokio::net::TcpStream;

#[cfg(feature = "agent-serve")]
use hyper_util::rt::{TokioExecutor, TokioIo};
#[cfg(feature = "agent-serve")]
use hyper_util::server::conn::auto::Builder as AutoConnBuilder;

use crate::rb_http_client::RbHttpClient;

#[cfg(feature = "agent-serve")]
use crate::remotebrowse::http_api;
#[cfg(feature = "agent-serve")]
use crate::remotebrowse_svc::MediaSvc;

pub async fn connect_h2(sock: TcpStream) -> Result<RbHttpClient, String> {
    connect_h2_io(sock).await
}

pub async fn connect_h2_io<I>(io: I) -> Result<RbHttpClient, String>
where
    I: AsyncRead + AsyncWrite + Unpin + Send + 'static,
{
    RbHttpClient::connect_io(io).await
}

#[cfg(feature = "agent-serve")]
pub async fn serve_h2(sock: TcpStream, sandbox: Option<String>, room: Option<String>) -> Result<(), String> {
    serve_h2_io(sock, sandbox, room).await
}

#[cfg(feature = "agent-serve")]
pub async fn serve_h2_io<I>(io: I, sandbox: Option<String>, room: Option<String>) -> Result<(), String>
where
    I: AsyncRead + AsyncWrite + Unpin + Send + 'static,
{
    tracing::info!(
        sandbox = sandbox.as_deref().unwrap_or(""),
        room = room.as_deref().unwrap_or(""),
        "p2p http2 serve"
    );
    let media = MediaSvc::new(sandbox, room);
    let mut app = http_api::router(media.clone());
    #[cfg(feature = "media")]
    {
        app = app.merge(crate::remotebrowse::media_http::router(media));
    }
    app = app.layer(crate::remotebrowse::http_log::layer());
    let svc = hyper_util::service::TowerToHyperService::new(app);
    AutoConnBuilder::new(TokioExecutor::new())
        .serve_connection(TokioIo::new(io), svc)
        .await
        .map_err(|e| e.to_string())?;
    tracing::info!("p2p http2 closed");
    Ok(())
}
