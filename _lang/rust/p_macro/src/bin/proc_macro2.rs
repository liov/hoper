extern crate proc_macro;
use proc_macro2::{Ident, Span};
use quote::quote;
/// A wrapper around the procedural macro API of the compiler's proc_macro crate. This library serves two purposes:
///1.Bring proc-macro-like functionality to other contexts like build.rs and main.rs.
///Types from proc_macro are entirely specific to procedural macros and cannot ever exist in code outside of a procedural macro.
///Meanwhile proc_macro2 types may exist anywhere including non-macro code.
///By developing foundational libraries like syn and quote against proc_macro2 rather than proc_macro, the procedural macro ecosystem becomes easily applicable to many other use cases and we avoid reimplementing non-macro equivalents of those libraries.
///2.Make procedural macros unit testable.
///As a consequence of being specific to procedural macros, nothing that uses proc_macro can be executed from a unit test.
///In order for helper libraries or components of a macro to be testable in isolation, they must be implemented using proc_macro2.

const IGNORE: &str = stringify! {
 #[proc_macro_derive(HelloMacro)]
};

#[proc_macro_derive(MyDerive)]
pub fn my_derive(input: proc_macro::TokenStream) -> proc_macro::TokenStream {
    let input = proc_macro2::TokenStream::from(input);
    let ast :&syn::DeriveInput= syn::parse2(input).unwrap();
    let name = &ast.ident;

    let output = quote! {
        impl HelloMacro for #name {
            fn foo() {
                println!("Hello, Macro! My name is {}", stringify!(#name));
            }
        }
    };
    gen.into()
}


pub trait HelloMacro {
    fn hello_macro();
}

//#[derive(HelloMacro)]
struct Pancakes;

fn main() {
    println!("目前找不到使用意义")
}
