use std::env;
use std::process::Command;
use std::str;

const PROTOPATH: &str = r"D:\code\hopeio\hoper\proto";
const PROJECTPATH: &str = r"D:\code\hopeio\hoper\server\go";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let tiga = get_include_path("github.com/hopeio/tiga");
    println!("tiga：{}", tiga);
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
                tiga + "/protobuf/_proto",
            ],
        )?;
    Ok(())
}

fn get_include_path(path: &str) -> String {
    let args = ["list", "-m", "-f", "{{.Dir}}", path];
    let stdout = Command::new("go").args(&args).current_dir(PROJECTPATH).output().expect("failed to execute process").stdout;
    str::from_utf8(&stdout).unwrap().strip_suffix("\n").unwrap().to_owned()
}