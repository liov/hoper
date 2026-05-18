use bytes::{BufMut, BytesMut};
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::TcpStream;

use crate::client::wire;

pub struct TcpWire {
    sock: TcpStream,
    buf: BytesMut,
    relay: bool,
}

impl TcpWire {
    pub fn direct(sock: TcpStream) -> Self {
        Self { sock, buf: BytesMut::with_capacity(4096), relay: false }
    }

    pub fn relay(sock: TcpStream) -> Self {
        Self { sock, buf: BytesMut::with_capacity(4096), relay: true }
    }

    pub async fn read_frame(&mut self) -> Result<(u8, Vec<u8>), String> {
        if self.relay {
            let payload = read_relay_payload(&mut self.sock).await?;
            return wire::decode_frame(&payload).map_err(|e| e.to_string());
        }
        loop {
            if self.buf.len() >= wire::HEADER_LEN {
                let n = u32::from_be_bytes(self.buf[2..6].try_into().unwrap()) as usize;
                let need = wire::HEADER_LEN + n;
                if self.buf.len() >= need {
                    let raw = self.buf.split_to(need).to_vec();
                    return wire::decode_frame(&raw).map_err(|e| e.to_string());
                }
            }
            let mut tmp = [0u8; 2048];
            let n = self.sock.read(&mut tmp).await.map_err(|e| e.to_string())?;
            if n == 0 {
                return Err("eof".into());
            }
            self.buf.put_slice(&tmp[..n]);
        }
    }

    pub async fn write_frame(&mut self, typ: u8, payload: &[u8]) -> Result<(), String> {
        let frame = wire::encode_frame(typ, payload);
        if self.relay {
            write_relay_payload(&mut self.sock, &frame).await?;
        } else {
            self.sock.write_all(&frame).await.map_err(|e| e.to_string())?;
        }
        Ok(())
    }
}

async fn read_relay_payload(sock: &mut TcpStream) -> Result<Vec<u8>, String> {
    let mut sz = [0u8; 4];
    sock.read_exact(&mut sz).await.map_err(|e| e.to_string())?;
    let n = u32::from_be_bytes(sz) as usize;
    if n > wire::MAX_PAYLOAD + wire::HEADER_LEN {
        return Err("relay frame too large".into());
    }
    let mut buf = vec![0u8; n];
    if n > 0 {
        sock.read_exact(&mut buf).await.map_err(|e| e.to_string())?;
    }
    Ok(buf)
}

async fn write_relay_payload(sock: &mut TcpStream, payload: &[u8]) -> Result<(), String> {
    sock.write_all(&u32::to_be_bytes(payload.len() as u32))
        .await
        .map_err(|e| e.to_string())?;
    if !payload.is_empty() {
        sock.write_all(payload).await.map_err(|e| e.to_string())?;
    }
    Ok(())
}
