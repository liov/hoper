use std::sync::{Arc, Mutex};

use crate::client::ice_common::IceInbox;
use crate::client::ice_handle::{IceRoleRuntime, IceRoleState};
use crate::client::ice_stream::IceWire;
use crate::client::wire;

pub struct AgentHandle {
    rt: IceRoleRuntime,
    ice: Arc<Mutex<IceInbox>>,
    state: Arc<Mutex<IceRoleState>>,
}

impl AgentHandle {
    pub fn new(timeout_ms: u32) -> Self {
        let ice = Arc::new(Mutex::new(IceInbox::default()));
        let state = Arc::new(Mutex::new(IceRoleState::Running));
        let rt = IceRoleRuntime::spawn_ice(ice.clone(), state.clone(), timeout_ms, false);
        Self { rt, ice, state }
    }

    pub fn push(&self, data: &[u8]) {
        self.ice
            .lock()
            .expect("lock")
            .inbox
            .push_back(data.to_vec());
    }

    pub fn poll_out(&self) -> Option<Vec<u8>> {
        self.ice.lock().expect("lock").outbox.pop_front()
    }

    pub fn state_code(&self) -> i32 {
        match &*self.state.lock().expect("lock") {
            IceRoleState::Running => 0,
            IceRoleState::Ready(_) => 1,
            IceRoleState::Failed(_) => -1,
        }
    }

    pub fn read_frame(&self, buf: &mut [u8]) -> Result<(u8, usize), &'static str> {
        let mut g = self.state.lock().expect("lock");
        let IceRoleState::Ready(ref mut ice) = *g else {
            return Err("not ready");
        };
        let (typ, payload) = self
            .rt
            .block_on(ice.read_frame())
            .map_err(|_| "read failed")?;
        let frame = wire::encode_frame(typ, &payload);
        if frame.len() > buf.len() {
            return Err("short buf");
        }
        buf[..frame.len()].copy_from_slice(&frame);
        Ok((typ, frame.len()))
    }

    pub fn write_frame(&self, typ: u8, payload: &[u8]) -> Result<(), &'static str> {
        let mut g = self.state.lock().expect("lock");
        let IceRoleState::Ready(ref mut ice) = *g else {
            return Err("not ready");
        };
        self.rt
            .block_on(ice.write_frame(typ, payload))
            .map_err(|_| "write failed")
    }

    pub fn try_take_wire(&self) -> Option<IceWire> {
        let mut g = self.state.lock().expect("lock");
        if let IceRoleState::Ready(_) = *g {
            if let IceRoleState::Ready(w) = std::mem::replace(&mut *g, IceRoleState::Running) {
                return Some(w);
            }
        }
        None
    }
}
