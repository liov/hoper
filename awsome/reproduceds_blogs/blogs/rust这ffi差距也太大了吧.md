```rust
use libc::c_uint;
use libc::c_ulonglong;
use std::time::{Duration, SystemTime};
#[link(name = "ffi")]
extern {
    fn fibonacci(n: c_uint) -> c_ulonglong;
}


//n=43
//C_fib:1388,701408733
//Rust_fib:4173,701408733
fn main(){
    c_bib();
    rust_fib();
}


fn c_bib(){
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
```