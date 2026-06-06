//! ICE Agent/Viewer 共用：在独立线程里 drop，避免 tokio 异步任务内销毁 Runtime panic。
use std::future::Future;
use std::sync::{Arc, Mutex};
use std::time::Duration;

use tokio::runtime::Runtime;
use tokio::task::JoinHandle;

use crate::client::ice_common::{IceInbox, run_ice};
use crate::client::ice_stream::IceWire;

pub(crate) enum IceRoleState {
    Running,
    Ready(IceWire),
    Failed(String),
}

pub(crate) struct IceRoleRuntime {
    owned_rt: Option<Runtime>,
    ice_task: Option<JoinHandle<()>>,
}

impl IceRoleRuntime {
    pub fn spawn_ice(
        ice: Arc<Mutex<IceInbox>>,
        state: Arc<Mutex<IceRoleState>>,
        timeout_ms: u32,
        controlling: bool,
    ) -> Self {
        let ms = timeout_ms;
        let fut = async move {
            match run_ice(ice, Duration::from_millis(ms as u64), controlling).await {
                Ok(stream) => *state.lock().expect("lock") = IceRoleState::Ready(stream),
                Err(e) => *state.lock().expect("lock") = IceRoleState::Failed(e),
            }
        };
        if let Ok(h) = tokio::runtime::Handle::try_current() {
            return Self {
                owned_rt: None,
                ice_task: Some(h.spawn(fut)),
            };
        }
        let rt = Runtime::new().expect("tokio runtime");
        let ice_task = rt.spawn(fut);
        Self {
            owned_rt: Some(rt),
            ice_task: Some(ice_task),
        }
    }

    pub fn block_on<F: Future>(&self, fut: F) -> F::Output {
        if let Some(rt) = &self.owned_rt {
            rt.block_on(fut)
        } else {
            tokio::runtime::Handle::current().block_on(fut)
        }
    }
}

impl Drop for IceRoleRuntime {
    fn drop(&mut self) {
        let owned_rt = self.owned_rt.take();
        let ice_task = self.ice_task.take();
        std::thread::spawn(move || {
            if let Some(h) = ice_task {
                h.abort();
                drop(h);
            }
            drop(owned_rt);
        });
    }
}
