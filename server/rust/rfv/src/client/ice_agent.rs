use std::sync::{Arc, Mutex};
use std::time::Duration;

use tokio::runtime::Runtime;

use crate::client::ice_common::{run_ice, IceInbox};
use crate::client::ice_stream::IceWire;
use crate::client::wire;

enum AgentState {
    Running,
    Ready(IceWire),
    Failed(String),
}

pub struct AgentHandle {
    runtime: Runtime,
    ice: Arc<Mutex<IceInbox>>,
    state: Arc<Mutex<AgentState>>,
}

impl AgentHandle {
    pub fn new(timeout_ms: u32) -> Self {
        let runtime = Runtime::new().expect("tokio runtime");
        let ice = Arc::new(Mutex::new(IceInbox::default()));
        let state = Arc::new(Mutex::new(AgentState::Running));
        let ice_bg = ice.clone();
        let state_bg = state.clone();
        runtime.spawn(async move {
            match run_ice(ice_bg, Duration::from_millis(timeout_ms as u64), false).await {
                Ok(stream) => *state_bg.lock().expect("lock") = AgentState::Ready(stream),
                Err(e) => *state_bg.lock().expect("lock") = AgentState::Failed(e),
            }
        });
        Self { runtime, ice, state }
    }

    pub fn push(&self, data: &[u8]) {
        self.ice.lock().expect("lock").inbox.push_back(data.to_vec());
    }

    pub fn poll_out(&self) -> Option<Vec<u8>> {
        self.ice.lock().expect("lock").outbox.pop_front()
    }

    pub fn state_code(&self) -> i32 {
        match &*self.state.lock().expect("lock") {
            AgentState::Running => 0,
            AgentState::Ready(_) => 1,
            AgentState::Failed(_) => -1,
        }
    }

    pub fn read_frame(&self, buf: &mut [u8]) -> Result<(u8, usize), &'static str> {
        let mut g = self.state.lock().expect("lock");
        let AgentState::Ready(ref mut ice) = *g else {
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
        let mut g = self.state.lock().expect("lock");
        let AgentState::Ready(ref mut ice) = *g else {
            return Err("not ready");
        };
        self.runtime
            .block_on(ice.write_frame(typ, payload))
            .map_err(|_| "write failed")
    }

    pub fn try_take_wire(&self) -> Option<IceWire> {
        let mut g = self.state.lock().expect("lock");
        if let AgentState::Ready(_) = *g {
            if let AgentState::Ready(w) = std::mem::replace(&mut *g, AgentState::Running) {
                return Some(w);
            }
        }
        None
    }
}