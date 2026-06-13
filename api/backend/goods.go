package backend

import "github.com/gogf/gf/v2/frame/g"

type GoodsCreateReq struct {
	g.Meta      `path:"/goods/create" method:"post" tags:"Backend Goods" summary:"Create goods"`
	CategoryId  uint   `json:"category_id" v:"required|min:1#Category is required|Category is invalid"`
	Name        string `json:"name" v:"required|length:1,128#Goods name is required|Goods name must be at most 128 characters"`
	ImageUrl    string `json:"image_url" v:"max-length:255#Goods image URL must be at most 255 characters"`
	PointsPrice uint   `json:"points_price" v:"required|min:1#Credits price is required|Credits price must be greater than 0"`
	Stock       uint   `json:"stock"`
	Description string `json:"description"`
}

type GoodsCreateRes struct {
	Id uint `json:"id"`
}

type GoodsListReq struct {
	g.Meta     `path:"/goods/list" method:"get" tags:"Backend Goods" summary:"Goods list"`
	Page       int    `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size       int    `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
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
	g.Meta `path:"/goods/detail" method:"get" tags:"Backend Goods" summary:"Goods detail"`
	Id     uint `json:"id" in:"query" v:"required|min:1#Goods ID is required|Goods ID is invalid"`
}

type GoodsDetailRes struct {
	Goods GoodsItem `json:"goods"`
}

type GoodsUpdateReq struct {
	g.Meta      `path:"/goods/update" method:"post" tags:"Backend Goods" summary:"Update goods"`
	Id          uint   `json:"id" v:"required|min:1#Goods ID is required|Goods ID is invalid"`
	CategoryId  uint   `json:"category_id" v:"required|min:1#Category is required|Category is invalid"`
	Name        string `json:"name" v:"required|length:1,128#Goods name is required|Goods name must be at most 128 characters"`
	ImageUrl    string `json:"image_url" v:"max-length:255#Goods image URL must be at most 255 characters"`
	PointsPrice uint   `json:"points_price" v:"required|min:1#Credits price is required|Credits price must be greater than 0"`
	Stock       uint   `json:"stock"`
	Description string `json:"description"`
}

type GoodsUpdateRes struct {
	Goods GoodsItem `json:"goods"`
}

type GoodsStatusReq struct {
	g.Meta `path:"/goods/status" method:"post" tags:"Backend Goods" summary:"Put goods on or off shelf"`
	Id     uint   `json:"id" v:"required|min:1#Goods ID is required|Goods ID is invalid"`
	Status string `json:"status" v:"required|in:0,1#Status is required|Status must be 0 or 1"`
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
