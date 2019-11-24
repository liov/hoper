async fn say_hello() {
    println!("Hello, world!");
}

#[async_std::main]
async fn main() {
    say_hello().await;
    let a = async { 1u8 };
    let b = async { 2u8 };
    assert_eq!(a.join(b).await, (1u8, 2u8));
}
