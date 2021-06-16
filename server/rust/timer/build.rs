use std::env;
use std::process::Command;
use std::str;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let current_dir = r"D:\hoper\server\go";
    let mut args = ["list", "-m","-f","{{.Dir}}",""];
    args[4] = "github.com/grpc-ecosystem/grpc-gateway/v2";
    let  stdout = Command::new("go").args(&args).current_dir(current_dir).output().expect("failed to execute process").stdout;
    let gateway = str::from_utf8(&stdout).unwrap().strip_suffix('\n').unwrap();
    println!("gateway：{}", gateway);
    let googleapis = gateway.to_owned()+r"\third_party\googleapis";
    println!("googleapis：{}", googleapis);
    args[4] = "github.com/alta/protopatch";
    let stdout = Command::new("go").args(&args).current_dir(current_dir).output().expect("failed to execute process").stdout;
    let protopatch = str::from_utf8(&stdout).unwrap().strip_suffix('\n').unwrap();
    println!("protopatch：{}", protopatch);
    args[4] = "google.golang.org/protobuf";
    let stdout = Command::new("go").args(&args).current_dir(current_dir).output().expect("failed to execute process").stdout;
    let protobuf = str::from_utf8(&stdout).unwrap().strip_suffix('\n').unwrap();
    println!("protobuf：{}", protobuf);
    let gopath =env::var("GOPATH").unwrap()+r"\src";
    println!("gopath：{}", gopath);
    tonic_build::configure()
        .compile(
            &[
                "D:/hoper/proto/user/user.service.proto",
                "D:/hoper/proto/user/user.model.proto",
                "D:/hoper/proto/user/user.enum.proto",
            ],
            &[
                "D:/hoper/proto",
                "D:/hoper/proto/utils/proto",
                gateway,&googleapis,protopatch,protobuf,&gopath
            ],
        )?;
    Ok(())
}