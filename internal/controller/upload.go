package controller

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"bit303_shop/api/backend"
	"bit303_shop/internal/consts"
	"bit303_shop/utility/upload"
)

type cUpload struct{}

var Upload = cUpload{}

func (c *cUpload) UploadImgToCloud(ctx context.Context, req *backend.UploadImgToCloudReq) (res *backend.UploadImgToCloudRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, consts.CodeMissingParameterMsg)
	}
	url, err := upload.UploadImgToCloud(ctx, req.File)
	if err != nil {
		g.Log(err.Error())
		return nil, err
	}
	return &backend.UploadImgToCloudRes{
		Url: url,
	}, nil
}
