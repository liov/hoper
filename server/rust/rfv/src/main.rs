mod file;

use axum::{
    routing::get,
    Router,
};
use tower::ServiceBuilder;
use tower_http::trace::TraceLayer;
use tracing_subscriber;
use ffmpeg_next as ffmpeg;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();
    ffmpeg::init().unwrap();
    let grpc_addr = std::env::var("RFV_GRPC_ADDR").unwrap_or_else(|_| "0.0.0.0:50051".into());
    tokio::spawn(async move {
        if let Err(e) = rfv::grpc_server::serve(grpc_addr).await {
            tracing::error!("rfv grpc: {e}");
        }
    });
    let app = Router::new()
        .route("/", get(|| async { "Hello, World!" }))
        .route("/list_files", get(file::list_files_handler))
        .route("/thumbnail/*path", get(file::file_thumbnail_handler))
        .layer(ServiceBuilder::new().layer(TraceLayer::new_for_http()));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app.into_make_service()).await.unwrap();
}
