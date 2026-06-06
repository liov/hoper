//! 转码档位（按需分片，不落盘整片）。
#[derive(Clone, Copy)]
pub struct TranscodePresetCfg {
    pub max_height: u32,
    pub video_bitrate: i64,
    pub audio_bitrate: i64,
}

pub fn preset_cfg(preset: &str) -> Option<TranscodePresetCfg> {
    match preset {
        "1080" => Some(TranscodePresetCfg { max_height: 1080, video_bitrate: 4_000_000, audio_bitrate: 160_000 }),
        "720" => Some(TranscodePresetCfg { max_height: 720, video_bitrate: 2_000_000, audio_bitrate: 128_000 }),
        "480" => Some(TranscodePresetCfg { max_height: 480, video_bitrate: 800_000, audio_bitrate: 96_000 }),
        "360" => Some(TranscodePresetCfg { max_height: 360, video_bitrate: 400_000, audio_bitrate: 64_000 }),
        _ => None,
    }
}

/// 每片时长（毫秒），与 m3u8 `#EXTINF` 一致。
pub const FRAGMENT_MS: i64 = 2000;
