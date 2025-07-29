use ffmpeg_next as ffmpeg;
use ffmpeg_next::software::scaling::Context;
use image::{DynamicImage, ImageFormat};
use std::io::Cursor;
use std::path::Path;
use ffmpeg_next::codec;
use ffmpeg_next::media::Type;

fn main() {
    ffmpeg::init().unwrap();
    let path = Path::new(r"D:\giphy.mp4");
    // 打开视频文件。
    let mut input =ffmpeg::format::input(&path).unwrap();
    let context = input.streams().best(Type::Video).unwrap();
    let video_stream_index = context.index();
    // 获取解码器上下文。
    let codec_context =
        codec::context::Context::from_parameters(context.parameters()).unwrap();
    let mut decoder = codec_context.decoder().video().unwrap();
    // 准备接收帧的变量。
    let mut key_frame_count = 0; // 用于计数关键帧的数量
    let mut frame_count = 0; // 用于计数帧的数量

    // 提取第一帧作为缩略图。
    const MAX_SIZE: u32 = 256;
    let mut frame = ffmpeg::frame::Video::empty();
    for (stream, packet) in input.packets() {
        if stream.index() == video_stream_index {
            decoder.send_packet(&packet).unwrap();
            while let Ok(_) = decoder.receive_frame(&mut frame) {
                if frame.is_key() {
                    // 检查当前帧是否是关键帧
                    key_frame_count += 1;
                    if key_frame_count == 3 {
                        // 只需要第二个关键帧
                        break;
                    }
                }
                frame_count += 1;
                if frame_count >= 10 {
                    break;
                }
            }
        }
    }
    decoder.send_eof().unwrap();
    unsafe {
        if frame.is_empty() {
            panic!("No frame found!")
        }
    }
    // 将帧转换为 RGB 格式。
    let mut rgb_frame = ffmpeg::frame::Video::empty();
    let owidth = frame.width();
    let oheight = frame.height();
    let (mut width, mut height) = (owidth, oheight);
    if owidth > oheight {
        width = MAX_SIZE;
        height = MAX_SIZE * oheight / owidth;
    } else {
        height = MAX_SIZE;
        width = MAX_SIZE * owidth / oheight;
    }
    let mut scaler = Context::get(
        frame.format(),
        frame.width(),
        frame.height(),
        ffmpeg::format::pixel::Pixel::RGB24,
        width,
        height,
        ffmpeg::software::scaling::flag::Flags::BILINEAR,
    )
    .unwrap();

    scaler.run(&frame, &mut rgb_frame).unwrap();
    // 将 RGB 数据转换为图像。
    let img = DynamicImage::ImageRgb8(
        image::RgbImage::from_raw(width, height, rgb_frame.data(0).to_vec()).unwrap(),
    );
    // 将缩略图编码为 WebP 字节流。
    img.save("result.webp").unwrap();
}
