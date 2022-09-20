///两数之和
use std::collections::HashMap;


//暴力法52ms，2MB
//一遍hash版本，0ms，2.7MB
pub fn two_sum1(nums: Vec<i32>, target: i32) -> Vec<i32> {
    for i in 0..nums.len() {
        for j in i + 1..nums.len() {
            if nums[j] == target - nums[i] {
                return vec![i as i32, j as i32];
            }
        }
    }
    panic!("不存在")
}

pub fn two_sum2(nums: Vec<i32>, target: i32) -> Vec<i32> {
    let mut map: HashMap<i32, usize> = HashMap::new();
    let mut index = 0;
    while index < nums.len() {
        if let Some(j) = map.get(&(target - nums[index])) {
            return vec![*j as i32, index as i32];
        }
        map.insert(nums[index], index);
        index = index + 1
    }
    panic!("不存在")
}


