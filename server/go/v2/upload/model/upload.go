package model

import (
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type FileUploadInfo struct {
	File
	UUID      string    `gorm:"type:varchar(100);unique;not null" json:"uuid"`
	UserId    uint64    `json:"upload_user_id"`
	CreatedAt time.Time `json:"created_at"`
	FilePath  string    `gorm:"type:varchar(100);not null" json:"upload_file_path"`
	Status    uint8     `gorm:"type:smallint;default:0" json:"status"`
}

type File struct {
	Id           uint64 `gorm:"primary_key" json:"id"`
	FileName     string `gorm:"type:varchar(100);not null" json:"file_name"`
	OriginalName string `gorm:"type:varchar(100);not null" json:"original_name"`
	URL          string `json:"url"`
	MD5          string `gorm:"type:varchar(32)" json:"md5"`
	Mime         string `json:"mime"`
	Size         uint64 `json:"size"`
}

func GenerateUploadedInfo(ext string) FileUploadInfo {

	sep := string(os.PathSeparator)
	uploadImgDir := initialize.Config.Server.UploadDir
	length := utf8.RuneCountInString(uploadImgDir)
	lastChar := uploadImgDir[length-1:]
	ymStr := utils.GetTodayYM(sep)

	var uploadDir string
	if lastChar != sep {
		uploadDir = uploadImgDir + sep + ymStr
	} else {
		uploadDir = uploadImgDir + ymStr
	}

	uuidName := uuid.NewV4().String()
	filename := uuidName + ext
	uploadFilePath := uploadDir + sep + filename
	fileURL := strings.Join([]string{
		initialize.Config.Server.UploadPath,
		ymStr,
		filename,
	}, "/")
	var fileUpload crm.FileUploadInfo

	fileUpload.FileName = filename
	fileUpload.File.URL = fileURL
	fileUpload.UUID = uuidName
	fileUpload.UploadFilePath = uploadFilePath

	/*	fileUpload = crm.FileUploadInfo{
		File:       model.File{FileName:filename,},
		FileURL:        fileURL,
		UUIDName:       uuidName,
		UploadDir:      uploadDir,
		UploadFilePath: uploadFilePath,
	}*/
	return fileUpload
}

func GetExt(file *multipart.FileHeader) (string, error) {
	var ext string
	var index = strings.LastIndex(file.Filename, ".")
	if index == -1 {
		return "", nil
	} else {
		ext = file.Filename[index:]
	}
	if len(ext) == 1 {
		return "", errors.New("无效的扩展名")
	}
	return ext, nil
}

func GetDirAndUrl(classify string, info *multipart.FileHeader) (string, string, error) {
	//sep := string(os.PathSeparator)
	sep := "/"
	var uploadDir, prefixUrl string
	ymdStr := utils.GetTodayYMD(sep)
	ext, err := GetExt(info)
	if err != nil {
		return "", "", err
	}

	if ext == "" {
		uploadDir = strings.Join([]string{initialize.Config.Server.UploadDir,
			"others",
			classify,
			ymdStr},
			"/")
		prefixUrl = strings.Join([]string{
			initialize.Config.Server.UploadPath,
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



func UploadMultiple(ctx iris.Context) {
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

func SaveUploadedFile(file *multipart.FileHeader, dir string, url string) (*FileUploadInfo, error) {
	uuidName := uuid.NewV4().String()
	ext, err := GetExt(file)
	filename := uuidName + ext
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(dir + filename)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	fileUpload := crm.FileUploadInfo{
		File: crm.File{
			FileName:     filename,
			OriginalName: file.Filename,
			URL:          url + filename,
			Mime:         mime.TypeByExtension(ext),
		},
		UUID:           uuidName,
		UploadFilePath: dir + filename,
		UploadAt:       time.Now(),
	}
	io.Copy(out, src)
	return &fileUpload, nil
}

func MD5(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint64)
	md5 := ctx.Params().Get("md5")
	var upI crm.FileUploadInfo
	var count int
	initialize.DB.Where("md5 = ?", md5).First(&upI).Count(&count)
	if count != 0 {
		upI.ID = 0
		upI.UploadUserID = userID
		upI.UUID = uuid.NewV4().String()
		upI.UploadAt = time.Now()
		if err := initialize.DB.Create(&upI).Error; err != nil {
			common.Response(ctx, err, "", e.ERROR)
		}
		common.Response(ctx, &upI, "", e.SUCCESS)
		return
	}
	common.Response(ctx, nil, "不存在", e.ERROR)
}
