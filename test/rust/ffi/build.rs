extern crate cpp_build;

fn main() {
    //build具体到bin中的某个文件无效，猜测可能是库crate才有效
    cpp_build::build("src/lib.rs");
}
