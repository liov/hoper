//! Viewer ICE：Rust 内 HTTP/2 + protobuf client（ice_link），FFI 回传 protobuf 字节。
use std::ffi::{c_char, c_void, CStr};
use std::ptr;
use std::sync::Mutex;

use prost::Message;

use crate::client::ice_viewer::ViewerHandle;
use crate::ice_link;
use crate::rb_http_client::RbHttpClient;
use crate::remotebrowse_svc::proto::{
    DeleteFileRequest, DeleteFileResponse, ListFilesRequest, ListFilesResponse, ListMediaRequest, ListMediaResponse,
    MediaListCursor, ReadFileRequest, ReadFileResponse, ThumbnailBatchRequest, ThumbnailBatchResponse,
};

struct IceHttpBridge {
    client: RbHttpClient,
}

static ICE_HTTP: Mutex<Option<IceHttpBridge>> = Mutex::new(None);

pub fn open_from_handle(h: &ViewerHandle) -> Result<(), String> {
    let wire = h.take_wire_for_grpc().ok_or("ice not ready")?;
    let client = h.block_on(ice_link::connect(wire))?;
    *ICE_HTTP.lock().expect("lock") = Some(IceHttpBridge { client });
    Ok(())
}

pub fn close_bridge() {
    *ICE_HTTP.lock().expect("lock") = None;
}

fn borrow_client() -> Result<RbHttpClient, String> {
    ICE_HTTP
        .lock()
        .expect("lock")
        .as_ref()
        .map(|b| b.client.fork())
        .ok_or_else(|| "ice http not open".into())
}

pub fn list_files(h: &ViewerHandle, root: &str) -> Result<Vec<u8>, String> {
    let root = root.to_string();
    let mut client = borrow_client()?;
    h.block_on(async move {
        let resp: ListFilesResponse = client
            .post_proto(
                "/rb/v1/list",
                &ListFilesRequest {
                    root_path: root.into(),
                    ..Default::default()
                },
            )
            .await?;
        Ok(resp.encode_to_vec())
    })
}

pub fn list_media(h: &ViewerHandle, root: &str, cursor: &[u8], page_size: u32) -> Result<Vec<u8>, String> {
    let root = root.to_string();
    let cursor = if cursor.is_empty() {
        None
    } else {
        Some(MediaListCursor::decode(cursor).map_err(|e| format!("bad MediaListCursor: {e}"))?)
    };
    let mut client = borrow_client()?;
    h.block_on(async move {
        let resp: ListMediaResponse = client
            .post_proto(
                "/rb/v1/list/media",
                &ListMediaRequest {
                    root_path: root.into(),
                    cursor,
                    page_size,
                },
            )
            .await?;
        Ok(resp.encode_to_vec())
    })
}

pub fn thumb_batch(h: &ViewerHandle, paths: Vec<String>, max_edge: u32) -> Result<Vec<u8>, String> {
    let mut client = borrow_client()?;
    h.block_on(async move {
        let resp: ThumbnailBatchResponse = client
            .post_proto(
                "/rb/v1/thumb/batch",
                &ThumbnailBatchRequest { paths, max_edge },
            )
            .await?;
        Ok(resp.encode_to_vec())
    })
}

pub fn read_file(h: &ViewerHandle, path: &str, max_bytes: usize) -> Result<Vec<u8>, String> {
    read_file_range(h, path, 0, 0, max_bytes)
}

pub fn read_file_range(h: &ViewerHandle, path: &str, offset: i64, length: u32, max_out: usize) -> Result<Vec<u8>, String> {
    let req_len = if length == 0 { 512 * 1024 } else { length.min(2 * 1024 * 1024) };
    let path = path.to_string();
    let mut client = borrow_client()?;
    h.block_on(async move {
        let c: ReadFileResponse = client
            .post_proto(
                "/rb/v1/read",
                &ReadFileRequest {
                    path: path.clone(),
                    offset,
                    length: req_len,
                    ..Default::default()
                },
            )
            .await?;
        if !c.error.is_empty() {
            return Err(c.error);
        }
        let mut buf = c.data;
        if buf.len() > max_out {
            buf.truncate(max_out);
        }
        Ok(ReadFileResponse {
            path,
            total_size: c.total_size,
            mime: c.mime,
            data: buf,
            eof: true,
            ..Default::default()
        }
        .encode_to_vec())
    })
}

pub fn delete_file(h: &ViewerHandle, path: &str) -> Result<Vec<u8>, String> {
    let path = path.to_string();
    let mut client = borrow_client()?;
    h.block_on(async move {
        let resp: DeleteFileResponse = client
            .post_proto("/rb/v1/delete", &DeleteFileRequest { path: path.into() })
            .await?;
        Ok(resp.encode_to_vec())
    })
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_open(h: *mut c_void) -> i32 {
    if h.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    match open_from_handle(handle) {
        Ok(()) => 0,
        Err(_) => -1,
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_close() {
    close_bridge();
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_list_files(
    h: *mut c_void,
    root: *const c_char,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
) -> i32 {
    ffi_write(h, root, out, cap, out_len, |handle, root| list_files(handle, root))
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_list_media(
    h: *mut c_void,
    root: *const c_char,
    cursor: *const u8,
    cursor_len: usize,
    page_size: u32,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
) -> i32 {
    if h.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let root = match unsafe { CStr::from_ptr(root) }.to_str() {
        Ok(s) => s,
        Err(_) => return -1,
    };
    let cursor = if cursor.is_null() || cursor_len == 0 {
        &[]
    } else {
        unsafe { std::slice::from_raw_parts(cursor, cursor_len) }
    };
    ffi_write_ptr(out, cap, out_len, || list_media(handle, root, cursor, page_size))
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_thumb_batch(
    h: *mut c_void,
    paths_nl: *const c_char,
    max_edge: u32,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
) -> i32 {
    if h.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let paths = if paths_nl.is_null() {
        vec![]
    } else {
        let s = match unsafe { CStr::from_ptr(paths_nl) }.to_str() {
            Ok(v) => v,
            Err(_) => return -1,
        };
        s.split('\n').filter(|p| !p.is_empty()).map(str::to_string).collect()
    };
    ffi_write_ptr(out, cap, out_len, || thumb_batch(handle, paths, max_edge))
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_read_file(
    h: *mut c_void,
    path: *const c_char,
    max_bytes: u32,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
) -> i32 {
    ffi_write(h, path, out, cap, out_len, |handle, path| read_file(handle, path, max_bytes as usize))
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_read_file_range(
    h: *mut c_void,
    path: *const c_char,
    offset: i64,
    length: u32,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
) -> i32 {
    ffi_write(h, path, out, cap, out_len, |handle, path| read_file_range(handle, path, offset, length, cap))
}

#[unsafe(no_mangle)]
pub extern "C" fn rb_ice_grpc_delete_file(
    h: *mut c_void,
    path: *const c_char,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
) -> i32 {
    ffi_write(h, path, out, cap, out_len, |handle, path| delete_file(handle, path))
}

fn ffi_write(
    h: *mut c_void,
    path: *const c_char,
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
    f: impl FnOnce(&ViewerHandle, &str) -> Result<Vec<u8>, String>,
) -> i32 {
    if h.is_null() {
        return -1;
    }
    let handle = unsafe { &*(h as *const ViewerHandle) };
    let path = match unsafe { CStr::from_ptr(path) }.to_str() {
        Ok(s) => s,
        Err(_) => return -1,
    };
    ffi_write_ptr(out, cap, out_len, || f(handle, path))
}

fn ffi_write_ptr(
    out: *mut u8,
    cap: usize,
    out_len: *mut usize,
    f: impl FnOnce() -> Result<Vec<u8>, String>,
) -> i32 {
    if out.is_null() || out_len.is_null() {
        return -1;
    }
    match f() {
        Ok(data) => {
            if data.len() > cap {
                return -2;
            }
            unsafe {
                ptr::copy_nonoverlapping(data.as_ptr(), out, data.len());
                *out_len = data.len();
            }
            0
        }
        Err(_) => -1,
    }
}
