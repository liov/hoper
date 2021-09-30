use std::any::{Any, TypeId};
use std::fmt::Debug;


//Trait/ &Trait是用更一致的Struct/ &Struct不是impl Trait/ &dyn Trait
//目前只找到dyn和impl的用法，至于为何这么用原因不清楚
fn load_config_impl(value:&(impl Any + Debug)) -> Vec<String>{
    let mut cfgs: Vec<String>= vec![];
    let value = value as &dyn Any;
    match value.downcast_ref::<String>() {
        Some(cfp) => cfgs.push(cfp.clone()),
        None => (),
    };
    cfgs
}
//dyn是为了区分struct和Trait
fn load_config_dyn(value: &dyn Any) -> Vec<String>{
    let mut cfgs: Vec<String>= vec![];
    match value.downcast_ref::<String>() {
        Some(cfp) => cfgs.push(cfp.clone()),
        None => (),
    };
    match value.downcast_ref::<Vec<String>>() {
        Some(v) => cfgs.extend_from_slice(&v),
        None =>(),
    }

    if cfgs.len() == 0 {
        panic!("No Config File");
    }
    cfgs
}

fn load_config_gen<T:Any+Debug>(value: &T) -> Vec<String>{
    let mut cfgs: Vec<String>= vec![];
    let value = value as &dyn Any;
    match value.downcast_ref::<String>() {
        Some(cfp) => cfgs.push(cfp.clone()),
        None => (),
    };
    cfgs
}

static PST:&Bar = &Bar{i: 0, s: String::new()};//calls in statics are limited to constant functions, tuple structs and tuple variants

fn main() {
    let cfp = "/etc/wayslog.conf".to_string();
    println!("{:?}", load_config_dyn(&cfp));
    println!("{:?}",load_config_impl(&cfp));
    let cfps = vec!["/etc/wayslog.conf".to_string(),
                    "/etc/wayslog_sec.conf".to_string()];
    println!("{:?}", load_config_dyn(&cfps));
    let mut foo = 1;
    test(&foo);
    let ptr = &mut foo;
    *ptr = 2;
    println!("{}",foo);
    let mut bar = Bar{ i: 0, s: "s".to_string()};
    let ptr_b = &mut  bar;
    //test(&ptr_b);//`bar` does not live long enough
    (*ptr_b).i = 2;
    println!("{:?}",bar);
    test(&PST)
}

fn test(v:&dyn Foo){
    v.bar::<i32>()
}
//这里需要的是Foo是static而不是&Foo，&Foo有正常的生命周期
pub trait Foo:'static {
    fn type_id(&self) -> TypeId;
}

impl<T: 'static + ?Sized> Foo for T {
    fn type_id(&self) -> TypeId {
        TypeId::of::<T>()
    }
}

//error[E0038]: the trait `Foo` cannot be made into an object
//为什么Any可以这样实现
//这样是可以实现的，研究研究
impl dyn Foo {
    pub fn bar<T: Foo>(&self){
        println!("实现dyn Trait")
    }
}

#[derive(Debug)]
struct Bar{
    i:i32,
    s:String,
}

/*
impl Trait
impl Foo for dyn Foo {
    fn bar(){println!("foo bar") }
}*/

/*
fn foo<T>(add: u8) -> impl Fn(u8) -> u8
{
    move |origin: u8| {
        origin + add
    }
}
该语法表示返回值是一个满足其指明的 trait 的约束的具体类型。另外，由于这个实现是该函数返回值自行指定，还解决了某些场景使用泛型时的一些问题，比如上面的代码例子中，我们使用了泛型，而泛型的实际类型是由调用者决定的，这在使用装箱语法时会报错，虽然你在返回时通过 where 指明了泛型 T 的约束，但那并不是指示泛型具体类型的。

而通过 impl Trait 则是一个具体类型，且由返回者指定。

不过我想看了这部分内容的，都可能还有点模糊的地方，就是 调用者指定类型 ，或者说有没有更直观例子来辨别，当然有，我们以官方对于这部分说明的例子来写：

trait Trait {}

fn foo<T: Trait>(arg: T) {
}

fn foo(arg: impl Trait) {
}
上述两种实现，前者是泛型，表示 T 泛型需要是一个 Trait 的实现，后者不是泛型但也是要求满足 Trait 约束，这两个乍一看是类似的，实际却大不一样。

使用泛型时我们说，其类型是调用者决定，具体代码上体现就是我们可以这样调用 foo::<usize>(1) 表示我们传入的参数是 usize 类型，亦或 let a: usize = 1 然后 foo(a) ，在编译时，编译器会将泛型转换为实际被调用者传入的类型： usize ，这就是所谓的调用者决定其类型。

而对于 impl Trait 这种形式，则无需调用者指定，仅需要保证满足约束即可。
dyn 来解决另外的歧义
我们在之前例子中，说过在没有 impl Trait 这种语法糖之前，需要靠装箱解决问题。这个地方其实还有一个问题，我们看代码：

fn my_function() -> Box<Foo> {
    // ...
}
上述代码存在一个歧义， Foo 到底是 trait 还是一个具体的类型？这两个是有明显区别的。通过新的关键字来明确两者差异。以下是官方例子：

trait Trait {}

impl Trait for i32 {}

// old
fn function1() -> Box<Trait> {
}

// new
fn function2() -> Box<dyn Trait> {
}
*/
