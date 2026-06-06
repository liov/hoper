//! 视频时长等元数据 SQLite 缓存，避免重复 ffprobe；支持删除同步与容量清理。
use std::path::{Path, PathBuf};
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Mutex, OnceLock};
use std::time::{SystemTime, UNIX_EPOCH};

use rusqlite::{params, Connection};

static DB: OnceLock<Mutex<Connection>> = OnceLock::new();
static LIST_OPS: AtomicU64 = AtomicU64::new(0);

const SCHEMA: &str = "
CREATE TABLE IF NOT EXISTS file_duration (
  abs_path TEXT PRIMARY KEY NOT NULL,
  size INTEGER NOT NULL,
  mtime_unix_ms INTEGER NOT NULL,
  duration_ms INTEGER NOT NULL,
  last_seen_ms INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_file_duration_last_seen ON file_duration(last_seen_ms);
CREATE TABLE IF NOT EXISTS file_motion (
  abs_path TEXT PRIMARY KEY NOT NULL,
  size INTEGER NOT NULL,
  mtime_unix_ms INTEGER NOT NULL,
  kind INTEGER NOT NULL,
  motion_offset INTEGER NOT NULL DEFAULT 0,
  motion_length INTEGER NOT NULL DEFAULT 0,
  motion_companion TEXT NOT NULL DEFAULT '',
  duration_ms INTEGER NOT NULL DEFAULT 0,
  last_seen_ms INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_file_motion_last_seen ON file_motion(last_seen_ms);
";

/// [file_motion.kind]：0=非动态照片，1=侧车视频，2=内嵌 MP4。
pub const MOTION_KIND_NONE: i32 = 0;
pub const MOTION_KIND_COMPANION: i32 = 1;
pub const MOTION_KIND_EMBEDDED: i32 = 2;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum MotionCacheHit {
    NotMotion,
    Companion {
        motion_companion: String,
        motion_length: u64,
        duration_ms: i64,
    },
    Embedded { motion_offset: u64, motion_length: u64 },
}

pub fn meta_db_path() -> PathBuf {
    let raw = std::env::var("RB_META_CACHE").unwrap_or_else(|_| ".rb_meta".into());
    let p = PathBuf::from(raw);
    if p.extension().is_some_and(|e| e == "db") {
        p
    } else {
        p.join("duration.db")
    }
}

fn now_ms() -> i64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_millis() as i64)
        .unwrap_or(0)
}

fn env_u64(key: &str, default: u64) -> u64 {
    std::env::var(key).ok().and_then(|s| s.parse().ok()).unwrap_or(default)
}

fn open_conn() -> Result<Connection, String> {
    let path = meta_db_path();
    if let Some(dir) = path.parent() {
        std::fs::create_dir_all(dir).map_err(|e| e.to_string())?;
    }
    let conn = Connection::open(&path).map_err(|e| e.to_string())?;
    conn.execute_batch(SCHEMA).map_err(|e| e.to_string())?;
    Ok(conn)
}

fn with_conn<R>(f: impl FnOnce(&Connection) -> Result<R, String>) -> Result<R, String> {
    let lock = DB.get_or_init(|| open_conn().map(Mutex::new).expect("meta db open"));
    let guard = lock.lock().map_err(|_| String::from("meta db lock poisoned"))?;
    f(&guard)
}

pub fn is_video_name(name: &str) -> bool {
    matches!(
        Path::new(name)
            .extension()
            .and_then(|s| s.to_str())
            .unwrap_or("")
            .to_ascii_lowercase()
            .as_str(),
        "mp4" | "flv" | "avi" | "rmvb" | "mov" | "mkv" | "webm" | "m4v"
    )
}

#[cfg(feature = "media")]
fn probe_duration_ms(path: &Path) -> i64 {
    use ffmpeg_next as ffmpeg;
    let Ok(ictx) = ffmpeg::format::input(path) else {
        return 0;
    };
    let d = ictx.duration();
    if d <= 0 {
        return 0;
    }
    d / 1000
}

#[cfg(not(feature = "media"))]
fn probe_duration_ms(_path: &Path) -> i64 {
    0
}

fn lookup_cached(conn: &Connection, key: &str, size: u64, mtime_unix_ms: i64) -> Result<Option<i64>, String> {
    let row: Result<(i64, i64, i64), rusqlite::Error> = conn.query_row(
        "SELECT size, mtime_unix_ms, duration_ms FROM file_duration WHERE abs_path = ?1",
        params![key],
        |r| Ok((r.get(0)?, r.get(1)?, r.get(2)?)),
    );
    match row {
        Ok((s, m, d)) if s == size as i64 && m == mtime_unix_ms => {
            let seen = now_ms();
            conn.execute(
                "UPDATE file_duration SET last_seen_ms = ?1 WHERE abs_path = ?2",
                params![seen, key],
            )
            .map_err(|e| e.to_string())?;
            Ok(Some(d))
        }
        Ok(_) => {
            conn.execute("DELETE FROM file_duration WHERE abs_path = ?1", params![key])
                .map_err(|e| e.to_string())?;
            Ok(None)
        }
        Err(rusqlite::Error::QueryReturnedNoRows) => Ok(None),
        Err(e) => Err(e.to_string()),
    }
}

fn upsert_row(conn: &Connection, key: &str, size: u64, mtime_unix_ms: i64, duration_ms: i64) -> Result<(), String> {
    let seen = now_ms();
    conn.execute(
        "INSERT INTO file_duration(abs_path, size, mtime_unix_ms, duration_ms, last_seen_ms)
         VALUES (?1, ?2, ?3, ?4, ?5)
         ON CONFLICT(abs_path) DO UPDATE SET
           size = excluded.size,
           mtime_unix_ms = excluded.mtime_unix_ms,
           duration_ms = excluded.duration_ms,
           last_seen_ms = excluded.last_seen_ms",
        params![key, size as i64, mtime_unix_ms, duration_ms, seen],
    )
    .map_err(|e| e.to_string())?;
    Ok(())
}

/// 缓存命中返回时长；未命中或 mtime/size 变化时 ffprobe 一次并落库。
pub fn duration_for_path(abs: &Path, size: u64, mtime_unix_ms: i64) -> i64 {
    let key = abs.to_string_lossy().into_owned();
    if let Ok(Some(d)) = with_conn(|c| lookup_cached(c, &key, size, mtime_unix_ms)) {
        return d;
    }
    let d = probe_duration_ms(abs);
    let _ = with_conn(|c| upsert_row(c, &key, size, mtime_unix_ms, d));
    d
}

pub fn purge_meta_path(abs: &Path) {
    let key = abs.to_string_lossy();
    let _ = with_conn(|c| {
        c.execute("DELETE FROM file_duration WHERE abs_path = ?1", params![key])
            .map_err(|e| e.to_string())?;
        c.execute("DELETE FROM file_motion WHERE abs_path = ?1", params![key])
            .map_err(|e| e.to_string())?;
        Ok(())
    });
}

pub fn purge_meta_tree(abs_dir: &Path) {
    let key = abs_dir.to_string_lossy().into_owned();
    let like = format!("{key}/%");
    let _ = with_conn(|c| {
        c.execute(
            "DELETE FROM file_duration WHERE abs_path = ?1 OR abs_path LIKE ?2",
            params![key, like],
        )
        .map_err(|e| e.to_string())?;
        c.execute(
            "DELETE FROM file_motion WHERE abs_path = ?1 OR abs_path LIKE ?2",
            params![key, like],
        )
        .map_err(|e| e.to_string())?;
        Ok(())
    });
}

fn motion_lookup_cached(conn: &Connection, key: &str, size: u64, mtime_unix_ms: i64) -> Result<Option<MotionCacheHit>, String> {
    let row: Result<(i64, i64, i32, i64, i64, String, i64), rusqlite::Error> = conn.query_row(
        "SELECT size, mtime_unix_ms, kind, motion_offset, motion_length, motion_companion, duration_ms
         FROM file_motion WHERE abs_path = ?1",
        params![key],
        |r| {
            Ok((
                r.get(0)?,
                r.get(1)?,
                r.get(2)?,
                r.get(3)?,
                r.get(4)?,
                r.get(5)?,
                r.get(6)?,
            ))
        },
    );
    match row {
        Ok((s, m, kind, off, len, comp, dur)) if s == size as i64 && m == mtime_unix_ms => {
            let seen = now_ms();
            conn.execute(
                "UPDATE file_motion SET last_seen_ms = ?1 WHERE abs_path = ?2",
                params![seen, key],
            )
            .map_err(|e| e.to_string())?;
            Ok(Some(match kind {
                MOTION_KIND_COMPANION => MotionCacheHit::Companion {
                    motion_companion: comp,
                    motion_length: len.max(0) as u64,
                    duration_ms: dur,
                },
                MOTION_KIND_EMBEDDED => MotionCacheHit::Embedded {
                    motion_offset: off.max(0) as u64,
                    motion_length: len.max(0) as u64,
                },
                _ => MotionCacheHit::NotMotion,
            }))
        }
        Ok(_) => {
            conn.execute("DELETE FROM file_motion WHERE abs_path = ?1", params![key])
                .map_err(|e| e.to_string())?;
            Ok(None)
        }
        Err(rusqlite::Error::QueryReturnedNoRows) => Ok(None),
        Err(e) => Err(e.to_string()),
    }
}

fn motion_upsert_row(
    conn: &Connection,
    key: &str,
    size: u64,
    mtime_unix_ms: i64,
    kind: i32,
    motion_offset: u64,
    motion_length: u64,
    motion_companion: &str,
    duration_ms: i64,
) -> Result<(), String> {
    let seen = now_ms();
    conn.execute(
        "INSERT INTO file_motion(abs_path, size, mtime_unix_ms, kind, motion_offset, motion_length, motion_companion, duration_ms, last_seen_ms)
         VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9)
         ON CONFLICT(abs_path) DO UPDATE SET
           size = excluded.size,
           mtime_unix_ms = excluded.mtime_unix_ms,
           kind = excluded.kind,
           motion_offset = excluded.motion_offset,
           motion_length = excluded.motion_length,
           motion_companion = excluded.motion_companion,
           duration_ms = excluded.duration_ms,
           last_seen_ms = excluded.last_seen_ms",
        params![
            key,
            size as i64,
            mtime_unix_ms,
            kind,
            motion_offset as i64,
            motion_length as i64,
            motion_companion,
            duration_ms,
            seen,
        ],
    )
    .map_err(|e| e.to_string())?;
    Ok(())
}

pub fn motion_lookup(abs: &Path, size: u64, mtime_unix_ms: i64) -> Option<MotionCacheHit> {
    let key = abs.to_string_lossy().into_owned();
    with_conn(|c| motion_lookup_cached(c, &key, size, mtime_unix_ms)).ok().flatten()
}

pub fn motion_store_none(abs: &Path, size: u64, mtime_unix_ms: i64) {
    let key = abs.to_string_lossy().into_owned();
    let _ = with_conn(|c| motion_upsert_row(c, &key, size, mtime_unix_ms, MOTION_KIND_NONE, 0, 0, "", 0));
}

pub fn motion_store_companion(
    abs: &Path,
    size: u64,
    mtime_unix_ms: i64,
    companion: &str,
    motion_length: u64,
    duration_ms: i64,
) {
    let key = abs.to_string_lossy().into_owned();
    let _ = with_conn(|c| {
        motion_upsert_row(
            c,
            &key,
            size,
            mtime_unix_ms,
            MOTION_KIND_COMPANION,
            0,
            motion_length,
            companion,
            duration_ms,
        )
    });
}

pub fn motion_store_embedded(abs: &Path, size: u64, mtime_unix_ms: i64, motion_offset: u64, motion_length: u64) {
    let key = abs.to_string_lossy().into_owned();
    let _ = with_conn(|c| {
        motion_upsert_row(
            c,
            &key,
            size,
            mtime_unix_ms,
            MOTION_KIND_EMBEDDED,
            motion_offset,
            motion_length,
            "",
            3000,
        )
    });
}

pub fn apply_motion_hit(entry: &mut super::RemoteFileEntry, hit: MotionCacheHit) {
    use super::FILE_ENTRY_FLAG_MOTION_PHOTO;
    match hit {
        MotionCacheHit::NotMotion => {}
        MotionCacheHit::Companion {
            motion_companion,
            motion_length,
            duration_ms,
        } => {
            entry.flags |= FILE_ENTRY_FLAG_MOTION_PHOTO;
            entry.motion_companion = motion_companion;
            entry.motion_offset = 0;
            entry.motion_length = motion_length;
            entry.duration_ms = duration_ms;
        }
        MotionCacheHit::Embedded { motion_offset, motion_length } => {
            entry.flags |= FILE_ENTRY_FLAG_MOTION_PHOTO;
            entry.motion_offset = motion_offset;
            entry.motion_length = motion_length;
            entry.motion_companion.clear();
            entry.duration_ms = 3000;
        }
    }
}

pub fn apply_motion_cache_to_entry(abs: &Path, entry: &mut super::RemoteFileEntry) {
    use super::FILE_ENTRY_FLAG_MOTION_PHOTO;
    if entry.flags & FILE_ENTRY_FLAG_MOTION_PHOTO != 0 {
        return;
    }
    let size = entry.size.max(0) as u64;
    let Some(hit) = motion_lookup(abs, size, entry.mtime_unix_ms) else {
        return;
    };
    apply_motion_hit(entry, hit);
}

fn delete_missing_from_table(conn: &Connection, table: &str, limit: u32) -> Result<u32, String> {
    let sql = format!("SELECT abs_path FROM {table} ORDER BY last_seen_ms ASC LIMIT ?1");
    let mut stmt = conn.prepare(&sql).map_err(|e| e.to_string())?;
    let paths: Vec<String> = stmt
        .query_map(params![limit], |r| r.get(0))
        .map_err(|e| e.to_string())?
        .filter_map(|r| r.ok())
        .collect();
    let mut n = 0u32;
    for p in paths {
        if Path::new(&p).exists() {
            continue;
        }
        let del = format!("DELETE FROM {table} WHERE abs_path = ?1");
        conn.execute(&del, params![p]).map_err(|e| e.to_string())?;
        n += 1;
    }
    Ok(n)
}

fn trim_table_to_max_rows(conn: &Connection, table: &str, max_rows: u64) -> Result<u32, String> {
    let count: i64 = conn
        .query_row(&format!("SELECT COUNT(*) FROM {table}"), [], |r| r.get(0))
        .map_err(|e| e.to_string())?;
    if count <= max_rows as i64 {
        return Ok(0);
    }
    let excess = (count - max_rows as i64) as u32;
    let sql = format!(
        "DELETE FROM {table} WHERE abs_path IN (
            SELECT abs_path FROM {table} ORDER BY last_seen_ms ASC LIMIT ?1
        )"
    );
    let n = conn.execute(&sql, params![excess]).map_err(|e| e.to_string())?;
    Ok(n as u32)
}

fn delete_missing_batch(conn: &Connection, limit: u32) -> Result<u32, String> {
    let half = limit.max(2) / 2;
    Ok(delete_missing_from_table(conn, "file_duration", half)? + delete_missing_from_table(conn, "file_motion", half)?)
}

fn trim_to_max_rows(conn: &Connection, max_rows: u64) -> Result<u32, String> {
    let half = max_rows / 2;
    Ok(trim_table_to_max_rows(conn, "file_duration", half)? + trim_table_to_max_rows(conn, "file_motion", half)?)
}

fn cleanup_stale() -> Result<(), String> {
    with_conn(|c| {
        let removed = delete_missing_batch(c, 512)?;
        let trimmed = trim_to_max_rows(c, env_u64("RB_META_MAX_ROWS", 200_000))?;
        if removed > 0 || trimmed > 0 {
            tracing::info!(removed, trimmed, "file_meta_db cleanup");
        }
        Ok(())
    })
}

/// 列举目录后调用：周期性清理不存在文件与超容量 LRU（后台线程，不阻塞 list）。
pub fn after_list_dir() {
    let every = env_u64("RB_META_CLEANUP_EVERY", 64).max(1);
    if LIST_OPS.fetch_add(1, Ordering::Relaxed) % every != 0 {
        return;
    }
    std::thread::spawn(|| {
        let _ = cleanup_stale();
    });
}
