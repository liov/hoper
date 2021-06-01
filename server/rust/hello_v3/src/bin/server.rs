use tonic::{transport::Server, Request, Response, Status};


pub mod user {
    tonic::include_proto!("user");
}

use user::{
    user_service_server::{UserService, UserServiceServer},
    SendVerifyCodeReq,SingUpVerifyReq,SignupReq,LoginRep,
    EditReq,LoginReq,ResetPasswordReq,UserRep,ActionLogListReq,
    ActionLogListRep,BaseListRep,FollowReq,ActiveReq,UserAuthInfo,
    BaseListReq
};

pub use hello_v3::{empty,response,request,any,oauth};


#[derive(Default)]
pub struct MyUserService {
    data: String,
}

#[tonic::async_trait]
impl UserService for MyUserService {
    async fn verify_code(
        &self,
        request: Request<empty::Empty>,
    ) -> Result<Response<String>, Status> {
        println!("Got a request: {:?}", request);

        let string = &self.data;

        println!("My data: {:?}", string);

        Ok(Response::new( "Zomg, it works!".into()))
    }
    async fn send_verify_code(
        &self,
        request: Request<SendVerifyCodeReq>,
    ) -> Result<Response<empty::Empty>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn signup_verify(
        &self,
        request: Request<SingUpVerifyReq>,
    ) -> Result<Response<empty::Empty>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn signup(
        &self,
        request: Request<SignupReq>,
    ) -> Result<Response<empty::Empty>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn easy_signup(
        &self,
        request: Request<SignupReq>,
    ) -> Result<Response<LoginRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn active(
        &self,
        request: Request<ActiveReq>,
    ) -> Result<Response<LoginRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn edit(
        &self,
        request: Request<EditReq>,
    ) -> Result<Response<empty::Empty>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn login(
        &self,
        request: Request<LoginReq>,
    ) -> Result<Response<LoginRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn logout(
        &self,
        request: Request<empty::Empty>,
    ) -> Result<Response<empty::Empty>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn auth_info(
        &self,
        request: Request<empty::Empty>,
    ) -> Result<Response<UserAuthInfo>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn forget_password(
        &self,
        request: Request<LoginReq>,
    ) -> Result<Response<response::TinyRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn reset_password(
        &self,
        request: Request<ResetPasswordReq>,
    ) -> Result<Response<response::TinyRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn info(
        &self,
        request: Request<request::Object>,
    ) -> Result<Response<UserRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn action_log_list(
        &self,
        request: Request<ActionLogListReq>,
    ) -> Result<Response<ActionLogListRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn base_list(
        &self,
        request: Request<BaseListReq>,
    ) -> Result<Response<BaseListRep>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn follow(
        &self,
        request: Request<FollowReq>,
    ) -> Result<Response<empty::Empty>, Status>{
        Ok(Response::new(Default::default()))
    }
    async fn del_follow(
        &self,
        request: Request<FollowReq>,
    ) -> Result<Response<BaseListRep>, Status>{
        Ok(Response::new(Default::default()))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse().unwrap();
    let greeter = MyUserService::default();

    Server::builder()
        .add_service(UserServiceServer::new(greeter))
        .serve(addr)
        .await?;

    Ok(())
}