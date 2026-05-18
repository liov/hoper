//! Agent 数据面 wire：本地列举 + rfv/ffmpeg 缩略图；对外业务 API 仍由 Go 提供。
use prost::Message;

#[cfg(feature = "media")]
use crate::grpc_server::proto::{ListFilesRequest, ListFilesResponse, ThumbnailRequest, ThumbnailResponse};

const TYPE_FILE_INDEX: u8 = 2;
const TYPE_THUMB_REQ: u8 = 3;
const TYPE_THUMB_DATA: u8 = 4;

pub async fn serve_tcp(mut link: crate::transport::tcp_wire::TcpWire, root: String) -> Result<(), String> {
    loop {
        let (typ, payload) = link.read_frame().await?;
        match typ {
            TYPE_FILE_INDEX => {
                let out = build_file_index(&root, &payload)?;
                link.write_frame(TYPE_FILE_INDEX, &out).await?;
            }
            TYPE_THUMB_REQ => {
                let out = build_thumb(&root, &payload)?;
                link.write_frame(TYPE_THUMB_DATA, &out).await?;
            }
            _ => {}
        }
    }
}

pub async fn serve_ice(mut ice: crate::client::ice_stream::IceWire, root: String) -> Result<(), String> {
    loop {
        let (typ, payload) = ice.read_frame().await?;
        match typ {
            TYPE_FILE_INDEX => {
                let out = build_file_index(&root, &payload)?;
                ice.write_frame(TYPE_FILE_INDEX, &out).await?;
            }
            TYPE_THUMB_REQ => {
                let out = build_thumb(&root, &payload)?;
                ice.write_frame(TYPE_THUMB_DATA, &out).await?;
            }
            _ => {}
        }
    }
}

#[cfg(feature = "media")]
fn build_file_index(root: &str, payload: &[u8]) -> Result<Vec<u8>, String> {
    let req = ListFilesRequest::decode(payload).unwrap_or_default();
    let path = if req.root_path.is_empty() { root } else { req.root_path.as_str() };
    let rows = crate::remotebrowse::list_remote_files(path)?;
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
fn build_file_index(_root: &str, _payload: &[u8]) -> Result<Vec<u8>, String> {
    Err("media feature required".into())
}

#[cfg(feature = "media")]
fn build_thumb(root: &str, payload: &[u8]) -> Result<Vec<u8>, String> {
    let req = ThumbnailRequest::decode(payload).unwrap_or_default();
    let path = std::path::PathBuf::from(resolve_path(root, &req.path));
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
fn build_thumb(_root: &str, _payload: &[u8]) -> Result<Vec<u8>, String> {
    Err("media feature required".into())
}

#[cfg(feature = "media")]
fn resolve_path(root: &str, path: &str) -> String {
    if path.is_empty() {
        return root.to_string();
    }
    if path.starts_with('/') || (cfg!(windows) && path.contains(':')) {
        return path.to_string();
    }
    format!("{}/{}", root.trim_end_matches('/'), path)
}
