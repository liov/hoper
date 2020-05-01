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
