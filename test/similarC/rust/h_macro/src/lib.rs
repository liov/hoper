//the `#[proc_macro_derive]` attribute is only usable with crates of the `proc-macro` crate type
//`proc-macro` crate types cannot export any items other than functions tagged with `#[proc_macro_derive]` currently
//宏属性必须extern crate
extern crate proc_macro;
use crate::proc_macro::TokenStream;

use quote::quote;
use syn;


#[proc_macro_derive(Foo)]
pub fn foo_derive(input: TokenStream) -> TokenStream {
    // 构建 Rust 代码所代表的语法树
    // 以便可以进行操作
    let ast = syn::parse(input).unwrap();

    // 构建 trait 实现
    impl_foo(&ast)
}

fn impl_foo(ast: &syn::DeriveInput) -> TokenStream {
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


/*#[proc_macro_attribute]
pub fn route(attr: TokenStream, item: TokenStream) -> TokenStream {

}

#[proc_macro]
pub fn sql(input: TokenStream) -> TokenStream {

}*/
