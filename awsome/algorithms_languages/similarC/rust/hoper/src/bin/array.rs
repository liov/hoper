

fn main(){
    let mut x:[i32;100]=[0;100];
    array1(x);
    dbg!(x[0]);
    let mut y:[i32;10000]=[0;10000];
    array2(y);
    dbg!(y[0]);
}

fn array1(mut arr:[i32;100]){
    arr[0] = 1;
    dbg!(arr[0]);
}

fn array2(mut arr:[i32;10000]){
    arr[1] = 1;
    dbg!(arr[0]);
}
