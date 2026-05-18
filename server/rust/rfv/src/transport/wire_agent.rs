//! Agent 数据面：列举/缩略图路径由 Viewer 在 wire 请求里指定；可选沙箱 `RB_AGENT_SANDBOX`。
use std::path::{Component, Path, PathBuf};

use prost::Message;

#[cfg(feature = "media")]
use crate::grpc_server::proto::{ListFilesRequest, ListFilesResponse, ThumbnailRequest, ThumbnailResponse};

const TYPE_FILE_INDEX: u8 = 2;
const TYPE_THUMB_REQ: u8 = 3;
const TYPE_THUMB_DATA: u8 = 4;

pub async fn serve_tcp(mut link: crate::transport::tcp_wire::TcpWire, sandbox: Option<String>) -> Result<(), String> {
    serve_link(&mut Link::Tcp(&mut link), sandbox).await
}

pub async fn serve_ice(mut ice: crate::client::ice_stream::IceWire, sandbox: Option<String>) -> Result<(), String> {
    serve_link(&mut Link::Ice(&mut ice), sandbox).await
}

enum Link<'a> {
    Tcp(&'a mut crate::transport::tcp_wire::TcpWire),
    Ice(&'a mut crate::client::ice_stream::IceWire),
}

impl Link<'_> {
    async fn read_frame(&mut self) -> Result<(u8, Vec<u8>), String> {
        match self {
            Link::Tcp(t) => t.read_frame().await,
            Link::Ice(i) => i.read_frame().await,
        }
    }

    async fn write_frame(&mut self, typ: u8, payload: &[u8]) -> Result<(), String> {
        match self {
            Link::Tcp(t) => t.write_frame(typ, payload).await,
            Link::Ice(i) => i.write_frame(typ, payload).await,
        }
    }
}

async fn serve_link(link: &mut Link<'_>, sandbox: Option<String>) -> Result<(), String> {
    loop {
        let (typ, payload) = link.read_frame().await?;
        match typ {
            TYPE_FILE_INDEX => {
                let out = build_file_index(sandbox.as_deref(), &payload)?;
                link.write_frame(TYPE_FILE_INDEX, &out).await?;
            }
            TYPE_THUMB_REQ => {
                let out = build_thumb(sandbox.as_deref(), &payload)?;
                link.write_frame(TYPE_THUMB_DATA, &out).await?;
            }
            _ => {}
        }
    }
}

#[cfg(feature = "media")]
fn build_file_index(sandbox: Option<&str>, payload: &[u8]) -> Result<Vec<u8>, String> {
    let req = ListFilesRequest::decode(payload).unwrap_or_default();
    let path = resolve_viewer_path(sandbox, &req.root_path)?;
    let rows = crate::remotebrowse::list_remote_files(path.to_str().ok_or("bad path")?)?;
    let entries = rows
        .into_iter()
        .map(|e| crate::grpc_server::proto::FileEntry {
            id: e.id,
            name: e.name,
            size: e.size,
            mtime_unix_ms: e.mtime_unix_ms,
            mime: e.mime,
            thumb_hash: e.thumb_hash,
            ..Default::default()
        })
        .collect();
    Ok(ListFilesResponse { entries, next_cursor: String::new() }.encode_to_vec())
}

#[cfg(not(feature = "media"))]
fn build_file_index(_sandbox: Option<&str>, _payload: &[u8]) -> Result<Vec<u8>, String> {
    Err("media feature required".into())
}

#[cfg(feature = "media")]
fn build_thumb(sandbox: Option<&str>, payload: &[u8]) -> Result<Vec<u8>, String> {
    let req = ThumbnailRequest::decode(payload).unwrap_or_default();
    let path = resolve_viewer_path(sandbox, &req.path)?;
    let max_edge = if req.max_edge == 0 { 256 } else { req.max_edge };
    let (data, hash, _) = crate::remotebrowse::ensure_thumbnail(&path, max_edge)?;
    Ok(ThumbnailResponse {
        data,
        mime: "image/webp".into(),
        thumb_hash: hash,
        ..Default::default()
    }
    .encode_to_vec())
}

#[cfg(not(feature = "media"))]
fn build_thumb(_sandbox: Option<&str>, _payload: &[u8]) -> Result<Vec<u8>, String> {
    Err("media feature required".into())
}

/// Viewer 在 ListFilesRequest.root_path / ThumbnailRequest.path 里指定目录或文件。
fn resolve_viewer_path(sandbox: Option<&str>, client: &str) -> Result<PathBuf, String> {
    let client = client.trim();
    if client.is_empty() {
        return Err("viewer must send path".into());
    }
    let mut path = PathBuf::from(client);
    if path.is_relative() {
        let base = sandbox.ok_or("relative path requires RB_AGENT_SANDBOX")?;
        path = PathBuf::from(base).join(path);
    }
    let path = normalize_path(&path);
    if let Some(sb) = sandbox {
        enforce_sandbox(&path, sb)?;
    }
    Ok(path)
}

fn normalize_path(path: &Path) -> PathBuf {
    let mut out = PathBuf::new();
    for c in path.components() {
        match c {
            Component::ParentDir => {
                out.pop();
            }
            Component::CurDir => {}
            other => out.push(other),
        }
    }
    out
}

fn enforce_sandbox(path: &Path, sandbox: &str) -> Result<(), String> {
    let sb = PathBuf::from(sandbox);
    let sb = sb.canonicalize().map_err(|e| format!("sandbox: {e}"))?;
    let p = path.canonicalize().map_err(|e| e.to_string())?;
    if !p.starts_with(&sb) {
        return Err("path outside RB_AGENT_SANDBOX".into());
    }
    Ok(())
}
