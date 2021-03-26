package upload

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"github.com/liov/hoper/go/v2/upload/config"
	"github.com/liov/hoper/go/v2/upload/dao"
	"github.com/liov/hoper/go/v2/upload/model"
	"github.com/liov/hoper/go/v2/utils/fs"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	timei "github.com/liov/hoper/go/v2/utils/time"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"net/http"
)

// Upload 文件上传
func Upload(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(config.Conf.Customize.UploadMaxSize)
	md5Var := req.FormValue("md5")
	file, info, err := req.FormFile("file")
	if err != nil {
		errorcode.ParamInvalid.Origin().Response(w)
		return
	}
	defer file.Close()
	ctxi := user.CtxFromContext(req.Context())
	auth, err := ctxi.GetAuthInfo(Auth)
	uploadDao := dao.GetDao(ctxi)
	db := dao.Dao.GetDB(ctxi.Logger)
	upload, err := uploadDao.UploadDB(db, md5Var, strconv.FormatInt(info.Size, 10))
	if err != nil {
		errorcode.DBError.Origin().Response(w)
		return
	}
	if upload != nil {
		upload.Id = 0
		upload.UserId = auth.Id
		upload.UUID = uuid.New().String()
		upload.CreatedAt = ctxi.RequestAt.Time
		if err := db.Create(&upload).Error; err != nil {
			ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}
		(&httpi.ResData{
			Code:    0,
			Message: "",
			Details: upload.URL,
		}).Response(w)
		return
	}

	upload = &model.FileUploadInfo{UserId: auth.Id}

	hash := md5.New()
	io.Copy(hash, file)
	md5Str := hex.EncodeToString(hash.Sum(nil))
	upload.MD5 = md5Str

	ext, err := fs.GetExt(info)
	if ext == "" || err != nil {
		errorcode.ParamInvalid.Origin().Message("无效的图片类型").Response(w)
		return
	}

	dir, url, err := GetDirAndUrl(classify, info)

	upInfo, err := SaveUploadedFile(info, dir, url)
	if err != nil {
		common.Response(ctx, nil, err.Error(), e.ERROR)
		return nil
	}

	upInfo.File.Size = uint64(info.Size)
	upInfo.UploadUserID = userID
	upInfo.Status = 1
	upInfo.MD5 = md5
	if err := initialize.DB.Create(upInfo).Error; err != nil {
		common.Response(ctx, nil, err.Error(), e.ERROR)
		return nil
	}
	common.Response(ctx, upInfo, "", e.SUCCESS)
	return upInfo
}

func MD5(w http.ResponseWriter, req *http.Request) {
	md5 := req.URL.Query().Get("md5")
	size := req.URL.Query().Get("size")
	ctxi := user.CtxFromContext(req.Context())
	auth, err := ctxi.GetAuthInfo(Auth)
	uploadDao := dao.GetDao(ctxi)
	db := dao.Dao.GetDB(ctxi.Logger)
	upload, err := uploadDao.UploadDB(db, md5, size)
	if err != nil {
		errorcode.DBError.Origin().Response(w)
		return
	}
	if upload != nil {
		upload.Id = 0
		upload.UserId = auth.Id
		upload.UUID = uuid.New().String()
		upload.CreatedAt = ctxi.RequestAt.Time
		if err := db.Create(&upload).Error; err != nil {
			ctxi.ErrorLog(errorcode.DBError, err, "Create")
		}
		(&httpi.ResData{
			Code:    0,
			Message: "",
			Details: upload.URL,
		}).Response(w)
		return
	}
}

func GetDirAndUrl(classify string, info *multipart.FileHeader) (string, string, error) {
	//sep := string(os.PathSeparator)
	sep := "/"
	var uploadDir, prefixUrl string
	ymdStr := timei.GetTodayYMD(sep)
	ext, err := fs.GetExt(info)
	if err != nil {
		return "", "", err
	}

	if ext == "" {
		uploadDir = strings.Join([]string{string(config.Conf.Customize.UploadDir),
			"others",
			classify,
			ymdStr},
			"/")
		prefixUrl = strings.Join([]string{
			config.Conf.Customize.UploadPath,
			"others",
			classify,
			ymdStr,
		}, "/")
		return uploadDir, prefixUrl, nil
	}

	var mimeType = mime.TypeByExtension(ext)
	if mimeType == "" && ext == ".jpeg" {
		mimeType = "image/jpeg"
	}

	dirType := strings.Split(mimeType, "/")

	uploadDir = strings.Join([]string{initialize.Config.Server.UploadDir,
		dirType[0] + "s",
		classify,
		ymdStr},
		"/")

	/*	length := utf8.RuneCountInString(uploadDir)
		lastChar := uploadDir[length-1:]

		if lastChar != sep {
			uploadDir = uploadDir + sep + ymdStr
		} else {
			uploadDir = uploadDir + ymdStr
		}
	*/

	prefixUrl = strings.Join([]string{
		initialize.Config.Server.UploadPath,
		dirType[0] + "s",
		classify,
		ymdStr,
	}, "/")

	if err := os.MkdirAll(uploadDir, 0777); err != nil {
		return uploadDir, prefixUrl, err
	}
	return uploadDir, prefixUrl, nil
}


func UploadMultiple(w http.ResponseWriter, req *http.Request) {
	userID := ctx.Values().Get("userID").(uint64)
	classify := ctx.Params().GetString("classify")
	//获取通过iris.WithPostMaxMemory获取的最大上传值大小。
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	var dir, url string
	form := ctx.Request().MultipartForm
	failures := 0
	var urls []string
	for _, file := range form.File {
		if dir == "" {
			dir, url, err = GetDirAndUrl(classify, file[0])
		}

		upInfo, err := SaveUploadedFile(file[0], dir, url)
		if err != nil {
			failures++
			common.Response(ctx, nil, file[0].Filename+"上传失败", e.ERROR)
		} else {
			upInfo.File.Size = uint64(file[0].Size)
			upInfo.UploadUserID = userID
			if err := initialize.DB.Create(&upInfo).Error; err != nil {
				common.Response(ctx, nil, err.Error(), e.ERROR)
			}
			urls = append(urls, upInfo.URL)
		}
	}

	common.Res(ctx, iris.Map{
		"errno": 0,
		"data":  urls,
	})
}
