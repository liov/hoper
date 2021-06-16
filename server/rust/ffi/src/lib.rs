#[macro_use]
extern crate cpp;

cpp!{{
    #include <stdint.h>
    #include <iostream>
}}


pub fn add(a: i32, b: i32) -> i32 {
    cpp!(unsafe [a as "int32_t", b as "int32_t"] -> i32 as "int32_t" {
        //printf不能用，原因未知
        std::cout << "adding " << a << " and " << b << std::endl;
        call_rust();
        return a + b;
    })
}

extern {
    pub fn call_rust();
}

cpp!{{
void call_rust() {
    rust!(cpp_call_rust [] {
        println!("This is in Rust!");
    });
}
}}



pub mod export;