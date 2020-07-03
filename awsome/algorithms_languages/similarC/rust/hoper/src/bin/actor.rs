use actix::prelude::*;
//在 Cargo.toml 中加入 actix_rt
// 依赖告诉了 Cargo 要从 https://crates.io 下载 actix_rt 和其依赖，
// 并使其可在项目代码中使用。



/// Define `Ping` message
struct Ping(usize);

impl Message for Ping {
    type Result = usize;
}

/// Actor
struct MyActor {
    count: usize,
}

/// Declare actor and its context
impl Actor for MyActor {
    type Context = Context<Self>;
}

/// Handler for `Ping` message
impl Handler<Ping> for MyActor {
    type Result = usize;

    fn handle(&mut self, msg: Ping, _: &mut Context<Self>) -> Self::Result {
        self.count += msg.0;
        self.count
    }
}

fn main() -> std::io::Result<()> {
    // start system, this is required step
    System::run(|| {
        // start new actor
        let addr = MyActor { count: 10 }.start();

        // send message and get future for result
        let res = addr.send(Ping(10));

        // handle() returns tokio handle
        actix_rt::spawn(
            res.map(|res| {
                println!("RESULT: {}", res == 20);

                // stop system and exit
                System::current().stop();
            })
                .map_err(|_| ()),
        );
    })
}
