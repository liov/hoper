mod test{
    #[test]
    fn iter(){
        let a =vec![1,2,3,4,5];
        let b =vec![1,2,3,4,5];
        let c: Vec<i32> = a.iter().
            zip(b.iter().skip(1)).
            map(|(x,y)|x+y).
            collect();
        println!("{:?}",c);
    }
}

mod hoper{
     mod hash{
         pub mod map{
             use std::collections::hash_map::RandomState;
             //什么操作啊
             use hashbrown::hash_map as base;
             pub struct HashMap<K, V, S = RandomState> {
                 base: base::HashMap<K, V, S>,
             }
         }
     }
    mod hash_map{
        pub use super::hash::map::*;
    }
}