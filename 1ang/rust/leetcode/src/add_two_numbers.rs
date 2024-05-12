use crate::list::ListNode;

///两数相加

//执行用时 : 4 ms, 在Add Two Numbers的Rust提交中击败了100.00% 的用户

//内存消耗 : 2 MB, 在Add Two Numbers的Rust提交中击败了100.00% 的用户


pub fn add_two_numbers(l1: Option<Box<ListNode>>, l2: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut next1 = l1;
    let mut next2 = l2;
    let mut result = Box::new(ListNode::new(0));
    let mut first = true;
    let mut carry = 0;
    let mut sum;
    loop {
        let lx = next1.unwrap_or(Box::new(ListNode::new(0)));
        let ly = next2.unwrap_or(Box::new(ListNode::new(0)));

        sum = lx.val + ly.val;

        if first {
            result.val = sum % 10;
            first = false;
            if sum >= 10 { carry = 1 } else { carry = 0 }
        } else {
            if sum % 10 + carry == 10 {
                result.push(0);
                carry = 1;
            } else {
                result.push(sum % 10 + carry);
                if sum >= 10 { carry = 1 } else { carry = 0 }
            }
        }

        next1 = lx.next;
        next2 = ly.next;
        if next1 == None && next2 == None && carry == 0 {
            break;
        }
    }
    Some(result)
}



