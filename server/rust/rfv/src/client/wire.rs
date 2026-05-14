pub const VERSION1: u8 = 1;
pub const HEADER_LEN: usize = 6;
pub const MAX_PAYLOAD: usize = 16 << 20;

pub fn encode_frame(typ: u8, payload: &[u8]) -> Vec<u8> {
    let mut out = Vec::with_capacity(HEADER_LEN + payload.len());
    out.push(VERSION1);
    out.push(typ);
    out.extend_from_slice(&(payload.len() as u32).to_be_bytes());
    out.extend_from_slice(payload);
    out
}

pub fn decode_frame(buf: &[u8]) -> Result<(u8, Vec<u8>), &'static str> {
    if buf.len() < HEADER_LEN {
        return Err("short header");
    }
    if buf[0] != VERSION1 {
        return Err("bad version");
    }
    let n = u32::from_be_bytes(buf[2..6].try_into().unwrap()) as usize;
    if n > MAX_PAYLOAD {
        return Err("frame too large");
    }
    if buf.len() < HEADER_LEN + n {
        return Err("short payload");
    }
    Ok((buf[1], buf[HEADER_LEN..HEADER_LEN + n].to_vec()))
}
