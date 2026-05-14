use std::sync::Arc;

use bytes::{BufMut, BytesMut};
use webrtc_util::conn::Conn;

use crate::wire;

pub struct IceWire {
    conn: Arc<dyn Conn + Send + Sync>,
    read_buf: BytesMut,
}

impl IceWire {
    pub fn new(conn: Arc<dyn Conn + Send + Sync>) -> Self {
        Self { conn, read_buf: BytesMut::with_capacity(4096) }
    }

    pub async fn read_frame(&mut self) -> Result<(u8, Vec<u8>), String> {
        loop {
            if self.read_buf.len() >= wire::HEADER_LEN {
                let n = u32::from_be_bytes(self.read_buf[2..6].try_into().unwrap()) as usize;
                let need = wire::HEADER_LEN + n;
                if self.read_buf.len() >= need {
                    let raw = self.read_buf.split_to(need).to_vec();
                    return wire::decode_frame(&raw).map_err(|e| e.to_string());
                }
            }
            let mut tmp = [0u8; 2048];
            let n = self.conn.recv(&mut tmp).await.map_err(|e| e.to_string())?;
            if n == 0 {
                return Err("eof".to_string());
            }
            self.read_buf.put_slice(&tmp[..n]);
        }
    }

    pub async fn write_frame(&mut self, typ: u8, payload: &[u8]) -> Result<(), String> {
        let frame = wire::encode_frame(typ, payload);
        self.conn.send(&frame).await.map_err(|e| e.to_string())?;
        Ok(())
    }
}

pub async fn connect(conn: Arc<dyn Conn + Send + Sync>) -> Result<IceWire, String> {
    Ok(IceWire::new(conn))
}
