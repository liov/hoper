//! Agent 同端口 HTTP Range：`GET /rb/v1/media`（与 HTTP/2 protobuf API 共用 TCP，见 `p2p_http::serve_h2_io`）。
use std::io::SeekFrom;

use axum::body::Body;
use axum::extract::{Query, State};
use axum::http::{header, HeaderMap, StatusCode};
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use axum::Router;
use serde::Deserialize;
use tokio::fs::File;
use tokio::io::AsyncSeekExt;
use tokio_util::io::ReaderStream;

use crate::remotebrowse_svc::MediaSvc;

const MEDIA_READ_CHUNK: usize = 2 * 1024 * 1024;

pub fn router(svc: MediaSvc) -> Router {
    Router::new()
        .route("/rb/v1/media", get(media_get).head(media_head))
        .merge(crate::remotebrowse::media_probe::router())
        .merge(crate::remotebrowse::media_transcode::router())
        .with_state(svc)
}

#[derive(Debug, Deserialize)]
struct MediaQuery {
    path: String,
    #[serde(default)]
    t: String,
    /// Motion Photo 内嵌 MP4 在容器文件内的字节偏移。
    #[serde(default)]
    off: u64,
    /// `1`：预览（bpp>0.3 或不可解码格式 → 0.3bpp WebP；否则原文件）。`2`：高清 WebP（原图按钮）。
    #[serde(default)]
    preview: u8,
}

async fn media_head(State(svc): State<MediaSvc>, Query(q): Query<MediaQuery>) -> Response {
    match media_meta(&svc, &q).await {
        Ok((mime, total)) => media_response_headers(StatusCode::OK, &mime, total, 0, total.saturating_sub(1), total, false),
        Err(code) => (code, ()).into_response(),
    }
}

async fn media_get(State(svc): State<MediaSvc>, Query(q): Query<MediaQuery>, headers: HeaderMap) -> Response {
    if q.preview == 2 {
        return media_get_hq(State(svc), Query(q), headers).await;
    }
    if q.preview == 1 {
        return media_get_preview(State(svc), Query(q), headers).await;
    }
    let (mime, total) = match media_meta(&svc, &q).await {
        Ok(v) => v,
        Err(code) => return (code, ()).into_response(),
    };
    if total == 0 {
        return (StatusCode::NOT_FOUND, ()).into_response();
    }
    let (start, end, partial) = parse_range(headers.get(header::RANGE).and_then(|v| v.to_str().ok()), total);
    let len = end - start + 1;
    let path = match resolve_media_path(&svc, &q) {
        Ok(p) => p,
        Err(code) => return (code, ()).into_response(),
    };
    let mut file = match File::open(&path).await {
        Ok(f) => f,
        Err(_) => return (StatusCode::NOT_FOUND, ()).into_response(),
    };
    let seek_to = q.off.saturating_add(start);
    if file.seek(SeekFrom::Start(seek_to)).await.is_err() {
        return (StatusCode::INTERNAL_SERVER_ERROR, ()).into_response();
    }
    use tokio::io::AsyncReadExt;
    let body = Body::from_stream(ReaderStream::with_capacity(file.take(len), MEDIA_READ_CHUNK));
    let mut resp = media_response_headers(
        if partial { StatusCode::PARTIAL_CONTENT } else { StatusCode::OK },
        &mime,
        total,
        start,
        end,
        len,
        partial,
    );
    *resp.body_mut() = body;
    resp
}

async fn media_get_preview(State(svc): State<MediaSvc>, Query(q): Query<MediaQuery>, headers: HeaderMap) -> Response {
    if !check_token(&svc, &q.t) {
        return (StatusCode::FORBIDDEN, ()).into_response();
    }
    let path = match resolve_media_path(&svc, &q) {
        Ok(p) => p,
        Err(code) => return (code, ()).into_response(),
    };
    let fname = path.file_name().and_then(|s| s.to_str()).unwrap_or("").to_string();
    let must_webp = crate::remotebrowse::preview_transcode_for_client(&fname);
    let webp = tokio::task::spawn_blocking(move || {
        let meta = std::fs::metadata(&path).map_err(|_| StatusCode::NOT_FOUND)?;
        if !meta.is_file() {
            return Err(StatusCode::NOT_FOUND);
        }
        match crate::remotebrowse::maybe_preview_webp(&path, meta.len()) {
            Ok(Some(v)) => Ok(Some(v)),
            Ok(None) => Ok(None),
            Err(e) => {
                if must_webp {
                    tracing::warn!(path = %path.display(), "media preview webp: {e}");
                    Err(StatusCode::UNPROCESSABLE_ENTITY)
                } else {
                    Ok(None)
                }
            }
        }
    })
    .await;
    let webp = match webp {
        Ok(Ok(v)) => v,
        Ok(Err(code)) => return (code, ()).into_response(),
        Err(_) => return (StatusCode::INTERNAL_SERVER_ERROR, ()).into_response(),
    };
    match webp {
        Some(data) if !data.is_empty() => serve_static_bytes(&data, "image/webp", &headers),
        _ => media_get(State(svc), Query(MediaQuery { preview: 0, ..q }), headers).await,
    }
}

async fn media_get_hq(State(svc): State<MediaSvc>, Query(q): Query<MediaQuery>, headers: HeaderMap) -> Response {
    if !check_token(&svc, &q.t) {
        return (StatusCode::FORBIDDEN, ()).into_response();
    }
    let path = match resolve_media_path(&svc, &q) {
        Ok(p) => p,
        Err(code) => return (code, ()).into_response(),
    };
    let webp = tokio::task::spawn_blocking(move || {
        let meta = std::fs::metadata(&path).map_err(|_| StatusCode::NOT_FOUND)?;
        if !meta.is_file() {
            return Err(StatusCode::NOT_FOUND);
        }
        crate::remotebrowse::maybe_hq_webp(&path, meta.len())
            .map_err(|e| {
                tracing::warn!(path = %path.display(), "media hq webp: {e}");
                StatusCode::UNPROCESSABLE_ENTITY
            })
    })
    .await;
    match webp {
        Ok(Ok(data)) if !data.is_empty() => serve_static_bytes(&data, "image/webp", &headers),
        Ok(Ok(_)) => (StatusCode::NOT_FOUND, ()).into_response(),
        Ok(Err(code)) => (code, ()).into_response(),
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, ()).into_response(),
    }
}

fn serve_static_bytes(data: &[u8], mime: &str, headers: &HeaderMap) -> Response {
    let total = data.len() as u64;
    if total == 0 {
        return (StatusCode::NOT_FOUND, ()).into_response();
    }
    let (start, end, partial) = parse_range(headers.get(header::RANGE).and_then(|v| v.to_str().ok()), total);
    let start = start as usize;
    let end = end as usize;
    let slice = &data[start..=end.min(data.len().saturating_sub(1))];
    let len = slice.len() as u64;
    let mut resp = media_response_headers(
        if partial { StatusCode::PARTIAL_CONTENT } else { StatusCode::OK },
        mime,
        total,
        start as u64,
        end as u64,
        len,
        partial,
    );
    *resp.body_mut() = Body::from(slice.to_vec());
    resp
}

async fn media_meta(svc: &MediaSvc, q: &MediaQuery) -> Result<(String, u64), StatusCode> {
    if !check_token(svc, &q.t) {
        return Err(StatusCode::FORBIDDEN);
    }
    let path = resolve_media_path(svc, q)?;
    let meta = tokio::fs::metadata(&path).await.map_err(|_| StatusCode::NOT_FOUND)?;
    if !meta.is_file() {
        return Err(StatusCode::NOT_FOUND);
    }
    let total = meta.len().saturating_sub(q.off);
    let mime = mime_guess::from_path(&path).first_or_octet_stream().to_string();
    Ok((mime, total))
}

fn resolve_media_path(svc: &MediaSvc, q: &MediaQuery) -> Result<std::path::PathBuf, StatusCode> {
    svc.resolve_path(&q.path).map_err(|_| StatusCode::BAD_REQUEST)
}

fn check_token(svc: &MediaSvc, t: &str) -> bool {
    let expect = svc.media_token();
    expect.is_empty() || t == expect
}

fn parse_range(header: Option<&str>, total: u64) -> (u64, u64, bool) {
    if total == 0 {
        return (0, 0, false);
    }
    let Some(h) = header else {
        return (0, total - 1, false);
    };
    let spec = h.strip_prefix("bytes=").unwrap_or(h).split(',').next().unwrap_or("").trim();
    let Some((a, b)) = spec.split_once('-') else {
        return (0, total - 1, false);
    };
    if a.is_empty() {
        let suffix: u64 = b.parse().unwrap_or(0);
        let start = total.saturating_sub(suffix);
        return (start, total - 1, true);
    }
    let start: u64 = a.parse().unwrap_or(0).min(total - 1);
    let end = if b.is_empty() {
        total - 1
    } else {
        b.parse().unwrap_or(total - 1).min(total - 1)
    };
    (start.max(0), end.max(start), true)
}

fn media_response_headers(
    status: StatusCode,
    mime: &str,
    total: u64,
    start: u64,
    end: u64,
    len: u64,
    partial: bool,
) -> Response {
    let mut r = Response::builder().status(status);
    r = r.header(header::ACCEPT_RANGES, "bytes");
    r = r.header(header::CONTENT_TYPE, mime);
    if partial {
        r = r.header(header::CONTENT_RANGE, format!("bytes {start}-{end}/{total}"));
    }
    r = r.header(header::CONTENT_LENGTH, len.to_string());
    r.body(Body::empty()).unwrap()
}
