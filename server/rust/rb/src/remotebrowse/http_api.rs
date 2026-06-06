//! 远程浏览数据面：HTTP/2 + `application/protobuf`（非 gRPC 帧）；载荷可选 zstd。
use axum::body::Bytes;
use axum::extract::State;
use axum::http::{header, HeaderMap, StatusCode};
use axum::response::Response;
use axum::routing::{get, post};
use axum::Router;
use prost::Message;

use crate::proto_zstd::{self, proto_err, proto_ok};
use crate::remotebrowse_svc::proto::{
    DeleteFileRequest, DeleteFileResponse, HealthResponse, ListFilesRequest, ListFilesResponse,
    ListMediaRequest, ListMediaResponse, ReadFileRequest, ReadFileResponse, ThumbnailBatchRequest,
    ThumbnailBatchResponse, ThumbnailRequest, ThumbnailResponse,
};
use crate::remotebrowse_svc::MediaSvc;

use crate::remotebrowse::thumb_meta::RB_PROTO;

pub fn router(svc: MediaSvc) -> Router {
    Router::new()
        .route("/rb/health", get(health))
        .route("/rb/v1/list/media", post(list_media))
        .route("/rb/v1/list", post(list_files))
        .route("/rb/v1/thumb", post(get_thumbnail))
        .route("/rb/v1/thumb/batch", post(thumb_batch))
        .route("/rb/v1/read", post(read_file))
        .route("/rb/v1/delete", post(delete_file))
        .with_state(svc)
}

async fn health(State(svc): State<MediaSvc>, headers: HeaderMap) -> Response {
    match svc.pb_health().await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_err(StatusCode::INTERNAL_SERVER_ERROR, e),
    }
}

async fn list_files(State(svc): State<MediaSvc>, headers: HeaderMap, body: Bytes) -> Response {
    let raw = match proto_zstd::decode_body(&headers, &body) {
        Ok(v) => v,
        Err(e) => return proto_err(StatusCode::BAD_REQUEST, e),
    };
    let req = match ListFilesRequest::decode(raw.as_slice()) {
        Ok(v) => v,
        Err(_) => return proto_err(StatusCode::BAD_REQUEST, "bad ListFilesRequest"),
    };
    match svc.pb_list_files(req).await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_err(StatusCode::BAD_REQUEST, e),
    }
}

async fn list_media(State(svc): State<MediaSvc>, headers: HeaderMap, body: Bytes) -> Response {
    let raw = match proto_zstd::decode_body(&headers, &body) {
        Ok(v) => v,
        Err(e) => return proto_err(StatusCode::BAD_REQUEST, e),
    };
    let req = match ListMediaRequest::decode(raw.as_slice()) {
        Ok(v) => v,
        Err(_) => return proto_err(StatusCode::BAD_REQUEST, "bad ListMediaRequest"),
    };
    match svc.pb_list_media(req).await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_err(StatusCode::BAD_REQUEST, e),
    }
}

async fn get_thumbnail(State(svc): State<MediaSvc>, headers: HeaderMap, body: Bytes) -> Response {
    let raw = match proto_zstd::decode_body(&headers, &body) {
        Ok(v) => v,
        Err(e) => return proto_err(StatusCode::BAD_REQUEST, e),
    };
    let req = match ThumbnailRequest::decode(raw.as_slice()) {
        Ok(v) => v,
        Err(_) => return proto_err(StatusCode::BAD_REQUEST, "bad ThumbnailRequest"),
    };
    match svc.pb_get_thumbnail(req).await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_err(StatusCode::BAD_REQUEST, e),
    }
}

async fn thumb_batch(State(svc): State<MediaSvc>, headers: HeaderMap, body: Bytes) -> Response {
    let raw = match proto_zstd::decode_body(&headers, &body) {
        Ok(v) => v,
        Err(e) => return proto_err(StatusCode::BAD_REQUEST, e),
    };
    let req = match ThumbnailBatchRequest::decode(raw.as_slice()) {
        Ok(v) => v,
        Err(_) => return proto_err(StatusCode::BAD_REQUEST, "bad ThumbnailBatchRequest"),
    };
    match svc.pb_thumb_batch(req).await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_err(StatusCode::INTERNAL_SERVER_ERROR, e),
    }
}

async fn read_file(State(svc): State<MediaSvc>, headers: HeaderMap, body: Bytes) -> Response {
    let raw = match proto_zstd::decode_body(&headers, &body) {
        Ok(v) => v,
        Err(e) => return proto_err(StatusCode::BAD_REQUEST, e),
    };
    let req = match ReadFileRequest::decode(raw.as_slice()) {
        Ok(v) => v,
        Err(e) => {
            tracing::warn!(len = raw.len(), "bad ReadFileRequest: {e}");
            return proto_err(StatusCode::BAD_REQUEST, format!("bad ReadFileRequest: {e}"));
        }
    };
    let path = req.path.clone();
    match svc.pb_read_file(req).await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_ok(
            &headers,
            ReadFileResponse {
                path,
                error: e,
                eof: true,
                ..Default::default()
            },
            RB_PROTO,
        ),
    }
}

async fn delete_file(State(svc): State<MediaSvc>, headers: HeaderMap, body: Bytes) -> Response {
    let raw = match proto_zstd::decode_body(&headers, &body) {
        Ok(v) => v,
        Err(e) => return proto_err(StatusCode::BAD_REQUEST, e),
    };
    let req = match DeleteFileRequest::decode(raw.as_slice()) {
        Ok(v) => v,
        Err(_) => return proto_err(StatusCode::BAD_REQUEST, "bad DeleteFileRequest"),
    };
    let path = req.path.clone();
    match svc.pb_delete_file(req).await {
        Ok(resp) => proto_ok(&headers, resp, RB_PROTO),
        Err(e) => proto_ok(
            &headers,
            DeleteFileResponse { path, error: e },
            RB_PROTO,
        ),
    }
}
