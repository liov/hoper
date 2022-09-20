// Definition for singly-linked list.
#[derive(PartialEq, Eq, Clone, Debug)]
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>,
}

impl ListNode {
    #[inline]
    pub fn new(val: i32) -> Self {
        ListNode {
            next: None,
            val,
        }
    }

    /*因为 Rust 是后向兼容的（backwards compatible），所以不会移除 ref 和 ref mut，同时它们在一些不明确的场景还有用，
    比如希望可变地借用结构体的部分值而可变地借用另一部分的情况。你可能会在老的 Rust 代码中看到它们，所以请记住它们仍有价值。*/
    pub fn push(&mut self, val: i32) {
        match self.next {
            Some(ref mut next) =>
                next.push(val),
            None => self.next = Some(Box::new(ListNode::new(val)))
        }
    }
}