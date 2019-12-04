pub mod utils;
pub mod math;
pub mod leetcode;

#[cfg(hoper)]
mod tests {
    use crate::math::add;

    #[hoper]
    fn add_two_a(){
        assert_eq!(4,add(1,3))
    }
}
