//! WebSocket路由处理器
//! 处理WebSocket连接升级和认证

use axum::{
    extract::{ws::{WebSocketUpgrade, WebSocket}, State, Query},
    response::IntoResponse,
    http::StatusCode,
};
use serde::Deserialize;
use std::sync::Arc;
use axum::body::Bytes;
use tracing::{info, warn, error};

use crate::websocket_manager::{WebSocketManager};

/// WebSocket连接参数
#[derive(Deserialize)]
pub struct WsQueryParams {
    /// 用户ID
    pub user_id: u64,
    /// 可选的认证令牌
    #[serde(default)]
    pub token: Option<String>,
}

/// WebSocket连接处理器
pub async fn websocket_handler(
    ws: WebSocketUpgrade,
    Query(params): Query<WsQueryParams>,
    State(manager): State<Arc<WebSocketManager>>,
) -> impl IntoResponse {
    info!("WebSocket connection attempt from user {}", params.user_id);
    
    // 这里可以添加认证逻辑
    if let Some(token) = params.token {
        if !validate_token(&token, params.user_id).await {
            warn!("Authentication failed for user {}", params.user_id);
            return StatusCode::UNAUTHORIZED.into_response();
        }
    }
    
    // 升级连接
    ws.on_upgrade(move |socket| handle_websocket(socket, manager, params.user_id))
}

/// 处理已升级的WebSocket连接
async fn handle_websocket(socket: WebSocket, manager: Arc<WebSocketManager>, user_id: u64) {
    info!("WebSocket connection established for user {}", user_id);
    
    match manager.add_connection(user_id, socket).await {
        Ok(()) => {
            info!("User {} successfully connected", user_id);
        }
        Err(e) => {
            error!("Failed to add connection for user {}: {}", user_id, e);
        }
    }
}

/// 简单的令牌验证函数
/// 实际应用中应该替换为真正的认证逻辑
async fn validate_token(token: &str, user_id: u64) -> bool {
    // 示例验证逻辑：检查令牌是否以"user_{id}_"开头
    token.starts_with(&format!("user_{}_", user_id))
}


/// 获取在线用户列表接口
pub async fn get_online_users(
    State(manager): State<Arc<WebSocketManager>>,
) -> impl IntoResponse {
    let users = manager.get_online_users().await;
    let json = serde_json::json!({
        "online_users": users,
        "count": users.len()
    });
    axum::Json(json)
}