
fn main(){
    let a = 5;
    println!("{:p}",&a);
    let b = Some(a);
    println!("{:p}",&b);
    if let Some(c) = b {
        println!("{:p}",&c)
    }
}
