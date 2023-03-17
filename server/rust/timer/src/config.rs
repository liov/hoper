use std::collections::HashMap;
use std::fs::File;
use std::io::Read;
use std::sync::{Arc,Mutex};
use lazy_static::lazy_static;
use once_cell::sync::Lazy;
use serde_derive::{Deserialize, Serialize};

static CONFIG1: Lazy<Mutex<Config>> = Lazy::new(|| {
    Mutex::new(Config::new())
});


lazy_static!{
   pub static ref CONFIG: Arc<Mutex<Config>> = Arc::new(Mutex::new(Config::new()));
}

#[derive(Serialize,Deserialize,Debug,Clone)]
pub struct Config {
    GORMDB: DataBase,
}

impl Config {
    fn new() -> Config {
        let file_path = "../go/mod/local.toml";
        let mut file = match File::open(file_path) {
            Ok(f) => f,
            Err(e) => panic!("no such file {} exception:{}", file_path, e)
        };
        let mut str_val = String::new();
        match file.read_to_string(&mut str_val) {
            Ok(s) => s
            ,
            Err(e) => panic!("Error Reading file: {}", e)
        };
        toml::from_str(&str_val).unwrap()
    }
}


#[derive(Serialize,Deserialize,Debug,Clone)]
struct DataBase {
    Type : String,
    User : String,
    Password : String,
    Host : String,
    Port : i32,
    Charset:String,
    Database:String,
    MaxIdleConns:i32,
    MaxOpenConns:i32,
}
