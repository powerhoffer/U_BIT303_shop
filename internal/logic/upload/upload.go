package upload

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

const goodsBizType = "goods"

type sUpload struct{}

func init() {
	service.RegisterUpload(New())
}

func New() *sUpload {
	return &sUpload{}
}

func (s *sUpload) GoodsImage(ctx context.Context, in model.UploadGoodsImageInput) (out model.UploadGoodsImageOutput, err error) {
	if in.File == nil {
		return out, errors.New("Image file is required")
	}

	ext := strings.ToLower(filepath.Ext(in.File.Filename))
	if !isAllowedImageExt(ext) {
		return out, errors.New("Only jpg, jpeg, png and webp images are allowed")
	}
	mimeType := in.File.Header.Get("Content-Type")
	if mimeType != "" && !strings.HasPrefix(strings.ToLower(mimeType), "image/") {
		return out, errors.New("Only image files are allowed")
	}

	maxSizeMB := g.Cfg().MustGet(ctx, "upload.maxSizeMB").Int64()
	if maxSizeMB <= 0 {
		maxSizeMB = 5
	}
	maxSize := maxSizeMB * 1024 * 1024
	if in.File.Size > maxSize {
		return out, errors.New("Image file must be at most 5MB")
	}

	uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
	urlPrefix := strings.TrimRight(g.Cfg().MustGet(ctx, "upload.urlPrefix").String(), "/")
	if uploadPath == "" || urlPrefix == "" {
		return out, errors.New("Upload configuration is missing")
	}

	dateDir := gtime.Now().Format("Ymd")
	saveDir := gfile.Join(uploadPath, goodsBizType, dateDir)
	fileName, err := in.File.Save(saveDir, true)
	if err != nil {
		return out, err
	}
	filePath := gfile.Join(saveDir, fileName)
	url := urlPrefix + "/" + goodsBizType + "/" + dateDir + "/" + fileName

	_, err = dao.UploadFile.Ctx(ctx).Data(do.UploadFile{
		AdminId:      in.AdminId,
		FileName:     fileName,
		OriginalName: in.File.Filename,
		FilePath:     filePath,
		Url:          url,
		FileSize:     uint64(in.File.Size),
		MimeType:     mimeType,
		FileExt:      strings.TrimPrefix(ext, "."),
		BizType:      goodsBizType,
	}).Insert()
	if err != nil {
		return out, err
	}

	out = model.UploadGoodsImageOutput{
		Url:          url,
		FileName:     fileName,
		OriginalName: in.File.Filename,
		FileSize:     uint64(in.File.Size),
	}
	return out, nil
}

func isAllowedImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}
