use std::os::raw::{c_char};
use std::ffi::{CString, CStr};


#[no_mangle]
pub extern "C" fn rust_greeting(to: *const c_char) -> *mut c_char {
    let c_str = unsafe { CStr::from_ptr(to) };
    let recipient = match c_str.to_str() {
        Err(_) => "there",
        Ok(string) => string,
    };

    CString::new("Hello ".to_owned() + recipient).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn rust_cstr_free(s: *mut c_char) {
    unsafe {
        if s.is_null() { return }
        CString::from_raw(s)
    };
}

#[repr(C)]
#[derive(Copy, Clone)]
pub struct NumPair {
    pub first: u64,
    pub second: usize,
}

#[no_mangle]
pub extern "C" fn process_pair(pair: NumPair) -> f64 {
    (pair.first as f64 * pair.second as f64) + 4.2
}
