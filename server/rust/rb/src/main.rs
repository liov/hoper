//! 被浏览主机：`rb <房间码>` 起 HTTP 服务 + 子进程 Agent；`rb <房间码> --agent-only` 仅跑信令/P2P。
use std::env;
use std::time::Duration;

use ffmpeg_next as ffmpeg;

#[tokio::main]
async fn main() {
    rb::tracing_init::init();
    ffmpeg::init().expect("ffmpeg init");
    if agent_only_mode() {
        run_agent_only(room_from_argv()).await;
        return;
    }
    run_host(room_from_argv()).await;
}

fn agent_only_mode() -> bool {
    env::args().any(|a| a == "--agent-only")
}

/// 房间码：`RB_ROOM` 或第一个非 flag 参数（如 `rb demo`）。
fn room_from_argv() -> Option<String> {
    if let Ok(s) = env::var("RB_ROOM") {
        let s = s.trim().to_string();
        if !s.is_empty() {
            return Some(s);
        }
    }
    for arg in env::args().skip(1) {
        if arg.starts_with('-') {
            continue;
        }
        return Some(arg);
    }
    None
}

async fn run_host(room: Option<String>) {
    let listen = env::var("RB_LISTEN")
        .or_else(|_| env::var("RB_GRPC_ADDR"))
        .unwrap_or_else(|_| "0.0.0.0:50051".into());
    tracing::info!(%listen, room = room.as_deref().unwrap_or(""), "rb host starting");
    if let Some(room) = room {
        spawn_agent_child_supervisor(room);
    }
    rb::grpc_server::serve_forever(listen).await;
}

#[cfg(feature = "transport")]
async fn run_agent_only(room: Option<String>) {
    let Some(room) = room else {
        tracing::error!("--agent-only 需要房间码，例如: rb demo --agent-only");
        return;
    };
    let sandbox = None::<String>;
    let signal_url =
        env::var("RB_SIGNAL_URL").unwrap_or_else(|_| "ws://127.0.0.1:8080/rb/signal".into());
    let ice_ms: u32 = env::var("RB_ICE_TIMEOUT_MS")
        .ok()
        .and_then(|s| s.parse().ok())
        .unwrap_or(15000);
    tracing::info!(%room, %signal_url, ice_ms, sandbox = sandbox.as_deref().unwrap_or(""), "rb agent-only");
    let mut backoff = Duration::from_secs(1);
    loop {
        rb::transport::run_agent(signal_url.clone(), room.clone(), sandbox.clone(), ice_ms).await;
        tracing::warn!(
            wait_secs = backoff.as_secs(),
            "agent-only returned, reconnect"
        );
        tokio::time::sleep(backoff).await;
        backoff = (backoff * 2).min(Duration::from_secs(30));
    }
}

#[cfg(not(feature = "transport"))]
async fn run_agent_only(_room: Option<String>) {}

/// 信令/P2P 在子进程：转码 segfault 不会拖死本机 HTTP 服务，子进程退出后自动拉起。
#[cfg(feature = "transport")]
fn spawn_agent_child_supervisor(room: String) {
    tokio::spawn(async move {
        let exe = match env::current_exe() {
            Ok(p) => p,
            Err(e) => {
                tracing::error!(%e, "agent child: current_exe failed");
                return;
            }
        };
        let mut backoff = Duration::from_secs(1);
        loop {
            tracing::info!(%room, "agent child spawn");
            let mut cmd = tokio::process::Command::new(&exe);
            cmd.arg(&room).arg("--agent-only");
            match cmd.status().await {
                Ok(st) => {
                    tracing::warn!(%st, wait_secs = backoff.as_secs(), "agent child exited, respawn")
                }
                Err(e) => tracing::error!(%e, "agent child wait failed"),
            }
            tokio::time::sleep(backoff).await;
            backoff = (backoff * 2).min(Duration::from_secs(30));
        }
    });
}

#[cfg(not(feature = "transport"))]
fn spawn_agent_child_supervisor(_room: String) {}
