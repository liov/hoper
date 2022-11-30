pub mod user {
    tonic::include_proto!("user");
}

use user::{user_service_client::UserServiceClient, SendVerifyCodeReq};
pub use timer::{empty,response,request,any,oauth};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = UserServiceClient::connect("http://[::1]:50051").await?;

    let request = empty::Empty{};

    let response = client.verify_code(request).await?;

    println!("RESPONSE={:?}", response);

    Ok(())
}