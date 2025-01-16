package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/protobuf/errcode"

	gormi "github.com/hopeio/utils/dao/database/gorm"
	errcode2 "github.com/hopeio/utils/errors/errcode"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/fs"
	timei "github.com/hopeio/utils/time"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/upload/data"
	"github.com/liov/hoper/server/go/upload/global"
	"github.com/liov/hoper/server/go/upload/model"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const errRep = "上传失败"
const sep = "/"

const (
	ApiExists = "/api/v1/exists"
	ApiUpload = "/api/v1/upload/"
)

// Upload 文件上传
func Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(global.Conf.Customize.UploadMaxSize)
	if err != nil {
		httpi.RespErrRep(w, errcode.ParamInvalid.Origin().Msg(errRep))
		return
	}

	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpi.RespErrRep(w, errcode.ParamInvalid.Origin().Msg(errRep))
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

	ctxi := httpctx.FromContextValue(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.ResAnyData{
			Code: errcode2.ErrCode(user.UserErrLogin),
			Msg:  errRep,
		}).Response(w, http.StatusOK)
		return
	}
	upload, err := save(ctxi, info, md5Str)
	if err != nil {
		httpi.RespErrRep(w, errcode.UploadFail.Origin().ErrRep())
		return
	}
	(&httpi.ResAnyData{Data: model.Rep{Id: upload.Id, URL: upload.Path}}).Response(w, http.StatusOK)

}

func Exists(w http.ResponseWriter, req *http.Request) {
	md5 := req.URL.Query().Get("md5")
	size := req.URL.Query().Get("size")
	exists(req.Context(), w, md5, size)
}

func exists(ctx context.Context, w http.ResponseWriter, md5, size string) {
	ctxi := httpctx.FromContextValue(ctx)
	auth, err := auth(ctxi, false)
	uploadDao := data.GetDao(ctxi)
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	upload, err := uploadDao.UploadDB(db, md5, size)
	if err != nil {
		httpi.RespErrRep(w, errcode.DBError.Origin().ErrRep())
		return
	}
	if upload != nil {
		uploadExt := model.UploadExt{
			UserId:    auth.Id,
			CreatedAt: ctxi.RequestAt.Time,
			UploadId:  upload.Id,
		}
		if err := db.Table(model.UploadExtTableName).Create(&uploadExt).Error; err != nil {
			ctxi.RespErrorLog(errcode.DBError, err, "Create")
		}
		(&httpi.ResAnyData{
			Code: 1,
			Msg:  "已存在",
			Data: model.Rep{Id: upload.Id, URL: upload.Path},
		}).Response(w, http.StatusOK)
		return
	}
	(&httpi.ResAnyData{Msg: "不存在"}).Response(w, http.StatusOK)
}

func save(ctx *httpctx.Context, info *multipart.FileHeader, md5Str string) (upload *model.UploadInfo, err error) {
	uploadDao := data.GetDao(ctx)
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx.Base(), ctx.TraceID())
	auth := ctx.AuthInfo.(*user.AuthBase)
	if md5Str != "" {
		upload, err = uploadDao.UploadDB(db, md5Str, strconv.FormatInt(info.Size, 10))
		if err != nil {
			return nil, err
		}
	}

	var file multipart.File
	if upload == nil {
		if info == nil {
			return nil, errcode.ParamInvalid
		}
		file, err = info.Open()
		if err != nil {
			return nil, ctx.RespErrorLog(errcode.IOError, err, "Open")
		}
		defer file.Close()
		hash := md5.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return nil, ctx.RespErrorLog(errcode.IOError, err, "Create")
		}
		md5Str = hex.EncodeToString(hash.Sum(nil))
		upload, err = uploadDao.UploadDB(db, md5Str, strconv.FormatInt(info.Size, 10))
		if err != nil {
			return nil, err
		}
	}
	if upload != nil {
		uploadExt := model.UploadExt{
			UserId:    auth.Id,
			CreatedAt: ctx.RequestAt.Time,
			UploadId:  upload.Id,
		}
		if err = db.Table(model.UploadExtTableName).Create(&uploadExt).Error; err != nil {
			return nil, ctx.RespErrorLog(errcode.DBError, err, "Create")
		}
		return
	}

	ymdStr := timei.GetYMD(ctx.RequestAt.Time, sep)

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
	dir := string(global.Conf.Customize.UploadDir) + uploadDir
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}

	fileName := md5Str + "_" + strconv.FormatInt(info.Size, 32) + ext
	out, err := os.Create(dir + fileName)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	file.Seek(0, io.SeekStart)
	_, err = io.Copy(out, file)
	if err != nil {
		return nil, err
	}
	fileUpload := model.UploadInfo{
		File: model.File{
			Name: info.Filename,
			MD5:  md5Str,
			Ext:  ext,
			Size: info.Size,
		},
		UserId:    auth.Id,
		Path:      uploadDir + fileName,
		CreatedAt: ctx.RequestAt.Time,
	}

	err = db.Table(model.UploadTableName).Create(&fileUpload).Error
	if err != nil {
		return nil, ctx.RespErrorLog(errcode.DBError, err, "Create")
	}
	return &fileUpload, nil
}

func MultiUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(global.Conf.Customize.UploadMaxSize)
	if err != nil {
		httpi.RespErrRep(w, errcode.ParamInvalid.Origin().Msg(errRep))
		return
	}
	ctxi := httpctx.FromContextValue(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.ResAnyData{
			Code: errcode2.ErrCode(user.UserErrLogin),
			Msg:  errRep,
		}).Response(w, http.StatusOK)
		return
	}
	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpi.RespErrRep(w, errcode.ParamInvalid.Origin().Msg(errRep))
		return
	}
	md5s := r.MultipartForm.Value["md5[]"]
	files := r.MultipartForm.File["file[]"]
	// 如果有md5
	if len(md5s) != 0 && len(md5s) != len(files) {
		httpi.RespErrRep(w, errcode.ParamInvalid.Origin().Msg(errRep))
		return
	}
	var urls = make([]model.MultiRep, len(files))
	var failures = make([]string, 0)
	for i, file := range files {
		upload, err := save(ctxi, file, md5s[i])
		if err != nil {
			failures = append(failures, file.Filename)
			httpi.RespErrRep(w, errcode.UploadFail.Origin().ErrRep())
			return
		}
		urls[i].URL = upload.Path
		urls[i].Success = true
	}
	(&httpi.ResAnyData{
		Msg:  strings.Join(failures, ",") + errRep,
		Data: urls,
	}).Response(w, http.StatusOK)
}
