//! 按需转码单个 HLS 分片：ffmpeg-next 解码/编码 + 内存 MPEG-TS，不写盘。
use std::path::Path;
use std::sync::{Mutex, Once, OnceLock};

use ffmpeg_next as ffmpeg;
use ffmpeg::codec::{self, encoder, Context};
use ffmpeg::format::context::Output;
use ffmpeg::format::{self, Pixel};
use ffmpeg::media::Type;
use ffmpeg::packet::Packet;
use ffmpeg::software::resampling::Context as AudioResampler;
use ffmpeg::software::scaling::{context::Context as Scaler, flag::Flags as ScaleFlags};
use ffmpeg::util::mathematics::rescale::Rescale;
use ffmpeg::util::rational::Rational;
use ffmpeg::{frame, ChannelLayout, Dictionary};

use crate::remotebrowse::mpegts_mem_out::{MpegtsMemOut, MpegtsStreamOut};
use crate::remotebrowse::transcode_preset::{TranscodePresetCfg, FRAGMENT_MS};
use crate::remotebrowse::transcode_vcodec::{find_video_encoder, video_encoder_dict, TranscodeVcodec};

static FFMPEG_INIT: Once = Once::new();
static TRANSCODE_MUX: OnceLock<Mutex<()>> = OnceLock::new();

pub(crate) fn transcode_mux() -> &'static Mutex<()> {
    TRANSCODE_MUX.get_or_init(|| Mutex::new(()))
}

pub fn transcode_fragment(
    input: &Path,
    start_ms: i64,
    preset: TranscodePresetCfg,
    vcodec: TranscodeVcodec,
) -> Result<Vec<u8>, String> {
    match transcode_fragments_serial(input, &[start_ms], preset, vcodec)?.into_iter().next() {
        Some((_, b)) => Ok(b),
        None => Err("empty".to_string()),
    }
}

/// 单次持锁连续转多片，避免逐片 spawn 与重复调度导致播放卡在 2s 边界。
pub fn transcode_fragments_serial(
    input: &Path,
    starts_ms: &[i64],
    preset: TranscodePresetCfg,
    vcodec: TranscodeVcodec,
) -> Result<Vec<(i64, Vec<u8>)>, String> {
    FFMPEG_INIT.call_once(|| {
        let _ = ffmpeg::init();
    });
    let _guard = transcode_mux().lock().map_err(|e| e.to_string())?;
    if starts_ms.is_empty() {
        return Ok(Vec::new());
    }
    let mut ictx = format::input(input).map_err(|e| e.to_string())?;
    let mut out = Vec::with_capacity(starts_ms.len());
    for &start_ms in starts_ms {
        let end_us = start_ms.saturating_mul(1000).saturating_add(FRAGMENT_MS.saturating_mul(1000));
        let bytes = transcode_mpegts_ictx(&mut ictx, start_ms, end_us, preset, vcodec, None)?;
        out.push((start_ms, bytes));
    }
    Ok(out)
}

/// 转码为 MPEG-TS；[end_us] 为排他上界，`i64::MAX` 表示直至文件结束（连续流）。
enum TsMux {
    Mem(MpegtsMemOut),
    Stream(MpegtsStreamOut),
}

impl TsMux {
    fn output_mut(&mut self) -> &mut Output {
        match self {
            Self::Mem(m) => m.output_mut(),
            Self::Stream(s) => s.output_mut(),
        }
    }

    fn mark_header_written(&mut self) {
        match self {
            Self::Mem(m) => m.mark_header_written(),
            Self::Stream(s) => s.mark_header_written(),
        }
    }

    fn emit_pending(&mut self, cb: &mut dyn FnMut(&[u8]) -> Result<(), String>) -> Result<(), String> {
        if let Self::Stream(s) = self {
            s.emit_pending(|b| cb(b))?;
        }
        Ok(())
    }

    fn finish(self, cb: Option<&mut dyn FnMut(&[u8]) -> Result<(), String>>) -> Result<Vec<u8>, String> {
        match self {
            Self::Mem(m) => m.into_bytes(),
            Self::Stream(s) => {
                if let Some(cb) = cb {
                    s.finish(|b| cb(b))?;
                }
                Ok(Vec::new())
            }
        }
    }
}

pub(crate) fn transcode_mpegts(
    input: &Path,
    start_ms: i64,
    end_us: i64,
    preset: TranscodePresetCfg,
    vcodec: TranscodeVcodec,
    on_chunk: Option<&mut dyn FnMut(&[u8]) -> Result<(), String>>,
) -> Result<Vec<u8>, String> {
    let mut ictx = format::input(input).map_err(|e| e.to_string())?;
    transcode_mpegts_ictx(&mut ictx, start_ms, end_us, preset, vcodec, on_chunk)
}

fn seek_ictx(ictx: &mut format::context::Input, start_ms: i64) -> Result<(), String> {
    let start_us = start_ms.saturating_mul(1000);
    if start_us <= 0 {
        return Ok(());
    }
    let span_us = FRAGMENT_MS.saturating_mul(1000);
    let min_ts = start_us.saturating_sub(span_us);
    ictx.seek(start_us, min_ts..start_us.saturating_add(1)).map_err(|e| e.to_string())
}

fn transcode_mpegts_ictx(
    ictx: &mut format::context::Input,
    start_ms: i64,
    end_us: i64,
    preset: TranscodePresetCfg,
    vcodec: TranscodeVcodec,
    mut on_chunk: Option<&mut dyn FnMut(&[u8]) -> Result<(), String>>,
) -> Result<Vec<u8>, String> {
    seek_ictx(ictx, start_ms)?;
    let mut mux = if on_chunk.is_some() {
        TsMux::Stream(MpegtsStreamOut::open()?)
    } else {
        TsMux::Mem(MpegtsMemOut::open()?)
    };
    let start_us = start_ms.saturating_mul(1000);
    let v_ist = ictx.streams().best(Type::Video).ok_or("no video stream")?;
    let v_idx = v_ist.index();
    let v_tb = v_ist.time_base();
    let a_ist = ictx.streams().best(Type::Audio);
    let a_idx = a_ist.as_ref().map(|s| s.index());
    let stream_mode = end_us == i64::MAX;
    let octx = mux.output_mut();
    let mut v = VideoLane::open(&v_ist, octx, preset, vcodec, start_us, end_us, stream_mode)?;
    let mut a = a_ist
        .as_ref()
        .map(|s| AudioLane::open(s, octx, preset, start_us, end_us))
        .transpose()?;
    octx.write_header().map_err(|e| e.to_string())?;
    mux.mark_header_written();
    let mut v_done = false;
    for (stream, packet) in ictx.packets() {
        let idx = stream.index();
        if idx == v_idx && !v_done {
            if packet_abs_us(v_tb, packet.pts()).is_some_and(|t| t >= end_us) {
                v_done = true;
                continue;
            }
            v.process_packet(&packet, mux.output_mut())?;
            if let Some(cb) = on_chunk.as_mut() {
                mux.emit_pending(cb)?;
            }
        } else if a_idx == Some(idx) {
            if let Some(lane) = a.as_mut() {
                let tb = stream.time_base();
                if packet_abs_us(tb, packet.pts()).is_some_and(|t| t >= end_us) {
                    continue;
                }
                lane.process_packet(&packet, mux.output_mut())?;
                if let Some(cb) = on_chunk.as_mut() {
                    mux.emit_pending(cb)?;
                }
            }
        }
    }
    let _ = v.flush(mux.output_mut());
    if let Some(lane) = a.as_mut() {
        let _ = lane.flush(mux.output_mut());
    }
    if let Some(cb) = on_chunk.as_mut() {
        mux.emit_pending(cb)?;
    }
    mux.finish(on_chunk)
}

fn packet_abs_us(tb: Rational, pts: Option<i64>) -> Option<i64> {
    pts.map(|p| p.rescale(tb, Rational(1, 1_000_000)))
}

fn even_dim(v: u32) -> u32 {
    v & !1
}

fn scaled_size(w: u32, h: u32, max_h: u32) -> (u32, u32) {
    if h == 0 || w == 0 {
        return (even_dim(w), even_dim(h.max(2)));
    }
    if h <= max_h {
        return (even_dim(w), even_dim(h));
    }
    let nh = max_h;
    let nw = (w as u64 * nh as u64 / h as u64) as u32;
    (even_dim(nw.max(2)), even_dim(nh))
}

fn gop_for_fragment(fps: Rational) -> usize {
    let n = fps.numerator().max(1) as i64;
    let d = fps.denominator().max(1) as i64;
    let frames = n * FRAGMENT_MS / (d * 1000);
    (frames.max(12).min(120)) as usize
}

fn gop_for_stream(fps: Rational) -> usize {
    let n = fps.numerator().max(1) as i64;
    let d = fps.denominator().max(1) as i64;
    let frames = n * FRAGMENT_MS * 4 / (d * 1000);
    (frames.max(24).min(300)) as usize
}

struct VideoLane {
    ost: usize,
    decoder: codec::decoder::Video,
    encoder: codec::encoder::Video,
    scaler: Scaler,
    enc_tb: Rational,
    stream_tb: Rational,
    dec_tb: Rational,
    start_us: i64,
    end_us: i64,
    out_w: u32,
    out_h: u32,
}

impl VideoLane {
    fn open(
        ist: &format::stream::Stream,
        octx: &mut Output,
        preset: TranscodePresetCfg,
        vcodec: TranscodeVcodec,
        start_us: i64,
        end_us: i64,
        stream_mode: bool,
    ) -> Result<Self, String> {
        let dec_tb = ist.time_base();
        let dec = Context::from_parameters(ist.parameters()).map_err(|e| e.to_string())?.decoder().video().map_err(|e| e.to_string())?;
        let (out_w, out_h) = scaled_size(dec.width(), dec.height(), preset.max_height);
        let enc_codec = find_video_encoder(vcodec).ok_or_else(|| format!("no encoder for {vcodec:?}"))?;
        let mut ost = octx.add_stream(enc_codec).map_err(|e| e.to_string())?;
        let mut enc = Context::new_with_codec(enc_codec).encoder().video().map_err(|e| e.to_string())?;
        let stream_tb = Rational(1, 90_000);
        let fps = dec.frame_rate().unwrap_or(Rational(24, 1));
        let gop = if stream_mode { gop_for_stream(fps) } else { gop_for_fragment(fps) };
        enc.set_width(out_w);
        enc.set_height(out_h);
        enc.set_aspect_ratio(dec.aspect_ratio());
        enc.set_format(Pixel::YUV420P);
        enc.set_time_base(stream_tb);
        enc.set_frame_rate(dec.frame_rate());
        enc.set_bit_rate(preset.video_bitrate as usize);
        enc.set_max_bit_rate(preset.video_bitrate as usize);
        enc.set_gop(gop as u32);
        let bufsize = preset.video_bitrate * 2;
        let opened = enc
            .open_with(video_encoder_dict(&enc_codec, vcodec, preset.video_bitrate, bufsize, gop as i32))
            .map_err(|e| e.to_string())?;
        ost.set_parameters(&opened);
        ost.set_time_base(stream_tb);
        let enc_tb = opened.time_base();
        let scaler = Scaler::get(dec.format(), dec.width(), dec.height(), Pixel::YUV420P, out_w, out_h, ScaleFlags::BILINEAR)
            .map_err(|e| e.to_string())?;
        Ok(Self {
            ost: ost.index(),
            decoder: dec,
            encoder: opened,
            scaler,
            enc_tb,
            stream_tb,
            dec_tb,
            start_us,
            end_us,
            out_w,
            out_h,
        })
    }

    fn process_packet(&mut self, packet: &Packet, octx: &mut Output) -> Result<(), String> {
        self.decoder.send_packet(packet).map_err(|e| e.to_string())?;
        self.drain_frames(octx)
    }

    fn drain_frames(&mut self, octx: &mut Output) -> Result<(), String> {
        let mut decoded = frame::Video::empty();
        let mut scaled = frame::Video::empty();
        while self.decoder.receive_frame(&mut decoded).is_ok() {
            let Some(p) = decoded.timestamp() else { continue };
            let abs_us = p.rescale(self.dec_tb, Rational(1, 1_000_000));
            if abs_us < self.start_us {
                continue;
            }
            if abs_us >= self.end_us {
                return Ok(());
            }
            scaled.set_width(self.out_w);
            scaled.set_height(self.out_h);
            scaled.set_format(Pixel::YUV420P);
            self.scaler.run(&decoded, &mut scaled).map_err(|e| e.to_string())?;
            let ts_90k = abs_us.rescale(Rational(1, 1_000_000), Rational(1, 90_000));
            scaled.set_pts(Some(ts_90k));
            self.encoder.send_frame(&scaled).map_err(|e| e.to_string())?;
            self.write_encoded(octx)?;
        }
        Ok(())
    }

    fn write_encoded(&mut self, octx: &mut Output) -> Result<(), String> {
        let mut pkt = Packet::empty();
        while self.encoder.receive_packet(&mut pkt).is_ok() {
            pkt.set_stream(self.ost);
            pkt.rescale_ts(self.enc_tb, self.stream_tb);
            pkt.write_interleaved(octx).map_err(|e| e.to_string())?;
        }
        Ok(())
    }

    fn flush(&mut self, octx: &mut Output) -> Result<(), String> {
        self.decoder.send_eof().map_err(|e| e.to_string())?;
        self.drain_frames(octx)?;
        self.encoder.send_eof().map_err(|e| e.to_string())?;
        self.write_encoded(octx)
    }
}

struct AudioLane {
    ost: usize,
    decoder: codec::decoder::Audio,
    encoder: codec::encoder::Audio,
    resampler: AudioResampler,
    enc_tb: Rational,
    stream_tb: Rational,
    dec_tb: Rational,
    start_us: i64,
    end_us: i64,
}

impl AudioLane {
    fn open(ist: &format::stream::Stream, octx: &mut Output, preset: TranscodePresetCfg, start_us: i64, end_us: i64) -> Result<Self, String> {
        let dec_tb = ist.time_base();
        let dec = Context::from_parameters(ist.parameters()).map_err(|e| e.to_string())?.decoder().audio().map_err(|e| e.to_string())?;
        let enc_codec = encoder::find(codec::Id::AAC).ok_or("no aac encoder")?;
        let mut ost = octx.add_stream(enc_codec).map_err(|e| e.to_string())?;
        let mut enc = Context::new_with_codec(enc_codec).encoder().audio().map_err(|e| e.to_string())?;
        let stream_tb = Rational(1, 48_000);
        let layout = ChannelLayout::STEREO;
        enc.set_rate(48_000);
        enc.set_channel_layout(layout);
        enc.set_format(ffmpeg::format::Sample::F32(ffmpeg::format::sample::Type::Planar));
        enc.set_bit_rate(preset.audio_bitrate as usize);
        enc.set_time_base(stream_tb);
        let mut opts = Dictionary::new();
        opts.set("profile", "aac_low");
        let opened = enc.open_with(opts).map_err(|e| e.to_string())?;
        ost.set_parameters(&opened);
        ost.set_time_base(stream_tb);
        let enc_tb = opened.time_base();
        let resampler = AudioResampler::get(
            dec.format(),
            dec.channel_layout(),
            dec.rate(),
            opened.format(),
            opened.channel_layout(),
            opened.rate(),
        )
        .map_err(|e| e.to_string())?;
        Ok(Self {
            ost: ost.index(),
            decoder: dec,
            encoder: opened,
            resampler,
            enc_tb,
            stream_tb,
            dec_tb,
            start_us,
            end_us,
        })
    }

    fn process_packet(&mut self, packet: &Packet, octx: &mut Output) -> Result<(), String> {
        self.decoder.send_packet(packet).map_err(|e| e.to_string())?;
        self.drain_frames(octx)
    }

    fn drain_frames(&mut self, octx: &mut Output) -> Result<(), String> {
        let mut decoded = frame::Audio::empty();
        let mut converted = frame::Audio::empty();
        while self.decoder.receive_frame(&mut decoded).is_ok() {
            let Some(p) = decoded.timestamp() else { continue };
            let abs_us = p.rescale(self.dec_tb, Rational(1, 1_000_000));
            if abs_us < self.start_us {
                continue;
            }
            if abs_us >= self.end_us {
                return Ok(());
            }
            self.resampler.run(&decoded, &mut converted).map_err(|e| e.to_string())?;
            let ts_48k = abs_us.rescale(Rational(1, 1_000_000), Rational(1, 48_000));
            converted.set_pts(Some(ts_48k));
            self.encoder.send_frame(&converted).map_err(|e| e.to_string())?;
            self.write_encoded(octx)?;
        }
        Ok(())
    }

    fn write_encoded(&mut self, octx: &mut Output) -> Result<(), String> {
        let mut pkt = Packet::empty();
        while self.encoder.receive_packet(&mut pkt).is_ok() {
            pkt.set_stream(self.ost);
            pkt.rescale_ts(self.enc_tb, self.stream_tb);
            pkt.write_interleaved(octx).map_err(|e| e.to_string())?;
        }
        Ok(())
    }

    fn flush(&mut self, octx: &mut Output) -> Result<(), String> {
        self.decoder.send_eof().map_err(|e| e.to_string())?;
        self.drain_frames(octx)?;
        self.encoder.send_eof().map_err(|e| e.to_string())?;
        self.write_encoded(octx)
    }
}

#[cfg(all(test, feature = "media"))]
mod tests {
    use std::path::Path;

    use super::*;
    use crate::remotebrowse::transcode_preset;

    fn smoke_path() -> &'static Path {
        Path::new("/Users/jyb/Downloads/WhatsApp Video 2026-03-25 at 5.48.42 PM (2).mp4")
    }

    #[test]
    fn fragment_hevc_smoke() {
        let p = smoke_path();
        if !p.exists() {
            return;
        }
        let preset = transcode_preset::preset_cfg("720").expect("720 preset");
        let out = transcode_fragment(p, 0, preset, TranscodeVcodec::Hevc).expect("transcode");
        assert!(out.len() > 1024, "ts too small: {}", out.len());
    }

    #[test]
    fn fragment_hevc_at_2s() {
        let p = smoke_path();
        if !p.exists() {
            return;
        }
        let preset = transcode_preset::preset_cfg("720").expect("720 preset");
        let out = transcode_fragment(p, 2000, preset, TranscodeVcodec::Hevc).expect("transcode");
        assert!(out.len() > 1024, "ts too small at 2s: {}", out.len());
    }

    #[test]
    fn fragment_hevc_at_26s() {
        let p = smoke_path();
        if !p.exists() {
            return;
        }
        let preset = transcode_preset::preset_cfg("720").expect("720 preset");
        let out = transcode_fragment(p, 26000, preset, TranscodeVcodec::Hevc).expect("transcode");
        assert!(out.len() > 1024, "ts too small at 26s: {}", out.len());
    }
}
