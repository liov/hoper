use std::sync::Arc;
use axum::{
    routing::{get, post},
    Router,
};
use tracing_subscriber;
use tracing::{info, error};

use message::{WebSocketManager, websocket_handler, get_online_users};

#[tokio::main]
async fn main() {
    // 初始化日志
    tracing_subscriber::fmt::init();
    
    // 创建WebSocket管理器 
    let ws_manager = Arc::new(WebSocketManager::new());
    
    // 构建路由
    let app = Router::new()
        .route("/ws", get(websocket_handler))
        .route("/online", get(get_online_users))
        .with_state(ws_manager);
    
    let addr = "0.0.0.0:3000";
    info!("🚀 Starting WebSocket server on {}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app.into_make_service()).await.unwrap();

}
