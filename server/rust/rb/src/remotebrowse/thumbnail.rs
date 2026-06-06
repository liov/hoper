//! 缩略图：生成后写入本地磁盘缓存，供 gRPC 与 HTTP 复用（仅 Agent / host）。
use image::GenericImageView;
use std::fs;
use std::io::Cursor;
use std::path::{Path, PathBuf};

use crate::remotebrowse::thumb_meta::{thumb_cache_dir, thumb_hash, DEFAULT_MAX_EDGE};

#[cfg(feature = "media")]
use ffmpeg_next as ffmpeg;
#[cfg(feature = "media")]
use ffmpeg_next::media::Type;
#[cfg(feature = "media")]
use ffmpeg_next::software::scaling::Context;

fn cache_file(hash: &str) -> PathBuf {
    thumb_cache_dir().join(format!("{hash}.webp"))
}

fn encode_image_thumb(path: &Path, max_edge: u32) -> Result<Vec<u8>, String> {
    let img = image::open(path).map_err(|e| e.to_string())?;
    resize_to_webp(&img, max_edge)
}

#[cfg(feature = "media")]
pub fn decode_still_frame_image(path: &Path) -> Result<image::DynamicImage, String> {
    let frame = ffmpeg_first_video_frame(path, 0)?;
    rgb24_frame_to_image(&frame)
}

#[cfg(feature = "media")]
fn ffmpeg_first_video_frame(path: &Path, seek_us: i64) -> Result<ffmpeg::frame::Video, String> {
    let mut input = ffmpeg::format::input(path).map_err(|e| e.to_string())?;
    if seek_us > 0 {
        let _ = input.seek(seek_us, seek_us..seek_us + 2_000_000);
    }
    let stream = input
        .streams()
        .best(Type::Video)
        .or_else(|| input.streams().find(|s| s.parameters().medium() == ffmpeg::media::Type::Attachment))
        .ok_or("no video/image stream")?;
    let idx = stream.index();
    let ctx = ffmpeg::codec::context::Context::from_parameters(stream.parameters()).map_err(|e| e.to_string())?;
    let mut decoder = ctx.decoder().video().map_err(|e| e.to_string())?;
    let mut frame = ffmpeg::frame::Video::empty();
    let mut packets = 0u32;
    for (s, pkt) in input.packets() {
        if s.index() != idx {
            continue;
        }
        packets += 1;
        if packets > 512 {
            break;
        }
        if decoder.send_packet(&pkt).is_err() {
            continue;
        }
        while decoder.receive_frame(&mut frame).is_ok() {
            if frame.width() > 0 && frame.height() > 0 {
                return Ok(frame);
            }
        }
    }
    let _ = decoder.send_eof();
    while decoder.receive_frame(&mut frame).is_ok() {
        if frame.width() > 0 && frame.height() > 0 {
            return Ok(frame);
        }
    }
    Err("no decodable frame".into())
}

#[cfg(feature = "media")]
fn rgb24_frame_to_image(frame: &ffmpeg::frame::Video) -> Result<image::DynamicImage, String> {
    let w = frame.width();
    let h = frame.height();
    if w == 0 || h == 0 {
        return Err("empty rgb frame".into());
    }
    let stride = frame.stride(0) as usize;
    let row = (w as usize) * 3;
    let mut buf = vec![0u8; row * h as usize];
    let plane = frame.data(0);
    for y in 0..h as usize {
        let off = y * stride;
        if off + row > plane.len() {
            return Err(format!("rgb stride overflow y={y} stride={stride}"));
        }
        buf[y * row..(y + 1) * row].copy_from_slice(&plane[off..off + row]);
    }
    let img = image::RgbImage::from_raw(w, h, buf).ok_or("rgb24 from_raw")?;
    Ok(image::DynamicImage::ImageRgb8(img))
}

#[cfg(feature = "media")]
fn frame_to_webp(frame: ffmpeg::frame::Video, max_edge: u32) -> Result<Vec<u8>, String> {
    let mut rgb = ffmpeg::frame::Video::empty();
    let (ow, oh) = (frame.width(), frame.height());
    let (tw, th) = fit_edge(ow, oh, max_edge);
    let mut scaler = Context::get(
        frame.format(),
        ow,
        oh,
        ffmpeg::format::pixel::Pixel::RGB24,
        tw,
        th,
        ffmpeg::software::scaling::flag::Flags::BILINEAR,
    )
    .map_err(|e| e.to_string())?;
    scaler.run(&frame, &mut rgb).map_err(|e| e.to_string())?;
    let img = rgb24_frame_to_image(&rgb)?;
    encode_dynamic_webp(&img, max_edge)
}

#[cfg(feature = "media")]
fn encode_dynamic_webp(img: &image::DynamicImage, max_edge: u32) -> Result<Vec<u8>, String> {
    let (w, h) = img.dimensions();
    if w <= max_edge && h <= max_edge {
        let mut out = Vec::new();
        img.write_to(&mut Cursor::new(&mut out), image::ImageFormat::WebP).map_err(|e| e.to_string())?;
        return Ok(out);
    }
    resize_to_webp(img, max_edge)
}

#[cfg(feature = "media")]
fn encode_still_via_ffmpeg(path: &Path, max_edge: u32) -> Result<Vec<u8>, String> {
    if let Ok(frame) = ffmpeg_first_video_frame(path, 0) {
        if let Ok(b) = frame_to_webp(frame, max_edge) {
            return Ok(b);
        }
    }
    Err("ffmpeg still decode failed".into())
}

#[cfg(feature = "media")]
fn encode_video_thumb(path: &Path, max_edge: u32) -> Result<Vec<u8>, String> {
    let dur = ffmpeg::format::input(path).ok().map(|i| i.duration()).unwrap_or(0);
    let seek_us = if dur > 0 { (dur / 10).clamp(0, 5_000_000) } else { 0 };
    let frame = ffmpeg_first_video_frame(path, seek_us).or_else(|_| ffmpeg_first_video_frame(path, 0))?;
    frame_to_webp(frame, max_edge)
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
    if w == 0 || h == 0 {
        return Err("empty image".into());
    }
    let (tw, th) = fit_edge(w, h, max_edge);
    let thumb = if tw == w && th == h {
        img.clone()
    } else {
        img.resize_exact(tw, th, image::imageops::FilterType::Triangle)
    };
    let mut out = Vec::new();
    thumb.write_to(&mut Cursor::new(&mut out), image::ImageFormat::WebP).map_err(|e| e.to_string())?;
    Ok(out)
}

fn file_ext(path: &Path) -> String {
    path.extension()
        .and_then(|s| s.to_str())
        .unwrap_or("")
        .to_ascii_lowercase()
}

fn is_video(path: &Path) -> bool {
    path.file_name()
        .and_then(|s| s.to_str())
        .map(crate::remotebrowse::media_ext::is_video_file_name)
        .unwrap_or(false)
}

fn is_still_photo(path: &Path) -> bool {
    path.file_name()
        .and_then(|s| s.to_str())
        .map(crate::remotebrowse::media_ext::is_image_file_name)
        .unwrap_or(false)
}

#[cfg(feature = "media")]
fn encode_still_thumb(path: &Path, max_edge: u32) -> Result<Vec<u8>, String> {
    let name = path.file_name().and_then(|s| s.to_str()).unwrap_or("");
    if crate::remotebrowse::media_ext::prefer_ffmpeg_still_decode(name) {
        if let Ok(b) = encode_still_via_ffmpeg(path, max_edge) {
            return Ok(b);
        }
    }
    match image::open(path) {
        Ok(img) => resize_to_webp(&img, max_edge),
        Err(e) => {
            tracing::debug!(path = %path.display(), err = %e, "image::open failed, try ffmpeg");
            encode_still_via_ffmpeg(path, max_edge)
        }
    }
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
    let path_key = crate::remotebrowse::path::display_path(abs_path);
    let hash = thumb_hash(&path_key, mtime, max_edge);
    let cache = cache_file(&hash);
    if cache.is_file() {
        let data = fs::read(&cache).map_err(|e| e.to_string())?;
        if !data.is_empty() {
            return Ok((data, hash, cache));
        }
    }
    if let Some(p) = cache.parent() {
        fs::create_dir_all(p).map_err(|e| e.to_string())?;
    }
    let data = if is_video(abs_path) {
        #[cfg(feature = "media")]
        {
            encode_video_thumb(abs_path, max_edge).map_err(|e| {
                tracing::warn!(path = %path_key, %e, "video thumbnail failed");
                e
            })?
        }
        #[cfg(not(feature = "media"))]
        {
            return Err("video thumb requires media feature".into());
        }
    } else if is_still_photo(abs_path) {
        #[cfg(feature = "media")]
        {
            encode_still_thumb(abs_path, max_edge)?
        }
        #[cfg(not(feature = "media"))]
        {
            encode_image_thumb(abs_path, max_edge)?
        }
    } else {
        encode_image_thumb(abs_path, max_edge)?
    };
    if data.is_empty() {
        return Err("empty thumbnail".into());
    }
    fs::write(&cache, &data).map_err(|e| e.to_string())?;
    Ok((data, hash, cache))
}
