//cargo run -p cmd --bin passbyvalue
use leetcode::list::ListNode;
use leetcode::rotate_right::rotate_right;

fn main() {
    let mut l1= Box::new(ListNode::new(1));
    println!("初始：{:p}",l1.as_mut());
    l1.push(2);
    l1.push(3);
    l1.push(4);
    l1.push(5);
    let mut t = rotate_right(Some(l1),2);

    unsafe {   println!("环:{:?}",t);}
}
