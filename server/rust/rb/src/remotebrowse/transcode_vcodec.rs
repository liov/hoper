//! HLS 分片视频编码（不用 H.264）：默认 HEVC，可选 AV1。

use ffmpeg_next::{self as ffmpeg, Dictionary};

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum TranscodeVcodec {
    Hevc,
    Av1,
}

impl TranscodeVcodec {
    pub fn parse(s: &str) -> Option<Self> {
        match s.trim().to_ascii_lowercase().as_str() {
            "hevc" | "h265" | "h.265" => Some(Self::Hevc),
            "av1" => Some(Self::Av1),
            "h264" | "avc" => None,
            _ => None,
        }
    }

    pub fn parse_or_default(s: &str) -> Self {
        Self::parse(s).unwrap_or(Self::Hevc)
    }
}

fn encoder_name(enc: &ffmpeg::codec::codec::Codec) -> &str {
    enc.name()
}

fn is_hevc_videotoolbox(enc: &ffmpeg::codec::codec::Codec) -> bool {
    encoder_name(enc).eq_ignore_ascii_case("hevc_videotoolbox")
}

fn is_hw_hevc(enc: &ffmpeg::codec::codec::Codec) -> bool {
    let n = encoder_name(enc).to_ascii_lowercase();
    n.contains("nvenc") || n.contains("_amf") || n.contains("_qsv") || n.contains("_mf") || n.contains("videotoolbox")
}

fn hevc_hw_encoder_candidates() -> &'static [&'static str] {
    #[cfg(target_os = "windows")]
    {
        &["hevc_nvenc", "hevc_amf", "hevc_qsv", "hevc_mf"]
    }
    #[cfg(target_os = "linux")]
    {
        &["hevc_nvenc", "hevc_vaapi", "hevc_qsv"]
    }
    #[cfg(not(any(target_os = "windows", target_os = "linux")))]
    {
        &[] as &[&str]
    }
}

pub fn find_video_encoder(vcodec: TranscodeVcodec) -> Option<ffmpeg::codec::codec::Codec> {
    match vcodec {
        TranscodeVcodec::Hevc => {
            #[cfg(target_os = "macos")]
            if let Some(c) = ffmpeg::encoder::find_by_name("hevc_videotoolbox") {
                tracing::info!("transcode encoder: hevc_videotoolbox");
                return Some(c);
            }
            for name in hevc_hw_encoder_candidates() {
                if let Some(c) = ffmpeg::encoder::find_by_name(name) {
                    tracing::info!(%name, "transcode encoder");
                    return Some(c);
                }
            }
            if let Some(c) = ffmpeg::encoder::find_by_name("libx265") {
                tracing::info!("transcode encoder: libx265");
                return Some(c);
            }
            ffmpeg::encoder::find(ffmpeg::codec::Id::HEVC)
        }
        TranscodeVcodec::Av1 => {
            ffmpeg::encoder::find_by_name("libsvtav1").or_else(|| ffmpeg::encoder::find(ffmpeg::codec::Id::AV1))
        }
    }
}

pub fn video_encoder_dict(
    enc: &ffmpeg::codec::codec::Codec,
    vcodec: TranscodeVcodec,
    vb: i64,
    bufsize: i64,
    gop: i32,
) -> Dictionary<'static> {
    let mut d = Dictionary::new();
    d.set("b:v", &vb.to_string());
    if is_hevc_videotoolbox(enc) {
        d.set("realtime", "1");
        d.set("allow_sw", "1");
        return d;
    }
    if is_hw_hevc(enc) {
        let n = encoder_name(enc).to_ascii_lowercase();
        if n.contains("nvenc") {
            d.set("preset", "p1");
            d.set("tune", "ll");
            d.set("rc", "vbr");
            d.set("g", &gop.to_string());
            return d;
        }
        d.set("g", &gop.to_string());
        return d;
    }
    d.set("maxrate", &vb.to_string());
    d.set("bufsize", &bufsize.to_string());
    match vcodec {
        TranscodeVcodec::Hevc => {
            d.set("preset", "ultrafast");
            let x265 = format!(
                "keyint={gop}:min-keyint={gop}:open-gop=0:bframes=0:frame-threads=1:pools=none"
            );
            d.set("x265-params", &x265);
        }
        TranscodeVcodec::Av1 => {
            d.set("preset", "12");
            d.set("g", &gop.to_string());
        }
    }
    d
}
