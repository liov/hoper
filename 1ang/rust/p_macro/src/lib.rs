//the `#[proc_macro_derive]` attribute is only usable with crates of the `proc-macro` crate type
//`proc-macro` crate types cannot export any items other than functions tagged with `#[proc_macro_derive]` currently
//宏属性必须extern crate
//error: `proc-macro` crate types currently cannot export any items other than functions tagged with `#[proc_macro]`, `#[proc_macro_derive]`, or `#[proc_macro_attribute]`
//不能导出非过程宏函数
extern crate proc_macro;
//实测下面的已经不需要了，但是ide不友好，不添加无法识别跳转
//extern crate proc_macro2;
extern crate syn;
extern crate quote;
use proc_macro::TokenStream;

use quote::quote;
use syn::{parse_macro_input, DeriveInput,LitStr};



#[proc_macro_derive(Foo)]
pub fn foo_derive(input: TokenStream) -> TokenStream {
    // 构建 Rust 代码所代表的语法树
    // 以便可以进行操作
    let ast = syn::parse(input).unwrap();

    // 构建 trait 实现
    impl_foo(&ast)
}

fn impl_foo(ast: &DeriveInput) -> TokenStream {
    let name = &ast.ident;
    let gen = quote! {
        impl Foo for #name {
            fn foo() {
                println!("Hello, Macro! My name is {}", stringify!(#name));
            }
        }
    };
    gen.into()
}


#[proc_macro_attribute]
pub fn inject(args: TokenStream, body: TokenStream) -> TokenStream {
    let mut param = parse_macro_input!(args as LitStr).value();
    if param.is_empty(){
        panic!("填写属性!")
    }
    match syn::parse::<syn::Item>(body).unwrap() {
        syn::Item::Fn(ref func) => {
            let syn::ItemFn {
                attrs, vis, sig, block,
            } = func;
            use std::ops::Deref;
            let syn::Signature {
                ident,generics, inputs, output, ..
            } = sig.deref();
            let gen = quote! {
                #(#attrs)* #vis
                fn #ident #generics (#inputs) #output {
                    println!("you want to print {}", #param);
                    let result = #block;
                    println!("end");
                    return result;
                }
            };
            gen.into()
        }
        _ => panic!("Only fn is allowed!"),
    }
}

#[proc_macro]
pub fn hello_world(_: TokenStream) -> TokenStream {
    r#"println!("Hello, World!");"#.parse().unwrap()
}

#[proc_macro]
pub fn count_tt(ts: TokenStream) -> TokenStream {
    ts.into_iter().count().to_string().parse().unwrap()
}
