use std::rc::Rc;
use std::cell::RefCell;
use std::io::Error;
use std::io::ErrorKind;


pub struct Node<T> {
    value: T,
    next: Option<Box<Node<T>>>,
}


pub enum EList<T>{
    Cons(Rc<RefCell<T>>, Rc<EList<T>>),
    Nil,
}

pub struct List<T> {
    value: Rc<RefCell<T>>,
    next: Option<Rc<List<T>>>,
}

impl<T> List<T> {
    fn new(elem: T) -> Self {
        List {
            value:Rc::new(RefCell::new(elem)),
            next: None,
        }
    }

    fn _new(elem:Rc<RefCell<T>>) -> Self {
        List {
            value:elem,
            next: None,
        }
    }

    fn add(&mut self){

    }

    fn write_fmt(&mut self) -> Result<(), Error>{
       Err(Error::new(ErrorKind::NotFound,"5"))
    }
}
