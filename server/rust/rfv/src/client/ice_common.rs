use std::collections::VecDeque;
use std::sync::{Arc, Mutex};
use std::time::Duration;

use prost::Message;
use rand::RngCore;
use tokio::sync::mpsc;
use webrtc_ice::agent::agent_config::AgentConfig;
use webrtc_ice::agent::Agent;
use webrtc_ice::candidate::candidate_base::unmarshal_candidate;
use webrtc_ice::candidate::Candidate;
use webrtc_ice::network_type::NetworkType;
use webrtc_ice::url::Url;

use crate::client::ice_stream::IceWire;
use crate::signal_proto::signal_envelope::Payload;
use crate::signal_proto::{IceCandidateInit, IceParameters, SignalEnvelope};

const STUN: &str = "stun:stun.l.google.com:19302";

pub(crate) struct IceInbox {
    pub outbox: VecDeque<Vec<u8>>,
    pub inbox: VecDeque<Vec<u8>>,
}

impl Default for IceInbox {
    fn default() -> Self {
        Self { outbox: VecDeque::new(), inbox: VecDeque::new() }
    }
}

pub async fn run_ice(
    inner: Arc<Mutex<IceInbox>>,
    timeout: Duration,
    controlling: bool,
) -> Result<IceWire, String> {
    let (ufrag, pwd) = new_credentials();
    let agent = Arc::new(build_agent(&ufrag, &pwd, controlling).await?);
    queue_out(&inner, ice_parameters_env(&ufrag, &pwd));
    bind_candidates(&agent, inner.clone());
    agent.gather_candidates().map_err(|e| e.to_string())?;
    let (remote_ufrag, remote_pwd) = wait_remote_credentials(&inner, &agent, timeout).await?;
    agent
        .set_remote_credentials(remote_ufrag.clone(), remote_pwd.clone())
        .await
        .map_err(|e| e.to_string())?;
    let (cancel_tx, cancel_rx) = mpsc::channel::<()>(1);
    if controlling {
        let conn = tokio::time::timeout(timeout, agent.dial(cancel_rx, remote_ufrag, remote_pwd))
            .await
            .map_err(|_| "ice timeout".to_string())?
            .map_err(|e| e.to_string())?;
        let _ = cancel_tx.send(()).await;
        return crate::client::ice_stream::connect(conn).await;
    }
    let conn = tokio::time::timeout(timeout, agent.accept(cancel_rx, remote_ufrag, remote_pwd))
        .await
        .map_err(|_| "ice timeout".to_string())?
        .map_err(|e| e.to_string())?;
    let _ = cancel_tx.send(()).await;
    crate::client::ice_stream::connect(conn).await
}

async fn build_agent(ufrag: &str, pwd: &str, controlling: bool) -> Result<Agent, String> {
    let urls = vec![Url::parse_url(STUN).map_err(|e| e.to_string())?];
    let cfg = AgentConfig {
        urls,
        local_ufrag: ufrag.to_string(),
        local_pwd: pwd.to_string(),
        network_types: vec![NetworkType::Udp4, NetworkType::Udp6],
        is_controlling: controlling,
        ..AgentConfig::default()
    };
    Agent::new(cfg).await.map_err(|e| e.to_string())
}

fn bind_candidates(agent: &Arc<Agent>, inner: Arc<Mutex<IceInbox>>) {
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
    inner: &Arc<Mutex<IceInbox>>,
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

async fn pop_inbox(inner: &Arc<Mutex<IceInbox>>, deadline: tokio::time::Instant) -> Result<SignalEnvelope, String> {
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

pub(crate) fn queue_out(inner: &Arc<Mutex<IceInbox>>, env: SignalEnvelope) {
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
