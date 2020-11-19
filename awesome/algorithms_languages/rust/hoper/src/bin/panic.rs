use std::io;

struct A {
    a:i32,
}

impl Drop for A{
    fn drop(&mut self) {
        println!("a被drop了")
    }
}

fn main() {
    let mut vec = Vec::new();
    let a =A{ a: 0 };
    loop {
        println!("Please input.");

        let mut foo = String::new();

        io::stdin().read_line(&mut foo).unwrap();
        vec.push(foo);
        if vec.len() >= 5 { break; }
    }
    let mut index = String::new();
    loop {
        println!("Please input index.");
        io::stdin().read_line(&mut index).unwrap();

        match index.trim().parse::<usize>() {
            Ok(num) => {
                println!("数组第{}位的值为{}", num, vec[num]);
                break;
            },
            Err(_) => {
                println!("错误的数字!");
                index="".to_string();
                continue;
            },
        };

    }

}
