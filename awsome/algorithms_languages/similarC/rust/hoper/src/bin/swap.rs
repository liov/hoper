fn main() {
    let mut t1 = ST::new(String::from("hello"));//这里有个拷贝过程
    t1.init();
    let mut t2 = ST::new(String::from("world"));
    t2.init();
    println!("{}-{};{}-{}", &t1.a, t1.b(), &t2.a, t2.b());
    std::mem::swap(&mut t1, &mut t2);
    println!("{}-{};{}-{}", &t1.a, t1.b(), &t2.a, t2.b())
}

struct ST {
    a: String,
    b: *const String,
}

impl ST {
    fn new(s: String) -> ST {
        ST { a: s, b: std::ptr::null() }
    }
    fn init(&mut self){
        let p = &(self.a) as *const String;
        self.b = p;
    }
    fn b(&self) -> &String {
        return unsafe { &*(self.b) };
    }
}