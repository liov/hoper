use actix::prelude::*;
use futures::{FutureExt, TryFutureExt};
use std::time::Duration;
//在 Cargo.toml 中加入 actix_rt
// 依赖告诉了 Cargo 要从 https://crates.io 下载 actix_rt 和其依赖，
// 并使其可在项目代码中使用。


/// Define `Ping` message
#[derive(Message)]
#[rtype(result = "()")]
struct Ping(usize);

/// Actor
#[derive(Debug)]
struct MyActor {
    count: usize,
    name: String,
    addr: Recipient<Ping>,
}

/// Declare actor and its context
impl Actor for MyActor {
    type Context = Context<Self>;
}

/// Handler for `Ping` message
impl Handler<Ping> for MyActor {
    type Result = ();

    fn handle(&mut self, msg: Ping, ctx: &mut Context<Self>) {
        self.count += 1;

        if self.count > 10 {
            System::current().stop();
        } else {
            println!("[{0}] Ping received {1}", self.name, msg.0);

            // wait 100 nanos
            ctx.run_later(Duration::new(0, 100), move |act, _| {
                act.addr.do_send(Ping ( msg.0 + 1));
            });
        }
    }
}


fn main() {
    let system = System::new("test");

    // To get a Recipient object, we need to use a different builder method
    // which will allow postponing actor creation
    let addr = MyActor::create(|ctx| {
        // now we can get an address of the first actor and create the second actor
        let addr = ctx.address();
        let addr2 = MyActor {
            count: 0,
            name: String::from("Game 2"),
            addr:addr.recipient()
        }.start();

        // let's start pings
        addr2.do_send(Ping(0));

        // now we can finally create first actor
        MyActor {
            count: 0,
            name: String::from("Game 1"),
            addr:addr2.recipient()
        }
    });

    system.run();
}
