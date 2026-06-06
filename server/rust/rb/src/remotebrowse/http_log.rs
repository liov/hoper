//! HTTP 请求访问日志（method / path / query）；敏感参数脱敏；不记录 response。
use axum::body::Body;
use axum::http::Request;
use tower_http::trace::{OnRequest, OnResponse, TraceLayer};

pub fn layer() -> TraceLayer<
    tower_http::classify::SharedClassifier<tower_http::classify::ServerErrorsAsFailures>,
    impl Fn(&Request<Body>) -> tracing::Span + Clone,
    impl OnRequest<Body> + Clone,
    impl OnResponse<Body> + Clone,
> {
    TraceLayer::new_for_http()
        .make_span_with(noop_span)
        .on_request(|req: &Request<Body>, _span: &tracing::Span| {
            tracing::info!(
                method = %req.method(),
                path = %req.uri().path(),
                query = %format_query_for_log(req.uri()),
                "http request"
            );
        })
        .on_response(|_response: &axum::http::Response<Body>, _latency: std::time::Duration, _span: &tracing::Span| {})
}

fn noop_span(_req: &Request<Body>) -> tracing::Span {
    tracing::Span::none()
}

fn format_query_for_log(uri: &axum::http::Uri) -> String {
    let Some(raw) = uri.query() else {
        return String::new();
    };
    if raw.is_empty() {
        return String::new();
    }
    raw.split('&')
        .map(redact_query_pair)
        .collect::<Vec<_>>()
        .join("&")
}

fn redact_query_pair(pair: &str) -> String {
    let Some((k, _)) = pair.split_once('=') else {
        return pair.to_string();
    };
    if is_sensitive_query_key(k) {
        return format!("{k}=***");
    }
    pair.to_string()
}

fn is_sensitive_query_key(key: &str) -> bool {
    let k = key.to_ascii_lowercase();
    matches!(
        k.as_str(),
        "t" | "token" | "access_token" | "auth" | "password" | "passwd" | "secret"
    )
}
