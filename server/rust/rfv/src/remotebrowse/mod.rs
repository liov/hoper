//! 远程浏览：目录列举与缩略图缓存（gRPC 数据面由 `grpc_server` 暴露）。
#[cfg(feature = "media")]
mod thumbnail;

#[cfg(feature = "media")]
pub use thumbnail::{ensure_thumbnail, thumb_hash, DEFAULT_MAX_EDGE};

use serde::Serialize;
use std::fs;
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize)]
pub struct RemoteFileEntry {
    pub id: String,
    pub name: String,
    pub size: i64,
    pub mtime_unix_ms: i64,
    pub mime: String,
    pub thumb_hash: String,
}

fn guess_mime(name: &str) -> String {
    mime_guess::from_path(name)
        .first_or_octet_stream()
        .to_string()
}

pub fn list_remote_files(root: &str) -> Result<Vec<RemoteFileEntry>, String> {
    let path = PathBuf::from(root);
    if !path.is_dir() {
        return Err("not a directory".into());
    }
    let mut out = Vec::new();
    for entry in fs::read_dir(&path).map_err(|e| e.to_string())? {
        let entry = entry.map_err(|e| e.to_string())?;
        let meta = entry.metadata().map_err(|e| e.to_string())?;
        if !meta.is_file() {
            continue;
        }
        let name = entry.file_name().to_string_lossy().into_owned();
        let mtime = meta
            .modified()
            .ok()
            .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
            .map(|d| d.as_millis() as i64)
            .unwrap_or(0);
        let id = format!("{root}:{name}");
        let abs = entry.path().to_string_lossy().into_owned();
        let hash = {
            #[cfg(feature = "media")]
            {
                thumb_hash(&abs, mtime, DEFAULT_MAX_EDGE)
            }
            #[cfg(not(feature = "media"))]
            {
                String::new()
            }
        };
        out.push(RemoteFileEntry {
            id,
            name: name.clone(),
            size: meta.len() as i64,
            mtime_unix_ms: mtime,
            mime: guess_mime(&name),
            thumb_hash: hash,
        });
    }
    out.sort_by(|a, b| a.name.cmp(&b.name));
    Ok(out)
}
