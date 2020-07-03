fn main(){
    let mut a = 2i32;
    let b = &mut a;
    //let c = b; //c,move occurs because `b` has type `&mut i32`,这里是模式匹配，转移所有权
    let c:&mut i32 = &mut (*b); //编译通过,种种迹象表明这里是借用b也是借用*b
    //let c = &mut (*b);
    //println!("{:p}",&b);//cannot borrow `b` as immutable because it is also borrowed as mutable
    //*b = 1;
    *c = 3;//种种迹象表明，c含着b的引用
    println!("{:p}",&b);
    *b = 5;
    /*let b  = &a;
    let c = b;//编译通过*/
    //println!("{:p},{:p}",&b,&c);//cannot assign to `*b` because it is borrowed
    println!("{}",a);
    //println!("{}",c);//cannot borrow `b` as immutable because it is also borrowed as mutable
    //cannot assign to `*b` because it is borrowed
    //唯一解释的通的就是
    let mut a = 2i32;
    let b = &mut a;
    println!("{:p}",b);
    let c = b;
    println!("{:p}",c);
    let mut a = 2i32;
    let b = &mut a;
    println!("{:p}",b);
    let c:&mut i32 = b;
    println!("{:p}",c);

    let mut a = 2;
    let b = &mut a;
    //let c = b;    //这里值move给了c 这里是模式匹配,下面是赋值?
    let c:&mut i32 = b;//说明啥,:不单单是类型说明
    println!("{}",b);//虽然这里不能同时输出c,b
}

