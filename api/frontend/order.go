package frontend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type OrderCreateReq struct {
	g.Meta `path:"/order/create" method:"post" tags:"Frontend Order" summary:"Create redemption order from cart"`
	Remark string `json:"remark" v:"max-length:255#Remark must be at most 255 characters"`
}

type OrderCreateRes struct {
	Order OrderItem `json:"order"`
}

type OrderListReq struct {
	g.Meta `path:"/order/list" method:"get" tags:"Frontend Order" summary:"Current employee order list"`
	Page   int `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size   int `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
}

type OrderListRes struct {
	List  []OrderItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type OrderDetailReq struct {
	g.Meta `path:"/order/detail" method:"get" tags:"Frontend Order" summary:"Current employee order detail"`
	Id     uint `json:"id" in:"query" v:"required|min:1#Order ID is required|Order ID is invalid"`
}

type OrderDetailRes struct {
	Order OrderDetail `json:"order"`
}

type OrderCancelReq struct {
	g.Meta `path:"/order/cancel" method:"post" tags:"Frontend Order" summary:"Cancel current employee order"`
	Id     uint `json:"id" v:"required|min:1#Order ID is required|Order ID is invalid"`
}

type OrderCancelRes struct {
	Order OrderItem `json:"order"`
}

type OrderItem struct {
	Id          uint      `json:"id"`
	OrderNo     string    `json:"order_no"`
	EmployeeId  uint      `json:"employee_id"`
	TotalPoints uint      `json:"total_points"`
	Status      int       `json:"status"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderDetail struct {
	OrderItem
	Items []OrderGoodsItem `json:"items"`
}

type OrderGoodsItem struct {
	Id            uint      `json:"id"`
	OrderId       uint      `json:"order_id"`
	GoodsId       uint      `json:"goods_id"`
	GoodsName     string    `json:"goods_name"`
	GoodsImageUrl string    `json:"goods_image_url"`
	PointsPrice   uint      `json:"points_price"`
	Count         uint      `json:"count"`
	TotalPoints   uint      `json:"total_points"`
	CreatedAt     time.Time `json:"created_at"`
}
