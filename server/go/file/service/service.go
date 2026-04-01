package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/hopeio/gox/log"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/file"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"time"

	"go.uber.org/zap"
)

type FileService struct {
	file.UnimplementedFileServiceServer
}

func (*FileService) GetUrls(ctx context.Context, req *file.GetUrlsReq) (*file.GetUrlsResp, error) {

	db := global.Dao.GORMDB.DB.WithContext(ctx)
	uploadDao := data.GetDao(db)

	files, err := uploadDao.GetUrls(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}

func (*FileService) GetUrlsByIdsStr(ctx context.Context, req *file.GetUrlsByIdsStrReq) (*file.GetUrlsResp, error) {

	db := global.Dao.GORMDB.DB.WithContext(ctx)
	uploadDao := data.GetDao(db)
	files, err := uploadDao.GetUrlsByStrId(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &file.GetUrlsResp{Files: files}, nil
}

func (*FileService) PreUpload(ctx context.Context, req *file.PreUploadReq) (*file.PreUploadResp, error) {
	auth, err := auth(ctx, false)
	db := global.Dao.GORMDB.DB.WithContext(ctx)
	uploadDao := data.GetDao(db)
	f, err := uploadDao.FileInfo(ctx, req.Md5, strconv.FormatInt(req.Size, 10))
	if err != nil {
		return nil, errcode.DBError
	}
	if f != nil {
		upload := model.UploadInfo{
			UserId:    auth.Id,
			CreatedAt: time.Now(),
			FileId:    f.Id,
		}
		if err := uploadDao.Table(model.TableNameUploadInfo).Create(&upload).Error; err != nil {
			log.Errorw("Exists", zap.Error(err))
		}
		return &file.PreUploadResp{PreUploadType: file.PreUploadType_PRE_UPLOAD_TYPE_EXISTS, File: &file.File{Id: f.Id, Url: f.Path}, Credentials: nil}, nil
	}
	if req.Size > 10*1024*1024*1024 {
		return nil, errcode.InvalidArgument.Msg("file size too large")
	}
	if req.Size < 100*1024*1024 {
		url, err := global.Dao.Minio.Client.PresignedPutObject(ctx, "my-bucket", fmt.Sprintf("uploads/%s", req.Name), time.Duration(15*60)*time.Second)
		if err != nil {
			log.Errorw("PresignedPutObject", zap.Error(err))
			return nil, errcode.IOError.Wrap(err)
		}
		return &file.PreUploadResp{PreUploadType: file.PreUploadType_PRE_UPLOAD_TYPE_UPLOAD_URL, UploadUrl: url.String()}, nil
	}
	if req.Size < 1024*1024*1024 {
		uploadId, presignedURLs, err := InitiateMultipartUpload(ctx, "my-bucket", fmt.Sprintf("uploads/%s", req.Name), 10)
		if err != nil {
			log.Errorw("InitiateMultipartUpload", zap.Error(err))
			return nil, errcode.IOError.Wrap(err)
		}
		return &file.PreUploadResp{
			PreUploadType: file.PreUploadType_PRE_UPLOAD_TYPE_MULTIPART_UPLOAD,
			MultipartUpload: &file.MultipartUpload{
				Id:   uploadId,
				Urls: presignedURLs,
			},
		}, nil
	}
	const policy = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:GetObject",
					"s3:PutObject",
					"s3:AbortMultipartUpload",
					"s3:ListMultipartUploadParts"
				],
				"Resource": [
					"arn:aws:s3:::test-bucket/uploads/*"
				]
			}
		]
	}`

	// 3. 调用 STS AssumeRole 获取临时凭证
	// 这里的 RoleArn 在 MinIO 中通常可以随便填，或者填 'arn:aws:iam::minio-user:sts'
	// 但为了规范，建议按照 AWS 格式填写
	stsProvider := &credentials.STSAssumeRole{
		STSEndpoint: global.Dao.Minio.Conf.Endpoint,
		Options: credentials.STSAssumeRoleOptions{
			AccessKey:       global.Dao.Minio.Conf.AccessKeyID,
			SecretKey:       global.Dao.Minio.Conf.SecretAccessKey,
			Policy:          policy,
			RoleARN:         "arn:aws:iam::123456789012:role/DummyRole",
			RoleSessionName: "FrontendSession",
		},
	}
	stsProvider.SetExpiration(time.Now().Add(time.Hour*1), 0)

	// 4. 获取凭证值
	creds, err := stsProvider.Retrieve()
	if err != nil {
		log.Fatalf("获取 STS 凭证失败: %v", err)
	}
	return &file.PreUploadResp{PreUploadType: file.PreUploadType_PRE_UPLOAD_TYPE_CREDENTIALS, File: nil, Credentials: &file.Credentials{
		AccessKeyId:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
		Expiration:      creds.Expiration.Format(time.RFC3339),
		SignerType:      int32(creds.SignerType),
	}}, nil
}

func InitiateMultipartUpload(ctx context.Context, bucketName, objectName string, totalParts int) (string, []string, error) {
	if totalParts <= 0 {
		return "", nil, errcode.InvalidArgument.Msg("totalParts must be greater than 0")
	}
	// 创建分片上传
	uploadId, err := global.Dao.MinioCore.NewMultipartUpload(ctx, bucketName, objectName, minio.PutObjectOptions{})
	if err != nil {
		return "", nil, errcode.IOError.Wrap(err)
	}
	presignedURLs := make([]string, totalParts)
	for i := 1; i <= totalParts; i++ {
		url, err := global.Dao.MinioCore.Presign(ctx, http.MethodPut, bucketName, objectName, time.Hour, url.Values{
			"uploadId":   []string{uploadId},
			"partNumber": []string{strconv.Itoa(i)},
		})
		if err != nil {
			if abortErr := global.Dao.MinioCore.AbortMultipartUpload(ctx, bucketName, objectName, uploadId); abortErr != nil {
				log.Errorw("AbortMultipartUpload", zap.Error(abortErr), zap.String("upload_id", uploadId))
			}
			return "", nil, errcode.IOError.Wrap(err)
		}
		presignedURLs[i-1] = url.String()
	}
	return uploadId, presignedURLs, nil
}

func CompleteMultipartUpload(ctx context.Context, bucketName, objectName, uploadId string, etags []file.UploadPart) error {
	sort.Slice(etags, func(i, j int) bool { return etags[i].PartNumber < etags[j].PartNumber })
	parts := make([]minio.CompletePart, len(etags))
	for i, part := range etags {
		parts[i] = minio.CompletePart{
			PartNumber: int(part.PartNumber),
			ETag:       part.Etag,
		}
	}

	_, err := global.Dao.MinioCore.CompleteMultipartUpload(ctx, bucketName, objectName, uploadId, parts, minio.PutObjectOptions{})
	if err != nil {
		return errcode.IOError.Wrap(err)
	}
	return nil
}
