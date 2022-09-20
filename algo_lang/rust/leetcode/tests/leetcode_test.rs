use leetcode::two_sum::*;
use leetcode::add_two_numbers::add_two_numbers;
use leetcode::multiply::multiply;
use leetcode::list::ListNode;
use leetcode::letter_combinations::letter_combinations;
use leetcode::rotate_right::rotate_right3;
use leetcode::h_trap_rain_water_ii::{trap_rain_water2, insert_sort};

#[test]
fn two_sum_test(){
    let nums = vec![0,4,3,0];
    let target = 0;
    assert_eq!(two_sum2(nums,target),[0,3])
}

#[test]
fn add_two_numbers_test(){
    let mut l1= Box::new(ListNode::new(2));
    l1.push(4);
    l1.push(3);
    let mut l2=Box::new(ListNode::new(5));
    l2.push(6);
    l2.push(4);
    let mut result=Box::new(ListNode::new(7));
    result.push(0);
    result.push(8);
    assert_eq!(add_two_numbers(Some(l1),Some(l2)),Some(result))
}

#[test]
fn multiply_test(){
    let s1 = String::from("9");
    let s2 = String::from("9");
    assert_eq!(multiply(s1,s2),  String::from("81"))
}

#[test]
fn multiply_test2(){
    let s1 = String::from("5");
    let s2 =s1.as_bytes();
    let b1:[u8;1]=[53];
    assert_eq!(*s2,b1 )
}


#[test]
fn letter_combinations_test(){
    assert_eq!(letter_combinations(String::from("23")),  vec![String::from("81")])
}

#[test]
fn list_len(){
    let mut l1= Box::new(ListNode::new(2));
    l1.push(4);
    l1.push(3);
    assert_eq!(l1.len(),  3)
}

#[test]
fn rotate_right_test(){
    let mut l1= Box::new(ListNode::new(1));
    l1.push(2);
    l1.push(3);
    l1.push(4);
    l1.push(5);
    let mut l2= Box::new(ListNode::new(3));
    l2.push(4);
    l2.push(5);
    l2.push(1);
    l2.push(2);
    assert_eq!(rotate_right3(Some(l1),3), Some(l2))
}

//leetcode上自带测试，这里就不写了

#[cfg(not(test))]
fn clone() -> i32 {
   5
}

#[cfg(test)]
fn clone() -> i32 {
    10
}

#[test]
fn clone_test(){
    assert_eq!(clone(), 10)
}

#[test]
fn trap_rain_water_test(){
    let t = vec![
        vec![1,4,3,1,3,2],
        vec![3,2,1,3,2,4],
        vec![2,3,3,2,3,1],
    ];
    assert_eq!(trap_rain_water2(t), 10)
}

#[test]
fn insert_sort_test(){
    let t1=vec![2,6,4];
    let mut t2 = vec![5,4,3,2,1];
    let t3= vec![6,6,5,4,3,2,2,1];
    insert_sort(&mut t2,t1);
    assert_eq!(t2, t3)
}
