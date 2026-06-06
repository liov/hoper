use std::ffi::{c_void, CStr};
use std::ptr;
use std::slice;
#[cfg(feature = "agent-serve")]
use std::ffi::c_char;
#[cfg(feature = "agent-serve")]
use std::sync::atomic::{AtomicBool, Ordering};
#[cfg(feature = "agent-serve")]
use std::thread;

#[cfg(feature = "agent-serve")]
static RB_AGENT_BUSY: AtomicBool = AtomicBool::new(false);

use crate::client::ice_viewer::ViewerHandle;
#[cfg(feature = "agent-serve")]
use crate::client::ice_agent::AgentHandle;

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_new(timeout_ms: u32) -> *mut c_void {
    let h = ViewerHandle::new(timeout_ms);
    Box::into_raw(Box::new(h)) as *mut c_void
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_push(h: *mut c_void, data: *const u8, len: usize) -> i32 {
    if h.is_null() || data.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let buf = unsafe { slice::from_raw_parts(data, len) };
    handle.push(buf);
    0
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_poll_out(h: *mut c_void, buf: *mut u8, cap: usize, out_len: *mut usize) -> i32 {
    if h.is_null() || buf.is_null() || out_len.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let Some(data) = handle.poll_out() else {
        unsafe { *out_len = 0 };
        return 0;
    };
    if data.len() > cap {
        return -2;
    }
    unsafe {
        ptr::copy_nonoverlapping(data.as_ptr(), buf, data.len());
        *out_len = data.len();
    }
    0
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_state(h: *mut c_void) -> i32 {
    if h.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    handle.state_code()
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_write(h: *mut c_void, typ: u8, data: *const u8, len: usize) -> i32 {
    if h.is_null() || (len > 0 && data.is_null()) {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let payload = unsafe { slice::from_raw_parts(data, len) };
    if handle.write_frame(typ, payload).is_err() {
        return -1;
    }
    0
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_read(h: *mut c_void, buf: *mut u8, cap: usize, out_len: *mut usize) -> i32 {
    if h.is_null() || buf.is_null() || out_len.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let slice = unsafe { slice::from_raw_parts_mut(buf, cap) };
    match handle.read_frame(slice) {
        Ok((_, n)) => {
            unsafe { *out_len = n };
            0
        }
        Err(_) => -1,
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_viewer_close(h: *mut c_void) {
    if h.is_null() {
        return;
    }
    unsafe {
        drop(Box::from_raw(h as *mut ViewerHandle));
    }
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_new(timeout_ms: u32) -> *mut c_void {
    Box::into_raw(Box::new(AgentHandle::new(timeout_ms))) as *mut c_void
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_push(h: *mut c_void, data: *const u8, len: usize) -> i32 {
    if h.is_null() || data.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const AgentHandle) };
    handle.push(unsafe { slice::from_raw_parts(data, len) });
    0
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_poll_out(h: *mut c_void, buf: *mut u8, cap: usize, out_len: *mut usize) -> i32 {
    if h.is_null() || buf.is_null() || out_len.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const AgentHandle) };
    let Some(data) = handle.poll_out() else {
        unsafe { *out_len = 0 };
        return 0;
    };
    if data.len() > cap {
        return -2;
    }
    unsafe {
        ptr::copy_nonoverlapping(data.as_ptr(), buf, data.len());
        *out_len = data.len();
    }
    0
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_state(h: *mut c_void) -> i32 {
    if h.is_null() {
        return -1;
    }
    unsafe { &*(h as *const AgentHandle) }.state_code()
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_write(h: *mut c_void, typ: u8, data: *const u8, len: usize) -> i32 {
    if h.is_null() || (len > 0 && data.is_null()) {
        return -1;
    }
    let handle = unsafe { &*(h as *const AgentHandle) };
    let payload = unsafe { slice::from_raw_parts(data, len) };
    if handle.write_frame(typ, payload).is_err() {
        return -1;
    }
    0
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_read(h: *mut c_void, buf: *mut u8, cap: usize, out_len: *mut usize) -> i32 {
    if h.is_null() || buf.is_null() || out_len.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const AgentHandle) };
    let slice = unsafe { slice::from_raw_parts_mut(buf, cap) };
    match handle.read_frame(slice) {
        Ok((_, n)) => {
            unsafe { *out_len = n };
            0
        }
        Err(_) => -1,
    }
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_close(h: *mut c_void) {
    if h.is_null() {
        return;
    }
    unsafe {
        drop(Box::from_raw(h as *mut AgentHandle));
    }
}

/// 后台运行 Agent；浏览路径由 Viewer wire 请求携带。`sandbox` 可空（同 `RB_AGENT_SANDBOX`）。
#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_agent_run(
    signal_url: *const c_char,
    room: *const c_char,
    sandbox: *const c_char,
    timeout_ms: u32,
) -> i32 {
    if signal_url.is_null() || room.is_null() {
        return -1;
    }
    if RB_AGENT_BUSY.swap(true, Ordering::SeqCst) {
        tracing::warn!("rb_agent_run: already running");
        return -2;
    }
    let signal_url = match unsafe { CStr::from_ptr(signal_url) }.to_str() {
        Ok(s) => s.to_string(),
        Err(_) => {
            RB_AGENT_BUSY.store(false, Ordering::SeqCst);
            return -1;
        }
    };
    let room = match unsafe { CStr::from_ptr(room) }.to_str() {
        Ok(s) => s.to_string(),
        Err(_) => {
            RB_AGENT_BUSY.store(false, Ordering::SeqCst);
            return -1;
        }
    };
    let sandbox = if sandbox.is_null() {
        None
    } else {
        match unsafe { CStr::from_ptr(sandbox) }.to_str() {
            Ok(s) if s.is_empty() => None,
            Ok(s) => Some(s.to_string()),
            Err(_) => {
                RB_AGENT_BUSY.store(false, Ordering::SeqCst);
                return -1;
            }
        }
    };
    tracing::info!(%signal_url, %room, timeout_ms, sandbox = sandbox.as_deref().unwrap_or(""), "rb_agent_run start");
    thread::spawn(move || {
        let mut backoff = std::time::Duration::from_secs(1);
        loop {
            let signal_url = signal_url.clone();
            let room = room.clone();
            let sandbox = sandbox.clone();
            let r = std::panic::catch_unwind(std::panic::AssertUnwindSafe(|| {
                let rt = tokio::runtime::Runtime::new().expect("tokio runtime");
                rt.block_on(crate::transport::run_agent(signal_url, room, sandbox, timeout_ms))
            }));
            match r {
                Ok(()) => tracing::warn!(wait_secs = backoff.as_secs(), "rb_agent_run returned, restart"),
                Err(_) => tracing::error!(wait_secs = backoff.as_secs(), "rb_agent_run panic, restart"),
            }
            std::thread::sleep(backoff);
            backoff = (backoff * 2).min(std::time::Duration::from_secs(30));
        }
    });
    0
}

#[cfg(feature = "agent-serve")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_agent_running() -> i32 {
    i32::from(RB_AGENT_BUSY.load(Ordering::SeqCst))
}

#[cfg(feature = "viewer-ffi")]
pub use crate::ice_grpc_ffi::{
    rb_ice_grpc_close, rb_ice_grpc_delete_file, rb_ice_grpc_list_files, rb_ice_grpc_list_media, rb_ice_grpc_open,
    rb_ice_grpc_read_file, rb_ice_grpc_read_file_range, rb_ice_grpc_thumb_batch,
};
