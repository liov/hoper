extern crate cpp_build;

fn main() {
    build();
    println!("build successfully");
    println!("cargo:rerun-if-changed=build.rs");
    //build具体到bin中的某个文件无效，猜测可能是库crate才有效
    cpp_build::build("src/lib.rs");
}

#[cfg(windows)]
fn build(){
    //gcc -shared -o ffi.dll ffi.c
    cc::Build::new()
        .file("c/ffi.c")
        .define("FOO", Some("bar"))
        .include("src")
        .shared_flag(true)
        .static_flag(true)
        .compile("ffi");
}