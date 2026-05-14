use std::ffi::c_void;
use std::ptr;
use std::slice;

use crate::ice_viewer::ViewerHandle;

#[no_mangle]
pub extern "C" fn rb_ice_viewer_new(timeout_ms: u32) -> *mut c_void {
    let h = ViewerHandle::new(timeout_ms);
    Box::into_raw(Box::new(h)) as *mut c_void
}

#[no_mangle]
pub extern "C" fn rb_ice_viewer_push(h: *mut c_void, data: *const u8, len: usize) -> i32 {
    if h.is_null() || data.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let buf = unsafe { slice::from_raw_parts(data, len) };
    handle.push(buf);
    0
}

#[no_mangle]
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

#[no_mangle]
pub extern "C" fn rb_ice_viewer_state(h: *mut c_void) -> i32 {
    if h.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    handle.state_code()
}

#[no_mangle]
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

#[no_mangle]
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

#[no_mangle]
pub extern "C" fn rb_ice_viewer_close(h: *mut c_void) {
    if h.is_null() {
        return;
    }
    unsafe {
        drop(Box::from_raw(h as *mut ViewerHandle));
    }
}
