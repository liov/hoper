package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/hopeio/lemon/context/http_context"
	"github.com/hopeio/lemon/protobuf/errorcode"
	httpi "github.com/hopeio/lemon/utils/net/http"
	"github.com/hopeio/lemon/utils/net/http/fs"
	timei "github.com/hopeio/lemon/utils/time"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/upload/confdao"
	"github.com/liov/hoper/server/go/upload/data"
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
	err := r.ParseMultipartForm(confdao.Conf.Customize.UploadMaxSize)
	if err != nil {
		errorcode.ParamInvalid.OriMessage(errRep).Response(w)
		return
	}

	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		errorcode.ParamInvalid.OriMessage(errRep).Response(w)
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

	ctxi := http_context.ContextFromContext(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.ResAnyData{
			Code:    errorcode.ErrCode(user.UserErrLogin),
			Message: errRep,
		}).Response(w, http.StatusOK)
		return
	}
	upload, err := save(ctxi, info, md5Str)
	if err != nil {
		errorcode.UploadFail.OriErrRep().Response(w)
		return
	}
	(&httpi.ResAnyData{Details: model.Rep{Id: upload.Id, URL: upload.Path}}).Response(w, http.StatusOK)

}

func Exists(w http.ResponseWriter, req *http.Request) {
	md5 := req.URL.Query().Get("md5")
	size := req.URL.Query().Get("size")
	exists(req.Context(), w, md5, size)
}

func exists(ctx context.Context, w http.ResponseWriter, md5, size string) {
	ctxi := http_context.ContextFromContext(ctx)
	auth, err := auth(ctxi, false)
	uploadDao := data.GetDao(ctxi)
	db := ctxi.NewDB(confdao.Dao.GORMDB.DB)
	upload, err := uploadDao.UploadDB(db, md5, size)
	if err != nil {
		errorcode.DBError.OriErrRep().Response(w)
		return
	}
	if upload != nil {
		uploadExt := model.UploadExt{
			UserId:    auth.Id,
			CreatedAt: ctxi.RequestAt.Time,
			UploadId:  upload.Id,
		}
		if err := db.Table(model.UploadExtTableName).Create(&uploadExt).Error; err != nil {
			ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}
		(&httpi.ResAnyData{
			Code:    1,
			Message: "已存在",
			Details: model.Rep{Id: upload.Id, URL: upload.Path},
		}).Response(w, http.StatusOK)
		return
	}
	(&httpi.ResAnyData{Message: "不存在"}).Response(w, http.StatusOK)
}

func save(ctx *http_context.Context, info *multipart.FileHeader, md5Str string) (upload *model.UploadInfo, err error) {
	uploadDao := data.GetDao(ctx)
	db := ctx.NewDB(confdao.Dao.GORMDB.DB)
	auth := ctx.AuthInfo.(*user.AuthInfo)
	if md5Str != "" {
		upload, err = uploadDao.UploadDB(db, md5Str, strconv.FormatInt(info.Size, 10))
		if err != nil {
			return nil, err
		}
	}

	var file multipart.File
	if upload == nil {
		if info == nil {
			return nil, errorcode.ParamInvalid
		}
		file, err = info.Open()
		if err != nil {
			return nil, ctx.ErrorLog(errorcode.IOError, err, "Open")
		}
		defer file.Close()
		hash := md5.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return nil, ctx.ErrorLog(errorcode.IOError, err, "Create")
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
			return nil, ctx.ErrorLog(errorcode.DBError, err, "Create")
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
	dir := string(confdao.Conf.Customize.UploadDir) + uploadDir
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
		return nil, ctx.ErrorLog(errorcode.DBError, err, "Create")
	}
	return &fileUpload, nil
}

func MultiUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(confdao.Conf.Customize.UploadMaxSize)
	if err != nil {
		errorcode.ParamInvalid.OriMessage(errRep).Response(w)
		return
	}
	ctxi := http_context.ContextFromContext(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.ResAnyData{
			Code:    errorcode.ErrCode(user.UserErrLogin),
			Message: errRep,
		}).Response(w, http.StatusOK)
		return
	}
	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		errorcode.ParamInvalid.OriMessage(errRep).Response(w)
		return
	}
	md5s := r.MultipartForm.Value["md5[]"]
	files := r.MultipartForm.File["file[]"]
	// 如果有md5
	if len(md5s) != 0 && len(md5s) != len(files) {
		errorcode.ParamInvalid.OriMessage(errRep).Response(w)
		return
	}
	var urls = make([]model.MultiRep, len(files))
	var failures = make([]string, 0)
	for i, file := range files {
		upload, err := save(ctxi, file, md5s[i])
		if err != nil {
			failures = append(failures, file.Filename)
			errorcode.UploadFail.OriErrRep().Response(w)
			return
		}
		urls[i].URL = upload.Path
		urls[i].Success = true
	}
	(&httpi.ResAnyData{
		Message: strings.Join(failures, ",") + errRep,
		Details: urls,
	}).Response(w, http.StatusOK)
}
