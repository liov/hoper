//! WebSocket连接管理模块
//! 提供用户连接管理、消息广播、定向发送等功能

use futures_util::stream::StreamExt;
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{broadcast, RwLock};
use axum::extract::ws::{Message, WebSocket};
use futures_util::SinkExt;
use tokio_util::bytes::Bytes;
use tracing::{debug, info, warn, error};

/// WebSocket会话信息
#[derive(Debug, Clone)]
pub struct WsSession {
    /// 用户ID
    pub user_id: u64,
    /// 发送消息的通道
    pub tx: broadcast::Sender<Bytes>,
    /// 连接时间戳
    pub connected_at: std::time::SystemTime,
}

/// WebSocket管理器
#[derive(Debug)]
pub struct WebSocketManager {
    /// 全局消息广播通道
    global_tx: broadcast::Sender<Bytes>,
    /// 用户会话映射 (user_id -> session)
    sessions: Arc<RwLock<HashMap<u64, WsSession>>>,
}

impl WebSocketManager {
    /// 创建新的WebSocket管理器
    pub fn new() -> Self {
        let (global_tx, _) = broadcast::channel(100);

        Self {
            global_tx,
            sessions: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    /// 获取当前连接数
    pub async fn connection_count(&self) -> usize {
        self.sessions.read().await.len()
    }
    

    /// 添加新连接
    pub async fn add_connection(&self, user_id: u64, socket: WebSocket) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {

        // 检查用户是否已连接
        {
            let sessions = self.sessions.read().await;
            if sessions.contains_key(&user_id) {
                return Err(format!("User {} already connected", user_id).into());
            }
        }

        // 创建用户专用通道
        let (user_tx, _) = broadcast::channel(100);

        // 创建会话
        let session = WsSession {
            user_id,
            tx: user_tx.clone(),
            connected_at: std::time::SystemTime::now(),
        };

        // 添加到会话列表
        {
            let mut sessions = self.sessions.write().await;
            sessions.insert(user_id, session);
        }

        info!("User {} connected, current connections: {}", user_id, self.connection_count().await);

        // 启动连接处理任务
        self.handle_connection(user_id, socket, user_tx).await;

        Ok(())
    }

    /// 移除连接
    pub async fn remove_connection(&self, user_id: u64) {
        let mut sessions = self.sessions.write().await;
        if sessions.remove(&user_id).is_some() {
            info!("User {} disconnected, current connections: {}", user_id, sessions.len());
        }
    }

    /// 处理单个连接
    async fn handle_connection(&self, user_id: u64, socket: WebSocket, user_tx: broadcast::Sender<Bytes>) {
        let manager = self.clone();
        let global_tx = self.global_tx.clone();

        tokio::spawn(async move {
            // 订阅全局消息
            let mut global_rx = global_tx.subscribe();

            // 分离读写
            let (mut sender, mut receiver) = socket.split();

            // 发送任务：转发全局消息到客户端
            let send_task = tokio::spawn(async move {
                loop {
                    match global_rx.recv().await {
                        Ok(msg) => {
                            if let Err(e) = sender.send(Message::Binary(msg)).await {
                                debug!("Failed to send message to user {}: {}", user_id, e);
                                break;
                            }
                        }
                        Err(broadcast::error::RecvError::Closed) => {
                            debug!("Global channel closed for user {}", user_id);
                            break;
                        }
                        Err(broadcast::error::RecvError::Lagged(_)) => {
                            warn!("Message lagged for user {}", user_id);
                            continue;
                        }
                    }
                }
            });

            // 接收任务：处理客户端消息
            let recv_task = tokio::spawn(async move {
                while let Some(result) = receiver.next().await {
                    match result {
                        Ok(Message::Binary(bin)) => {
                            debug!("Received message from user {}: {:?}", user_id, bin);
                            // 广播消息到所有用户
                            if let Err(e) = global_tx.send(bin) {
                                error!("Failed to broadcast message: {}", e);
                            }
                        }
                        Ok(Message::Close(frame)) => {
                            debug!("User {} sent close frame: {:?}", user_id, frame);
                            break;
                        }
                        Ok(_) => {
                            // 忽略非文本消息
                            continue;
                        }
                        Err(e) => {
                            warn!("Error receiving message from user {}: {}", user_id, e);
                            break;
                        }
                    }
                }
            });

            // 等待任一任务完成
            tokio::select! {
                _ = send_task => debug!("Send task completed for user {}", user_id),
                _ = recv_task => debug!("Receive task completed for user {}", user_id),
            }

            // 清理连接
            manager.remove_connection(user_id).await;
        });
    }

    /// 广播消息给所有连接的用户
    pub async fn broadcast(&self, message: Bytes) -> Result<usize, broadcast::error::SendError<Bytes>> {
        let count = self.global_tx.send(message)?;
        debug!("Broadcast message to {} users", count);
        Ok(count)
    }

    /// 向特定用户发送消息
    pub async fn send_to_user(&self, user_id: u64, message: Bytes) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
        let sessions = self.sessions.read().await;
        if let Some(session) = sessions.get(&user_id) {
            session.tx.send(message)?;
            Ok(())
        } else {
            Err(format!("User {} not found", user_id).into())
        }
    }

    /// 获取在线用户列表
    pub async fn get_online_users(&self) -> Vec<u64> {
        let sessions = self.sessions.read().await;
        sessions.keys().cloned().collect()
    }

    /// 获取用户会话信息
    pub async fn get_user_session(&self, user_id: u64) -> Option<WsSession> {
        let sessions = self.sessions.read().await;
        sessions.get(&user_id).cloned()
    }

    /// 检查用户是否在线
    pub async fn is_user_online(&self, user_id: u64) -> bool {
        let sessions = self.sessions.read().await;
        sessions.contains_key(&user_id)
    }
}

// 手动实现Clone trait，因为RwLock不能自动派生Clone
impl Clone for WebSocketManager {
    fn clone(&self) -> Self {
        Self {
            global_tx: self.global_tx.clone(),
            sessions: Arc::clone(&self.sessions),
        }
    }
}

// 为WebSocket实现next方法替代futures_util
trait WebSocketExt {
    async fn next(&mut self) -> Option<Result<Message, axum::Error>>;
}

impl WebSocketExt for WebSocket {
    async fn next(&mut self) -> Option<Result<Message, axum::Error>> {
        self.recv().await
    }
}