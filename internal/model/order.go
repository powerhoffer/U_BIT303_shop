package model

import "time"

type OrderCreateInput struct {
	EmployeeId uint
	Remark     string
}

type OrderCreateOutput struct {
	Order OrderBase
}

type OrderListInput struct {
	EmployeeId uint
	Page       int
	Size       int
}

type OrderListOutput struct {
	List  []OrderBase
	Total int
	Page  int
	Size  int
}

type OrderDetailInput struct {
	EmployeeId uint
	Id         uint
}

type OrderDetailOutput struct {
	Order OrderDetail
}

type OrderCancelInput struct {
	EmployeeId uint
	Id         uint
}

type OrderCancelOutput struct {
	Order OrderBase
}

type OrderBase struct {
	Id          uint
	OrderNo     string
	EmployeeId  uint
	TotalPoints uint
	Status      int
	Remark      string
	CreatedAt   time.Time
}

type OrderDetail struct {
	OrderBase
	Items []OrderGoodsItem
}

type OrderGoodsItem struct {
	Id            uint
	OrderId       uint
	EmployeeId    uint
	GoodsId       uint
	GoodsName     string
	GoodsImageUrl string
	PointsPrice   uint
	Count         uint
	TotalPoints   uint
	CreatedAt     time.Time
}
