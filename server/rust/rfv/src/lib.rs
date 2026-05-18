pub mod remotebrowse;

#[cfg(feature = "media")]
pub mod grpc_server;

#[cfg(feature = "client")]
pub mod client;

#[cfg(any(feature = "client", feature = "transport"))]
mod ffi;

#[cfg(feature = "daemon")]
pub mod daemon;

#[cfg(any(feature = "client", feature = "daemon", feature = "transport"))]
pub mod signal_proto;

#[cfg(feature = "transport")]
pub mod transport;

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
