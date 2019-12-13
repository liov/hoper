use std::env;
use std::thread;

use hello::protobuf::helloworld::*;
use hello::protobuf::helloworld_grpc::*;

use grpc::RequestOptions;
use grpc::SingleResponse;

struct GreeterImpl;

impl Greeter for GreeterImpl {
    fn say_hello(&self, _: RequestOptions, req: HelloRequest) -> SingleResponse<HelloReply> {
        let mut r = HelloReply::new();
        let name = if req.get_name().is_empty() {
            "world"
        } else {
            req.get_name()
        };
        println!("greeting request from {}", name);
        r.set_message(format!("Rust {}", name));
        SingleResponse::completed(r)
    }
}


fn is_tls() -> bool {
    env::args().any(|a| a == "--tls")
}

fn main() {
    let tls = is_tls();

    let port = if !tls { 50051 } else { 50052 };

    let mut server = grpc::ServerBuilder::new_plain();
    server.http.set_port(port);
    server.add_service(GreeterServer::new_service_def(GreeterImpl));
    //server.http.set_cpu_pool_threads(4);
    if tls {
        println!("tls")
    }
    let _server = server.build().expect("server");

    println!(
        "greeter server started on port {} {}",
        port,
        if tls { "with tls" } else { "without tls" }
    );

    loop {
        thread::park();
    }
}