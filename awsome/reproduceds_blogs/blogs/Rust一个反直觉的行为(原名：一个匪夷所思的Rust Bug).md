做这个测试源于群里有个小伙伴问了一个问题
```rust
let a = 2;
let b = &mut a;
let c = b;
```
代码有点小问题，a应该是可变的，上面代码b不再可用，当你试图打印b的时候
```rust
println!("{}",b);
```
会报错，告诉你，

borrow of moved value: `b`,move occurs because `b` has type `&mut i32`，&mut i32 没有实现Copy
但是，如果是下面这样，b仍然可用
```rust
let a = 2;
let b = &a;
let c = b;
```
为此还引起了重借（reborrow）的讨论

包括出现了这样的代码，第三行和第四行会出现不同的结果
```rust
let a = 2;
let b = &mut a;
let c = b;
let c:&mut i32 = b;
```
好奇心驱使，测试

代码如下
```rust
let mut a = 2;
let b = &mut a;
let c = b;
let c:&mut i32 = b;
println!("{}",b);
```
注释第三行，可以正常输入b，注释第四行，编译报错，b已被move

继续测试
```rust
let mut a = 2;
let b = &mut a;
let c:&mut i32 = b;
*c = 3;
println!("{}",c);
*b = 5;
println!("{}",a);
```
上面代码分别输出3,5，为什么上面输出c，下面输出a，这很好理解，第二行b拥有a的可变借用，一直到第6行，所以五行不能用不可变借用，只能输出指针指向的值了

这段代码的问题在于，b和c都可以改变a，意味着b和c同时拥有了a的可变引用，但是这违背了rust的借用规则，变量的可变引用只能有一个

不行，这段违反直觉的代码的错误过于明显，我需要找出问题所在，继续测试，输出b和c的地址
```rust
println!("{:p},{:p}",&b,&c);
```
问题就出在了这里，编译不通过cannot borrow `b` as immutable because it is also borrowed as mutable

第三行，b被c可变借用，输出b的地址又需要不可变借用，于是报错

那么问题很明显了，let c：&mut i32 =b 这样的写法其实是借用了b，也就是获取到的是b的引用，为了验证，将i32改为usize，报错，想法验证失败

看下编译器报错吧，E502，看官方文档，例子很简单
```rust
This error indicates that you are trying to borrow a variable as mutable when it has already been borrowed as immutable.

Example of erroneous code:

ⓘ
fn bar(x: &mut i32) {}
fn foo(a: &mut i32) {
    let ref y = a; // a is borrowed as immutable.
    bar(a); // error: cannot borrow `*a` as mutable because `a` is also borrowed
            //        as immutable
}Run
To fix this error, ensure that you don't have any other references to the variable before trying to access it mutably:

fn bar(x: &mut i32) {}
fn foo(a: &mut i32) {
    bar(a);
    let ref y = a; // ok!
}Run
For more information on the rust ownership system, take a look at https://doc.rust-lang.org/book/ch04-02-references-and-borrowing.html.
```
一目了然，let c：&mut i32的写法应当与 let ref y 相当，而 y的类型是&&mut i32，很明显了，此处取的是b的引用，验证一下
```rust
let mut a = 2i32;
let b = &mut a;
let c:&mut i32 = b;
let ref d = b;
println!("{}",d==c);
```
完美的编译报错，can't compare `&mut i32` with `i32`

c取的并不是b的引用，当然了，从代码上看，c也不可能是b的引用，这样就奇怪了，为啥会报502

继续打印地址，打印地址由于b502，需要一个中间变量存地址
```rust
let mut a = 2i32;
let b = &mut a;
let p = &b as *const &mut i32 as usize;
let c = b;
println!("0x{:x},{:p}",p,&c);
```
结果0x80fbc8,0x80fbd8
```rust
let mut a = 2i32;
let b = &mut a;
let p = &b as *const &mut i32 as usize;
let c:&mut i32 = b;
println!("0x{:x},{:p}",p,&c);
```
结果0x80fbc8,0x80fbd8

很遗憾，并看不出来有什么问题，只能归结为，bug

----------------------------------------------------------------------------------------------------------------------------------------

好吧，继续看了看，上边提到了同时存在可变借用是错的，由于NLL的缘故，其实可变借用并不是同时存在的，之所以貌似存在是因为先*c之后没用c所以可以*b，如果先*b，c拿着b的可变借用编译报错

let c:&mut i32 = b;  - borrow of `*b` occurs here

唯一解释的通的就是，c有*b的引用，&（*b），强制解引用多态。&mut i32 -> &mut i32

经验证let c : &mut i32 = b 与 let c : &mut i32 = &mut (*b)行为一致，也与let c = &mut (*b)行为一致

借用*b的同时也借用着b