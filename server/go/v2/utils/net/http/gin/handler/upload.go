package handler

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/structure/uuid"
	"github.com/liov/hoper/go/v2/utils/time"

)

type UploadFile interface {
	Exist()
	Create()
}

type FileUploadInfo struct {
	fs.File
	UUID           string    `gorm:"type:varchar(100);unique;not null" json:"uuid"`
	UploadUserID   uint64    `json:"upload_user_id"`
	UploadAt       time.Time `json:"upload_at"`
	UploadFilePath string    `gorm:"type:varchar(100);not null" json:"upload_file_path"`
	Status         uint8     `gorm:"type:smallint;default:0" json:"status"`
}

func PathMerge(path1, path2 string) string {
	sep := string(os.PathSeparator)
	length := utf8.RuneCountInString(path1)
	lastChar := path1[length-1:]
	if lastChar != sep {
		path1 = path1 + sep + path2
	} else {
		path1 = path1 + path2
	}
	return path1
}

func GenerateUploadedInfo(uploadDir, ext string) FileUploadInfo {

	ymStr := timei.GetTodayYM(string(os.PathSeparator))

	uuidName := uuid.NewV4().String()
	filename := uuidName + ext
	uploadFilePath := PathMerge(uploadDir, filename)
	fileURL := strings.Join([]string{
		uploadDir, ymStr, filename}, "/")

	var fileUpload FileUploadInfo

	fileUpload.FileName = filename
	fileUpload.File.URL = fileURL
	fileUpload.UUID = uuidName
	fileUpload.UploadFilePath = uploadFilePath

	/*	fileUpload = FileUploadInfo{
		File:       File{FileName:filename,},
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

func GetDirAndUrl(uploadDir string, info *multipart.FileHeader) (string, error) {
	//sep := string(os.PathSeparator)
	sep := "/"
	ymdStr := timei.GetTodayYMD(sep)
	ext, err := GetExt(info)
	if err != nil {
		return "", err
	}

	if ext == "" {
		uploadDir = strings.Join([]string{
			uploadDir, ymdStr}, "/")
		return uploadDir, nil
	}

	var mimeType = mime.TypeByExtension(ext)
	if mimeType == "" && ext == ".jpeg" {
		mimeType = "image/jpeg"
	}

	dirType := strings.Split(mimeType, "/")

	uploadDir = strings.Join([]string{uploadDir,
		dirType[0] + "s",
		ymdStr},
		"/")

	if err := os.MkdirAll(uploadDir, 0777); err != nil {
		return uploadDir, err
	}
	return uploadDir, nil
}

// Upload 文件上传
func Upload(ctx *gin.Context) *FileUploadInfo {
	userID := ctx.Value("userID").(uint64)
	classify, _ := ctx.Params.Get("classify")
	info, err := ctx.FormFile("file")
	md5 := ctx.PostForm("md5")
	/*	var upI FileUploadInfo
		var count int
		initialize.DB.Where("md5 = ?", md5).First(&upI).Count(&count)
		if count != 0 {
			upI.ID = 0
			upI.UploadUserID = userID
			upI.UUID = uuid.NewV4().String()
			upI.UploadAt = time.Now()
			if err := initialize.DB.Create(&upI).Error; err != nil {
				common.Response(ctx, err.Error())
				return nil
			}
			common.Response(ctx, &upI)
			return &upI
		}*/
	/*	md5 := md5.New()
		io.Copy(md5,file)
		MD5Str := hex.EncodeToString(md5.Sum(nil))*/

	if err != nil {
		Response(ctx, nil, "参数无效", errorcode.Canceled)
		return nil
	}

	ext, err := GetExt(info)
	if ext == "" || err != nil {
		Response(ctx, nil, "无效的图片类型", errorcode.Canceled)
		return nil
	}

	dir, err := GetDirAndUrl(classify, info)

	upInfo, err := SaveUploadedFile(info, dir)
	if err != nil {
		Response(ctx, nil, err.Error(), errorcode.Canceled)
		return nil
	}

	upInfo.File.Size = uint64(info.Size)
	upInfo.UploadUserID = userID
	upInfo.Status = 1
	upInfo.MD5 = md5
	/*	if err := initialize.DB.Create(upInfo).Error; err != nil {
		Response(ctx,nil, err.Error(), errorcode.Canceled)
		return nil
	}*/
	Response(ctx, upInfo, "", errorcode.SUCCESS)
	return upInfo
}

func UploadMultiple(ctx *gin.Context) {
	userID := ctx.Value("userID").(uint64)
	classify, _ := ctx.Params.Get("classify")
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Writer.WriteString(err.Error())
		return
	}

	var dir string
	failures := 0
	var urls []string
	for _, file := range form.File {
		if dir == "" {
			dir, err = GetDirAndUrl(classify, file[0])
		}

		upInfo, err := SaveUploadedFile(file[0], dir)
		if err != nil {
			failures++
			Response(ctx, nil, file[0].Filename+"上传失败", errorcode.Canceled)
		} else {
			upInfo.File.Size = uint64(file[0].Size)
			upInfo.UploadUserID = userID
			/*			if err := initialize.DB.Create(&upInfo).Error; err != nil {
						Response(ctx, nil, err.Error(), errorcode.Canceled)
					}*/
			urls = append(urls, upInfo.URL)
		}
	}

	Res(ctx, errorcode.SUCCESS, "", urls)
}

func SaveUploadedFile(file *multipart.FileHeader, dir string) (*FileUploadInfo, error) {
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

	fileUpload := FileUploadInfo{
		File: fs.File{
			FileName:     filename,
			OriginalName: file.Filename,
			URL:          dir + filename,
			Mime:         mime.TypeByExtension(ext),
		},
		UUID:           uuidName,
		UploadFilePath: dir + filename,
		UploadAt:       time.Now(),
	}
	io.Copy(out, src)
	return &fileUpload, nil
}

func MD5(ctx *gin.Context) {
	userID := ctx.Value("userID").(uint64)
	//md5,_ := ctx.Params.Get("md5")
	var upI FileUploadInfo
	var count int
	//initialize.DB.Where("md5 = ?", md5).First(&upI).Count(&count)
	if count != 0 {
		upI.ID = 0
		upI.UploadUserID = userID
		upI.UUID = uuid.NewV4().String()
		upI.UploadAt = time.Now()
		/*		if err := initialize.DB.Create(&upI).Error; err != nil {
				Response(ctx, err, "", errorcode.Canceled)
			}*/
		Response(ctx, &upI, "", errorcode.SUCCESS)
		return
	}
	Response(ctx, nil, "不存在", errorcode.Canceled)
}
