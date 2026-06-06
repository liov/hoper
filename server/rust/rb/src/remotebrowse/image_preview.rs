//! 图片预览：源 bpp≤0.3 且 Flutter 可解码 → 直传原文件；否则 0.3bpp WebP。
//! `#rb-original` / `preview=2`：全分辨率高清 WebP（目标 bpp 2.0），非磁盘原文件。
use std::fs;
use std::path::Path;

use image::GenericImageView;

use crate::remotebrowse::media_ext::{is_image_file_name, preview_transcode_for_client, prefer_ffmpeg_still_decode};
use crate::remotebrowse::thumb_meta::{thumb_cache_dir, thumb_hash};

pub const PREVIEW_BPP_THRESHOLD: f64 = 0.3;
pub const PREVIEW_TARGET_BPP: f64 = 0.3;
pub const PREVIEW_MIN_WEBP_QUALITY: f32 = 75.0;
pub const HQ_TARGET_BPP: f64 = 2.0;
pub const HQ_MIN_WEBP_QUALITY: f32 = 88.0;
/// 与 [thumb_hash] 的 max_edge 区分：0=预览 WebP 缓存，1=高清 WebP 缓存。
pub const PREVIEW_CACHE_EDGE: u32 = 0;
pub const HQ_CACHE_EDGE: u32 = 1;

fn file_bpp(file_size: u64, width: u32, height: u32) -> f64 {
    let pixels = (width as u64).saturating_mul(height as u64);
    if pixels == 0 {
        return 0.0;
    }
    file_size as f64 * 8.0 / pixels as f64
}

fn target_bytes(width: u32, height: u32, bpp: f64) -> usize {
    ((bpp * width as f64 * height as f64) / 8.0).ceil() as usize
}

fn file_name(path: &Path) -> &str {
    path.file_name().and_then(|s| s.to_str()).unwrap_or("")
}

fn preview_cache_path(hash: &str) -> std::path::PathBuf {
    thumb_cache_dir().join(format!("{hash}.webp"))
}

fn image_dimensions(path: &Path) -> Result<(u32, u32), String> {
    if let Ok(d) = image::image_dimensions(path) {
        return Ok(d);
    }
    #[cfg(feature = "media")]
    {
        let img = crate::remotebrowse::thumbnail::decode_still_frame_image(path)?;
        return Ok(img.dimensions());
    }
    #[cfg(not(feature = "media"))]
    {
        Err("cannot read image dimensions".into())
    }
}

fn encode_webp_quality(img: &image::DynamicImage, quality: f32) -> Result<Vec<u8>, String> {
    let rgb = img.to_rgb8();
    let (w, h) = rgb.dimensions();
    let enc = webp::Encoder::from_rgb(rgb.as_raw(), w, h);
    Ok(enc.encode(quality.clamp(1.0, 100.0)).to_vec())
}

fn encode_webp_target_bpp(img: &image::DynamicImage, target_bpp: f64, min_q: f32) -> Result<Vec<u8>, String> {
    let (w, h) = img.dimensions();
    if w == 0 || h == 0 {
        return Err("empty image".into());
    }
    let want = target_bytes(w, h, target_bpp).max(256);
    let floor = min_q.clamp(1.0, 100.0);
    let mut lo = floor;
    let mut hi = 100.0f32;
    let mut best = encode_webp_quality(img, hi)?;
    if best.len() <= want {
        return Ok(best);
    }
    for _ in 0..8 {
        let q = (lo + hi) * 0.5;
        let cur = encode_webp_quality(img, q)?;
        if cur.len() > want {
            hi = q;
        } else {
            best = cur;
            lo = q;
        }
    }
    if best.len() > want {
        return encode_webp_quality(img, floor);
    }
    Ok(best)
}

fn load_still(path: &Path) -> Result<image::DynamicImage, String> {
    let name = file_name(path);
    #[cfg(feature = "media")]
    if prefer_ffmpeg_still_decode(name) {
        if let Ok(img) = crate::remotebrowse::thumbnail::decode_still_frame_image(path) {
            return Ok(img);
        }
    }
    if let Ok(img) = image::open(path) {
        return Ok(img);
    }
    #[cfg(feature = "media")]
    {
        return crate::remotebrowse::thumbnail::decode_still_frame_image(path);
    }
    #[cfg(not(feature = "media"))]
    {
        Err("unsupported image".into())
    }
}

fn needs_preview_webp(name: &str, bpp: f64) -> bool {
    preview_transcode_for_client(name) || bpp > PREVIEW_BPP_THRESHOLD
}

fn read_webp_cache(path: &Path, cache_edge: u32) -> Result<Option<Vec<u8>>, String> {
    let meta = fs::metadata(path).map_err(|e| e.to_string())?;
    let mtime = meta
        .modified()
        .ok()
        .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
        .map(|d| d.as_millis() as i64)
        .unwrap_or(0);
    let path_key = crate::remotebrowse::path::display_path(path);
    let hash = thumb_hash(&path_key, mtime, cache_edge);
    let cache = preview_cache_path(&hash);
    if !cache.is_file() {
        return Ok(None);
    }
    let data = fs::read(&cache).map_err(|e| e.to_string())?;
    if data.is_empty() {
        Ok(None)
    } else {
        Ok(Some(data))
    }
}

fn write_webp_cache(path: &Path, cache_edge: u32, data: &[u8]) -> Result<(), String> {
    let meta = fs::metadata(path).map_err(|e| e.to_string())?;
    let mtime = meta
        .modified()
        .ok()
        .and_then(|t| t.duration_since(std::time::UNIX_EPOCH).ok())
        .map(|d| d.as_millis() as i64)
        .unwrap_or(0);
    let hash = thumb_hash(&crate::remotebrowse::path::display_path(path), mtime, cache_edge);
    let cache = preview_cache_path(&hash);
    if let Some(p) = cache.parent() {
        fs::create_dir_all(p).map_err(|e| e.to_string())?;
    }
    fs::write(&cache, data).map_err(|e| e.to_string())
}

fn encode_cached_webp(path: &Path, cache_edge: u32, target_bpp: f64, min_q: f32) -> Result<Vec<u8>, String> {
    if let Some(cached) = read_webp_cache(path, cache_edge)? {
        return Ok(cached);
    }
    let img = load_still(path)?;
    let data = encode_webp_target_bpp(&img, target_bpp, min_q)?;
    if data.is_empty() {
        return Err("empty webp".into());
    }
    let _ = write_webp_cache(path, cache_edge, &data);
    Ok(data)
}

fn encode_preview_webp(path: &Path) -> Result<Vec<u8>, String> {
    encode_cached_webp(path, PREVIEW_CACHE_EDGE, PREVIEW_TARGET_BPP, PREVIEW_MIN_WEBP_QUALITY)
}

fn encode_hq_webp(path: &Path) -> Result<Vec<u8>, String> {
    encode_cached_webp(path, HQ_CACHE_EDGE, HQ_TARGET_BPP, HQ_MIN_WEBP_QUALITY)
}

/// 需 Agent 侧重编码时返回全分辨率 WebP；否则 `None`（ReadFile 直传原图字节）。
pub fn maybe_preview_webp(path: &Path, file_size: u64) -> Result<Option<Vec<u8>>, String> {
    if !path.is_file() || !is_image_file_name(file_name(path)) {
        return Ok(None);
    }
    let name = file_name(path);
    if preview_transcode_for_client(name) {
        return encode_preview_webp(path).map(Some);
    }
    let (w, h) = match image_dimensions(path) {
        Ok(v) => v,
        Err(_) => return Ok(None),
    };
    let bpp = file_bpp(file_size, w, h);
    if !needs_preview_webp(name, bpp) {
        return Ok(None);
    }
    encode_preview_webp(path).map(Some)
}

/// 「原图」：全分辨率高清 WebP（`#rb-original` / `preview=2`）。
pub fn maybe_hq_webp(path: &Path, file_size: u64) -> Result<Vec<u8>, String> {
    if !path.is_file() || !is_image_file_name(file_name(path)) {
        return Err("not an image".into());
    }
    let _ = file_size;
    encode_hq_webp(path)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn bpp_calc() {
        assert!((file_bpp(1_000_000, 1000, 1000) - 8.0).abs() < 0.001);
        assert!(file_bpp(37_500, 1000, 1000) <= PREVIEW_BPP_THRESHOLD);
    }

    #[test]
    fn avif_always_transcode() {
        assert!(needs_preview_webp("photo.avif", 0.1));
        assert!(!needs_preview_webp("photo.jpg", 0.2));
        assert!(needs_preview_webp("photo.jpg", 0.5));
    }

    #[test]
    fn client_formats_force_transcode_flag() {
        assert!(preview_transcode_for_client("a.heic"));
        assert!(preview_transcode_for_client("b.avif"));
        assert!(preview_transcode_for_client("c.jxl"));
        assert!(preview_transcode_for_client("scan.tiff"));
        assert!(!preview_transcode_for_client("d.jpg"));
    }
}
