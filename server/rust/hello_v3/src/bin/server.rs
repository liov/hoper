use tonic::{transport::Server, Request, Response, Status};

pub mod user {
    tonic::include_proto!("user");
}

pub mod empty {
    tonic::include_proto!("empty");
}

use user::{
    user_service_server::{UserService, UserServiceServer},
    HelloReply, HelloRequest,
};
use empty::Empty;

#[derive(Default)]
pub struct MyUserService {
    data: String,
}

#[tonic::async_trait]
impl UserService for MyUserService {
    async fn verify_code(
        &self,
        request: tonic::Request<Empty>,
    ) -> Result<tonic::Response<::std::string::String>, tonic::Status> {
        println!("Got a request: {:?}", request);

        let string = &self.data;

        println!("My data: {:?}", string);

        let reply = hello_world::HelloReply {
            message: "Zomg, it works!".into(),
        };
        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse().unwrap();
    let greeter = MyGreeter::default();

    Server::builder()
        .serve(addr, GreeterServer::new(greeter))
        .await?;

    Ok(())
}