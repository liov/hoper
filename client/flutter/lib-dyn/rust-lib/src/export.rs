use std::os::raw::{c_char, c_int};
use std::ffi::{CString, CStr};
use axum::{
    routing::{get, post},
    http::StatusCode,
    response::IntoResponse,
    Json, Router,
};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;
#[cfg(target_os="android")]
use android_logger::Config;
use axum::extract::Path;
use lazy_static::lazy_static;
use slog::{Drain, Logger, o, slog_o};
use log::{info, Level, LevelFilter, Metadata, Record};

lazy_static! {
    static ref LOG:Logger = {
        let decorator = slog_term::TermDecorator::new().build();
        let drain = slog_term::FullFormat::new(decorator).build().fuse();
        let drain = slog_async::Async::new(drain).build().fuse();

        slog::Logger::root(drain, o!())
        };
}

#[cfg(target_os="android")]
fn log_init() {
    android_logger::init_once(
        Config::default()
            .with_min_level(log::Level::Trace)
            .format(|f, record| write!(f, "native: {}", record.args()))
    );
    info!("log init");
}

struct SimpleLogger;

impl log::Log for SimpleLogger {
    fn enabled(&self, metadata: &Metadata) -> bool {
        metadata.level() <= Level::Info
    }

    fn log(&self, record: &Record) {
        if self.enabled(record.metadata()) {
            println!("{} - {}", record.level(), record.args());
        }
    }

    fn flush(&self) {}
}

static LOGGER: SimpleLogger = SimpleLogger;
#[cfg(target_os="windows")]
fn log_init(){
    log::set_logger(&LOGGER)
        .map(|()| log::set_max_level(LevelFilter::Info));
    info!("log init");
}


#[no_mangle]
pub extern "C" fn server(port: c_int) {
    log_init();
    // initialize tracing
    //tracing_subscriber::fmt::init();
    info!("tracing_subscriber init");
    tokio::runtime::Builder::new_current_thread()
        .enable_all()
        .build()
        .unwrap()
        .block_on(start(port as u16));
}

#[no_mangle]
pub extern "C" fn test(port: c_int) -> c_int {
    port + 1
}


pub async fn start(port: u16) {
    // build our application with a route
    let app = Router::new()
        // `GET /` goes to `root`
        .route("/", get(root))
        // `POST /users` goes to `create_user`
        .route("/users", post(create_user))
        .route("/user/:id", get(get_user));

    // run our app with hyper
    // `axum::Server` is a re-export of `hyper::Server`
    let addr = SocketAddr::from(([127, 0, 0, 1], port));
    info!("listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

// basic handler that responds with a static string
async fn root() -> &'static str {
    "Hello, World!"
}

async fn create_user(
    // this argument tells axum to parse the request body
    // as JSON into a `CreateUser` type
    Json(payload): Json<CreateUser>,
) -> impl IntoResponse {
    // insert your application logic here
    let user = User {
        id: 1337,
        username: payload.username,
    };

    // this will be converted into a JSON response
    // with a status code of `201 Created`
    (StatusCode::CREATED, Json(user))
}

async fn get_user(Path(user_id): Path<u64>) -> Json<User> {
    let user = User {
        id: user_id,
        username: "随机的".to_string(),
    };
    Json(user)
}

// the input to our `create_user` handler
#[derive(Deserialize)]
struct CreateUser {
    username: String,
}

// the output to our `create_user` handler
#[derive(Serialize)]
struct User {
    id: u64,
    username: String,
}