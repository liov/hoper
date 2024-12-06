use std::fs;
use std::io;
use std::path::Path;
use axum::{
    http::StatusCode,
    extract::{Query,Path as AxumPath},
    response::{Response, IntoResponse},
    Json,
};
use std::path::{PathBuf};
use serde::{Deserialize, Serialize};
use crate::file::FileType::File;

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


pub async fn list_files_handler(Query(params): Query<ListFilesParams>) -> Result<Json<Vec<FileInfo>>, (StatusCode, String)> {
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
                } else  {
                    FileType::Directory as i32
                };
                FileInfo {
                    name: entry.file_name().to_string_lossy().to_string(),
                    typ: file_type,
                }
            })
            .collect(),
        Err(e) =>  return Err((StatusCode::INTERNAL_SERVER_ERROR, format!("Error reading directory: {}", e))),
    };

    Ok(Json(files))
}

pub async fn file_thumbnail_handler(AxumPath(path):AxumPath<String>) -> Result<Response, (StatusCode, String)> {
    let path_buf = PathBuf::from(path.clone());
    if !path_buf.exists() || !path_buf.is_file() {
        return Err((StatusCode::NOT_FOUND, "File not found".to_string()));
    }

    // 读取文件内容
    match fs::read(&path_buf) {
        Ok(content) => {
            // 设置适当的头部信息
            let disposition = format!("attachment; filename=\"{}\"", path);
            let headers = [
                (axum::http::header::CONTENT_TYPE, "application/octet-stream"),
                (axum::http::header::CONTENT_DISPOSITION, disposition.as_str()),
            ];

            // 返回带有头部信息和文件内容的响应
            Ok((headers, content).into_response())
        }
        Err(_) => Err((StatusCode::INTERNAL_SERVER_ERROR, "Error reading file".to_string())),
    }
}