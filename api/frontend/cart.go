package frontend

import "github.com/gogf/gf/v2/frame/g"

type CartAddReq struct {
	g.Meta  `path:"/cart/add" method:"post" tags:"前台购物车" summary:"加入购物车"`
	GoodsId uint `json:"goods_id" v:"required|min:1#商品ID不能为空|商品ID错误"`
	Count   int  `json:"count" v:"required|min:1#商品数量不能为空|商品数量必须大于0"`
}

type CartAddRes struct {
	Id uint `json:"id"`
}

type CartListReq struct {
	g.Meta `path:"/cart/list" method:"get" tags:"前台购物车" summary:"购物车列表"`
	Page   int `json:"page" in:"query" d:"1" v:"min:1#分页页码错误"`
	Size   int `json:"size" in:"query" d:"10" v:"max:50#分页数量最多50条"`
}

type CartListRes struct {
	List  []CartItem `json:"list"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type CartUpdateReq struct {
	g.Meta `path:"/cart/update" method:"post" tags:"前台购物车" summary:"更新购物车商品数量"`
	Id     uint `json:"id" v:"required|min:1#购物车ID不能为空|购物车ID错误"`
	Count  int  `json:"count" v:"required|min:0#商品数量不能为空|商品数量不能小于0"`
}

type CartUpdateRes struct {
	Id uint `json:"id"`
}

type CartRemoveReq struct {
	g.Meta `path:"/cart/remove" method:"post" tags:"前台购物车" summary:"移除购物车商品"`
	Id     uint `json:"id" v:"required|min:1#购物车ID不能为空|购物车ID错误"`
}

type CartRemoveRes struct {
	Id uint `json:"id"`
}

type CartItem struct {
	Id           uint   `json:"id"`
	GoodsId      uint   `json:"goods_id"`
	CategoryId   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	GoodsName    string `json:"goods_name"`
	ImageUrl     string `json:"image_url"`
	PointsPrice  uint   `json:"points_price"`
	Stock        uint   `json:"stock"`
	Count        uint   `json:"count"`
	TotalPoints  uint   `json:"total_points"`
}
