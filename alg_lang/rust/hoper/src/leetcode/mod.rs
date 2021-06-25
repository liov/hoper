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


///两数相加

//执行用时 : 4 ms, 在Add Two Numbers的Rust提交中击败了100.00% 的用户

//内存消耗 : 2 MB, 在Add Two Numbers的Rust提交中击败了100.00% 的用户
// Definition for singly-linked list.
#[derive(PartialEq, Eq, Clone, Debug)]
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>,
}

impl ListNode {
    #[inline]
    pub fn new(val: i32) -> Self {
        ListNode {
            next: None,
            val,
        }
    }

    /*因为 Rust 是后向兼容的（backwards compatible），所以不会移除 ref 和 ref mut，同时它们在一些不明确的场景还有用，
    比如希望可变地借用结构体的部分值而可变地借用另一部分的情况。你可能会在老的 Rust 代码中看到它们，所以请记住它们仍有价值。*/
    pub fn push(&mut self, val: i32) {
        match self.next {
            Some(ref mut next) =>
                next.push(val),
            None => self.next = Some(Box::new(ListNode::new(val)))
        }
    }
}

pub fn add_two_numbers(l1: Option<Box<ListNode>>, l2: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
    let mut next1 = l1;
    let mut next2 = l2;
    let mut result = Box::new(ListNode::new(0));
    let mut first = true;
    let mut carry = 0;
    let mut sum;
    loop {
        let lx = next1.unwrap_or(Box::new(ListNode::new(0)));
        let ly = next2.unwrap_or(Box::new(ListNode::new(0)));

        sum = lx.val + ly.val;

        if first {
            result.val = sum % 10;
            first = false;
            if sum >= 10 { carry = 1 } else { carry = 0 }
        } else {
            if sum % 10 + carry == 10 {
                result.push(0);
                carry = 1;
            } else {
                result.push(sum % 10 + carry);
                if sum >= 10 { carry = 1 } else { carry = 0 }
            }
        }

        next1 = lx.next;
        next2 = ly.next;
        if next1 == None && next2 == None && carry == 0 {
            break;
        }
    }
    Some(result)
}

///字符串相乘

//执行用时 : 0 ms, 在Multiply Strings的Rust提交中击败了100.00% 的用户
//内存消耗 : 2 MB, 在Multiply Strings的Rust提交中击败了100.00% 的用户
/*m 位数和n位数相乘，结果位数为m+n-1或m + n，因此存进位数据数组大小申请为m+n位
对应位相乘后的结果与进位数据数组对应位置相加，十位数存入进位数组下一位，个位数留在该位
存进位数据数组在转为字符串返回的时候，要把前导零去掉*/
pub fn multiply(num1: String, num2: String) -> String {
    let num_vec1 = num1.as_bytes();
    let num_vec2 = num2.as_bytes();
    let cap = num_vec1.len() + num_vec2.len();
    let mut result = Vec::with_capacity(cap);
    unsafe { result.set_len(cap); }
    let con = 48;
    let mut product; //乘积
    let mut decade = 0; //十位

    if (num_vec1.len() == 1 && num_vec1[0] == con) || (num_vec2.len() == 1 && num_vec2[0] == con) {
        return String::from("0");
    }

    for i in 0..num_vec1.len() {
        for j in 0..num_vec2.len() {
            if i == 0 { result[cap - i - j - 1] = 0 } else if j == 0 { result[cap - i - num_vec2.len() - 1] = 0 }
            product = (num_vec1[num_vec1.len() - i - 1] - con) * (num_vec2[num_vec2.len() - j - 1] - con) + decade;
            decade = product / 10 + (product % 10 + result[cap - i - j - 1]) / 10;
            result[cap - i - j - 1] = (product % 10 + result[cap - i - j - 1]) % 10;

            if i == num_vec1.len() - 1 {
                result[cap - i - j - 1] = result[cap - i - j - 1] + con
            } else if j == 0 {
                result[cap - i - j - 1] = result[cap - i - j - 1] + con
            }
        }
        result[cap - i - num_vec2.len() - 1] = decade;
        decade = 0;
    }

    if result[0] != 0 { result[0] = result[0] + con } else { result.remove(0); }
    String::from_utf8(result).unwrap()
}

///两数相除

pub fn divide(dividend: i32, divisor: i32) -> i32 {
    if dividend == -0x80000000 && divisor == -1 {
        return 0x7fffffff;
    }
    if dividend == 0 {
        return 0;
    }
    if divisor == 1 { return dividend; }
    if divisor == -1 { return -dividend; }
    let mut result = 0;
    let mut dividend_copy = dividend;
    let mut divisor_copy = divisor;
    if dividend_copy > 0 { dividend_copy = -dividend_copy }
    if divisor_copy > 0 { divisor_copy = -divisor_copy }
    while dividend_copy <= divisor_copy {
        dividend_copy = dividend_copy - divisor_copy;
        result += 1;
    }
    if (dividend > 0 && divisor < 0) || (dividend < 0 && divisor > 0) {
        return -result;
    }
    result
}

///电话号码的字母组合

//Y 组合子是 lambda 演算中的一个概念，是任意函数的不动点，在函数式编程中主要作用是 提供一种匿名函数的递归方式。
//"λ" 字符可以看作 function 声明，"."字符前为参数列表，"."字符后为函数体。
//不动点是函数的一个特征：对于函数 f(x)，如果有变量  a 使得  f(a)=a 成立，则称 a 是函数 f 上的一个不动点。

//注意结尾没有分号的那一行 x+1，与你见过的大部分代码行不同。表达式的结尾没有分号。如果在表达式的结尾加上分号，
// 它就变成了语句，而语句不会返回值。在接下来探索具有返回值的函数和表达式时要谨记这一点。

pub fn letter_combinations(digits: String) -> Vec<String> {
    let combination = vec![
        vec![b'a', b'b', b'c'],
        vec![b'd', b'e', b'f'],
        vec![b'g', b'h', b'i'],
        vec![b'j', b'k', b'l'],
        vec![b'm', b'n', b'o'],
        vec![b'p', b'q', b'r', b's'],
        vec![b't', b'u', b'v'],
        vec![b'w', b'x', b'y', b'z'],
    ];

    let mut result = Vec::new();

    let cap = digits.len();
    let mut middle: Vec<u8> = Vec::with_capacity(cap);
    unsafe { middle.set_len(cap); }

    let digits_vec = digits.as_bytes();
    //闭包递归实在不会，改成了函数递归，消耗几何未知，有个clone()难受啊
    fn get(n: usize, res: &mut Vec<String>, mid: &mut Vec<u8>, com: &Vec<Vec<u8>>, dig: &[u8]) {
        if n < mid.len() {
            for i in &(com[dig[n] as usize - 50]) {
                mid[n] = *i;
                get(n + 1, res, mid, com, dig);
                if n == mid.len() - 1 {
                    res.push(String::from_utf8(mid.clone()).unwrap());
                }
            }
        }
    }


    //let Y = |y|(|x|y(x(x)))(|x|y(|n|x(x)(n)));

    get(0, &mut result, &mut middle, &combination, &digits_vec);
    result
}


///旋转链表
use std::mem;

impl ListNode {
    pub fn len(&self) -> usize {
        fn get_len(list: &ListNode, len: usize) -> usize {
            if let Some(ref next) = list.next {
                get_len(next, len + 1)
            } else {
                return len;
            }
        }
        get_len(self, 1)
    }
}

pub fn rotate_right(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if k == 0 { return head; }
    let mut l = head.unwrap();

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> (Option<*mut Box<ListNode>>, Option<*mut Box<ListNode>>, i32) {
        return if let Some(ref mut next) = list.next {
            let (p2, p3, size) = get(next, len + 1, k);
            if len == size - k + 1 {
                return (Some(list), p3, size);
            } else if len == size - k {
                let mut t = list.clone();
                t.next.take();
                mem::replace(list, t);
            }
            (p2, p3, size)
        } else {
            (None, Some(list), len)
        }
    }

    let (p2, p3, len) = get(&mut l, 1, k);

    if k > len {
        return rotate_right(Some(l), k % len);
    };

    unsafe {
        mem::replace(&mut (*p3.unwrap()).next, Some(l.clone()));
        Some((*p2.unwrap()).clone())
    }
}

pub fn rotate_right2(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if head == None || k == 0 { return head; }
    let mut l = head.unwrap();

    let len = l.len() as i32;
    if len == 1 { return Some(l); };
    let mut k2 = k;
    if k >= len {
        k2 = k % len;
        if k2 == 0 {
            return Some(l);
        }
    };

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> (Option<*mut Box<ListNode>>, Option<*mut Box<ListNode>>) {
        if len == k {
            let mut next = list.next.take();

            if let Some(ref mut next) = next {
                return get(next, len + 1, k);
            }
        }
        if len == k + 1 {
            let (_, p3) = get(list, len + 1, k);
            return (Some(list), p3);
        }

        return if let Some(ref mut next) = list.next {
            get(next, len + 1, k)
        } else {
            (None, Some(list))
        }

    }


    let (p2, p3) = get(&mut l, 1, len - k2);

    let (p2, p3) = (p2.unwrap(), p3.unwrap());
    unsafe {
        mem::replace(&mut ((*p3).next), Some(l));
        Some((*p2).clone())
    }
}

pub fn rotate_right3(head: Option<Box<ListNode>>, k: i32) -> Option<Box<ListNode>> {
    if head == None || k == 0 { return head; }
    let mut l = head.unwrap();

    let len = l.len() as i32;
    if len == 1 { return Some(l); };
    let mut k2 = k;
    if k >= len {
        k2 = k % len;
        if k2 == 0 {
            return Some(l);
        }
    };

    fn get(list: &mut Box<ListNode>, len: i32, k: i32) -> Box<ListNode> {
        if len == k {
            let mut next = list.next.take();

            if let Some(ref mut next) = next {
                return get(next, len + 1, k);
            }
        }
        if len == k + 1 {
            return list.clone();
        }

        if let Some(ref mut next) = list.next {
            return get(next, len + 1, k);
        }

        panic!("错误")
    }


    let mut list = get(&mut l, 1, len - k2);
    fn set(list: &mut Box<ListNode>, list2: Box<ListNode>) {
        if let Some(ref mut next) = list.next {
            set(next, list2);
        } else {
            list.next = Some(list2.clone());
        }
    }
    set(&mut list, l);
    Some(list)
}

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

///接雨水

pub fn trap(height: Vec<i32>) -> i32 {
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

///接雨水 II
use std::collections::{HashSet, BTreeMap};

#[derive(PartialEq, Eq, PartialOrd, Ord, Debug, Hash, Clone, Copy)]
struct Point(i32, usize, usize);


pub fn trap_rain_water(height_map: Vec<Vec<i32>>) -> i32 {
    if height_map.len() < 3 || height_map[0].len() < 3 { return 0; };
    let m = height_map.len();
    let n = height_map[0].len();
    let mut result = 0;
    let mut side = BTreeMap::new();
    let mut drop = HashSet::new();
    for x in 0..m {
        for y in 0..n {
            if x == 0 || x == m - 1 || y == 0 || y == n - 1 {
                side.insert(Point(height_map[x][y], x, y), false);
                drop.insert((x, y));
            }
        }
    }

    let round: [[i32; 2]; 4] = [[0, -1], [0, 1], [-1, 0], [1, 0]];
    let mut x = 0;
    let mut y = 0;

    let mut point_iter = side.iter();
    while let Some((point, yet)) = point_iter.next() {
        if *yet { continue; };
        let point = *point;
        for j in round.iter() {
            x = (point.1 as i32 + j[0]) as usize;
            y = (point.2 as i32 + j[1]) as usize;
            if x >= m - 1 || y >= n - 1 {
                continue;
            }
            if drop.get(&(x, y)) != None {
                continue;
            }
            drop.insert((x, y));
            if height_map[x][y] <= point.0 {
                result = result + (point.0 - height_map[x][y]);
                side.insert(Point(point.0, x, y), false);
            } else {
                side.insert(Point(height_map[x][y], x, y), false);
            }
            point_iter = side.iter();
        }
        side.insert(point, true);
        point_iter = side.iter();
    }
    result
}


//用Vec插入排开销太大，放弃
use std::fmt::{Debug};


#[derive(PartialEq, Eq, PartialOrd, Ord, Debug, Hash, Clone, Copy)]
struct Point2(i32, usize, usize, bool);


pub fn insert_sort<T>(ord_vec:&mut Vec<T>, mut other:Vec<T>) where T: Ord+Copy+Debug {
    other.sort_by(|a, b| b.cmp(a));
    let len = ord_vec.len();
    for i in 0..other.len() {
        ord_vec.push(other[i]);
    }

    if other[0]<=ord_vec[len-1]{return;}

    for i in (0..len).rev() {
        let other_len =other.len();
        if other_len == 0 { break; }
        for j in 0..other_len {
            if other[j] <= ord_vec[i] {
                ord_vec[i + j] = ord_vec[i];
                for x in (j..other_len).rev() {
                    ord_vec[i + x + 1] = other[x];
                    other.pop();
                }
                break;
            }
        }
        if other.len() == other_len{
            ord_vec[i+other_len]=ord_vec[i];
        }
    }
    if other.len()>0{
        for j in 0..other.len(){
            ord_vec[j+other.len()]=ord_vec[j];
            ord_vec[j]=other[j];
        }
    }
}


pub fn trap_rain_water2(height_map: Vec<Vec<i32>>) -> i32 {
    if height_map.len() < 3 || height_map[0].len() < 3 { return 0; };
    let m = height_map.len();
    let n = height_map[0].len();
    let mut result = 0;
    let mut side = Vec::with_capacity(m * n);
    let mut drop = HashSet::new();
    for x in 0..m {
        for y in 0..n {
            if x == 0 || x == m - 1 || y == 0 || y == n - 1 {
                side.push(Point2(height_map[x][y], x, y, false));
                drop.insert((x, y));
            }
        }
    }
    side.sort_by(|a, b| b.cmp(a));
    let round: [[i32; 2]; 4] = [[0, -1], [0, 1], [-1, 0], [1, 0]];
    let mut x = 0;
    let mut y = 0;
    let mut i = side.len() - 1;
    loop {
        if side[i].3 {
            i = i - 1;
            continue;
        }
        let mut sub_side = Vec::with_capacity(3);

        for j in round.iter() {
            x = (side[i].1 as i32 + j[0]) as usize;
            y = (side[i].2 as i32 + j[1]) as usize;
            if  x >= m - 1  || y >= n - 1 {
                continue;
            }
            println!("({:?},{:?}):{:?}:{:?}", x, y, side[i], height_map[x][y]);
            if drop.get(&(x, y)) != None {
                continue;
            }
            drop.insert((x, y));
            if height_map[x][y] <= side[i].0 {
                result = result + (side[i].0 - height_map[x][y]);
                sub_side.push(Point2(side[i].0, x, y, false));
            } else {
                sub_side.push(Point2(height_map[x][y], x, y, false))
            }
            //这里应该有个else，插入排大于边的新边
        }
        side[i].3 = true;
        if sub_side.len() > 0 {
            insert_sort(&mut side,sub_side);
            i = side.len();
        }

        if i == 0 { break; }
        i = i - 1;
    }
    result
}
