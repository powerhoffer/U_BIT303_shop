package backend

import "github.com/gogf/gf/v2/frame/g"

type GoodsCreateReq struct {
	g.Meta      `path:"/goods/create" method:"post" tags:"后台商品" summary:"新增商品"`
	CategoryId  uint   `json:"category_id" v:"required|min:1#商品分类不能为空|商品分类错误"`
	Name        string `json:"name" v:"required|length:1,128#商品名称不能为空|商品名称最多128位"`
	ImageUrl    string `json:"image_url" v:"max-length:255#商品图片最多255位"`
	PointsPrice uint   `json:"points_price" v:"required|min:1#兑换积分不能为空|兑换积分必须大于0"`
	Stock       uint   `json:"stock"`
	Description string `json:"description"`
}

type GoodsCreateRes struct {
	Id uint `json:"id"`
}

type GoodsListReq struct {
	g.Meta     `path:"/goods/list" method:"get" tags:"后台商品" summary:"商品列表"`
	Page       int    `json:"page" in:"query" d:"1" v:"min:1#分页页码错误"`
	Size       int    `json:"size" in:"query" d:"10" v:"max:50#分页数量最多50条"`
	CategoryId uint   `json:"category_id" in:"query"`
	Name       string `json:"name" in:"query"`
	Status     int    `json:"status" in:"query" d:"-1"`
}

type GoodsListRes struct {
	List  []GoodsItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type GoodsDetailReq struct {
	g.Meta `path:"/goods/detail" method:"get" tags:"后台商品" summary:"商品详情"`
	Id     uint `json:"id" in:"query" v:"required|min:1#商品ID不能为空|商品ID错误"`
}

type GoodsDetailRes struct {
	Goods GoodsItem `json:"goods"`
}

type GoodsUpdateReq struct {
	g.Meta      `path:"/goods/update" method:"post" tags:"后台商品" summary:"编辑商品"`
	Id          uint   `json:"id" v:"required|min:1#商品ID不能为空|商品ID错误"`
	CategoryId  uint   `json:"category_id" v:"required|min:1#商品分类不能为空|商品分类错误"`
	Name        string `json:"name" v:"required|length:1,128#商品名称不能为空|商品名称最多128位"`
	ImageUrl    string `json:"image_url" v:"max-length:255#商品图片最多255位"`
	PointsPrice uint   `json:"points_price" v:"required|min:1#兑换积分不能为空|兑换积分必须大于0"`
	Stock       uint   `json:"stock"`
	Description string `json:"description"`
}

type GoodsUpdateRes struct {
	Goods GoodsItem `json:"goods"`
}

type GoodsStatusReq struct {
	g.Meta `path:"/goods/status" method:"post" tags:"后台商品" summary:"上架或下架商品"`
	Id     uint   `json:"id" v:"required|min:1#商品ID不能为空|商品ID错误"`
	Status string `json:"status" v:"required|in:0,1#状态不能为空|状态只能是0或1"`
}

type GoodsStatusRes struct {
	Message string `json:"message"`
}

type GoodsItem struct {
	Id           uint   `json:"id"`
	CategoryId   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	PointsPrice  uint   `json:"points_price"`
	Stock        uint   `json:"stock"`
	Description  string `json:"description"`
	Status       int    `json:"status"`
}
