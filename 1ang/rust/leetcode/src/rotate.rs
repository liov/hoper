///旋转数组

pub fn rotate(nums: &mut Vec<i32>, k: i32) {
    let len = nums.len();
    if len == 0 || len == 1 || k == 0 { return; }
    let mut k2 = k as usize;
    if k2 >= len {
        k2 = k2 % len;
        if k2 == 0 {
            return;
        }
    };

    let mut idx;
    let mut tmp = nums[k2];

    //求最大公因数
    let mut ii = 1;
    let mut times = len;
    let mut kk = k2;
    while ii != 0 {
        ii = times % kk;
        times = kk;
        kk = ii;
    }

    for i in 0..times {
        idx = i;
        loop {
            idx = idx + k2;
            if idx >= len {
                idx = idx - len
            }
            tmp = nums[idx];
            nums[idx] = nums[i];
            nums[i] = tmp;
            if idx == i { break; }
        }
    }
}
