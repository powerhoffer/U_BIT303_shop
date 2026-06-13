package frontend

import "github.com/gogf/gf/v2/frame/g"

type CartAddReq struct {
	g.Meta  `path:"/cart/add" method:"post" tags:"Frontend Cart" summary:"Add to cart"`
	GoodsId uint `json:"goods_id" v:"required|min:1#Goods ID is required|Goods ID is invalid"`
	Count   int  `json:"count" v:"required|min:1#Goods count is required|Goods count must be greater than 0"`
}

type CartAddRes struct {
	Id uint `json:"id"`
}

type CartListReq struct {
	g.Meta `path:"/cart/list" method:"get" tags:"Frontend Cart" summary:"Cart list"`
	Page   int `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size   int `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
}

type CartListRes struct {
	List  []CartItem `json:"list"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type CartUpdateReq struct {
	g.Meta `path:"/cart/update" method:"post" tags:"Frontend Cart" summary:"Update cart item count"`
	Id     uint `json:"id" v:"required|min:1#Cart ID is required|Cart ID is invalid"`
	Count  int  `json:"count" v:"required|min:0#Goods count is required|Goods count cannot be negative"`
}

type CartUpdateRes struct {
	Id uint `json:"id"`
}

type CartRemoveReq struct {
	g.Meta `path:"/cart/remove" method:"post" tags:"Frontend Cart" summary:"Remove cart item"`
	Id     uint `json:"id" v:"required|min:1#Cart ID is required|Cart ID is invalid"`
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
