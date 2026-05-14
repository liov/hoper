pub mod remotebrowse;

#[cfg(feature = "media")]
pub mod grpc_server;

#[cfg(feature = "client")]
pub mod client;

#[cfg(feature = "client")]
mod ffi;

#[cfg(feature = "daemon")]
pub mod daemon;

#[cfg(any(feature = "client", feature = "daemon"))]
pub mod signal_proto;

pub fn add(left: u64, right: u64) -> u64 {
    left + right
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(add(2, 2), 4);
    }
}
