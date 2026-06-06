//! 远程浏览：目录列举与缩略图缓存（HTTP/2 + protobuf 数据面）。
#[cfg(any(feature = "agent-serve", feature = "media"))]
pub mod http_api;
#[cfg(any(feature = "agent-serve", feature = "media"))]
pub mod http_log;
#[cfg(feature = "media")]
pub mod media_http;
#[cfg(feature = "media")]
pub mod media_transcode;
#[cfg(feature = "media")]
pub mod media_probe;
#[cfg(feature = "media")]
mod mpegts_mem_out;
#[cfg(feature = "media")]
mod transcode_fragment;
#[cfg(feature = "media")]
mod transcode_stream;
#[cfg(feature = "media")]
mod media_transcode_stream;
#[cfg(feature = "media")]
mod transcode_preset;
#[cfg(feature = "media")]
mod transcode_vcodec;
#[cfg(any(feature = "agent-serve", feature = "media"))]
mod motion_photo;
mod media_ext;
#[cfg(any(feature = "agent-serve", feature = "media"))]
mod file_meta_db;
#[cfg(any(feature = "agent-serve", feature = "media"))]
pub use file_meta_db::{after_list_dir, duration_for_path, is_video_name, meta_db_path, purge_meta_path, purge_meta_tree};
#[cfg(any(feature = "agent-serve", feature = "media"))]
mod list_media;
#[cfg(any(feature = "agent-serve", feature = "media"))]
pub use list_media::{list_media_page, MediaListCursorState};
#[cfg(any(feature = "agent-serve", feature = "media"))]
pub mod path;
#[cfg(feature = "rb-core")]
pub mod thumb_meta;
#[cfg(feature = "rb-thumb")]
mod thumbnail;
#[cfg(feature = "rb-thumb")]
mod image_preview;

#[cfg(feature = "rb-core")]
pub use thumb_meta::{thumb_cache_dir, thumb_hash, DEFAULT_MAX_EDGE, RB_PROTO};
#[cfg(feature = "rb-core")]
pub use thumb_meta::purge_thumbnail_cache;
#[cfg(feature = "rb-thumb")]
pub use thumbnail::ensure_thumbnail;
#[cfg(feature = "rb-thumb")]
pub use image_preview::{maybe_hq_webp, maybe_preview_webp};
pub use media_ext::preview_transcode_for_client;

#[cfg(any(feature = "agent-serve", feature = "media"))]
use serde::{Deserialize, Serialize};
#[cfg(any(feature = "agent-serve", feature = "media"))]
use std::fs;
#[cfg(any(feature = "agent-serve", feature = "media"))]
use std::path::PathBuf;

/// `FileEntry.flags` bit0：目录
pub const FILE_ENTRY_FLAG_DIRECTORY: u32 = 1;
/// bit1：Live Photo / Motion Photo（静图 + 短视频）
pub const FILE_ENTRY_FLAG_MOTION_PHOTO: u32 = 1 << 1;

#[cfg(any(feature = "agent-serve", feature = "media"))]
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct RemoteFileEntry {
    pub id: String,
    pub name: String,
    pub size: i64,
    pub mtime_unix_ms: i64,
    pub mime: String,
    pub thumb_hash: String,
    pub flags: u32,
    pub duration_ms: i64,
    pub motion_offset: u64,
    pub motion_length: u64,
    pub motion_companion: String,
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn guess_mime(name: &str) -> String {
    mime_guess::from_path(name)
        .first_or_octet_stream()
        .to_string()
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
pub struct ListRemoteFilesResult {
    pub entries: Vec<RemoteFileEntry>,
    pub hidden_skipped: u32,
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
use std::path::Path;

#[cfg(any(feature = "agent-serve", feature = "media"))]
pub fn media_entry_from_abs(abs: &Path) -> Option<RemoteFileEntry> {
    let meta = fs::metadata(abs).ok()?;
    if !meta.is_file() {
        return None;
    }
    let name = abs.file_name()?.to_string_lossy().into_owned();
    let mtime = meta
        .modified()
        .ok()
        .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
        .map(|d| d.as_millis() as i64)
        .unwrap_or(0);
    let abs_key = path::display_path(&abs);
    let mut entry = RemoteFileEntry {
        id: name.clone(),
        name,
        size: meta.len() as i64,
        mtime_unix_ms: mtime,
        mime: guess_mime(abs.file_name()?.to_str().unwrap_or("")),
        thumb_hash: thumb_hash(&abs_key, mtime, DEFAULT_MAX_EDGE),
        ..Default::default()
    };
    file_meta_db::apply_motion_cache_to_entry(abs, &mut entry);
    Some(entry)
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
pub fn list_remote_files(root: &str) -> Result<ListRemoteFilesResult, String> {
    let path = path::normalize_client_path(root)?;
    if !path.is_dir() {
        tracing::warn!(%root, "list_remote_files: not a directory");
        return Err("not a directory".into());
    }
    let mut raw = Vec::new();
    for entry in fs::read_dir(&path).map_err(|e| e.to_string())? {
        let entry = entry.map_err(|e| e.to_string())?;
        let meta = entry.metadata().map_err(|e| e.to_string())?;
        let name = entry.file_name().to_string_lossy().into_owned();
        let mtime = meta
            .modified()
            .ok()
            .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
            .map(|d| d.as_millis() as i64)
            .unwrap_or(0);
        if meta.is_dir() {
            raw.push(motion_photo::Listed {
                entry: RemoteFileEntry {
                    id: name.clone(),
                    name,
                    mtime_unix_ms: mtime,
                    mime: "inode/directory".into(),
                    flags: FILE_ENTRY_FLAG_DIRECTORY,
                    ..Default::default()
                },
                abs: entry.path(),
            });
            continue;
        }
        if !meta.is_file() {
            continue;
        }
        let abs = entry.path();
        let abs_key = path::display_path(&abs);
        let hash = thumb_hash(&abs_key, mtime, DEFAULT_MAX_EDGE);
        let file_entry = RemoteFileEntry {
            id: name.clone(),
            name: name.clone(),
            size: meta.len() as i64,
            mtime_unix_ms: mtime,
            mime: guess_mime(&name),
            thumb_hash: hash,
            ..Default::default()
        };
        // 列表不 probe 时长（大目录会卡）；预览/按需再查 SQLite/ffprobe。
        raw.push(motion_photo::Listed { entry: file_entry, abs });
    }
    let probe_embedded = std::env::var("RB_MOTION_PROBE_LIST").ok().as_deref() == Some("1");
    let mut out = motion_photo::merge_motion_photos_with_options(raw, probe_embedded);
    out.sort_by(|a, b| {
        let da = a.flags & FILE_ENTRY_FLAG_DIRECTORY != 0;
        let db = b.flags & FILE_ENTRY_FLAG_DIRECTORY != 0;
        match (da, db) {
            (true, false) => std::cmp::Ordering::Less,
            (false, true) => std::cmp::Ordering::Greater,
            _ => a.name.cmp(&b.name),
        }
    });
    tracing::info!(%root, count = out.len(), "list_remote_files");
    after_list_dir();
    Ok(ListRemoteFilesResult { entries: out, hidden_skipped: 0 })
}
