#[macro_use]
extern crate log;

use std::sync::Arc;

use grpcio::{ChannelBuilder, EnvBuilder};
use hello_v2::log_util;
use hello_v2::protobuf::helloworld::HelloRequest;
use hello_v2::protobuf::helloworld_grpc::GreeterClient;


fn main() {
    let _guard = log_util::init_log(None);
    let env = Arc::new(EnvBuilder::new().build());
    let ch = ChannelBuilder::new(env).connect("localhost:50051");
    let client = GreeterClient::new(ch);

    let mut req = HelloRequest::default();
    req.set_name("world".to_owned());
    let reply = client.say_hello(&req).expect("rpc");
    info!("Greeter received: {}", reply.get_message());
}