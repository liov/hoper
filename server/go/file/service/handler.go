package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/scaffold/errcode"
	gormi "github.com/hopeio/utils/dao/database/gorm"
	errcode2 "github.com/hopeio/utils/errors/errcode"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/fs"
	timei "github.com/hopeio/utils/time"
	"github.com/liov/hoper/server/go/file/data"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/liov/hoper/server/go/protobuf/user"
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
		httpi.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}

	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpi.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
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

	ctxi, _ := httpctx.FromContextValue(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.RespAnyData{
			Code: errcode2.ErrCode(user.UserErrLogin),
			Msg:  errRep,
		}).Response(w)
		return
	}
	upload, err := save(ctxi, info, md5Str)
	if err != nil {
		httpi.RespErrRep(w, errcode.UploadFail.ErrRep())
		return
	}
	(&httpi.RespAnyData{Data: model.Rep{Id: upload.File.Id, URL: upload.File.Path}}).Response(w)

}

func Exists(w http.ResponseWriter, req *http.Request) {
	md5 := req.URL.Query().Get("md5")
	size := req.URL.Query().Get("size")
	exists(req.Context(), w, md5, size)
}

func exists(ctx context.Context, w http.ResponseWriter, md5, size string) {
	ctxi, _ := httpctx.FromContextValue(ctx)
	auth, err := auth(ctxi, false)
	uploadDao := data.GetDao(ctxi)
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx, ctxi.TraceID())
	file, err := uploadDao.FileInfo(db, md5, size)
	if err != nil {
		httpi.RespErrRep(w, errcode.DBError.ErrRep())
		return
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
		(&httpi.RespAnyData{
			Code: 1,
			Msg:  "已存在",
			Data: model.Rep{Id: file.Id, URL: file.Path},
		}).Response(w)
		return
	}
	(&httpi.RespAnyData{Msg: "不存在"}).Response(w)
}

func save(ctx *httpctx.Context, info *multipart.FileHeader, md5Str string) (upload *model.UploadInfo, err error) {
	uploadDao := data.GetDao(ctx)
	db := gormi.NewTraceDB(global.Dao.GORMDB.DB, ctx.Base(), ctx.TraceID())
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
	err := r.ParseMultipartForm(global.Conf.Customize.UploadMaxSize)
	if err != nil {
		httpi.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}
	ctxi, _ := httpctx.FromContextValue(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.RespAnyData{
			Code: errcode2.ErrCode(user.UserErrLogin),
			Msg:  errRep,
		}).Response(w)
		return
	}
	if r.MultipartForm == nil || (r.MultipartForm.Value == nil && r.MultipartForm.File == nil) {
		httpi.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}
	md5s := r.MultipartForm.Value["md5[]"]
	multipartFiles := r.MultipartForm.File["file[]"]
	// 如果有md5
	if len(md5s) != 0 && len(md5s) != len(multipartFiles) {
		httpi.RespErrRep(w, errcode.InvalidArgument.Msg(errRep))
		return
	}
	var urls = make([]model.MultiRep, len(multipartFiles))
	var failures = make([]string, 0)
	for i, multipartFile := range multipartFiles {
		upload, err := save(ctxi, multipartFile, md5s[i])
		if err != nil {
			failures = append(failures, multipartFile.Filename)
			httpi.RespErrRep(w, errcode.UploadFail.ErrRep())
			return
		}
		urls[i].URL = upload.File.Path
		urls[i].Success = true
	}
	(&httpi.RespAnyData{
		Msg:  strings.Join(failures, ",") + errRep,
		Data: urls,
	}).Response(w)
}
