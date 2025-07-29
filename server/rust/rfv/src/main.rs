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
    // build our application with a single route
    let app = Router::new().route("/", get(|| async { "Hello, World!" })).
        route("/list_files", get(file::list_files_handler)).
        route("/thumbnail/*path", get(file::file_thumbnail_handler)).
        layer(
        ServiceBuilder::new()
            .layer(TraceLayer::new_for_http())
        // 可以添加其他中间件
    );

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app.into_make_service()).await.unwrap();
}