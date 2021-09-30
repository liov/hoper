#![feature(fn_traits,unboxed_closures)]


struct Fact(usize);

impl<'a> FnOnce<()> for Fact {
    type Output = usize;

    extern "rust-call" fn call_once(self, _args: () )-> Self::Output {
        let val = self.0;
        if val==1{
            return 1;
        };
        val*Self(val - 1).call_once(())
    }
}

impl<'a> FnMut<()> for Fact {
    extern "rust-call" fn call_mut(&mut self, _args: ()) -> Self::Output {
        let val = self.0;
        if val==1{
            return 1;
        };
        self.0 -=1;
        val *self()
    }
}

impl<'a> Fn<()> for Fact {
    extern "rust-call" fn call(&self, _args: ()) -> Self::Output {
        let val = self.0;
        if val==1{
            return 1;
        };
        val*Self(val - 1).call(())
    }
}

fn main() {
    let mut five = Fact(5);

    println!("{}",five());

    let fib = fix(|f, x| {
        if x == 0 || x == 1 {
            x
        } else {
            // `f` is `fib`
            f(x - 1) + f(x - 2)
        }
    });
    for n in 0..=10 {
        println!("{}: {}", n, fib(n));
    }
}


pub fn fix<A, B, F>(f: F) -> impl Fn(A) -> B
    where F: Fn(&dyn Fn(A) -> B, A) -> B
{
    use std::cell::Cell;

    move |a: A| {
        // Hopefully optimized away. Can probably use some unsafe magic to help the optimizer...
        let tmp_fn = |_: A| -> B { panic!("Hmm... not good.") };
        let (fun_holder, fun);
        fun_holder = Cell::new(&tmp_fn as &dyn Fn(A) -> B);
        fun = |ai: A| { f(fun_holder.get(), ai) };
        fun_holder.set(&fun as &dyn Fn(A) -> B);
        f(fun_holder.get(), a)
    }
}


