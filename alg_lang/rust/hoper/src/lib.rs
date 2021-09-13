pub mod utils;
pub mod math;
pub mod timer_future;

mod test {
    #[test]
    fn iter() {
        let a = vec![1, 2, 3, 4, 5];
        let b = vec![1, 2, 3, 4, 5];
        let c: Vec<i32> = a.iter().
            zip(b.iter().skip(1)).
            map(|(x, y)| x + y).
            collect();
        println!("{:?}", c);
    }
}


mod hash {
    pub mod map {
        use std::collections::hash_map::RandomState;
        //什么操作啊
        use hashbrown::hash_map as base;

        pub struct HashMap<K, V, S = RandomState> {
            base: base::HashMap<K, V, S>,
        }
    }
}

mod hash_map {
    pub use super::hash::map::*;
}


#[cfg(hoper)]
mod tests {
    use crate::math::add;

    #[hoper]
    fn add_two_a() {
        assert_eq!(4, add(1, 3))
    }
}
