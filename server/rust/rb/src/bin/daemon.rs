#[tokio::main]
async fn main() {
    rb::tracing_init::init();
    rb::daemon::run().await;
}
