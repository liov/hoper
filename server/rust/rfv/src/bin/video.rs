use std::io::Cursor;
use std::path::Path;
use ffmpeg_next as ffmpeg;
use ffmpeg_next::software::scaling::Context;
use image::{DynamicImage, ImageFormat};

fn main(){
    let path = Path::new("");
    // 打开视频文件。
    let mut context =  ffmpeg::format::input(&path).unwrap();
    // 提取第一帧作为缩略图。
    const max_size = 256i64;
    let mut frame = ffmpeg::frame::Video::empty();
    for (stream, packet) in context.packets() {
        if stream.parameters().medium() == ffmpeg::media::Type::Video {
            // 获取解码器上下文。
            let codec_context = ffmpeg::codec::context::Context::from_parameters(
                stream.parameters(),
            ).unwrap();
            let mut decoder = codec_context.decoder().video().unwrap();
            // 准备接收帧的变量。
            let mut key_frame_count = 0; // 用于计数关键帧的数量
            let mut frame_count = 0; // 用于计数帧的数量
            decoder.send_packet(&packet).unwrap();

            while let Ok(_) =  decoder.receive_frame(&mut frame) {

                        if frame.is_key() { // 检查当前帧是否是关键帧
                            key_frame_count += 1;
                            if key_frame_count == 3 { // 只需要第二个关键帧
                                break;
                            }
                        }
                        frame_count += 1;
                        if frame_count >= 10 {
                            break;
                        }
                }
            }
            break;
    }
    // 将帧转换为 RGB 格式。
    let mut rgb_frame = ffmpeg::frame::Video::empty();
    let mut width = frame.width();
    let mut height = frame.height();
    if width > height {
        width=max_size;
        height=max_size*height/width;
    } else {
        height = max_size;
        width=max_size*width/height;
    }
    let mut scaler = Context::get(
        frame.format(),
        frame.width(),
        frame.height(),
        ffmpeg::format::pixel::Pixel::RGB24,
        width,
        height,
        ffmpeg::software::scaling::flag::Flags::BILINEAR,
    ).unwrap();

    scaler.run(&frame, &mut rgb_frame).unwrap();
    // 将 RGB 数据转换为图像。
    let img = DynamicImage::ImageRgb8(image::RgbImage::from_raw(
        width,
        height,
        rgb_frame.data(0).to_vec(),
    ).unwrap());
    // 将缩略图编码为 WebP 字节流。
    let mut thumb_bytes: Vec<u8> = Vec::new();
    img.write_to(&mut Cursor::new(&mut thumb_bytes), ImageFormat::WebP).unwrap();
}