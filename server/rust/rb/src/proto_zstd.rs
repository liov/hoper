//! HTTP protobuf 载荷 zstd 编解码（与 Flutter `RbZstdCodec` / `Accept-Encoding: zstd` 对齐）。
use axum::http::{header, HeaderMap, HeaderValue, StatusCode};
use axum::response::{IntoResponse, Response};

pub const ZSTD: &str = "zstd";
const MIN_BYTES: usize = 64;
const LEVEL: i32 = 3;

pub fn accepts_zstd(headers: &HeaderMap) -> bool {
    headers
        .get(header::ACCEPT_ENCODING)
        .and_then(|v| v.to_str().ok())
        .is_some_and(|s| {
            s.split(',')
                .any(|p| p.trim().split(';').next().unwrap_or("").eq_ignore_ascii_case(ZSTD))
        })
}

pub fn request_is_zstd(headers: &HeaderMap) -> bool {
    headers
        .get(header::CONTENT_ENCODING)
        .and_then(|v| v.to_str().ok())
        .is_some_and(|e| e.eq_ignore_ascii_case(ZSTD))
}

pub fn decode_body(headers: &HeaderMap, body: &[u8]) -> Result<Vec<u8>, String> {
    if request_is_zstd(headers) {
        zstd::decode_all(body).map_err(|e| format!("zstd decode: {e}"))
    } else {
        Ok(body.to_vec())
    }
}

/// 压缩请求体；仅当更短且超过 [MIN_BYTES] 时返回 `true`。
pub fn encode_request(plain: &[u8]) -> (Vec<u8>, bool) {
    if plain.len() < MIN_BYTES {
        return (plain.to_vec(), false);
    }
    match zstd::encode_all(plain, LEVEL) {
        Ok(c) if c.len() < plain.len() => (c, true),
        _ => (plain.to_vec(), false),
    }
}

pub fn encode_response(req_headers: &HeaderMap, plain: Vec<u8>) -> (Vec<u8>, bool) {
    if plain.len() < MIN_BYTES || !accepts_zstd(req_headers) {
        return (plain, false);
    }
    match zstd::encode_all(&plain[..], LEVEL) {
        Ok(c) if c.len() < plain.len() => (c, true),
        _ => (plain, false),
    }
}

pub fn proto_ok<M: prost::Message>(req_headers: &HeaderMap, msg: M, content_type: &'static str) -> Response {
    let plain = msg.encode_to_vec();
    let (body, zstd) = encode_response(req_headers, plain);
    let mut resp = (
        [(header::CONTENT_TYPE, HeaderValue::from_static(content_type))],
        body,
    )
        .into_response();
    if zstd {
        resp.headers_mut().insert(
            header::CONTENT_ENCODING,
            HeaderValue::from_static(ZSTD),
        );
    }
    resp
}

pub fn proto_err(status: StatusCode, msg: impl Into<String>) -> Response {
    (status, msg.into()).into_response()
}

pub fn decode_http_body(headers: &http::HeaderMap, body: &[u8]) -> Result<Vec<u8>, String> {
    let enc = headers.get(http::header::CONTENT_ENCODING).and_then(|v| v.to_str().ok());
    if enc.is_some_and(|e| e.eq_ignore_ascii_case(ZSTD)) {
        zstd::decode_all(body).map_err(|e| format!("zstd decode: {e}"))
    } else {
        Ok(body.to_vec())
    }
}

pub fn accepts_zstd_http(headers: &http::HeaderMap) -> bool {
    headers
        .get(http::header::ACCEPT_ENCODING)
        .and_then(|v| v.to_str().ok())
        .is_some_and(|s| {
            s.split(',')
                .any(|p| p.trim().split(';').next().unwrap_or("").eq_ignore_ascii_case(ZSTD))
        })
}
