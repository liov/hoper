#![allow(unused_variables)]

use futures::{
    executor::block_on,
    future::{Fuse, FusedFuture, FutureExt},
    stream::{FusedStream, FuturesUnordered, Stream, StreamExt},
    pin_mut,
    select,
};
use std::thread::sleep;
use std::time::Duration;
use async_test::timer_future::TimerFuture;

async fn hello_world() {
    println!("hello, world!");
}

fn main() {
    block_on(sing_and_dance()); // `future` is run and "hello, world!" is printed
}

type Song<'a> =  &'a str;

async fn learn_song() -> Song<'static> {
    "song"
}

async fn sing_song(song: Song<'static>) {
    println!("{}","sing".to_owned() + song);
}

async fn dance() {
    println!("dance");
}

async fn learn_and_sing() {
    TimerFuture::new(Duration::new(2, 0)).await;
    let song = learn_song().await;
    sing_song(song).await;
}

async fn sing_and_dance() {
    let f1 = learn_and_sing();
    let f2 = dance();
    futures::join!(f1, f2);
}
