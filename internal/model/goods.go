package model

type GoodsCreateInput struct {
	OperatorAdminId uint
	CategoryId      uint
	Name            string
	ImageUrl        string
	PointsPrice     uint
	Stock           uint
	Description     string
}

type GoodsCreateOutput struct {
	Id uint
}

type GoodsListInput struct {
	Page       int
	Size       int
	CategoryId uint
	Name       string
	Status     int
}

type GoodsListOutput struct {
	List  []GoodsItem
	Total int
	Page  int
	Size  int
}

type GoodsDetailOutput struct {
	Goods GoodsItem
}

type GoodsUpdateInput struct {
	OperatorAdminId uint
	Id              uint
	CategoryId      uint
	Name            string
	ImageUrl        string
	PointsPrice     uint
	Stock           uint
	Description     string
}

type GoodsUpdateOutput struct {
	Goods GoodsItem
}

type GoodsStatusInput struct {
	Id     uint
	Status int
}

type FrontendGoodsListInput struct {
	Page       int
	Size       int
	CategoryId uint
	Name       string
}

type FrontendGoodsListOutput struct {
	List  []FrontendGoodsListItem
	Total int
	Page  int
	Size  int
}

type FrontendGoodsDetailOutput struct {
	Goods FrontendGoodsDetailItem
}

type GoodsItem struct {
	Id           uint
	CategoryId   uint
	CategoryName string
	Name         string
	ImageUrl     string
	PointsPrice  uint
	Stock        uint
	Description  string
	Status       int
}

type FrontendGoodsListItem struct {
	Id           uint
	CategoryId   uint
	CategoryName string
	Name         string
	ImageUrl     string
	PointsPrice  uint
	Stock        uint
}

type FrontendGoodsDetailItem struct {
	Id           uint
	CategoryId   uint
	CategoryName string
	Name         string
	ImageUrl     string
	PointsPrice  uint
	Stock        uint
	Description  string
}
