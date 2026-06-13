package model

type CartAddInput struct {
	EmployeeId uint
	GoodsId    uint
	Count      int
}

type CartAddOutput struct {
	Id uint
}

type CartListInput struct {
	EmployeeId uint
	Page       int
	Size       int
}

type CartListOutput struct {
	List  []CartItem
	Total int
	Page  int
	Size  int
}

type CartUpdateInput struct {
	EmployeeId uint
	Id         uint
	Count      int
}

type CartUpdateOutput struct {
	Id uint
}

type CartRemoveInput struct {
	EmployeeId uint
	Id         uint
}

type CartRemoveOutput struct {
	Id uint
}

type CartItem struct {
	Id           uint
	GoodsId      uint
	CategoryId   uint
	CategoryName string
	GoodsName    string
	ImageUrl     string
	PointsPrice  uint
	Stock        uint
	Count        uint
	TotalPoints  uint
}
