//! 连续 MPEG-TS 转码流（单连接 chunked HTTP，无 HLS 分片边界）。
use std::path::Path;

use crate::remotebrowse::transcode_fragment::{transcode_mpegts, transcode_mux};
use crate::remotebrowse::transcode_preset::TranscodePresetCfg;
use crate::remotebrowse::transcode_vcodec::TranscodeVcodec;

static FFMPEG_INIT: std::sync::Once = std::sync::Once::new();

pub fn transcode_stream(
    input: &Path,
    start_ms: i64,
    preset: TranscodePresetCfg,
    vcodec: TranscodeVcodec,
    mut on_chunk: impl FnMut(&[u8]) -> Result<(), String>,
) -> Result<(), String> {
    FFMPEG_INIT.call_once(|| {
        let _ = ffmpeg_next::init();
    });
    let _guard = transcode_mux().lock().map_err(|e| e.to_string())?;
    transcode_mpegts(input, start_ms, i64::MAX, preset, vcodec, Some(&mut on_chunk)).map(|_| ())
}
