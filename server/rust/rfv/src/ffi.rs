use std::ffi::{c_char, c_void, CStr};
use std::ptr;
use std::slice;
use std::sync::atomic::{AtomicBool, Ordering};
use std::thread;

#[cfg(feature = "transport")]
static RB_AGENT_BUSY: AtomicBool = AtomicBool::new(false);

use crate::client::ice_agent::AgentHandle;
use crate::client::ice_viewer::ViewerHandle;

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

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_new(timeout_ms: u32) -> *mut c_void {
    Box::into_raw(Box::new(AgentHandle::new(timeout_ms))) as *mut c_void
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_push(h: *mut c_void, data: *const u8, len: usize) -> i32 {
    if h.is_null() || data.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const AgentHandle) };
    handle.push(unsafe { slice::from_raw_parts(data, len) });
    0
}

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

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_state(h: *mut c_void) -> i32 {
    if h.is_null() {
        return -1;
    }
    unsafe { &*(h as *const AgentHandle) }.state_code()
}

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

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_agent_close(h: *mut c_void) {
    if h.is_null() {
        return;
    }
    unsafe {
        drop(Box::from_raw(h as *mut AgentHandle));
    }
}

/// 后台运行完整 Agent（信令 + 选路 + ffmpeg 数据面）。0=已启动；-1=参数错；-2=已在运行。
#[cfg(feature = "transport")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_agent_run(
    signal_url: *const c_char,
    room: *const c_char,
    root: *const c_char,
    timeout_ms: u32,
) -> i32 {
    if signal_url.is_null() || room.is_null() || root.is_null() {
        return -1;
    }
    if RB_AGENT_BUSY.swap(true, Ordering::SeqCst) {
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
    let root = match unsafe { CStr::from_ptr(root) }.to_str() {
        Ok(s) => s.to_string(),
        Err(_) => {
            RB_AGENT_BUSY.store(false, Ordering::SeqCst);
            return -1;
        }
    };
    thread::spawn(move || {
        let rt = tokio::runtime::Runtime::new().expect("tokio runtime");
        if let Err(e) = rt.block_on(crate::transport::run_agent(signal_url, room, root, timeout_ms)) {
            eprintln!("rb_agent_run: {e}");
        }
        RB_AGENT_BUSY.store(false, Ordering::SeqCst);
    });
    0
}

#[cfg(feature = "transport")]
#[unsafe(no_mangle)]
pub extern "C" fn rb_agent_running() -> i32 {
    i32::from(RB_AGENT_BUSY.load(Ordering::SeqCst))
}
