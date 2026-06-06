//! P2P / 本机共用的 RemoteBrowse 数据面（HTTP/2 + protobuf）。
#[cfg(any(feature = "agent-serve", feature = "media"))]
use std::path::PathBuf;
#[cfg(any(feature = "agent-serve", feature = "media"))]
use std::sync::Arc;

#[cfg(any(feature = "agent-serve", feature = "media"))]
use tokio::sync::Semaphore;

pub mod proto {
    include!(concat!(env!("OUT_DIR"), "/remotebrowse.rs"));
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
use proto::{
    DeleteFileRequest, DeleteFileResponse, FileEntry, HealthResponse, ListFilesRequest, ListFilesResponse,
    ListMediaRequest, ListMediaResponse, MediaEntry, MediaListCursor, ReadFileRequest, ReadFileResponse,
    ThumbnailBatchItem, ThumbnailBatchRequest, ThumbnailBatchResponse, ThumbnailRequest, ThumbnailResponse,
};

#[cfg(any(feature = "agent-serve", feature = "media"))]
pub use crate::remotebrowse::path::{default_browse_root, display_path, resolve_viewer_path};
pub use crate::remotebrowse::MediaListCursorState;

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn media_cursor_from_proto(c: Option<MediaListCursor>) -> MediaListCursorState {
    let Some(c) = c.filter(|c| !c.stack.is_empty() || !c.pending.is_empty() || c.scanned_dirs > 0) else {
        return MediaListCursorState::default();
    };
    MediaListCursorState {
        stack: c.stack,
        pending: c.pending,
        scanned_dirs: c.scanned_dirs,
    }
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn media_cursor_to_proto(c: &MediaListCursorState) -> Option<MediaListCursor> {
    if c.is_empty() {
        return None;
    }
    Some(MediaListCursor {
        stack: c.stack.clone(),
        pending: c.pending.clone(),
        scanned_dirs: c.scanned_dirs,
    })
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
#[derive(Clone)]
pub struct MediaSvc {
    sandbox: Option<Arc<str>>,
    media_token: Arc<str>,
    thumb_sem: Arc<Semaphore>,
    read_sem: Arc<Semaphore>,
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
impl Default for MediaSvc {
    fn default() -> Self {
        Self::new(None, None)
    }
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
impl MediaSvc {
    pub fn new(sandbox: Option<String>, room: Option<String>) -> Self {
        let sandbox = sandbox.filter(|s| !s.is_empty()).map(|s| Arc::<str>::from(s.as_str()));
        let media_token = Arc::<str>::from(media_token_for_room(room.as_deref()).as_str());
        Self {
            sandbox,
            media_token,
            thumb_sem: Arc::new(Semaphore::new(8)),
            read_sem: Arc::new(Semaphore::new(12)),
        }
    }

    pub fn media_token(&self) -> &str {
        &self.media_token
    }

    pub fn resolve_path(&self, client: &str) -> Result<PathBuf, String> {
        resolve_viewer_path(self.sandbox.as_deref(), client)
    }
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
impl MediaSvc {
    pub async fn pb_health(&self) -> Result<HealthResponse, String> {
        Ok(HealthResponse {
            signal_ws: "/rb/signal".into(),
            relay_tcp: String::new(),
            rb_grpc: std::env::var("RB_GRPC").unwrap_or_else(|_| "127.0.0.1:50051".into()),
            thumb_cache: crate::remotebrowse::thumb_cache_dir().to_string_lossy().into_owned(),
        })
    }

    pub async fn pb_list_files(&self, req: ListFilesRequest) -> Result<ListFilesResponse, String> {
        let root = req.root_path;
        tracing::debug!(%root, "http list_files");
        let sandbox = self.sandbox.clone();
        tokio::task::spawn_blocking(move || {
            let path = resolve_viewer_path(sandbox.as_deref(), &root)?;
            let listed = crate::remotebrowse::list_remote_files(path.to_str().ok_or_else(|| "bad path".to_string())?)?;
            let resolved = display_path(&path);
            let entries = listed
                .entries
                .into_iter()
                .map(|e| FileEntry {
                    id: e.id,
                    name: e.name,
                    size: e.size,
                    mtime_unix_ms: e.mtime_unix_ms,
                    mime: e.mime,
                    thumb_hash: e.thumb_hash,
                    flags: e.flags,
                    duration_ms: e.duration_ms,
                    motion_offset: e.motion_offset as i64,
                    motion_length: e.motion_length as i64,
                    motion_companion: e.motion_companion,
                    ..Default::default()
                })
                .collect();
            Ok(ListFilesResponse {
                entries,
                next_cursor: String::new(),
                resolved_root_path: resolved,
                hidden_skipped: listed.hidden_skipped,
            })
        })
        .await
        .map_err(|e| e.to_string())?
    }

    pub async fn pb_list_media(&self, req: ListMediaRequest) -> Result<ListMediaResponse, String> {
        let root = req.root_path;
        let cursor = media_cursor_from_proto(req.cursor);
        let page_size = if req.page_size == 0 { 64 } else { req.page_size };
        tracing::debug!(
            %root,
            cursor_stack = cursor.stack.len(),
            cursor_pending = cursor.pending.len(),
            page = page_size,
            "http list_media"
        );
        let sandbox = self.sandbox.clone();
        tokio::task::spawn_blocking(move || {
            let path = resolve_viewer_path(sandbox.as_deref(), &root)?;
            let page = crate::remotebrowse::list_media_page(&path, &cursor, page_size)?;
            let items = page
                .items
                .into_iter()
                .map(|(rel_path, e)| MediaEntry {
                    rel_path,
                    entry: Some(FileEntry {
                        id: e.id,
                        name: e.name,
                        size: e.size,
                        mtime_unix_ms: e.mtime_unix_ms,
                        mime: e.mime,
                        thumb_hash: e.thumb_hash,
                        flags: e.flags,
                        duration_ms: e.duration_ms,
                        motion_offset: e.motion_offset as i64,
                        motion_length: e.motion_length as i64,
                        motion_companion: e.motion_companion,
                        ..Default::default()
                    }),
                })
                .collect();
            Ok(ListMediaResponse {
                items,
                next_cursor: media_cursor_to_proto(&page.next_cursor),
                done: page.done,
                resolved_root_path: page.resolved_root_path,
                scanned_dirs: page.scanned_dirs,
            })
        })
        .await
        .map_err(|e| e.to_string())?
    }

    pub async fn pb_get_thumbnail(&self, inner: ThumbnailRequest) -> Result<ThumbnailResponse, String> {
        let path = resolve_viewer_path(self.sandbox.as_deref(), &inner.path)?;
        let max_edge = if inner.max_edge == 0 {
            crate::remotebrowse::DEFAULT_MAX_EDGE
        } else {
            inner.max_edge
        };
        let _permit = self.thumb_sem.acquire().await.map_err(|e| e.to_string())?;
        let (data, hash, cache) = tokio::task::spawn_blocking(move || crate::remotebrowse::ensure_thumbnail(&path, max_edge))
            .await
            .map_err(|e| e.to_string())?
            .map_err(|e| e)?;
        Ok(ThumbnailResponse {
            data,
            mime: "image/webp".into(),
            thumb_hash: hash,
            cache_path: cache.to_string_lossy().into_owned(),
        })
    }

    pub async fn pb_thumb_batch(&self, inner: ThumbnailBatchRequest) -> Result<ThumbnailBatchResponse, String> {
        let max_edge = if inner.max_edge == 0 { 256 } else { inner.max_edge };
        let sb = self.sandbox.clone();
        let paths = inner.paths;
        let items = tokio::task::spawn_blocking(move || thumb_batch_sync(sb.as_deref(), &paths, max_edge))
            .await
            .map_err(|e| e.to_string())?;
        Ok(ThumbnailBatchResponse { items })
    }

    pub async fn pb_read_file(&self, inner: ReadFileRequest) -> Result<ReadFileResponse, String> {
        let (rel_ref, want_original) = crate::remotebrowse::path::strip_read_original_suffix(inner.path.trim());
        let rel = rel_ref.to_string();
        let path = resolve_viewer_path(self.sandbox.as_deref(), &rel)?;
        let offset = inner.offset.max(0) as u64;
        let cap = if inner.length == 0 {
            2 * 1024 * 1024
        } else {
            inner.length.min(2 * 1024 * 1024)
        };
        let _permit = self.read_sem.clone().acquire_owned().await.map_err(|e| format!("read busy: {e}"))?;
        let path_buf = path.clone();
        let rel_s = rel.clone();
        tokio::task::spawn_blocking(move || read_file_once_sync(&path_buf, &rel_s, offset, cap, want_original))
            .await
            .map_err(|e| e.to_string())?
    }

    pub async fn pb_delete_file(&self, inner: DeleteFileRequest) -> Result<DeleteFileResponse, String> {
        let rel = inner.path.clone();
        let path = resolve_viewer_path(self.sandbox.as_deref(), &rel)?;
        let meta = std::fs::symlink_metadata(&path).map_err(|e| e.to_string())?;
        if meta.is_dir() {
            crate::remotebrowse::purge_meta_tree(&path);
            std::fs::remove_dir_all(&path).map_err(|e| e.to_string())?;
        } else {
            let mtime = meta
                .modified()
                .ok()
                .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
                .map(|d| d.as_millis() as i64)
                .unwrap_or(0);
            crate::remotebrowse::purge_thumbnail_cache(&path, mtime);
            crate::remotebrowse::purge_meta_path(&path);
            std::fs::remove_file(&path).map_err(|e| e.to_string())?;
        }
        Ok(DeleteFileResponse { path: rel, ..Default::default() })
    }
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
const THUMB_BATCH_PARALLEL: usize = 4;

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn thumb_batch_sync(sandbox: Option<&str>, paths: &[String], max_edge: u32) -> Vec<ThumbnailBatchItem> {
    let n = paths.len();
    if n == 0 {
        return vec![];
    }
    let mut items = vec![ThumbnailBatchItem::default(); n];
    let sb = sandbox.map(str::to_string);
    let mut offset = 0usize;
    while offset < n {
        let end = (offset + THUMB_BATCH_PARALLEL).min(n);
        let batch = &paths[offset..end];
        std::thread::scope(|scope| {
            let handles: Vec<_> = batch
                .iter()
                .enumerate()
                .map(|(j, rel)| {
                    let rel = rel.clone();
                    let sb = sb.clone();
                    let idx = offset + j;
                    scope.spawn(move || (idx, thumb_one(sb.as_deref(), rel, max_edge)))
                })
                .collect();
            for h in handles {
                if let Ok((idx, item)) = h.join() {
                    items[idx] = item;
                }
            }
        });
        offset = end;
    }
    items
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn thumb_one(sandbox: Option<&str>, rel: String, max_edge: u32) -> ThumbnailBatchItem {
    let path = match resolve_viewer_path(sandbox, &rel) {
        Ok(p) => p,
        Err(e) => return ThumbnailBatchItem { path: rel, error: e, ..Default::default() },
    };
    let out = std::panic::catch_unwind(|| crate::remotebrowse::ensure_thumbnail(&path, max_edge));
    match out {
        Ok(Ok((data, hash, _))) => ThumbnailBatchItem { path: rel, data, thumb_hash: hash, ..Default::default() },
        Ok(Err(e)) => ThumbnailBatchItem { path: rel, error: e, ..Default::default() },
        Err(_) => {
            tracing::warn!(%rel, "thumbnail panic (recovered)");
            ThumbnailBatchItem { path: rel, error: "thumbnail encode panic".into(), ..Default::default() }
        }
    }
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn read_file_webp_range(rel: &str, webp: &[u8], offset: u64, cap: u32) -> Result<ReadFileResponse, String> {
    let total = webp.len() as i64;
    if offset >= webp.len() as u64 {
        return Ok(ReadFileResponse {
            path: rel.into(),
            total_size: total,
            mime: "image/webp".into(),
            eof: true,
            ..Default::default()
        });
    }
    let start = offset as usize;
    let end = (start + cap as usize).min(webp.len());
    Ok(ReadFileResponse {
        path: rel.into(),
        total_size: total,
        mime: "image/webp".into(),
        data: webp[start..end].to_vec(),
        eof: end >= webp.len(),
        ..Default::default()
    })
}

#[cfg(any(feature = "agent-serve", feature = "media"))]
fn read_file_once_sync(
    path: &std::path::Path,
    rel: &str,
    mut offset: u64,
    cap: u32,
    want_original: bool,
) -> Result<ReadFileResponse, String> {
    use std::io::{Read, Seek, SeekFrom};
    let meta = std::fs::metadata(path).map_err(|e| e.to_string())?;
    if !meta.is_file() {
        return Err("not a file".into());
    }
    let file_len = meta.len();
    let fname = path.file_name().and_then(|s| s.to_str()).unwrap_or("");
    if want_original {
        let webp = crate::remotebrowse::maybe_hq_webp(path, file_len)?;
        return read_file_webp_range(rel, &webp, offset, cap);
    }
    if crate::remotebrowse::preview_transcode_for_client(fname) {
        return match crate::remotebrowse::maybe_preview_webp(path, file_len) {
            Ok(Some(webp)) => read_file_webp_range(rel, &webp, offset, cap),
            Ok(None) => Err(format!("preview transcode unavailable for {fname}")),
            Err(e) => Err(format!("preview transcode failed for {fname}: {e}")),
        };
    }
    if let Ok(Some(webp)) = crate::remotebrowse::maybe_preview_webp(path, file_len) {
        return read_file_webp_range(rel, &webp, offset, cap);
    }
    let total = file_len as i64;
    let mime = mime_guess::from_path(path).first_or_octet_stream().to_string();
    if offset >= file_len {
        return Ok(ReadFileResponse { path: rel.into(), total_size: total, mime, eof: true, ..Default::default() });
    }
    let mut f = std::fs::File::open(path).map_err(|e| e.to_string())?;
    if offset > 0 {
        f.seek(SeekFrom::Start(offset)).map_err(|e| e.to_string())?;
    }
    let to_read = ((meta.len() - offset) as usize).min(cap as usize);
    let mut buf = vec![0u8; to_read];
    let n = f.read(&mut buf).map_err(|e| e.to_string())?;
    buf.truncate(n);
    offset += n as u64;
    Ok(ReadFileResponse {
        path: rel.into(),
        total_size: total,
        mime,
        data: buf,
        eof: offset >= meta.len(),
        ..Default::default()
    })
}

/// 与 Flutter [rbMediaTokenForRoom] 一致。
#[cfg(any(feature = "agent-serve", feature = "media"))]
pub fn media_token_for_room(room: Option<&str>) -> String {
    let room = room.unwrap_or("").trim();
    if room.is_empty() {
        return String::new();
    }
    use sha2::{Digest, Sha256};
    let mut h = Sha256::new();
    h.update(b"rb-media:");
    h.update(room.as_bytes());
    hex::encode(h.finalize())[..32].to_string()
}
