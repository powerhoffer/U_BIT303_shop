package model

import "time"

type StockAdjustInput struct {
	GoodsId         uint
	Action          string
	Quantity        uint
	Remark          string
	OperatorAdminId uint
}

type StockAdjustOutput struct {
	GoodsId uint
	Stock   uint
}

type StockRecordInput struct {
	GoodsId        uint
	GoodsName      string
	ChangeType     int
	ChangeQuantity int
	BeforeStock    uint
	AfterStock     uint
	BizType        string
	BizId          uint
	OperatorType   int
	OperatorId     uint
	Remark         string
}

type StockRecordsInput struct {
	Page       int
	Size       int
	GoodsId    uint
	ChangeType int
	StartTime  string
	EndTime    string
}

type StockRecordsOutput struct {
	List  []StockRecordItem
	Total int
	Page  int
	Size  int
}

type StockRecordItem struct {
	Id             uint
	GoodsId        uint
	GoodsName      string
	ChangeType     int
	ChangeQuantity int
	BeforeStock    uint
	AfterStock     uint
	BizType        string
	BizId          uint
	OperatorType   int
	OperatorId     uint
	Remark         string
	CreatedAt      time.Time
}
