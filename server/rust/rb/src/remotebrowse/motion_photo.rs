//! Live Photo（HEIC+MOV）与 Android Motion Photo（内嵌 MP4 / JPG+MP4 侧车）列举合并。
use std::collections::{HashMap, HashSet};
use std::fs::File;
use std::io::{Read, Seek, SeekFrom};
use std::path::{Path, PathBuf};

use super::file_meta_db::{
    apply_motion_hit, motion_lookup, motion_store_companion, motion_store_embedded, motion_store_none,
    MotionCacheHit,
};
use super::{RemoteFileEntry, FILE_ENTRY_FLAG_DIRECTORY, FILE_ENTRY_FLAG_MOTION_PHOTO};

const STILL_EXTS: &[&str] = &["heic", "heif", "jpg", "jpeg"];
const COMPANION_EXTS: &[&str] = &["mov", "mp4"];

#[derive(Clone)]
pub(crate) struct Listed {
    pub entry: RemoteFileEntry,
    pub abs: PathBuf,
}

fn ext_lower(name: &str) -> String {
    Path::new(name)
        .extension()
        .and_then(|s| s.to_str())
        .unwrap_or("")
        .to_ascii_lowercase()
}

fn stem_lower(name: &str) -> String {
    Path::new(name)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or(name)
        .to_ascii_lowercase()
}

fn is_still_ext(ext: &str) -> bool {
    STILL_EXTS.contains(&ext)
}

fn is_companion_ext(ext: &str) -> bool {
    COMPANION_EXTS.contains(&ext)
}

fn xmp_u64(blob: &str, key: &str) -> Option<u64> {
    let needle = format!("{key}=\"");
    let i = blob.find(&needle)?;
    let rest = &blob[i + needle.len()..];
    let end = rest.find('"')?;
    rest[..end].parse().ok()
}

/// Google Motion Photo：XMP 标注 + 文件尾 MP4。
fn probe_embedded_motion(path: &Path, file_size: u64) -> Option<(u64, u64)> {
    if file_size < 4096 {
        return None;
    }
    let mut head = vec![0u8; 512 * 1024];
    let head_len = File::open(path).ok()?.read(&mut head).unwrap_or(0);
    head.truncate(head_len);
    let head_s = String::from_utf8_lossy(&head);
    let tail_take = file_size.min(8 * 1024 * 1024) as usize;
    let mut tail = vec![0u8; tail_take];
    let mut f = File::open(path).ok()?;
    f.seek(SeekFrom::End(-(tail_take as i64))).ok()?;
    let tail_len = f.read(&mut tail).unwrap_or(0);
    tail.truncate(tail_len);
    if let Some(off_end) = xmp_u64(&head_s, "GCamera:MicroVideo") {
        let len = xmp_u64(&head_s, "GCamera:MicroVideoLength")
            .or_else(|| xmp_u64(&head_s, "GCamera:MicroVideoSize"))
            .or_else(|| xmp_u64(&head_s, "Item:Length"))
            .unwrap_or(0);
        if off_end > 0 && off_end < file_size {
            let start = file_size.saturating_sub(off_end);
            let len = if len > 0 && start + len <= file_size {
                len
            } else {
                file_size.saturating_sub(start)
            };
            if len >= 1024 {
                return Some((start, len));
            }
        }
    }
    if let Some(len) = xmp_u64(&head_s, "Item:Length") {
        if let Some(pos) = tail.windows(4).rposition(|w| w == b"ftyp") {
            let start = file_size.saturating_sub(tail.len() as u64) + pos as u64;
            let start = start.saturating_sub(4);
            let len = if len > 0 && start + len <= file_size { len } else { file_size.saturating_sub(start) };
            if len >= 1024 {
                return Some((start, len));
            }
        }
    }
    if let Some(pos) = tail.windows(4).rposition(|w| w == b"ftyp") {
        let start = file_size.saturating_sub(tail.len() as u64) + pos as u64;
        let start = start.saturating_sub(4);
        let len = file_size.saturating_sub(start);
        if len >= 32 * 1024 && len < file_size / 2 {
            return Some((start, len));
        }
    }
    None
}

/// 仅读文件头判断是否需要 tail 探测，避免对大目录每张 JPG 读 8MB 尾。
fn head_hints_embedded_motion(path: &Path) -> bool {
    let mut head = vec![0u8; 64 * 1024];
    let Ok(mut f) = File::open(path) else {
        return false;
    };
    let n = f.read(&mut head).unwrap_or(0);
    if n == 0 {
        return false;
    }
    let s = String::from_utf8_lossy(&head[..n]);
    s.contains("GCamera:MicroVideo")
        || s.contains("MicroVideoOffset")
        || s.contains("MotionPhoto")
        || s.contains("Container:Directory")
}

pub fn merge_motion_photos(items: Vec<Listed>) -> Vec<RemoteFileEntry> {
    merge_motion_photos_with_options(items, true)
}

/// [probe_embedded]：为 false 时不对未缓存 JPG 读盘；命中 [file_meta_db::file_motion] 仍恢复内嵌/侧车标记。
/// 全量首次扫描可设环境变量 `RB_MOTION_PROBE_LIST=1`（结果写入 SQLite，同目录再次列举不再读盘）。
pub fn merge_motion_photos_with_options(items: Vec<Listed>, probe_embedded: bool) -> Vec<RemoteFileEntry> {
    let mut stills: HashMap<String, usize> = HashMap::new();
    let mut companions: HashMap<String, Vec<usize>> = HashMap::new();
    for (i, it) in items.iter().enumerate() {
        if it.entry.flags & FILE_ENTRY_FLAG_DIRECTORY != 0 {
            continue;
        }
        let ext = ext_lower(&it.entry.name);
        let stem = stem_lower(&it.entry.name);
        if is_still_ext(&ext) {
            stills.insert(stem, i);
        } else if is_companion_ext(&ext) {
            companions.entry(stem).or_default().push(i);
        }
    }
    let mut hide = HashSet::new();
    let mut out: Vec<Option<RemoteFileEntry>> = items.iter().map(|it| Some(it.entry.clone())).collect();
    for (stem, still_i) in &stills {
        let Some(comp_idxs) = companions.get(stem) else {
            continue;
        };
        let ext = ext_lower(&items[*still_i].entry.name);
        if ext != "jpg" && ext != "jpeg" && ext != "heic" && ext != "heif" {
            continue;
        }
        let Some(comp_i) = comp_idxs
            .iter()
            .min_by_key(|&&i| {
                let e = ext_lower(&items[i].entry.name);
                if e == "mov" {
                    0
                } else {
                    1
                }
            })
            .copied()
        else {
            continue;
        };
        hide.insert(comp_i);
        let still = &mut out[*still_i].as_mut().unwrap();
        let comp = &items[comp_i];
        still.flags |= FILE_ENTRY_FLAG_MOTION_PHOTO;
        still.motion_companion = comp.entry.name.clone();
        still.motion_offset = 0;
        still.motion_length = comp.entry.size.max(0) as u64;
        still.duration_ms = comp.entry.duration_ms;
        motion_store_companion(
            &items[*still_i].abs,
            still.size.max(0) as u64,
            still.mtime_unix_ms,
            &still.motion_companion,
            still.motion_length,
            still.duration_ms,
        );
    }
    for (still_i, it) in items.iter().enumerate() {
        if hide.contains(&still_i) || out[still_i].is_none() {
            continue;
        }
        let ext = ext_lower(&it.entry.name);
        if ext != "jpg" && ext != "jpeg" {
            continue;
        }
        if out[still_i].as_ref().unwrap().flags & FILE_ENTRY_FLAG_MOTION_PHOTO != 0 {
            continue;
        }
        let size = it.entry.size.max(0) as u64;
        let mtime = it.entry.mtime_unix_ms;
        if let Some(hit) = motion_lookup(&it.abs, size, mtime) {
            if hit != MotionCacheHit::NotMotion {
                apply_motion_hit(out[still_i].as_mut().unwrap(), hit);
            }
            continue;
        }
        if !probe_embedded {
            continue;
        }
        if !head_hints_embedded_motion(&it.abs) {
            motion_store_none(&it.abs, size, mtime);
            continue;
        }
        let Some((start, len)) = probe_embedded_motion(&it.abs, size) else {
            motion_store_none(&it.abs, size, mtime);
            continue;
        };
        let still = out[still_i].as_mut().unwrap();
        still.flags |= FILE_ENTRY_FLAG_MOTION_PHOTO;
        still.motion_offset = start;
        still.motion_length = len;
        still.motion_companion.clear();
        still.duration_ms = 3000;
        motion_store_embedded(&it.abs, size, mtime, start, len);
    }
    out.into_iter()
        .enumerate()
        .filter_map(|(i, e)| if hide.contains(&i) { None } else { e })
        .collect()
}
