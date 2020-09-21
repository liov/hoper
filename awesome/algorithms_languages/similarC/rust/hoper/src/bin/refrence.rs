#[derive(Debug)]
struct Foo<'a> {
    a:&'a String
}

fn main(){
    let mut s1 = String::from("aaa");
    let mut f = Foo{a:&s1};
    println!("{:?}",f);
    let s2 = String::from("bbb");
    f.a = &s2;
    println!("{:?}",f);
    println!("{:?}",s1);
    s1 = String::from("ccc");
}