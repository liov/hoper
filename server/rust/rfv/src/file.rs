use crate::file::FileType::File;
use axum::http::header;
use axum::{
    extract::{Path as AxumPath, Query},
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use ffmpeg_next as ffmpeg;
use ffmpeg_next::decoder::Decoder;
use ffmpeg_next::media::Type;
use ffmpeg_next::software::scaling::Context;
use image::{DynamicImage, GenericImageView, ImageFormat, ImageReader};
use serde::{Deserialize, Serialize};
use std::fs;
use std::io;
use std::io::Cursor;
use std::path::Path;
use std::path::PathBuf;
use axum::body::Body;
use tracing::debug;
use tokio_util::codec::{BytesCodec, FramedRead};

#[derive(Deserialize)]
pub struct ListFilesParams {
    path: String,
}

#[derive(Serialize)]
enum FileType {
    File = 0,
    Directory = 1,
}

#[derive(Serialize)]
pub struct FileInfo {
    name: String,
    typ: i32,
}

pub async fn list_files_handler(
    Query(params): Query<ListFilesParams>,
) -> Result<Json<Vec<FileInfo>>, (StatusCode, String)> {
    let path = PathBuf::from(&params.path);

    // 检查路径是否存在并且是一个目录
    if !path.exists() || !path.is_dir() {
        return Err((StatusCode::NOT_FOUND, "Directory not found".to_string()));
    }

    // 读取目录中的文件列表
    let files = match fs::read_dir(&path) {
        Ok(entries) => entries
            .filter_map(|entry| entry.ok()) // 忽略错误的条目
            .map(|entry| {
                let file_type = if entry.path().is_file() {
                    FileType::File as i32
                } else {
                    FileType::Directory as i32
                };
                FileInfo {
                    name: entry.file_name().to_string_lossy().to_string(),
                    typ: file_type,
                }
            })
            .collect(),
        Err(e) => {
            return Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("Error reading directory: {}", e),
            ))
        }
    };

    Ok(Json(files))
}

const MAX_SIZE: u32 = 256;
const THUMBNAIL_DIR: &str = ".thumbnails";

pub async fn file_thumbnail_handler(
    AxumPath(path): AxumPath<String>,
) -> Result<Response, (StatusCode, String)> {

    let mut thumbnail_path_buf = PathBuf::from(format!("{}/{}", THUMBNAIL_DIR, path));
    thumbnail_path_buf.set_extension("webp");
    if thumbnail_path_buf.exists()  {
        // 打开文件进行异步读取。
        let file = match tokio::fs::File::open(path).await {
            Ok(file) => file,
            Err(_) => return Err((StatusCode::NOT_FOUND, "File not found".to_string())),
        };

        // 将文件转换为流，以便可以被发送到客户端。
        let stream = FramedRead::new(file, BytesCodec::new());

        // 从流创建一个Body，然后构建响应。
        let body = Body::from_stream(stream);

       return  Ok(Response::builder()
            .header("Content-Type", "application/octet-stream")
            .body(body)
            .unwrap())
    }

    let path_buf = PathBuf::from(format!("D:/{}", path));
    if !path_buf.exists() || !path_buf.is_file() {
        return Err((StatusCode::NOT_FOUND, "File not found".to_string()));
    }
    fs::create_dir_all(thumbnail_path_buf.parent().unwrap()).unwrap();
    if let Some(ext) = path_buf.extension() {
        if let Some(ext_str) = ext.to_str() {
            let lower_ext = ext_str.to_lowercase();
            if ["png", "jpg", "jpeg", "gif", "avif", "webp", "heic", "heif"]
                .contains(&lower_ext.as_str())
            {
                let img = match image::open(&path_buf) {
                    Ok(img) => img,
                    Err(e) => {
                        return Err((
                            StatusCode::INTERNAL_SERVER_ERROR,
                            format!("Error decoding image: {}", e),
                        ))
                    }
                };
                // 计算新的尺寸，保持宽高比。
                let (width, height) = img.dimensions();
                let thumbnail = if width > height {
                    img.resize(
                        MAX_SIZE,
                        MAX_SIZE * height / width,
                        image::imageops::FilterType::Triangle,
                    )
                } else {
                    img.resize(
                        MAX_SIZE * width / height,
                        MAX_SIZE,
                        image::imageops::FilterType::Triangle,
                    )
                };
                let mut thumb_bytes: Vec<u8> = Vec::new();
                thumbnail
                    .write_to(&mut Cursor::new(&mut thumb_bytes), ImageFormat::WebP)
                    .unwrap();
                fs::write(thumbnail_path_buf, &thumb_bytes).unwrap();
                // 构建 HTTP 响应。
                return Ok((
                    [
                        (header::CONTENT_TYPE, "image/webp"),
                        (header::CACHE_CONTROL, "public, max-age=31536000"),
                    ],
                    thumb_bytes,
                )
                    .into_response());
            } else if ["mp4", "flv", "avi", "rmvb"].contains(&lower_ext.as_str()) {
                // 打开视频文件。
                let mut input = match ffmpeg::format::input(&path_buf) {
                    Ok(ctx) => ctx,
                    Err(_) => {
                        return Err((StatusCode::NOT_FOUND, "Video file not found".to_string()))
                    }
                };
                let context = match input.streams().best(Type::Video) {
                    Some(ctx) => ctx,
                    None => {
                        return Err((StatusCode::NOT_FOUND, "Video file not found".to_string()))
                    }
                };
                let video_stream_index = context.index();
                // 提取第一帧作为缩略图。
                // 获取解码器上下文。
                let codec_context =
                    ffmpeg::codec::context::Context::from_parameters(context.parameters()).unwrap();
                let mut decoder = codec_context.decoder().video().unwrap();
                // 准备接收帧的变量。
                let mut key_frame_count = 0; // 用于计数关键帧的数量
                let mut frame_count = 0; // 用于计数帧的数量
                let mut frame = ffmpeg::frame::Video::empty();
                'lable: for (stream, packet) in input.packets() {
                    if stream.index() == video_stream_index {
                        debug!("stream: {:?}", stream);
                        decoder.send_packet(&packet).unwrap();
                        while decoder.receive_frame(&mut frame).is_ok() {
                            if frame.is_key() {
                                // 检查当前帧是否是关键帧
                                key_frame_count += 1;
                                if key_frame_count == 10 {
                                    // 只需要第二个关键帧
                                    break 'lable;
                                }
                            }
                            frame_count += 1;
                            if frame_count >= 100 {
                                break 'lable;
                            }
                        }
                    }
                }
                decoder.send_eof().unwrap();
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
                let mut thumb_bytes: Vec<u8> = Vec::new();
                img.write_to(&mut Cursor::new(&mut thumb_bytes), ImageFormat::WebP)
                    .unwrap();
                fs::write(thumbnail_path_buf, &thumb_bytes).unwrap();
                // 构建 HTTP 响应。
                return Ok((
                    [
                        (header::CONTENT_TYPE, "image/webp"),
                        (header::CACHE_CONTROL, "public, max-age=31536000"),
                    ],
                    thumb_bytes,
                )
                    .into_response());
            }
        }
    }
    // 读取文件内容
    match fs::read(&path_buf) {
        Ok(content) => {
            let mime_type = mime_guess::from_path(&path)
                .first_or_octet_stream()
                .to_string();
            // 设置适当的头部信息
            let headers = if mime_type == "application/octet-stream" {
                [
                    (header::CONTENT_TYPE, mime_type),
                    (
                        header::CONTENT_DISPOSITION,
                        format!(
                            "attachment; filename=\"{}\"",
                            path_buf.file_name().unwrap().to_str().unwrap()
                        ),
                    ),
                ]
            } else {
                [
                    (header::CONTENT_TYPE, mime_type),
                    (
                        header::CACHE_CONTROL,
                        "public, max-age=31536000".parse().unwrap(),
                    ),
                ]
            };

            // 返回带有头部信息和文件内容的响应
            Ok((headers, content).into_response())
        }
        Err(_) => Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            "Error reading file".to_string(),
        )),
    }
}
