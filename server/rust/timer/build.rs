use std::env;
use std::process::Command;
use std::str;

const PROTOPATH: &str = r"D:\code\hoper\proto";
const PROJECTPATH: &str = r"D:\code\hoper\server\go\lib";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let gateway = get_include_path("github.com/grpc-ecosystem/grpc-gateway/v2");
    println!("gateway：{}", gateway);
    let googleapis = get_include_path("github.com/googleapis/googleapis");
    println!("googleapis：{}", googleapis);
    let protopatch = get_include_path("github.com/alta/protopatch");
    println!("protopatch：{}", protopatch);
    let protobuf = get_include_path("google.golang.org/protobuf");
    println!("protobuf：{}", protobuf);
    let gopath = env::var("GOPATH").unwrap() + r"\src";
    println!("gopath：{}", gopath);
    tonic_build::configure()
        .server_mod_attribute("attrs", "#[cfg(feature = \"server\")]")
        .server_attribute("Echo", "#[derive(PartialEq)]")
        .client_mod_attribute("attrs", "#[cfg(feature = \"client\")]")
        .client_attribute("Echo", "#[derive(PartialEq)]")
        .compile(
            &[
                "/user/user.service.proto",
                "/user/user.model.proto",
                "/user/user.enum.proto",
            ].map(|v| PROTOPATH.to_owned() + v),
            &[
                PROTOPATH.to_owned(),
                PROJECTPATH.to_owned() + "/protobuf",
                PROJECTPATH.to_owned() + "/protobuf/third",
               // gateway, googleapis, protopatch, protobuf, gopath
            ],
        )?;
    Ok(())
}

fn get_include_path(path: &str) -> String {
    let args = ["list", "-m", "-f", "{{.Dir}}", path];
    let stdout = Command::new("go").args(&args).current_dir(PROJECTPATH).output().expect("failed to execute process").stdout;
    str::from_utf8(&stdout).unwrap().to_owned()
}