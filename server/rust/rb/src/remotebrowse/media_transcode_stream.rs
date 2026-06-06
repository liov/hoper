//! `GET /rb/v1/transcode/stream.ts`：从 `start_ms` 起连续输出 MPEG-TS（chunked）。
use std::convert::Infallible;
use axum::body::Body;
use axum::extract::{Query, State};
use axum::http::{header, StatusCode};
use axum::response::{IntoResponse, Response};
use axum::routing::get;
use axum::Router;
use tokio::sync::mpsc;
use tokio_stream::{wrappers::ReceiverStream, StreamExt};

use crate::remotebrowse::media_transcode::{check_token, preset_params, TranscodeQuery};
use crate::remotebrowse::transcode_stream;
use crate::remotebrowse::transcode_vcodec::TranscodeVcodec;
use crate::remotebrowse_svc::MediaSvc;

const STREAM_CHUNK_CHAN: usize = 24;

pub fn stream_route() -> Router<MediaSvc> {
    Router::new().route("/rb/v1/transcode/stream.ts", get(transcode_stream_http))
}

async fn transcode_stream_http(State(svc): State<MediaSvc>, Query(q): Query<TranscodeQuery>) -> Response {
    if q.off > 0 {
        return (StatusCode::NOT_IMPLEMENTED, ()).into_response();
    }
    if preset_params(&q.preset).is_none() {
        return (StatusCode::BAD_REQUEST, ()).into_response();
    }
    if !check_token(&svc, &q.t) {
        return (StatusCode::FORBIDDEN, ()).into_response();
    }
    let input = match svc.resolve_path(&q.path) {
        Ok(p) => p,
        Err(_) => return (StatusCode::BAD_REQUEST, ()).into_response(),
    };
    let start_ms = stream_start_ms(&q);
    let preset = match preset_params(&q.preset) {
        Some(p) => p,
        None => return (StatusCode::BAD_REQUEST, ()).into_response(),
    };
    let vcodec = TranscodeVcodec::parse_or_default(&q.vcodec);
    let path_log = q.path.clone();
    let (tx, rx) = mpsc::channel::<Vec<u8>>(STREAM_CHUNK_CHAN);
    tokio::task::spawn_blocking(move || {
        let r = transcode_stream::transcode_stream(&input, start_ms, preset, vcodec, |chunk| {
            if tx.blocking_send(chunk.to_vec()).is_err() {
                return Err("cancelled".to_string());
            }
            Ok(())
        });
        if let Err(e) = r {
            tracing::warn!(%e, start_ms, path = %path_log, "transcode stream");
        }
    });
    let body = Body::from_stream(ReceiverStream::new(rx).map(|c| Ok::<_, Infallible>(c)));
    Response::builder()
        .status(StatusCode::OK)
        .header(header::CONTENT_TYPE, "video/MP2T")
        .header(header::CACHE_CONTROL, "no-store")
        .header(header::CONNECTION, "keep-alive")
        .body(body)
        .unwrap()
}

fn stream_start_ms(q: &TranscodeQuery) -> i64 {
    let raw = q.start_ms.or(q.begin_ms).unwrap_or(0).max(0);
    let step = crate::remotebrowse::transcode_preset::FRAGMENT_MS;
    (raw / step) * step
}
