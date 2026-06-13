package frontend

import "github.com/gogf/gf/v2/frame/g"

type GoodsListReq struct {
	g.Meta     `path:"/goods/list" method:"get" tags:"前台商品" summary:"前台商品列表"`
	Page       int    `json:"page" in:"query" d:"1" v:"min:1#分页页码错误"`
	Size       int    `json:"size" in:"query" d:"10" v:"max:50#分页数量最多50条"`
	CategoryId uint   `json:"category_id" in:"query"`
	Name       string `json:"name" in:"query"`
}

type GoodsListRes struct {
	List  []GoodsListItem `json:"list"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}

type GoodsDetailReq struct {
	g.Meta `path:"/goods/detail" method:"get" tags:"前台商品" summary:"前台商品详情"`
	Id     uint `json:"id" in:"query" v:"required|min:1#商品ID不能为空|商品ID错误"`
}

type GoodsDetailRes struct {
	Goods GoodsDetailItem `json:"goods"`
}

type GoodsListItem struct {
	Id           uint   `json:"id"`
	CategoryId   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	PointsPrice  uint   `json:"points_price"`
	Stock        uint   `json:"stock"`
}

type GoodsDetailItem struct {
	Id           uint   `json:"id"`
	CategoryId   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	PointsPrice  uint   `json:"points_price"`
	Stock        uint   `json:"stock"`
	Description  string `json:"description"`
}
