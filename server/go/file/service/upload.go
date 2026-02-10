package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	errcodex "github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	timex "github.com/hopeio/gox/time"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/file/api/request"
	"github.com/liov/hoper/server/go/file/api/response"
	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.uber.org/zap"
)

const errResp = "上传失败"
const sep = "/"

const (
	ApiExists = "/api/exists"
	ApiUpload = "/api/upload/"
)

// Upload 文件上传
func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(global.Conf.Upload.UploadMaxSize)
	if err != nil {
		httpx.ServeError(w, r, errcode.InvalidArgument.Msg(errResp))
		return
	}

	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpx.ServeError(w, r, errcode.InvalidArgument.Msg(errResp))
		return
	}
	md5Str := r.RequestURI[len(ApiUpload):]
	if md5Str == "" {
		md5Str = r.FormValue("md5")
	}

	var info *multipart.FileHeader
	if fhs := r.MultipartForm.File["file"]; len(fhs) > 0 {
		info = fhs[0]
	}

	ctx, span := Tracer.Start(r.Context(), "Upload")
	defer span.End()
	_, err = auth(ctx, false)
	if err != nil {
		httpx.ServeError(w, r, user.UserErrLogin.Msg(errResp))
		return
	}
	upload, err := save(ctx, info, md5Str)
	if err != nil {
		httpx.ServeError(w, r, errcode.UploadFail.ErrResp())
		return
	}
	(&httpx.CommonAnyResp{Data: response.File{Id: upload.File.Id, URL: upload.File.Path}}).ServeHTTP(w, r)

}

func (*FileService) Exists(ctx context.Context, req *request.Exists) (*response.File, error) {

	ctx, span := Tracer.Start(ctx, "Exists")
	defer span.End()
	auth, err := auth(ctx, false)
	uploadDao := data.GetDao(ctx, global.Dao.GORMDB.DB)
	file, err := uploadDao.FileInfo(req.Md5, req.Size)
	if err != nil {
		return nil, errcode.DBError
	}
	if file != nil {
		upload := model.UploadInfo{
			UserId:    auth.Id,
			CreatedAt: time.Now(),
			FileId:    file.Id,
		}
		if err := uploadDao.Table(model.TableNameUploadInfo).Create(&upload).Error; err != nil {
			log.Errorw("Exists", zap.Error(err))
		}
		return &response.File{Id: file.Id, URL: file.Path}, nil
	}
	return nil, errcode.NotFound
}

func save(ctx context.Context, info *multipart.FileHeader, md5Str string) (upload *model.UploadInfo, err error) {
	uploadDao := data.GetDao(ctx, global.Dao.GORMDB.DB)

	auth, _ := auth(ctx, false)
	var file *model.FileInfo
	if md5Str != "" {
		file, err = uploadDao.FileInfo(md5Str, strconv.FormatInt(info.Size, 10))
		if err != nil {
			return nil, err
		}
	}

	var multipartFile multipart.File
	if file == nil {
		if info == nil {
			return nil, errcode.InvalidArgument
		}
		multipartFile, err = info.Open()
		if err != nil {
			return nil, errcode.IOError.Wrap(err)
		}
		defer multipartFile.Close()
		hash := md5.New()
		_, err = io.Copy(hash, multipartFile)
		if err != nil {
			return nil, errcode.IOError.Wrap(err)
		}
		md5Str = hex.EncodeToString(hash.Sum(nil))
		file, err = uploadDao.FileInfo(md5Str, strconv.FormatInt(info.Size, 10))
		if err != nil {
			return nil, err
		}
	}
	if file != nil {
		upload = &model.UploadInfo{
			UserId:    auth.Id,
			CreatedAt: time.Now(),
			FileId:    file.Id,
		}
		if err = uploadDao.Table(model.TableNameUploadInfo).Create(upload).Error; err != nil {
			return nil, errcode.DBError.Wrap(err)
		}
		return
	}

	ymdStr := timex.GetYMD(time.Now(), sep)

	ext, err := httpx.GetFileExt(info)
	if err != nil {
		return nil, err
	}

	mimeType := mime.TypeByExtension(ext)
	dirType := strings.Split(mimeType, "/")[0]
	if ext == "" {
		dirType = "other"
	}

	uploadDir := dirType + sep + ymdStr + sep
	dir := string(global.Conf.Upload.UploadDir) + uploadDir
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	fileName := md5Str + "_" + strconv.FormatInt(info.Size, 32) + ext
	out, err := os.Create(dir + fileName)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	multipartFile.Seek(0, io.SeekStart)
	_, err = io.Copy(out, multipartFile)
	if err != nil {
		return nil, err
	}
	file = &model.FileInfo{
		Name: info.Filename,
		MD5:  md5Str,
		Size: info.Size,
		Path: uploadDir + fileName,
	}
	err = uploadDao.Table(model.TableNameFileInfo).Create(file).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	upload = &model.UploadInfo{
		FileId:    file.Id,
		UserId:    auth.Id,
		CreatedAt: time.Now(),
	}

	err = uploadDao.Table(model.TableNameUploadInfo).Create(upload).Error
	if err != nil {
		return nil, errcode.DBError.Wrap(err)
	}
	return upload, nil
}

func MultiUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(global.Conf.Upload.UploadMaxSize)
	if err != nil {
		httpx.ServeError(w, r, errcode.InvalidArgument.Msg(errResp))
		return
	}
	ctx, span := Tracer.Start(r.Context(), "MultiUpload")
	defer span.End()
	_, err = auth(ctx, false)
	if err != nil {
		(&httpx.CommonAnyResp{
			Code: errcodex.ErrCode(user.UserErrLogin),
			Msg:  errResp,
		}).ServeHTTP(w, r)
		return
	}
	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpx.ServeError(w, r, errcode.InvalidArgument.Msg(errResp))
		return
	}
	md5s := r.MultipartForm.Value["md5[]"]
	multipartFiles := r.MultipartForm.File["file[]"]
	// 如果有md5
	if len(md5s) != 0 && len(md5s) != len(multipartFiles) {
		httpx.ServeError(w, r, errcode.InvalidArgument.Msg(errResp))
		return
	}
	var urls = make([]response.UploadRes, len(multipartFiles))
	var failures = make([]string, 0)
	for i, multipartFile := range multipartFiles {
		upload, err := save(ctx, multipartFile, md5s[i])
		if err != nil {
			failures = append(failures, multipartFile.Filename)
			httpx.ServeError(w, r, errcode.UploadFail.ErrResp())
			return
		}
		urls[i].Name = multipartFile.Filename
		urls[i].Path = upload.File.Path
	}
	(&httpx.CommonAnyResp{
		Msg:  strings.Join(failures, ",") + errResp,
		Data: urls,
	}).ServeHTTP(w, r)
}
