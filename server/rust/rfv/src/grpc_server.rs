//! gRPC：Go `webrtc` 包通过本服务访问 rfv 媒体能力。
use std::path::PathBuf;

use tokio::sync::mpsc;
use tokio_stream::wrappers::ReceiverStream;
use tonic::{Request, Response, Status, Streaming};

pub mod proto {
    tonic::include_proto!("remotebrowse");
}

use proto::remote_browse_service_server::{RemoteBrowseService, RemoteBrowseServiceServer};
use proto::{
    FileEntry, HealthResponse, ListFilesRequest, ListFilesResponse, ThumbnailChunk, ThumbnailRequest,
    ThumbnailResponse,
};

#[derive(Default)]
struct MediaSvc;

#[tonic::async_trait]
impl RemoteBrowseService for MediaSvc {
    async fn list_files(
        &self,
        req: Request<ListFilesRequest>,
    ) -> Result<Response<ListFilesResponse>, Status> {
        let root = req.into_inner().root_path;
        let rows = crate::remotebrowse::list_remote_files(&root).map_err(Status::invalid_argument)?;
        let entries = rows
            .into_iter()
            .map(|e| FileEntry {
                id: e.id,
                name: e.name,
                size: e.size,
                mtime_unix_ms: e.mtime_unix_ms,
                mime: e.mime,
                thumb_hash: e.thumb_hash,
                ..Default::default()
            })
            .collect();
        Ok(Response::new(ListFilesResponse {
            entries,
            next_cursor: String::new(),
        }))
    }

    async fn get_thumbnail(
        &self,
        req: Request<ThumbnailRequest>,
    ) -> Result<Response<ThumbnailResponse>, Status> {
        let inner = req.into_inner();
        let max_edge = if inner.max_edge == 0 {
            crate::remotebrowse::DEFAULT_MAX_EDGE
        } else {
            inner.max_edge
        };
        let path = PathBuf::from(inner.path);
        let (data, hash, cache) = tokio::task::spawn_blocking(move || {
            crate::remotebrowse::ensure_thumbnail(&path, max_edge)
        })
        .await
        .map_err(|e| Status::internal(e.to_string()))?
        .map_err(Status::invalid_argument)?;
        Ok(Response::new(ThumbnailResponse {
            data,
            mime: "image/webp".into(),
            thumb_hash: hash,
            cache_path: cache.to_string_lossy().into_owned(),
        }))
    }

    type ThumbnailPipeStream = ReceiverStream<Result<ThumbnailChunk, Status>>;

    async fn thumbnail_pipe(
        &self,
        req: Request<Streaming<ThumbnailRequest>>,
    ) -> Result<Response<Self::ThumbnailPipeStream>, Status> {
        let mut inbound = req.into_inner();
        let (tx, rx) = mpsc::channel(4);
        tokio::spawn(async move {
            while let Ok(Some(job)) = inbound.message().await {
                let path = PathBuf::from(job.path);
                let edge = if job.max_edge == 0 {
                    crate::remotebrowse::DEFAULT_MAX_EDGE
                } else {
                    job.max_edge
                };
                let out = tokio::task::spawn_blocking(move || {
                    crate::remotebrowse::ensure_thumbnail(&path, edge)
                })
                .await;
                let chunk = match out {
                    Ok(Ok((data, hash, _))) => ThumbnailChunk {
                        data,
                        thumb_hash: hash,
                        mime: "image/webp".into(),
                        done: true,
                        ..Default::default()
                    },
                    Ok(Err(e)) => ThumbnailChunk {
                        error: e,
                        done: true,
                        ..Default::default()
                    },
                    Err(e) => ThumbnailChunk {
                        error: e.to_string(),
                        done: true,
                        ..Default::default()
                    },
                };
                if tx.send(Ok(chunk)).await.is_err() {
                    break;
                }
            }
        });
        Ok(Response::new(ReceiverStream::new(rx)))
    }

    async fn get_health(
        &self,
        _req: Request<()>,
    ) -> Result<Response<HealthResponse>, Status> {
        Ok(Response::new(HealthResponse {
            signal_ws: "/rb/signal".into(),
            relay_tcp: String::new(),
            rfv_grpc: std::env::var("RFV_GRPC").unwrap_or_else(|_| "127.0.0.1:50051".into()),
            thumb_cache: crate::remotebrowse::thumb_cache_dir().to_string_lossy().into_owned(),
        }))
    }
}

pub async fn serve(addr: String) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    tonic::transport::Server::builder()
        .add_service(RemoteBrowseServiceServer::new(MediaSvc::default()))
        .serve(addr.parse()?)
        .await?;
    Ok(())
}
