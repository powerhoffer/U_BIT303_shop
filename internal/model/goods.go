package model

type GoodsCreateInput struct {
	CategoryId  uint
	Name        string
	ImageUrl    string
	PointsPrice uint
	Stock       uint
	Description string
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
	Id          uint
	CategoryId  uint
	Name        string
	ImageUrl    string
	PointsPrice uint
	Stock       uint
	Description string
}

type GoodsUpdateOutput struct {
	Goods GoodsItem
}

type GoodsStatusInput struct {
	Id     uint
	Status int
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
