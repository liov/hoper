use std::sync::{mpsc, Arc, Mutex, Condvar};
use std::{thread, fmt};
use std::rc::Rc;
use std::fmt::Display;
use std::time::Duration;
use std::sync::atomic::{AtomicUsize,Ordering};

pub struct Foo<T>{
    id:Rc<T>
}

impl<T:Display> Display for Foo<T> {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "foo {}", self.id)
    }
}

unsafe impl<T> Send for Foo<T> {}

// 线程数量
const THREAD_COUNT :i32 = 10;

fn main(){

    // 创建一个通道
    let (tx, rx)= mpsc::sync_channel(0);
    let new_thread = thread::spawn(move || {
        // 创建线程用于发送消息
        for id in 0..THREAD_COUNT {
            // 注意Sender是可以clone的，这样就可以支持多个发送者
            let thread_tx = tx.clone();
            thread::spawn(move || {
                // 发送一个消息，此处是数字id
                thread_tx.send(id + 1).unwrap();
                println!("send {}", id + 1);
            });
            //顺序打印
            thread::sleep(Duration::from_millis(100));
        }
    });

    // 在主线程中接收子线程发送的消息并输出
    for _ in 0..THREAD_COUNT {
        println!("receive {}",rx.recv().unwrap());
    }

    new_thread.join().unwrap();

    let var : Arc<i32> = Arc::new(5);
    let share_var = var.clone();

    let new_thread = thread::spawn(move || {
        println!("share value in new thread: {}, address: {:p}", share_var, &*share_var);
    });

    new_thread.join().unwrap();
    println!("share value in main thread: {}, address: {:p}", var, &*var);

    let pair = Arc::new((Mutex::new(false), Condvar::new()));
    let pair2 = pair.clone();

    // 创建一个新线程
    thread::spawn(move|| {
        let &(ref lock, ref cvar) = &*pair2;
        let mut started = lock.lock().unwrap();
        *started = true;
        cvar.notify_one();
        println!("notify main thread");
    });

    // 等待新线程先运行
    let &(ref lock, ref cvar) = &*pair;
    let mut started = lock.lock().unwrap();
    while !*started {
        println!("before wait");
        started = cvar.wait(started).unwrap();
        println!("after wait");
    }

    let var : Arc<AtomicUsize> = Arc::new(AtomicUsize::new(5));
    let share_var = var.clone();

    // 创建一个新线程
    let new_thread = thread::spawn(move|| {
        println!("share value in new thread: {}", share_var.load(Ordering::SeqCst));
        // 修改值
        share_var.store(9, Ordering::SeqCst);
    });

    // 等待新建线程先执行
    new_thread.join().unwrap();
    println!("share value in main thread: {}", var.load(Ordering::SeqCst));
}
