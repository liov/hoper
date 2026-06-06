//! 缩略图元数据（hash / 缓存路径），不含 image 编解码。
use sha2::{Digest, Sha256};
use std::fs;
use std::path::{Path, PathBuf};

pub const DEFAULT_MAX_EDGE: u32 = 256;

pub const RB_PROTO: &str = "application/protobuf";

pub fn thumb_cache_dir() -> PathBuf {
    std::env::var("RB_THUMB_CACHE")
        .map(PathBuf::from)
        .unwrap_or_else(|_| PathBuf::from(".thumbnails"))
}

pub fn thumb_hash(path: &str, mtime_unix_ms: i64, max_edge: u32) -> String {
    let mut h = Sha256::new();
    h.update(path.as_bytes());
    h.update(mtime_unix_ms.to_le_bytes());
    h.update(max_edge.to_le_bytes());
    hex::encode(h.finalize())
}

fn cache_file(hash: &str) -> PathBuf {
    thumb_cache_dir().join(format!("{hash}.webp"))
}

/// 删除源文件前调用：按绝对路径与 mtime 清除各档缩略图磁盘缓存。
pub fn purge_thumbnail_cache(abs_path: &Path, mtime_unix_ms: i64) {
    const EDGES: [u32; 9] = [0, 64, 80, 96, 128, 256, 512, 1024, DEFAULT_MAX_EDGE];
    let path_key = crate::remotebrowse::path::display_path(abs_path);
    let mut removed = 0u32;
    for edge in EDGES {
        let cache = cache_file(&thumb_hash(&path_key, mtime_unix_ms, edge));
        if cache.is_file() && fs::remove_file(&cache).is_ok() {
            removed += 1;
        }
    }
    if removed > 0 {
        tracing::debug!(%path_key, removed, "purge thumbnail cache");
    }
}
