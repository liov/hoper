//! 媒体探测：时长、是否可原画 Range 直播（仅元数据，不占视频副本空间）。
use std::path::Path;
use std::sync::{Mutex, OnceLock};

use axum::extract::{Query, State};
use axum::http::{header, StatusCode};
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use axum::Router;
use ffmpeg::codec::Id;
use ffmpeg::format;
use ffmpeg::media::Type;
use ffmpeg_next as ffmpeg;
use serde::{Deserialize, Serialize};

use crate::remotebrowse_svc::MediaSvc;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MediaProbeInfo {
    pub duration_ms: i64,
    pub direct_play: bool,
    pub video_codec: String,
    pub audio_codec: String,
}

#[derive(Debug, Deserialize)]
pub struct MediaMetaQuery {
    pub path: String,
    #[serde(default)]
    pub t: String,
}

static PROBE_CACHE: OnceLock<Mutex<std::collections::HashMap<String, MediaProbeInfo>>> = OnceLock::new();

fn cache() -> &'static Mutex<std::collections::HashMap<String, MediaProbeInfo>> {
    PROBE_CACHE.get_or_init(|| Mutex::new(std::collections::HashMap::new()))
}

pub fn router() -> Router<MediaSvc> {
    Router::new().route("/rb/v1/media/meta", get(media_meta))
}

pub fn probe_file(path: &Path) -> Result<MediaProbeInfo, String> {
    ffmpeg::init().map_err(|e| e.to_string())?;
    let mut ictx = format::input(path).map_err(|e| e.to_string())?;
    let dur_ms = ictx.duration().max(0) / 1_000;
    let v = ictx.streams().best(Type::Video);
    let a = ictx.streams().best(Type::Audio);
    let video_codec = v.as_ref().map(|s| s.parameters().id().name()).unwrap_or("").to_string();
    let audio_codec = a.as_ref().map(|s| s.parameters().id().name()).unwrap_or("").to_string();
    let v_ok = v.as_ref().map_or(false, |s| is_direct_video(s.parameters().id()));
    let a_ok = a.as_ref().map_or(true, |s| is_direct_audio(s.parameters().id()));
    let direct_play = v_ok && a_ok;
    Ok(MediaProbeInfo { duration_ms: dur_ms, direct_play, video_codec, audio_codec })
}

fn is_direct_video(id: Id) -> bool {
    matches!(id, Id::H264 | Id::HEVC | Id::H265 | Id::AV1 | Id::VP8 | Id::VP9)
}

fn is_direct_audio(id: Id) -> bool {
    matches!(id, Id::AAC | Id::MP3 | Id::OPUS | Id::VORBIS)
}

pub fn probe_cached(svc: &MediaSvc, rel_path: &str) -> Result<MediaProbeInfo, String> {
    let path = svc.resolve_path(rel_path)?;
    let key = cache_key(&path);
    if let Ok(g) = cache().lock() {
        if let Some(hit) = g.get(&key) {
            return Ok(hit.clone());
        }
    }
    let info = probe_file(&path)?;
    if let Ok(mut g) = cache().lock() {
        g.insert(key, info.clone());
    }
    Ok(info)
}

fn cache_key(path: &Path) -> String {
    let mtime = std::fs::metadata(path)
        .ok()
        .and_then(|m| m.modified().ok())
        .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
        .map(|d| d.as_millis())
        .unwrap_or(0);
    format!("{}|{mtime}", path.display())
}

async fn media_meta(State(svc): State<MediaSvc>, Query(q): Query<MediaMetaQuery>) -> Response {
    if !check_token(&svc, &q.t) {
        return (StatusCode::FORBIDDEN, ()).into_response();
    }
    match probe_cached(&svc, &q.path) {
        Ok(info) => {
            let body = serde_json::to_string(&info).unwrap_or_else(|_| "{}".into());
            Response::builder()
                .status(StatusCode::OK)
                .header(header::CONTENT_TYPE, "application/json")
                .header(header::CACHE_CONTROL, "no-store")
                .body(axum::body::Body::from(body))
                .unwrap()
        }
        Err(_) => (StatusCode::NOT_FOUND, ()).into_response(),
    }
}

fn check_token(svc: &MediaSvc, t: &str) -> bool {
    let expect = svc.media_token();
    expect.is_empty() || t == expect
}
