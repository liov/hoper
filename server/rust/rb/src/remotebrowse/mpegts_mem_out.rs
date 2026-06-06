//! MPEG-TS 复用到临时文件再读入内存（避免自定义 AVIO 与 HEVC 编码器组合时段错误）。
use std::io::{Read, Seek, SeekFrom};
use std::path::{Path, PathBuf};
use std::sync::atomic::{AtomicU64, Ordering};

use ffmpeg_next as ffmpeg;
use ffmpeg::format::context::Output;

static TS_SEQ: AtomicU64 = AtomicU64::new(0);

pub struct MpegtsMemOut {
    output: Output,
    path: PathBuf,
    header_written: bool,
}

impl MpegtsMemOut {
    pub fn open() -> Result<Self, String> {
        let _ = ffmpeg::init();
        let n = TS_SEQ.fetch_add(1, Ordering::Relaxed);
        let path = std::env::temp_dir().join(format!("rb-frag-{}-{n}.ts", std::process::id()));
        let output = ffmpeg::format::output_as(&path, "mpegts").map_err(|e| e.to_string())?;
        Ok(Self { output, path, header_written: false })
    }

    pub fn path(&self) -> &Path {
        &self.path
    }

    pub fn output_mut(&mut self) -> &mut Output {
        &mut self.output
    }

    pub fn mark_header_written(&mut self) {
        self.header_written = true;
    }

    pub fn write_trailer(&mut self) -> Result<(), String> {
        if self.header_written {
            self.output.write_trailer().map_err(|e| e.to_string())?;
        }
        Ok(())
    }

    pub fn into_bytes(mut self) -> Result<Vec<u8>, String> {
        self.write_trailer()?;
        drop(self.output);
        let bytes = std::fs::read(&self.path).map_err(|e| e.to_string())?;
        let _ = std::fs::remove_file(&self.path);
        if bytes.is_empty() {
            return Err("empty mpegts".into());
        }
        Ok(bytes)
    }
}

/// 边转码边从临时 TS 文件增量读出，供 HTTP chunked 连续流。
pub struct MpegtsStreamOut {
    inner: MpegtsMemOut,
    read_off: u64,
}

impl MpegtsStreamOut {
    pub fn open() -> Result<Self, String> {
        Ok(Self { inner: MpegtsMemOut::open()?, read_off: 0 })
    }

    pub fn output_mut(&mut self) -> &mut Output {
        self.inner.output_mut()
    }

    pub fn mark_header_written(&mut self) {
        self.inner.mark_header_written();
    }

    pub fn emit_pending(&mut self, mut emit: impl FnMut(&[u8]) -> Result<(), String>) -> Result<(), String> {
        let path = self.inner.path().to_path_buf();
        let len = std::fs::metadata(&path).map_err(|e| e.to_string())?.len();
        if len <= self.read_off {
            return Ok(());
        }
        let mut f = std::fs::File::open(&path).map_err(|e| e.to_string())?;
        f.seek(SeekFrom::Start(self.read_off)).map_err(|e| e.to_string())?;
        let mut left = (len - self.read_off) as usize;
        let mut buf = [0u8; 256 * 1024];
        while left > 0 {
            let cap = left.min(buf.len());
            let n = f.read(&mut buf[..cap]).map_err(|e| e.to_string())?;
            if n == 0 {
                break;
            }
            emit(&buf[..n])?;
            self.read_off += n as u64;
            left -= n;
        }
        Ok(())
    }

    pub fn finish(mut self, mut emit: impl FnMut(&[u8]) -> Result<(), String>) -> Result<(), String> {
        let path = self.inner.path().to_path_buf();
        self.inner.write_trailer()?;
        drop(self.inner);
        let len = std::fs::metadata(&path).map_err(|e| e.to_string())?.len();
        if len > self.read_off {
            let mut f = std::fs::File::open(&path).map_err(|e| e.to_string())?;
            f.seek(SeekFrom::Start(self.read_off)).map_err(|e| e.to_string())?;
            let mut left = (len - self.read_off) as usize;
            let mut buf = [0u8; 256 * 1024];
            while left > 0 {
                let cap = left.min(buf.len());
            let n = f.read(&mut buf[..cap]).map_err(|e| e.to_string())?;
                if n == 0 {
                    break;
                }
                emit(&buf[..n])?;
                left -= n;
            }
        }
        let _ = std::fs::remove_file(&path);
        Ok(())
    }
}
