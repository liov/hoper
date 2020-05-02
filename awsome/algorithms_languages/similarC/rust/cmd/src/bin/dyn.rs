#![feature(box_syntax)]

trait Fly {
    fn fly(&self) -> bool;
}

struct Duck;

impl Fly for Duck{
    fn fly(&self) ->bool{
        true
    }
}

/*
fn foo<T>()->&'static T where T:Fly{
    &Duck
}
*/

fn foo()->&'static dyn Fly {
    &Duck
}

fn foo1()->impl Fly{
    Duck
}

fn foo2()-> Box<dyn Fly>{
    box Duck
}

fn foo3<T:Fly>(t:T)->T{t}

// fn foo<T:Fly>(t:T)->T{Duck} expected type parameter `T`, found struct `Duck`

fn  main(){
    //foo3(Duck) expected `()`, found struct `Duck`
    let a:Duck= foo3(Duck);
}
