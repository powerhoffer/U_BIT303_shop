package model

import "github.com/gogf/gf/v2/net/ghttp"

type UploadGoodsImageInput struct {
	AdminId uint
	File    *ghttp.UploadFile
}

type UploadGoodsImageOutput struct {
	Url          string
	FileName     string
	OriginalName string
	FileSize     uint64
}
