use libc::c_int;
use libc::c_uint;
use libc::c_ulonglong;
use libc::size_t;
use std::ffi::{CStr,CString, c_void};
use std::os::raw::c_char;
use std::time::{Duration, SystemTime};
#[link(name = "ffi")]
extern {
    fn fibonacci(n: c_uint) -> c_ulonglong;
    fn your_func(arg1: c_int, arg2: *mut c_void) -> size_t; // 声明ffi函数
    fn your_func2(arg1: c_int, arg2: *mut c_void) -> size_t;
    static FFI_GLOBAL: c_int; // 声明ffi全局变量
    fn run_callback(data: i32, cb: extern fn(i32));
    fn char_func() -> *mut c_char;
    fn my_printer(s: *const c_char);
}

#[repr(C)]
struct RustObject {
    a: c_int,
    // other members
}

extern "C" fn callback(a: c_int) { // 这个函数是传给c调用的
    println!("hello {}!", a);
}

fn get_string() -> String {
    unsafe {
        let raw_string: *mut c_char = char_func();
        let cstr = CStr::from_ptr(raw_string);
        cstr.to_string_lossy().into_owned()
    }
}
//n=43
//C_fib:1388,701408733
//Rust_fib:4173,701408733
fn main(){
    /*let result: size_t = unsafe {
        your_func(1 as c_int, Box::into_raw(Box::new(3)) as *mut c_void)
    };
    let c_to_print = CString::new("Hello, world!").unwrap();
    unsafe {
        my_printer(c_to_print.as_ptr()); // 使用 as_ptr 将CString转化成char指针传给c函数
    }*/
    c_fib();
    rust_fib();
    call_dynamic();
}


fn c_fib(){
    let now = SystemTime::now();
    let value = unsafe{
      fibonacci(43)
    };
    println!("C_fib:{:?},{:?}", SystemTime::now().duration_since(now).unwrap().as_millis(),value);
}

fn rust_fib(){
    let now = SystemTime::now();
    let value = unsafe{
        rust_fibonacci(43)
    };
    println!("Rust_fib:{:?},{:?}", SystemTime::now().duration_since(now).unwrap().as_millis(),value);
}

fn rust_fibonacci(n:usize) -> u64{
    if n<2{return 1};
    rust_fibonacci(n-1)+rust_fibonacci(n-2)
}

fn call_dynamic() {
    let now = SystemTime::now();
    unsafe {
        let lib = libloading::Library::new("./ffi/c/ffi.o").unwrap();
        let func: libloading::Symbol<unsafe extern fn(u32) -> u64> = lib.get(b"fibonacci").unwrap();
        let value = func(43);
        println!("C_fib:{:?},{:?}", SystemTime::now().duration_since(now).unwrap().as_millis(),value);
    }
}