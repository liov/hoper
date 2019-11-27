#![allow(unused_variables)]

use futures::executor::block_on;
use std::thread::sleep;
use std::time::Duration;

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
    sleep(Duration::from_secs(2));
    println!("{}","sing".to_owned() + song);
}

async fn dance() {
    println!("dance");
}

async fn learn_and_sing() {
    let song = learn_song().await;
    sing_song(song).await;
}

async fn sing_and_dance() {
    let f1 = learn_and_sing();
    let f2 = dance();
    futures::join!(f1, f2);
}
