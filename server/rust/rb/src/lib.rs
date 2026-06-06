#[cfg(any(feature = "rb-core", feature = "rb-thumb", feature = "agent-serve", feature = "media"))]
pub mod remotebrowse;

#[cfg(any(feature = "media", feature = "daemon", feature = "agent-serve"))]
pub mod tracing_init;

#[cfg(feature = "media")]
pub mod file;

#[cfg(feature = "media")]
pub mod grpc_server;

#[cfg(any(feature = "rb-core", feature = "viewer-ffi"))]
pub mod remotebrowse_svc;

#[cfg(any(feature = "viewer-ffi", feature = "agent-serve"))]
pub mod p2p_http;
#[cfg(feature = "viewer-ffi")]
pub mod rb_http_client;

#[cfg(any(feature = "media", feature = "viewer-ffi"))]
pub mod proto_zstd;

#[cfg(any(feature = "viewer-ffi", feature = "agent-serve"))]
pub mod ice_link;

#[cfg(feature = "viewer-ffi")]
mod ice_grpc_ffi;

#[cfg(feature = "client")]
pub mod client;

#[cfg(feature = "viewer-ffi")]
mod ffi;

#[cfg(feature = "daemon")]
pub mod daemon;

#[cfg(any(feature = "client", feature = "daemon", feature = "viewer-ffi", feature = "agent-serve"))]
pub mod signal_proto;

#[cfg(feature = "agent-serve")]
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
