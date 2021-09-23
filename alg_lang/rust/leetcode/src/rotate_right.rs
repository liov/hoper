///旋转链表
use std::mem;
use crate::list::ListNode;

impl ListNode {
    pub fn len(&self) -> usize {
        fn get_len(list: &ListNode, len: usize) -> usize {
            if let Some(ref next) = list.next {
                get_len(next, len + 1)
            } else {
                return len;
            }
        }
        get_len(self, 1)
    }
}

pub fn rotate_right(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if k == 0 { return head; }
    let mut l = head.unwrap();

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> (Option<*mut Box<ListNode>>, Option<*mut Box<ListNode>>, i32) {
        return if let Some(ref mut next) = list.next {
            let (p2, p3, size) = get(next, len + 1, k);
            if len == size - k + 1 {
                return (Some(list), p3, size);
            } else if len == size - k {
                let mut t = list.clone();
                t.next.take();
                mem::replace(list, t);
            }
            (p2, p3, size)
        } else {
            (None, Some(list), len)
        }
    }

    let (p2, p3, len) = get(&mut l, 1, k);

    if k > len {
        return rotate_right(Some(l), k % len);
    };

    unsafe {
        mem::replace(&mut (*p3.unwrap()).next, Some(l.clone()));
        Some((*p2.unwrap()).clone())
    }
}

pub fn rotate_right2(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if head == None || k == 0 { return head; }
    let mut l = head.unwrap();

    let len = l.len() as i32;
    if len == 1 { return Some(l); };
    let mut k2 = k;
    if k >= len {
        k2 = k % len;
        if k2 == 0 {
            return Some(l);
        }
    };

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> (Option<*mut Box<ListNode>>, Option<*mut Box<ListNode>>) {
        if len == k {
            let mut next = list.next.take();

            if let Some(ref mut next) = next {
                return get(next, len + 1, k);
            }
        }
        if len == k + 1 {
            let (_, p3) = get(list, len + 1, k);
            return (Some(list), p3);
        }

        return if let Some(ref mut next) = list.next {
            get(next, len + 1, k)
        } else {
            (None, Some(list))
        }

    }


    let (p2, p3) = get(&mut l, 1, len - k2);

    let (p2, p3) = (p2.unwrap(), p3.unwrap());
    unsafe {
        mem::replace(&mut ((*p3).next), Some(l));
        Some((*p2).clone())
    }
}

pub fn rotate_right3(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if head == None || k == 0 { return head; }
    let mut l = head.unwrap();

    let len = l.len() as i32;
    if len == 1 { return Some(l); };
    let mut k2 = k;
    if k >= len {
        k2 = k % len;
        if k2 == 0 {
            return Some(l);
        }
    };

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> Box<ListNode> {
        if len == k {
            let mut next = list.next.take();

            if let Some(ref mut next) = next {
                return get(next, len + 1, k);
            }
        }
        if len == k + 1 {
            return list.clone();
        }

        if let Some(ref mut next) = list.next {
            return get(next, len + 1, k);
        }

        panic!("错误")
    }


    let mut list = get(&mut l, 1, len - k2);
    fn set(list: &mut Box<ListNode>, list2: Box<ListNode>) {
        if let Some(ref mut next) = list.next {
            set(next, list2);
        } else {
            list.next = Some(list2.clone());
        }
    }
    set(&mut list, l);
    Some(list)
}
