use std::cmp::PartialOrd;
use std::fmt::{Debug, Display};

//本人二叉树
#[derive(Debug)]
pub struct MyTree<T: PartialOrd + Debug> {
    pub data: Option<Box<T>>,
    pub left: Option<Box<MyTree<T>>>,
    pub right: Option<Box<MyTree<T>>>,
}

impl<T: PartialOrd + Debug> MyTree<T> {
    pub fn insert(&mut self, data: T) {
        match self.data {
            Some(ref mut rdata) => {
                let node = if data < **rdata { &mut self.left } else { &mut self.right };
                match node {
                    Some(ref mut node) =>
                        node.insert(data),
                    None =>
                        {
                            *node = Some(Box::new(MyTree::new()));
                            node.as_mut().unwrap().insert(data)
                        }
                }
            }
            None => self.data = Some(Box::new(data)),
        }
    }

    pub fn new() -> MyTree<T> {
        MyTree {
            data: None,
            left: None,
            right: None,
        }
    }
    pub fn peek(&self) {
        match self.data {
            Some(ref data) => {
                if let Some(ref left) = self.left {
                    left.peek()
                }
                println!("{:?}", **data);
                if let Some(ref right) = self.right {
                    right.peek()
                }
            }
            None => {}
        }
    }
}


struct Node<T: PartialOrd + Debug> {
    elem: T,
    left: Option<Box<Node<T>>>,
    right: Option<Box<Node<T>>>,
}

struct Tree<T: PartialOrd + Debug> {
    root: Option<Box<Node<T>>>,
    length: usize,
}

impl<T: PartialOrd + Debug> Tree<T> {
    pub fn insert(&mut self, data: T) {
        match self.root {
            Some(ref mut node) =>
                node.insert(data),
            None => self.root = Some(Box::new(Node::new(data))),
        }
        self.length += 1;
    }

    fn new() -> Tree<T> {
        Tree { root: None, length: 0 }
    }
}

impl<T: PartialOrd + Debug> Node<T> {
    pub fn new(t: T) -> Self {
        Node {
            elem: t,
            left: None,
            right: None,
        }
    }

    pub fn insert(&mut self, data: T) {
        if self.elem > data {
            match self.left {
                Some(ref mut node) =>
                    node.insert(data),
                None => self.left = Some(Box::new(Node::new(data))),
            }
        } else {
            match self.right {
                Some(ref mut node) =>
                    node.insert(data),
                None => self.right = Some(Box::new(Node::new(data))),
            }
        }
    }
}

#[derive(Debug)]
pub struct TNode<T> {
    elem: T,
    left: TTree<T>,
    right: TTree<T>,
}

impl<T> TNode<T> {
    fn new(t: T) -> TNode<T> {
        TNode {
            elem: t,
            left: None,
            right: None,
        }
    }
}


pub type TTree<T> = Option<Box<TNode<T>>>;

pub trait TreeT<T> {
    fn new() -> TTree<T>;
    fn insert(&mut self, data: T);
}

/*E0117
This error indicates a violation of one of Rust's orphan rules for trait implementations.
The rule prohibits any implementation of a foreign trait (a trait defined in another crate) where

the type that is implementing the trait is foreign
all of the parameters being passed to the trait (if there are any) are also foreign.
Here's one example of this error:*/
//1.类型别名不能实现方法,不能为crate外的类型实现方法，必须通过trait
impl<T: PartialOrd + Debug + Display> TreeT<T> for Option<Box<TNode<T>>> {
    fn new() -> TTree<T> {
        None
    }

    //E0449A visibility qualifier was used when it was unnecessary. Erroneous code examples:
    fn insert(&mut self, data: T) {
        match self {
            Some(ref mut node) =>
                if node.elem > data {
                    TreeT::insert(&mut node.left, data)
                } else {
                    TreeT::insert(&mut node.right, data)
                },
            None => *self = Some(Box::new(TNode::new(data))),
        }
    }
}
