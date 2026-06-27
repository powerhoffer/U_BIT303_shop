package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

var Upload = cUpload{}

type cUpload struct{}

func (c *cUpload) GoodsImage(ctx context.Context, req *backendApi.UploadGoodsImageReq) (res *backendApi.UploadGoodsImageRes, err error) {
	out, err := service.Upload().GoodsImage(ctx, model.UploadGoodsImageInput{
		AdminId: currentAdminId(ctx),
		File:    req.File,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.UploadGoodsImageRes{
		Url:          out.Url,
		FileName:     out.FileName,
		OriginalName: out.OriginalName,
		FileSize:     out.FileSize,
	}, nil
}
