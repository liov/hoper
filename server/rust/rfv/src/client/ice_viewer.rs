use std::collections::VecDeque;
use std::sync::{Arc, Mutex};
use std::time::Duration;

use prost::Message;
use rand::RngCore;
use tokio::runtime::Runtime;
use tokio::sync::mpsc;
use webrtc_ice::agent::agent_config::AgentConfig;
use webrtc_ice::candidate::candidate_base::unmarshal_candidate;
use webrtc_ice::candidate::Candidate;
use webrtc_ice::network_type::NetworkType;
use webrtc_ice::url::Url;
use webrtc_ice::agent::Agent;

use crate::client::ice_stream::IceWire;
use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::{IceCandidateInit, IceParameters, SignalEnvelope};
use crate::client::wire;

const STUN: &str = "stun:stun.l.google.com:19302";

pub struct ViewerHandle {
    runtime: Runtime,
    inner: Arc<Mutex<Inner>>,
}

struct Inner {
    state: ViewerState,
    outbox: VecDeque<Vec<u8>>,
    inbox: VecDeque<Vec<u8>>,
}

enum ViewerState {
    Running,
    Ready(IceWire),
    Failed(String),
}

impl ViewerHandle {
    pub fn new(timeout_ms: u32) -> Self {
        let runtime = Runtime::new().expect("tokio runtime");
        let inner = Arc::new(Mutex::new(Inner {
            state: ViewerState::Running,
            outbox: VecDeque::new(),
            inbox: VecDeque::new(),
        }));
        let bg = inner.clone();
        let fail = inner.clone();
        runtime.spawn(async move {
            if let Err(e) = run_viewer(bg, Duration::from_millis(timeout_ms as u64)).await {
                fail.lock().expect("lock").state = ViewerState::Failed(e);
            }
        });
        Self { runtime, inner }
    }

    pub fn push(&self, data: &[u8]) {
        self.inner.lock().expect("lock").inbox.push_back(data.to_vec());
    }

    pub fn poll_out(&self) -> Option<Vec<u8>> {
        self.inner.lock().expect("lock").outbox.pop_front()
    }

    pub fn state_code(&self) -> i32 {
        match &self.inner.lock().expect("lock").state {
            ViewerState::Running => 0,
            ViewerState::Ready(_) => 1,
            ViewerState::Failed(_) => -1,
        }
    }

    pub fn read_frame(&self, buf: &mut [u8]) -> Result<(u8, usize), &'static str> {
        let mut g = self.inner.lock().expect("lock");
        let ViewerState::Ready(ref mut ice) = g.state else {
            return Err("not ready");
        };
        let (typ, payload) = self.runtime.block_on(ice.read_frame()).map_err(|_| "read failed")?;
        let frame = wire::encode_frame(typ, &payload);
        if frame.len() > buf.len() {
            return Err("short buf");
        }
        buf[..frame.len()].copy_from_slice(&frame);
        Ok((typ, frame.len()))
    }

    pub fn write_frame(&self, typ: u8, payload: &[u8]) -> Result<(), &'static str> {
        let mut g = self.inner.lock().expect("lock");
        let ViewerState::Ready(ref mut ice) = g.state else {
            return Err("not ready");
        };
        self.runtime
            .block_on(ice.write_frame(typ, payload))
            .map_err(|_| "write failed")
    }
}

async fn run_viewer(inner: Arc<Mutex<Inner>>, timeout: Duration) -> Result<(), String> {
    let (ufrag, pwd) = new_credentials();
    let agent = Arc::new(build_agent(&ufrag, &pwd).await?);
    queue_out(&inner, ice_parameters_env(&ufrag, &pwd));
    bind_candidates(&agent, inner.clone());
    agent.gather_candidates().map_err(|e| e.to_string())?;
    let (remote_ufrag, remote_pwd) = wait_remote_credentials(&inner, &agent, timeout).await?;
    agent
        .set_remote_credentials(remote_ufrag.clone(), remote_pwd.clone())
        .await
        .map_err(|e| e.to_string())?;
    let (cancel_tx, cancel_rx) = mpsc::channel::<()>(1);
    let dial = agent.dial(cancel_rx, remote_ufrag, remote_pwd);
    let conn = tokio::time::timeout(timeout, dial)
        .await
        .map_err(|_| "ice timeout".to_string())?
        .map_err(|e| e.to_string())?;
    let _ = cancel_tx.send(()).await;
    let stream = crate::client::ice_stream::connect(conn).await?;
    inner.lock().expect("lock").state = ViewerState::Ready(stream);
    Ok(())
}

async fn build_agent(ufrag: &str, pwd: &str) -> Result<Agent, String> {
    let urls = vec![Url::parse_url(STUN).map_err(|e| e.to_string())?];
    let cfg = AgentConfig {
        urls,
        local_ufrag: ufrag.to_string(),
        local_pwd: pwd.to_string(),
        network_types: vec![NetworkType::Udp4, NetworkType::Udp6],
        is_controlling: true,
        ..AgentConfig::default()
    };
    Agent::new(cfg).await.map_err(|e| e.to_string())
}

fn bind_candidates(agent: &Arc<Agent>, inner: Arc<Mutex<Inner>>) {
    agent.on_candidate(Box::new(move |c: Option<Arc<dyn Candidate + Send + Sync>>| {
        let inner2 = inner.clone();
        Box::pin(async move {
            if let Some(c) = c {
                let raw = c.marshal();
                let cand = raw.strip_prefix("candidate:").unwrap_or(&raw).to_string();
                queue_out(&inner2, ice_candidate_env(&cand));
                return;
            }
            queue_out(&inner2, ice_complete_env());
        })
    }));
}

async fn wait_remote_credentials(
    inner: &Arc<Mutex<Inner>>,
    agent: &Arc<Agent>,
    timeout: Duration,
) -> Result<(String, String), String> {
    let deadline = tokio::time::Instant::now() + timeout;
    let mut ufrag = String::new();
    let mut pwd = String::new();
    while ufrag.is_empty() || pwd.is_empty() {
        let msg = pop_inbox(inner, deadline).await?;
        if let Some(p) = env_parameters(&msg) {
            ufrag = p.ufrag.clone();
            pwd = p.pwd.clone();
        }
        if let Some(c) = env_candidate(&msg) {
            add_remote_candidate(agent, c)?;
        }
    }
    Ok((ufrag, pwd))
}

fn add_remote_candidate(agent: &Arc<Agent>, c: &IceCandidateInit) -> Result<(), String> {
    let raw = if c.candidate.starts_with("candidate:") {
        c.candidate.clone()
    } else {
        format!("candidate:{}", c.candidate)
    };
    let cand = unmarshal_candidate(&raw).map_err(|e| e.to_string())?;
    let arc: Arc<dyn Candidate + Send + Sync> = Arc::new(cand);
    agent.add_remote_candidate(&arc).map_err(|e| e.to_string())?;
    Ok(())
}

async fn pop_inbox(inner: &Arc<Mutex<Inner>>, deadline: tokio::time::Instant) -> Result<SignalEnvelope, String> {
    loop {
        if let Some(b) = inner.lock().expect("lock").inbox.pop_front() {
            return SignalEnvelope::decode(b.as_slice()).map_err(|e| e.to_string());
        }
        if tokio::time::Instant::now() >= deadline {
            return Err("signal timeout".to_string());
        }
        tokio::time::sleep(Duration::from_millis(10)).await;
    }
}

fn queue_out(inner: &Arc<Mutex<Inner>>, env: SignalEnvelope) {
    inner.lock().expect("lock").outbox.push_back(env.encode_to_vec());
}

fn ice_parameters_env(ufrag: &str, pwd: &str) -> SignalEnvelope {
    SignalEnvelope {
        payload: Some(Payload::IceParameters(IceParameters {
            ufrag: ufrag.to_string(),
            pwd: pwd.to_string(),
        })),
        ..Default::default()
    }
}

fn ice_candidate_env(candidate: &str) -> SignalEnvelope {
    SignalEnvelope {
        payload: Some(Payload::IceCandidate(IceCandidateInit {
            candidate: candidate.to_string(),
            ..Default::default()
        })),
        ..Default::default()
    }
}

fn ice_complete_env() -> SignalEnvelope {
    SignalEnvelope {
        payload: Some(Payload::IceComplete(true)),
        ..Default::default()
    }
}

fn env_parameters(env: &SignalEnvelope) -> Option<&IceParameters> {
    match env.payload.as_ref()? {
        Payload::IceParameters(p) => Some(p),
        _ => None,
    }
}

fn env_candidate(env: &SignalEnvelope) -> Option<&IceCandidateInit> {
    match env.payload.as_ref()? {
        Payload::IceCandidate(c) => Some(c),
        _ => None,
    }
}

fn new_credentials() -> (String, String) {
    let mut raw = [0u8; 16];
    rand::thread_rng().fill_bytes(&mut raw);
    (hex::encode(&raw[..4]), hex::encode(&raw[4..]))
}
