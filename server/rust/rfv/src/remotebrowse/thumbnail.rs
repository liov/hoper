//! 缩略图：生成后写入本地磁盘缓存，供 gRPC 与 HTTP 复用。
use image::GenericImageView;
use sha2::{Digest, Sha256};
use std::fs;
use std::io::Cursor;
use std::path::{Path, PathBuf};

#[cfg(feature = "media")]
use ffmpeg_next as ffmpeg;
#[cfg(feature = "media")]
use ffmpeg_next::decoder::Decoder;
#[cfg(feature = "media")]
use ffmpeg_next::media::Type;
#[cfg(feature = "media")]
use ffmpeg_next::software::scaling::Context;

pub const DEFAULT_MAX_EDGE: u32 = 256;

pub fn thumb_cache_dir() -> PathBuf {
    std::env::var("RFV_THUMB_CACHE")
        .map(PathBuf::from)
        .unwrap_or_else(|_| PathBuf::from(".thumbnails"))
}

pub fn thumb_hash(path: &str, mtime_unix_ms: i64, max_edge: u32) -> String {
    let mut h = Sha256::new();
    h.update(path.as_bytes());
    h.update(mtime_unix_ms.to_le_bytes());
    h.update(max_edge.to_le_bytes());
    hex::encode(h.finalize())
}

fn cache_file(hash: &str) -> PathBuf {
    thumb_cache_dir().join(format!("{hash}.webp"))
}

fn encode_image_thumb(path: &Path, max_edge: u32) -> Result<Vec<u8>, String> {
    let img = image::open(path).map_err(|e| e.to_string())?;
    resize_to_webp(&img, max_edge)
}

#[cfg(feature = "media")]
fn encode_video_thumb(path: &Path, max_edge: u32) -> Result<Vec<u8>, String> {
    let mut input = ffmpeg::format::input(path).map_err(|e| e.to_string())?;
    let stream = input.streams().best(Type::Video).ok_or("no video stream")?;
    let idx = stream.index();
    let ctx = ffmpeg::codec::context::Context::from_parameters(stream.parameters()).map_err(|e| e.to_string())?;
    let mut decoder = ctx.decoder().video().map_err(|e| e.to_string())?;
    let mut frame = ffmpeg::frame::Video::empty();
    for (s, pkt) in input.packets() {
        if s.index() != idx {
            continue;
        }
        decoder.send_packet(&pkt).map_err(|e| e.to_string())?;
        if decoder.receive_frame(&mut frame).is_ok() {
            break;
        }
    }
    let mut rgb = ffmpeg::frame::Video::empty();
    let (ow, oh) = (frame.width(), frame.height());
    let (w, h) = fit_edge(ow, oh, max_edge);
    let mut scaler = Context::get(frame.format(), ow, oh, ffmpeg::format::pixel::Pixel::RGB24, w, h, ffmpeg::software::scaling::flag::Flags::BILINEAR)
        .map_err(|e| e.to_string())?;
    scaler.run(&frame, &mut rgb).map_err(|e| e.to_string())?;
    let img = image::RgbImage::from_raw(w, h, rgb.data(0).to_vec()).ok_or("bad rgb frame")?;
    resize_to_webp(&image::DynamicImage::ImageRgb8(img), max_edge)
}

fn fit_edge(w: u32, h: u32, max_edge: u32) -> (u32, u32) {
    if w > h {
        (max_edge, max_edge * h / w)
    } else {
        (max_edge * w / h, max_edge)
    }
}

fn resize_to_webp(img: &image::DynamicImage, max_edge: u32) -> Result<Vec<u8>, String> {
    let (w, h) = img.dimensions();
    let thumb = if w > h {
        img.resize(max_edge, max_edge * h / w, image::imageops::FilterType::Triangle)
    } else {
        img.resize(max_edge * w / h, max_edge, image::imageops::FilterType::Triangle)
    };
    let mut out = Vec::new();
    thumb.write_to(&mut Cursor::new(&mut out), image::ImageFormat::WebP).map_err(|e| e.to_string())?;
    Ok(out)
}

fn is_video(path: &Path) -> bool {
    let ext = path
        .extension()
        .and_then(|s| s.to_str())
        .unwrap_or("")
        .to_ascii_lowercase();
    matches!(ext.as_str(), "mp4" | "flv" | "avi" | "rmvb" | "mov" | "mkv")
}

/// 命中缓存则读盘；否则生成 WebP 并落盘。返回 (字节, thumb_hash, 缓存路径)。
pub fn ensure_thumbnail(abs_path: &Path, max_edge: u32) -> Result<(Vec<u8>, String, PathBuf), String> {
    if !abs_path.is_file() {
        return Err("not a file".into());
    }
    let meta = fs::metadata(abs_path).map_err(|e| e.to_string())?;
    let mtime = meta
        .modified()
        .ok()
        .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
        .map(|d| d.as_millis() as i64)
        .unwrap_or(0);
    let path_key = abs_path.to_string_lossy();
    let hash = thumb_hash(&path_key, mtime, max_edge);
    let cache = cache_file(&hash);
    if cache.is_file() {
        let data = fs::read(&cache).map_err(|e| e.to_string())?;
        return Ok((data, hash, cache));
    }
    if let Some(p) = cache.parent() {
        fs::create_dir_all(p).map_err(|e| e.to_string())?;
    }
    let data = if is_video(abs_path) {
        #[cfg(feature = "media")]
        {
            encode_video_thumb(abs_path, max_edge)?
        }
        #[cfg(not(feature = "media"))]
        {
            return Err("video thumb requires media feature".into());
        }
    } else {
        encode_image_thumb(abs_path, max_edge)?
    };
    fs::write(&cache, &data).map_err(|e| e.to_string())?;
    Ok((data, hash, cache))
}
