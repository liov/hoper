use futures::{executor, future, future::FutureExt, pin_mut, select};

fn main() {
    executor::block_on(count());
    executor::block_on(race_tasks());
}

async fn count() {
    let mut a_fut = future::ready(4);
    let mut b_fut = future::ready(6);
    let mut total = 0;

    loop {
        select! {
            a = a_fut => total += a,
            b = b_fut => total += b,
            complete => break,
            default => unreachable!(), // never runs (futures are ready, then complete)
        }
    }
    println!("{}",total);
}

async fn task_one() {  println!("task_one 执行了"); }

async fn task_two() { println!("task_two 执行了"); }

async fn race_tasks() {
    let t1 = task_one().fuse();  
    let t2 = task_two().fuse();

    pin_mut!(t1, t2);

    select! {
        () = t1 => println!("task one completed first"),
        () = t2 => println!("task two completed first"),
    }
}
