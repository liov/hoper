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

	"github.com/hopeio/context/httpctx"
	gormx "github.com/hopeio/gox/database/sql/gorm"
	errcodex "github.com/hopeio/gox/errors"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/fs"
	timex "github.com/hopeio/gox/time"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/file/api/request"
	"github.com/liov/hoper/server/go/file/api/response"
	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/user"
)

const errRep = "上传失败"
const sep = "/"

const (
	ApiExists = "/api/v1/exists"
	ApiUpload = "/api/v1/upload/"
)

// Upload 文件上传
func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(global.Conf.Upload.UploadMaxSize)
	if err != nil {
		httpx.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}

	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpx.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
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

	ctxi, _ := httpctx.FromContext(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpx.RespAnyData{
			Code: errcodex.ErrCode(user.UserErrLogin),
			Msg:  errRep,
		}).Response(w)
		return
	}
	upload, err := save(ctxi, info, md5Str)
	if err != nil {
		httpx.RespErrRep(w, errcode.UploadFail.ErrRep())
		return
	}
	(&httpx.RespAnyData{Data: response.File{Id: upload.File.Id, URL: upload.File.Path}}).Response(w)

}

func (*FileService) Exists(ctx context.Context, req *request.Exists) (*response.File, error) {

	ctxi, _ := httpctx.FromContext(ctx)
	auth, err := auth(ctxi, false)
	uploadDao := data.GetDao(ctxi)
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	file, err := uploadDao.FileInfo(db, req.Md5, req.Size)
	if err != nil {
		return nil, errcode.DBError
	}
	if file != nil {
		upload := model.UploadInfo{
			UserId:    auth.Id,
			CreatedAt: ctxi.RequestAt.Time,
			FileId:    file.Id,
		}
		if err := db.Table(model.TableNameUploadInfo).Create(&upload).Error; err != nil {
			ctxi.RespErrorLog(errcode.DBError, err, "Create")
		}
		return &response.File{Id: file.Id, URL: file.Path}, nil
	}
	return nil, errcode.NotFound
}

func save(ctx *httpctx.Context, info *multipart.FileHeader, md5Str string) (upload *model.UploadInfo, err error) {
	uploadDao := data.GetDao(ctx)
	db := gormx.NewTraceDB(global.Dao.GORMDB.DB, ctx.Base(), ctx.TraceID())
	auth := ctx.AuthInfo.(*user.AuthBase)
	var file *model.FileInfo
	if md5Str != "" {
		file, err = uploadDao.FileInfo(db, md5Str, strconv.FormatInt(info.Size, 10))
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
			return nil, ctx.RespErrorLog(errcode.IOError, err, "Open")
		}
		defer multipartFile.Close()
		hash := md5.New()
		_, err = io.Copy(hash, multipartFile)
		if err != nil {
			return nil, ctx.RespErrorLog(errcode.IOError, err, "Create")
		}
		md5Str = hex.EncodeToString(hash.Sum(nil))
		file, err = uploadDao.FileInfo(db, md5Str, strconv.FormatInt(info.Size, 10))
		if err != nil {
			return nil, err
		}
	}
	if file != nil {
		upload = &model.UploadInfo{
			UserId:    auth.Id,
			CreatedAt: ctx.RequestAt.Time,
			FileId:    file.Id,
		}
		if err = db.Table(model.TableNameUploadInfo).Create(upload).Error; err != nil {
			return nil, ctx.RespErrorLog(errcode.DBError, err, "Create")
		}
		return
	}

	ymdStr := timex.GetYMD(ctx.RequestAt.Time, sep)

	ext, err := fs.GetExt(info)
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
	err = db.Table(model.TableNameFileInfo).Create(file).Error
	if err != nil {
		return nil, ctx.RespErrorLog(errcode.DBError, err, "Create")
	}
	upload = &model.UploadInfo{
		FileId:    file.Id,
		UserId:    auth.Id,
		CreatedAt: ctx.RequestAt.Time,
	}

	err = db.Table(model.TableNameUploadInfo).Create(upload).Error
	if err != nil {
		return nil, ctx.RespErrorLog(errcode.DBError, err, "Create")
	}
	return upload, nil
}

func MultiUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(global.Conf.Upload.UploadMaxSize)
	if err != nil {
		httpx.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}
	ctxi, _ := httpctx.FromContext(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpx.RespAnyData{
			Code: errcodex.ErrCode(user.UserErrLogin),
			Msg:  errRep,
		}).Response(w)
		return
	}
	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpx.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}
	md5s := r.MultipartForm.Value["md5[]"]
	multipartFiles := r.MultipartForm.File["file[]"]
	// 如果有md5
	if len(md5s) != 0 && len(md5s) != len(multipartFiles) {
		httpx.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}
	var urls = make([]model.MultiRep, len(multipartFiles))
	var failures = make([]string, 0)
	for i, multipartFile := range multipartFiles {
		upload, err := save(ctxi, multipartFile, md5s[i])
		if err != nil {
			failures = append(failures, multipartFile.Filename)
			httpx.RespErrRep(w, errcode.UploadFail.ErrRep())
			return
		}
		urls[i].URL = upload.File.Path
		urls[i].Success = true
	}
	(&httpx.RespAnyData{
		Msg:  strings.Join(failures, ",") + errRep,
		Data: urls,
	}).Response(w)
}
