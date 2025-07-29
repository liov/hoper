package service

import (
	"context"
	"fmt"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/datastructure/idgen/id"
	"github.com/hopeio/utils/datax/database/datatypes"
	"github.com/hopeio/utils/log"
	"github.com/hopeio/utils/os/fs"
	timei "github.com/hopeio/utils/time"
	"github.com/liov/hoper/server/go/file/global"
	"github.com/liov/hoper/server/go/file/model"
	"github.com/tus/tusd/v2/pkg/handler"
	"gorm.io/gorm"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var defaultFilePerm = os.FileMode(0664)
var defaultDirectoryPerm = os.FileMode(0754)

// See the handler.DataStore interface for documentation about the different
// methods.
type FileStore struct {
	// Relative or absolute path to store files in. FileStore does not check
	// whether the path exists, use os.MkdirAll in this case on your own.
	Dir string
}

// NewFileStore creates a new file based storage backend. The directory specified will
// be used as the only storage entry. This method does not check
// whether the path exists, use os.MkdirAll to ensure.
// In addition, a locking mechanism is provided.
func NewFileStore(dir string) FileStore {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}
	err = os.MkdirAll(absDir, defaultDirectoryPerm)
	if err != nil {
		log.Fatal(err)
	}
	return FileStore{absDir}
}

// UseIn sets this store as the core data store in the passed composer and adds
// all possible extension to it.
func (store FileStore) UseIn(composer *handler.StoreComposer) {
	composer.UseCore(store)
	composer.UseTerminater(store)
	composer.UseConcater(store)
	composer.UseLengthDeferrer(store)
	composer.UseContentServer(store)
}

func (store FileStore) NewUpload(ctx context.Context, info handler.FileInfo) (handler.Upload, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	authInfo, _ := auth(ctxi, false)
	if info.ID == "" {
		info.ID = id.UniqueID()
	}
	name := info.MetaData["filename"]
	if name == "" {
		name = info.ID
	}
	md5 := info.MetaData["md5"]
	path := filepath.Join(timei.GetYMD(time.Now(), sep), info.ID+"_"+name)
	completePath := filepath.Join(store.Dir, path)

	// Create binary file with no content
	if err := createFile(completePath, nil); err != nil {
		return nil, err
	}
	upload := &fileUpload{
		FileInfo: model.FileInfo{
			Id:             info.ID,
			Name:           name,
			MD5:            md5,
			Path:           path,
			Size:           info.Size,
			SizeIsDeferred: info.SizeIsDeferred,
			Offset:         info.Offset,
			MetaData:       datatypes.MapJson[string](info.MetaData),
			IsPartial:      info.IsPartial,
			IsFinal:        info.IsFinal,
			PartialUploads: info.PartialUploads,
			Storage:        info.Storage,
		},
		CompletePath: completePath,
	}

	if authInfo != nil {
		err := global.Dao.GORMDB.Create(&model.UploadInfo{
			FileId: info.ID,
			UserId: authInfo.Id,
		}).Error
		if err != nil {
			return nil, err
		}
	}

	// writeInfo creates the file by itself if necessary
	err := global.Dao.GORMDB.Table(model.TableNameFileInfo).Create(&upload.FileInfo).Error
	if err != nil {
		return nil, err
	}

	return upload, nil
}

func (store FileStore) GetUpload(ctx context.Context, id string) (handler.Upload, error) {
	info := model.FileInfo{}
	err := global.Dao.GORMDB.Table(model.TableNameFileInfo).Where("id = ?", id).Scan(&info).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Interpret os.ErrNotExist as 404 Not Found
			err = handler.ErrNotFound
		}
		return nil, err
	}
	completePath := filepath.Join(store.Dir, info.Path)
	stat, err := os.Stat(completePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = handler.ErrNotFound
		}
		return nil, err
	}

	info.Offset = stat.Size()

	return &fileUpload{
		FileInfo:     info,
		CompletePath: completePath,
	}, nil
}

func (store FileStore) AsTerminatableUpload(upload handler.Upload) handler.TerminatableUpload {
	return upload.(*fileUpload)
}

func (store FileStore) AsLengthDeclarableUpload(upload handler.Upload) handler.LengthDeclarableUpload {
	return upload.(*fileUpload)
}

func (store FileStore) AsConcatableUpload(upload handler.Upload) handler.ConcatableUpload {
	return upload.(*fileUpload)
}

func (store FileStore) AsServableUpload(upload handler.Upload) handler.ServableUpload {
	return upload.(*fileUpload)
}

var _ handler.Upload = &fileUpload{}

type fileUpload struct {
	model.FileInfo
	md5Hash      hash.Hash
	CompletePath string
}

func (upload *fileUpload) GetInfo(ctx context.Context) (handler.FileInfo, error) {
	return handler.FileInfo{
		ID:             upload.Id,
		Size:           upload.Size,
		SizeIsDeferred: upload.SizeIsDeferred,
		Offset:         upload.Offset,
		MetaData:       handler.MetaData(upload.MetaData),
		IsPartial:      upload.IsPartial,
		IsFinal:        upload.IsFinal,
		PartialUploads: upload.PartialUploads,
		Storage:        upload.Storage,
	}, nil
}

func (upload *fileUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	file, err := os.OpenFile(upload.CompletePath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err := io.Copy(file, src)

	upload.Offset += n
	return n, err
}

func (upload *fileUpload) GetReader(ctx context.Context) (io.ReadCloser, error) {
	return os.Open(upload.CompletePath)
}

func (upload *fileUpload) Terminate(ctx context.Context) error {
	if err := global.Dao.GORMDB.Table(model.TableNameFileInfo).UpdateColumn("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	if err := os.Remove(upload.CompletePath); err != nil {
		return err
	}
	return nil
}

func (upload *fileUpload) ConcatUploads(ctx context.Context, uploads []handler.Upload) (err error) {
	file, err := os.OpenFile(upload.CompletePath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, partialUpload := range uploads {
		fileUpload := partialUpload.(*fileUpload)

		src, err := os.Open(fileUpload.CompletePath)
		if err != nil {
			return err
		}

		if _, err = io.Copy(file, src); err != nil {
			return err
		}

		if err = src.Close(); err != nil {
			return err
		}
	}

	return
}

func (upload *fileUpload) DeclareLength(ctx context.Context, length int64) error {
	upload.Size = length
	upload.SizeIsDeferred = false
	return global.Dao.GORMDB.Table(model.TableNameFileInfo).Where("id = ?", upload.Id).Updates(map[string]any{
		"size": length, "size_is_deferred": false,
	}).Error
}

// writeInfo updates the entire information. Everything will be overwritten.

func (upload *fileUpload) FinishUpload(ctx context.Context) error {
	if upload.MD5 == "" {
		var err error
		upload.MD5, err = fs.Md5(upload.CompletePath)
		if err != nil {
			return err
		}
	}
	return global.Dao.GORMDB.Table(model.TableNameFileInfo).Where("id = ?", upload.Id).Updates(map[string]any{
		"is_final": true, "finished_at": time.Now(), "md5": upload.MD5,
	}).Error
}

func (upload *fileUpload) ServeContent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, upload.CompletePath)

	return nil
}

// createFile creates the file with the content. If the corresponding directory does not exist,
// it is created. If the file already exists, its content is removed.
func createFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, defaultFilePerm)
	if err != nil {
		if os.IsNotExist(err) {
			// An upload ID containing slashes is mapped onto different directories on disk,
			// for example, `myproject/uploadA` should be put into a folder called `myproject`.
			// If we get an error indicating that a directory is missing, we try to create it.
			if err := os.MkdirAll(filepath.Dir(path), defaultDirectoryPerm); err != nil {
				return fmt.Errorf("failed to create directory for %s: %s", path, err)
			}

			// Try creating the file again.
			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, defaultFilePerm)
			if err != nil {
				// If that still doesn't work, error out.
				return err
			}
		} else {
			return err
		}
	}

	if content != nil {
		if _, err := file.Write(content); err != nil {
			return err
		}
	}

	return file.Close()
}
