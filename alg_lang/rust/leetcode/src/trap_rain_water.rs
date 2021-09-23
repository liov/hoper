///接雨水

pub fn trap_rain_water(height: Vec<i32>) -> i32 {
    if height.len() < 2 { return 0; }
    let mut result = 0;
    let mut sub = 0;
    let mut num = 0;
    let mut first = true;
    let mut left = 0;
    let mut idx = 0;
    let mut i = 0;
    let mut add = true;
    loop {
        if add && idx == height.len() - 1 {
            return result;
        }
        if height[i] >= left && height[i] != 0 {
            if first {
                left = height[i];
                first = false;
            } else {
                /*语义逻辑是这样，但是不好写
                if num == 0 && height[i] == left {
                    continue;
                }*/
                if !(num == 0 && height[i] == left) {
                    result = result + left * num - sub;
                    left = height[i];
                    sub = 0;
                    num = 0;
                    if add { idx = i; }
                }
            }
        } else if !first {
            num = num + 1;
            sub = sub + height[i];
        }
        if add && i == height.len() - 1 {
            sub = 0;
            num = 0;
            first = true;
            left = 0;
            add = false;
            continue;
        }

        if !add && i == idx {
            return result;
        }
        if add { i = i + 1; } else { i = i - 1; }
    }
}