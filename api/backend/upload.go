package backend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UploadGoodsImageReq struct {
	g.Meta `path:"/upload/goods-image" method:"post" mime:"multipart/form-data" tags:"Backend Upload" summary:"Upload goods image"`
	File   *ghttp.UploadFile `json:"file" type:"file" v:"required#Image file is required"`
}

type UploadGoodsImageRes struct {
	Url          string `json:"url"`
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FileSize     uint64 `json:"file_size"`
}
