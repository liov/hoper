//! 深度优先递归列举媒体（分页 cursor，避免客户端逐目录 listFiles）。
use std::path::Path;

use super::media_ext::is_media_file_name;
use super::path::{display_path, normalize_client_path};
use super::{list_remote_files, media_entry_from_abs, RemoteFileEntry, FILE_ENTRY_FLAG_DIRECTORY, FILE_ENTRY_FLAG_MOTION_PHOTO};

/// DFS 续扫状态（与 proto [MediaListCursor] 字段一致）。
#[derive(Debug, Default, Clone, PartialEq, Eq)]
pub struct MediaListCursorState {
    pub stack: Vec<String>,
    pub pending: Vec<String>,
    pub scanned_dirs: u32,
}

impl MediaListCursorState {
    pub fn is_empty(&self) -> bool {
        self.stack.is_empty() && self.pending.is_empty()
    }
}

pub struct ListMediaPage {
    pub items: Vec<(String, RemoteFileEntry)>,
    pub next_cursor: MediaListCursorState,
    pub done: bool,
    pub resolved_root_path: String,
    pub scanned_dirs: u32,
}

fn is_media_entry(e: &RemoteFileEntry) -> bool {
    e.flags & FILE_ENTRY_FLAG_DIRECTORY == 0
        && (e.flags & FILE_ENTRY_FLAG_MOTION_PHOTO != 0 || is_media_file_name(&e.name))
}

fn resolve_pending_paths(paths: &[String]) -> Vec<(String, RemoteFileEntry)> {
    paths
        .iter()
        .filter_map(|rel| {
            let p = normalize_client_path(rel).ok()?;
            let entry = media_entry_from_abs(&p)?;
            if !is_media_entry(&entry) {
                return None;
            }
            Some((rel.clone(), entry))
        })
        .collect()
}

fn take_pending(cur: &mut MediaListCursorState, limit: usize) -> Vec<(String, RemoteFileEntry)> {
    let n = limit.min(cur.pending.len());
    let chunk: Vec<String> = cur.pending.drain(..n).collect();
    resolve_pending_paths(&chunk)
}

fn push_subdirs(stack: &mut Vec<String>, dir: &Path, names: &mut Vec<String>) {
    names.sort();
    for name in names.iter().rev() {
        stack.push(display_path(&dir.join(name)));
    }
}

fn collect_dir_media(dir: &Path, out: &mut Vec<(String, RemoteFileEntry)>) -> Result<Vec<String>, String> {
    let listed = list_remote_files(dir.to_str().ok_or_else(|| "bad path".to_string())?)?;
    let mut subdirs = Vec::new();
    for e in listed.entries {
        if e.flags & FILE_ENTRY_FLAG_DIRECTORY != 0 {
            subdirs.push(e.name);
            continue;
        }
        if !is_media_entry(&e) {
            continue;
        }
        let abs = dir.join(&e.name);
        out.push((display_path(&abs), e));
    }
    Ok(subdirs)
}

/// DFS 分页：每页最多 [page_size] 个媒体项（1..=256）。
pub fn list_media_page(root: &Path, cursor: &MediaListCursorState, page_size: u32) -> Result<ListMediaPage, String> {
    let limit = page_size.clamp(1, 256) as usize;
    let resolved_root_path = display_path(root);
    let mut cur = cursor.clone();
    if cur.is_empty() {
        cur.stack.push(resolved_root_path.clone());
    }
    let mut items = take_pending(&mut cur, limit);
    while items.len() < limit && !cur.stack.is_empty() {
        let dir_s = cur.stack.pop().unwrap();
        let dir = match normalize_client_path(&dir_s) {
            Ok(p) => p,
            Err(_) => continue,
        };
        if !dir.is_dir() {
            continue;
        }
        cur.scanned_dirs = cur.scanned_dirs.saturating_add(1);
        let mut pending = Vec::new();
        let mut subdirs = collect_dir_media(&dir, &mut pending)?;
        push_subdirs(&mut cur.stack, &dir, &mut subdirs);
        for (rel_path, entry) in pending {
            if items.len() < limit {
                items.push((rel_path, entry));
            } else {
                cur.pending.push(rel_path);
            }
        }
    }
    let done = cur.is_empty();
    let scanned_dirs = cur.scanned_dirs;
    let next_cursor = if done { MediaListCursorState::default() } else { cur };
    Ok(ListMediaPage {
        items,
        next_cursor,
        done,
        resolved_root_path,
        scanned_dirs,
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn cursor_empty_when_no_stack_or_pending() {
        assert!(MediaListCursorState::default().is_empty());
        assert!(!MediaListCursorState { stack: vec!["/a".into()], ..Default::default() }.is_empty());
    }
}
