package backend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type StockAdjustReq struct {
	g.Meta   `path:"/stock/adjust" method:"post" tags:"Backend Stock" summary:"Adjust goods stock"`
	GoodsId  uint   `json:"goods_id" v:"required|min:1#Goods ID is required|Goods ID is invalid"`
	Action   string `json:"action" v:"required|in:increase,decrease#Action is required|Action must be increase or decrease"`
	Quantity uint   `json:"quantity" v:"required|min:1#Quantity is required|Quantity must be greater than 0"`
	Remark   string `json:"remark" v:"max-length:255#Remark must be at most 255 characters"`
}

type StockAdjustRes struct {
	GoodsId uint `json:"goods_id"`
	Stock   uint `json:"stock"`
}

type StockRecordsReq struct {
	g.Meta     `path:"/stock/records" method:"get" tags:"Backend Stock" summary:"Goods stock records"`
	Page       int    `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size       int    `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
	GoodsId    uint   `json:"goods_id" in:"query"`
	ChangeType int    `json:"change_type" in:"query" v:"in:0,1,2,3,4,5#Change type is invalid"`
	StartTime  string `json:"start_time" in:"query" description:"Start date in YYYY-MM-DD format"`
	EndTime    string `json:"end_time" in:"query" description:"End date in YYYY-MM-DD format"`
}

type StockRecordsRes struct {
	List  []StockRecordItem `json:"list"`
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

type StockRecordItem struct {
	Id             uint      `json:"id"`
	GoodsId        uint      `json:"goods_id"`
	GoodsName      string    `json:"goods_name"`
	ChangeType     int       `json:"change_type"`
	ChangeQuantity int       `json:"change_quantity"`
	BeforeStock    uint      `json:"before_stock"`
	AfterStock     uint      `json:"after_stock"`
	BizType        string    `json:"biz_type"`
	BizId          uint      `json:"biz_id"`
	OperatorType   int       `json:"operator_type"`
	OperatorId     uint      `json:"operator_id"`
	Remark         string    `json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}
