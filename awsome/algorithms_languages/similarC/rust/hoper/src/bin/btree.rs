use std::collections::BTreeSet;

fn main() {
    let mut set: BTreeSet<usize> = [3, 1, 2].iter().cloned().collect();
    let mut set_iter = set.iter();
    let mut i = 0;
    while let Some(t) = set_iter.next() {
        println!("{:?}", set);
        i+=1;
        set.insert(i);
        set_iter = set.iter();
    }
}
