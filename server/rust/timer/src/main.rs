use std::fs::File;
use std::io::Read;
use std::ops::Deref;
use std::sync::{Arc, Mutex, MutexGuard};
use axum::{routing::get, Router, Json};
use serde_json::{Value, json};
use timer::config::{CONFIG,Config};

#[tokio::main]
async fn main() {

    // build our application with a single route
    let app = Router::new().route("/", get(root))
        .route("/json", get(json).post(json))
        .route("/plain_text", get(plain_text));

    // run it with hyper on localhost:3000
    axum::Server::bind(&"0.0.0.0:3000".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn root() {}


async fn plain_text() -> &'static str {
    "foo"
}

async fn json() -> Json<&'static Mutex<Config>> {
    Json(&CONFIG)
}
