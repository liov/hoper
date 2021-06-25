
#[cfg(test)]
mod tests {
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}

//test --package cmd --bin macro tt
//Tests passed: 1 (moments ago)
mod tt {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}

//test --package cmd --bin macro t
//Tests failed: 1, passed: 6 (moments ago)
mod t {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}

//test --package cmd --bin macro m
//Tests passed: 1 (moments ago)
mod m {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}

//Tests passed: 1 (moments ago)
#[test]
fn it_works() {
    assert_eq!(3 + 2, 4);
}

fn main() {}