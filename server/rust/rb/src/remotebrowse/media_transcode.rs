//! 按需分片转码：动态 m3u8 + 按 GET 现编分片（仅内存，可选 RAM LRU，不落盘整片）。
use std::collections::HashMap;
use std::path::PathBuf;
use std::sync::{Mutex, OnceLock};
use std::time::Instant;

use axum::body::Body;
use axum::extract::{Query, State};
use axum::http::{header, StatusCode};
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use axum::Router;
use serde::Deserialize;

use crate::remotebrowse::media_probe;
use crate::remotebrowse::transcode_fragment;
use crate::remotebrowse::transcode_preset::{self, FRAGMENT_MS};
use crate::remotebrowse::transcode_vcodec::TranscodeVcodec;
use crate::remotebrowse_svc::MediaSvc;

/// 分片编码参数变更时递增，避免 RAM 缓存旧格式 TS。
const FRAG_CACHE_REV: u64 = 12;

/// m3u8 分片 URL 自带 prefetch，播放器逐片 GET 时 Agent 批量预编后续片。
const M3U8_FRAGMENT_PREFETCH: i64 = 10;

#[derive(Debug, Clone, Deserialize)]
pub struct TranscodeQuery {
    pub path: String,
    pub preset: String,
    #[serde(default)]
    pub t: String,
    #[serde(default)]
    pub off: u64,
    #[serde(default)]
    pub start_ms: Option<i64>,
    /// 播放列表从该时间点（毫秒）对应的分片起，避免客户端先拉 0 秒再 seek 造成音画 PTS 回跳。
    #[serde(default)]
    pub begin_ms: Option<i64>,
    /// 转码视频编码：`hevc`（默认）、`av1`；不支持 `h264`。
    #[serde(default)]
    pub vcodec: String,
    /// 拉取分片时同步预热后续分片数（含当前片），上限 16。
    #[serde(default)]
    pub prefetch: Option<i64>,
}

struct FragHit {
    bytes: Vec<u8>,
    at: Instant,
}

struct FragRamCache {
    map: HashMap<(u64, String, String, String, u64), FragHit>,
    bytes: usize,
}

impl FragRamCache {
    const MAX_BYTES: usize = 48 << 20;

    fn contains(&self, key: &(u64, String, String, String, u64)) -> bool {
        self.map.contains_key(key)
    }

    fn get(&mut self, key: &(u64, String, String, String, u64)) -> Option<Vec<u8>> {
        self.map.get(key).map(|h| h.bytes.clone())
    }

    fn put(&mut self, key: (u64, String, String, String, u64), bytes: Vec<u8>) {
        let n = bytes.len();
        if let Some(old) = self.map.insert(key, FragHit { bytes, at: Instant::now() }) {
            self.bytes = self.bytes.saturating_sub(old.bytes.len());
        }
        self.bytes += n;
        while self.bytes > Self::MAX_BYTES {
            let Some(victim) = self.map.iter().min_by_key(|(_, v)| v.at).map(|(k, _)| k.clone()) else {
                break;
            };
            if let Some(rem) = self.map.remove(&victim) {
                self.bytes = self.bytes.saturating_sub(rem.bytes.len());
            }
        }
    }
}

static FRAG_CACHE: OnceLock<Mutex<FragRamCache>> = OnceLock::new();

fn frag_cache() -> &'static Mutex<FragRamCache> {
    FRAG_CACHE.get_or_init(|| Mutex::new(FragRamCache { map: HashMap::new(), bytes: 0 }))
}

pub fn router() -> Router<MediaSvc> {
    Router::new()
        .route("/rb/v1/transcode/index.m3u8", get(vod_playlist))
        .route("/rb/v1/transcode/fragment", get(transcode_fragment_http))
        .route("/rb/v1/transcode/fragment.ts", get(transcode_fragment_http))
        .merge(crate::remotebrowse::media_transcode_stream::stream_route())
}

pub fn preset_params(preset: &str) -> Option<transcode_preset::TranscodePresetCfg> {
    transcode_preset::preset_cfg(preset)
}

async fn vod_playlist(State(svc): State<MediaSvc>, Query(q): Query<TranscodeQuery>) -> Response {
    if q.off > 0 {
        return (StatusCode::NOT_IMPLEMENTED, ()).into_response();
    }
    if preset_params(&q.preset).is_none() || !vcodec_ok(&q.vcodec) {
        return (StatusCode::BAD_REQUEST, ()).into_response();
    }
    if !check_token(&svc, &q.t) {
        return (StatusCode::FORBIDDEN, ()).into_response();
    }
    let info = match media_probe::probe_cached(&svc, &q.path) {
        Ok(v) => v,
        Err(_) => return (StatusCode::NOT_FOUND, ()).into_response(),
    };
    if info.duration_ms <= 0 {
        return (StatusCode::NOT_FOUND, ()).into_response();
    }
    let begin_ms = q.begin_ms.unwrap_or(0).max(0);
    let begin_idx = begin_ms / FRAGMENT_MS;
    let body = build_vod_m3u8(&q, info.duration_ms, begin_idx);
    let ahead = fragment_prefetch_count(&q).max(M3U8_FRAGMENT_PREFETCH);
    prefetch_fragments(&svc, &q, begin_idx, ahead, info.duration_ms);
    Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "application/vnd.apple.mpegurl")
        .header(header::CACHE_CONTROL, "no-store")
        .body(Body::from(body))
        .unwrap()
}

fn build_vod_m3u8(q: &TranscodeQuery, duration_ms: i64, begin_idx: i64) -> String {
    let n = (duration_ms + FRAGMENT_MS - 1) / FRAGMENT_MS;
    let begin_idx = begin_idx.clamp(0, n.saturating_sub(1).max(0));
    let mut s = format!(
        "#EXTM3U\n#EXT-X-VERSION:7\n#EXT-X-INDEPENDENT-SEGMENTS\n#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-MEDIA-SEQUENCE:{begin_idx}\n#EXT-X-TARGETDURATION:3\n",
    );
    for i in begin_idx..n {
        let start = i * FRAGMENT_MS;
        s.push_str("#EXTINF:2.000,\n");
        s.push_str(&fragment_url(q, start));
        s.push('\n');
    }
    s.push_str("#EXT-X-ENDLIST\n");
    s
}

fn fragment_url(q: &TranscodeQuery, start_ms: i64) -> String {
    let vc = vcodec_id(&q.vcodec);
    let pf = fragment_prefetch_count(q).max(M3U8_FRAGMENT_PREFETCH);
    let mut qs = format!(
        "fragment.ts?path={}&preset={}&vcodec={vc}&start_ms={start_ms}&prefetch={pf}",
        urlencoding_path(&q.path),
        q.preset
    );
    if !q.t.is_empty() {
        qs.push_str(&format!("&t={}", q.t));
    }
    qs
}

fn vcodec_id(raw: &str) -> &'static str {
    match TranscodeVcodec::parse_or_default(raw) {
        TranscodeVcodec::Hevc => "hevc",
        TranscodeVcodec::Av1 => "av1",
    }
}

fn vcodec_ok(raw: &str) -> bool {
    raw.trim().is_empty() || TranscodeVcodec::parse(raw).is_some()
}

async fn ensure_fragments(svc: &MediaSvc, q: &TranscodeQuery, from_idx: i64, count: i64, duration_ms: i64) -> Result<(), String> {
    let preset = preset_params(&q.preset).ok_or_else(|| "preset".to_string())?;
    let vcodec = TranscodeVcodec::parse_or_default(&q.vcodec);
    let input = svc.resolve_path(&q.path).map_err(|_| "path".to_string())?;
    let n = (duration_ms + FRAGMENT_MS - 1) / FRAGMENT_MS;
    if from_idx >= n {
        return Err("range".to_string());
    }
    let end = (from_idx + count).min(n);
    let mut todo: Vec<i64> = Vec::new();
    for i in from_idx..end {
        let start_ms = i * FRAGMENT_MS;
        let ck = frag_cache_key(&q.path, &q.preset, &q.vcodec, start_ms as u64);
        if frag_cache().lock().ok().is_some_and(|g| g.contains(&ck)) {
            continue;
        }
        todo.push(start_ms);
    }
    if todo.is_empty() {
        return Ok(());
    }
    let path = input.clone();
    let pairs = tokio::task::spawn_blocking(move || transcode_fragment::transcode_fragments_serial(&path, &todo, preset, vcodec))
        .await
        .map_err(|_| "join".to_string())??;
    let mut g = frag_cache().lock().map_err(|_| "lock".to_string())?;
    for (start_ms, bytes) in pairs {
        g.put(frag_cache_key(&q.path, &q.preset, &q.vcodec, start_ms as u64), bytes);
    }
    Ok(())
}

fn fragment_prefetch_count(q: &TranscodeQuery) -> i64 {
    q.prefetch.unwrap_or(M3U8_FRAGMENT_PREFETCH).clamp(1, 16)
}

fn prefetch_fragments(svc: &MediaSvc, q: &TranscodeQuery, from_idx: i64, count: i64, duration_ms: i64) {
    let q = q.clone();
    let svc = svc.clone();
    tokio::spawn(async move {
        let _ = ensure_fragments(&svc, &q, from_idx, count, duration_ms).await;
    });
}

async fn transcode_fragment_http(State(svc): State<MediaSvc>, Query(q): Query<TranscodeQuery>) -> Response {
    if q.off > 0 {
        return (StatusCode::NOT_IMPLEMENTED, ()).into_response();
    }
    let start_ms = q.start_ms.unwrap_or(0).max(0);
    if preset_params(&q.preset).is_none() {
        return (StatusCode::BAD_REQUEST, ()).into_response();
    }
    if !check_token(&svc, &q.t) {
        return (StatusCode::FORBIDDEN, ()).into_response();
    }
    if resolve_path(&svc, &q.path).is_err() {
        return (StatusCode::BAD_REQUEST, ()).into_response();
    }
    let ck = frag_cache_key(&q.path, &q.preset, &q.vcodec, start_ms as u64);
    if let Ok(info) = media_probe::probe_cached(&svc, &q.path) {
        let idx = start_ms / FRAGMENT_MS;
        if let Ok(mut g) = frag_cache().lock() {
            if let Some(hit) = g.get(&ck) {
                let ahead = fragment_prefetch_count(&q).saturating_sub(1).max(4);
                prefetch_fragments(&svc, &q, idx + 1, ahead, info.duration_ms);
                return ts_response(hit);
            }
        }
    } else if let Ok(mut g) = frag_cache().lock() {
        if let Some(hit) = g.get(&ck) {
            return ts_response(hit);
        }
    }
    if let Ok(info) = media_probe::probe_cached(&svc, &q.path) {
        let n = (info.duration_ms + FRAGMENT_MS - 1) / FRAGMENT_MS;
        let idx = start_ms / FRAGMENT_MS;
        if idx >= n {
            tracing::warn!(start_ms, duration_ms = info.duration_ms, path = %q.path, "fragment out of range");
            return (StatusCode::NOT_FOUND, ()).into_response();
        }
        let ensure_n = fragment_prefetch_count(&q);
        if let Err(e) = ensure_fragments(&svc, &q, idx, ensure_n, info.duration_ms).await {
            tracing::warn!(%e, start_ms, ensure_n, path = %q.path, "fragment ensure");
            return (StatusCode::INTERNAL_SERVER_ERROR, ()).into_response();
        }
        if let Ok(mut g) = frag_cache().lock() {
            if let Some(hit) = g.get(&ck) {
                let ahead = ensure_n.saturating_sub(1).max(4);
                prefetch_fragments(&svc, &q, idx + ensure_n, ahead, info.duration_ms);
                return ts_response(hit);
            }
        }
    }
    tracing::warn!(start_ms, path = %q.path, "fragment cache miss after ensure");
    (StatusCode::INTERNAL_SERVER_ERROR, ()).into_response()
}

fn ts_response(bytes: Vec<u8>) -> Response {
    Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "video/MP2T")
        .header(header::CACHE_CONTROL, "no-store")
        .header(header::ACCEPT_RANGES, "bytes")
        .body(Body::from(bytes))
        .unwrap()
}

pub(crate) fn check_token(svc: &MediaSvc, t: &str) -> bool {
    let expect = svc.media_token();
    expect.is_empty() || t == expect
}

fn resolve_path(svc: &MediaSvc, path: &str) -> Result<PathBuf, StatusCode> {
    svc.resolve_path(path).map_err(|_| StatusCode::BAD_REQUEST)
}

fn frag_cache_key(path: &str, preset: &str, vcodec: &str, start_ms: u64) -> (u64, String, String, String, u64) {
    (FRAG_CACHE_REV, path.to_string(), preset.to_string(), vcodec_id(vcodec).to_string(), start_ms)
}

fn urlencoding_path(p: &str) -> String {
    p.chars()
        .map(|c| match c {
            'A'..='Z' | 'a'..='z' | '0'..='9' | '-' | '_' | '.' | '~' => c.to_string(),
            _ => format!("%{:02X}", c as u8),
        })
        .collect()
}
