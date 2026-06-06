//! 媒体扩展名：`.ts` 与 TypeScript 区分，仅 HLS/MPEG-TS 分片命名按视频处理。
use std::path::Path;

fn ext_lower(name: &str) -> String {
    Path::new(name)
        .extension()
        .and_then(|s| s.to_str())
        .unwrap_or("")
        .to_ascii_lowercase()
}

fn stem_lower(name: &str) -> String {
    Path::new(name)
        .file_stem()
        .and_then(|s| s.to_str())
        .unwrap_or("")
        .to_ascii_lowercase()
}

/// `.ts` 多为 TypeScript；纯数字或 segment/chunk 等前缀视为 MPEG-TS 分片。
pub fn is_likely_mpeg_ts_segment_name(name: &str) -> bool {
    if ext_lower(name) != "ts" {
        return false;
    }
    let stem = stem_lower(name);
    if stem.is_empty() {
        return false;
    }
    if stem.chars().all(|c| c.is_ascii_digit()) {
        return true;
    }
    ["segment", "chunk", "stream", "media_", "video_"]
        .iter()
        .any(|p| stem.starts_with(p))
}

pub fn is_video_file_name(name: &str) -> bool {
    match ext_lower(name).as_str() {
        "mp4" | "m4v" | "mov" | "mkv" | "webm" | "flv" | "avi" | "rmvb" | "3gp" | "m2ts" | "mts" => true,
        "ts" => is_likely_mpeg_ts_segment_name(name),
        _ => false,
    }
}

pub fn is_image_file_name(name: &str) -> bool {
    matches!(
        ext_lower(name).as_str(),
        "jpg" | "jpeg" | "jfif" | "png" | "gif" | "webp" | "bmp" | "heic" | "heif" | "hif" | "tif" | "tiff" | "avif"
            | "jxl"
    )
}

/// 静图预览/缩略图：优先 ffmpeg 解码（HEIC/AVIF/JXL 等）。
pub fn prefer_ffmpeg_still_decode(name: &str) -> bool {
    matches!(
        ext_lower(name).as_str(),
        "heic" | "heif" | "hif" | "avif" | "jxl" | "tif" | "tiff" | "jpeg" | "jpg" | "jfif"
    )
}

/// Flutter 无法原生解码的静图：预览固定 0.3bpp WebP；点「原图」为高清 WebP。
pub fn preview_transcode_for_client(name: &str) -> bool {
    matches!(
        ext_lower(name).as_str(),
        "heic" | "heif" | "hif" | "avif" | "jxl" | "tif" | "tiff"
    )
}

pub fn is_media_file_name(name: &str) -> bool {
    is_image_file_name(name) || is_video_file_name(name)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn ts_typescript_not_video() {
        assert!(!is_video_file_name("index.ts"));
        assert!(!is_video_file_name("utils.ts"));
        assert!(!is_media_file_name("foo.ts"));
    }

    #[test]
    fn ts_hls_segment_is_video() {
        assert!(is_video_file_name("00001.ts"));
        assert!(is_video_file_name("segment012.ts"));
        assert!(is_media_file_name("chunk_1.ts"));
    }

    #[test]
    fn m2ts_always_video() {
        assert!(is_video_file_name("clip.m2ts"));
    }
}
