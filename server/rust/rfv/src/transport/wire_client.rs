//! Viewer 侧 wire：向 Agent 请求列表/缩略图。
use prost::Message;

#[cfg(feature = "media")]
use crate::grpc_server::proto::{FileEntry, ListFilesRequest, ListFilesResponse, ThumbnailRequest, ThumbnailResponse};

use crate::transport::link::ViewerLink;

const TYPE_FILE_INDEX: u8 = 2;
const TYPE_THUMB_REQ: u8 = 3;
const TYPE_THUMB_DATA: u8 = 4;

pub struct ListedFile {
    pub name: String,
    pub size: i64,
    pub id: String,
    pub thumb_hash: String,
}

pub async fn list_files(mut link: ViewerLink, root: &str) -> Result<Vec<ListedFile>, String> {
    let req = ListFilesRequest { root_path: root.into(), ..Default::default() };
    write_link(&mut link, TYPE_FILE_INDEX, &req.encode_to_vec()).await?;
    let (typ, payload) = read_link(&mut link).await?;
    if typ != TYPE_FILE_INDEX {
        return Err(format!("unexpected wire type {typ}"));
    }
    let resp = ListFilesResponse::decode(payload.as_slice()).map_err(|e| e.to_string())?;
    Ok(resp
        .entries
        .into_iter()
        .map(|e: FileEntry| ListedFile {
            name: e.name,
            size: e.size,
            id: e.id,
            thumb_hash: e.thumb_hash,
        })
        .collect())
}

pub async fn fetch_thumb(mut link: ViewerLink, path: &str, max_edge: u32) -> Result<Vec<u8>, String> {
    let req = ThumbnailRequest { path: path.into(), max_edge, ..Default::default() };
    write_link(&mut link, TYPE_THUMB_REQ, &req.encode_to_vec()).await?;
    let (typ, payload) = read_link(&mut link).await?;
    if typ != TYPE_THUMB_DATA {
        return Err(format!("unexpected wire type {typ}"));
    }
    let resp = ThumbnailResponse::decode(payload.as_slice()).map_err(|e| e.to_string())?;
    Ok(resp.data)
}

async fn write_link(link: &mut ViewerLink, typ: u8, payload: &[u8]) -> Result<(), String> {
    match link {
        ViewerLink::Tcp(t) => t.write_frame(typ, payload).await,
        ViewerLink::Ice(i) => i.write_frame(typ, payload).await,
    }
}

async fn read_link(link: &mut ViewerLink) -> Result<(u8, Vec<u8>), String> {
    match link {
        ViewerLink::Tcp(t) => t.read_frame().await,
        ViewerLink::Ice(i) => i.read_frame().await,
    }
}
