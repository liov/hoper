fn main() {
    build();
    println!("build successfully");
    println!("cargo:rerun-if-changed=build.rs");
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