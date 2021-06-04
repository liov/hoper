use hoper::utils::tree::{MyTree, TTree, TreeT};
//如果项目是二进制 crate 并且只包含 src/main.rs 而没有 src/lib.rs，
// 这样就不可能在 tests 目录创建集成测试并使用 extern crate 导入 src/main.rs 中定义的函数。
// 只有库 crate 才会向其他 crate 暴露了可供调用和使用的函数；二进制 crate 只意在单独运行。

//为什么 Rust 二进制项目的结构明确采用 src/main.rs 调用 src/lib.rs 中的逻辑的方式？
// 因为通过这种结构，集成测试 就可以 通过 extern crate 测试库 crate 中的主要功能了，
// 而如果这些重要的功能没有问题的话，src/main.rs 中的少量代码也就会正常工作且不需要测试。
#[test]
fn add_tree() {
        let x = std :: f64 :: consts :: PI;
        let mut t = MyTree::new();
        t.insert(x);
        println!("{:?}",t);
        t.insert(3f64);
        println!("{:?}",t);
        t.insert(5f64);
        println!("{:?}",t);
        t.insert(9f64);
        println!("{:?}",t);
        t.insert(66f64);
        println!("{:?}",t);
        t.insert(18f64);
        println!("{:?}",t);
        t.insert(12f64);
        println!("{:?}",t);
        t.insert(111f64);
        println!("{:?}",t);
        t.insert(12f64);
        println!("{:?}",t);
        t.insert(2f64);
        println!("{:?}",t);
        t.peek();
}

#[test]
fn t_tree(){
        let mut a = TTree::new();
        a.insert(5);
        //打印--nocapture
        println!("{:?}",a)
}
