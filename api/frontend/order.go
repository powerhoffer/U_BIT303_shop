package frontend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type OrderAddReq struct {
	g.Meta `path:"/add/order" method:"post" tags:"前台订单" summary:"创建订单"`
	//主订单维度
	Price            int    `json:"price"            description:"订单金额 单位分"`
	CouponPrice      int    `json:"coupon_price"      description:"优惠券金额 单位分"`
	ActualPrice      int    `json:"actual_price"      description:"实际支付金额 单位分"`
	ConsigneeName    string `json:"consignee_name"    description:"收货人姓名"`
	ConsigneePhone   string `json:"consignee_phone"   description:"收货人手机号"`
	ConsigneeAddress string `json:"consignee_address" description:"收货人详细地址"`
	Remark           string `json:"remark"           description:"备注"`
	AddressId        uint   `json:"address_id"       description:"地址ID"`
	// 兼容不同调用方式
	OrderAddGoodsInfos []*OrderAddGoodsInfo `json:"order_add_goods_infos"`
	GoodsList          []OrderAddGoodsInfo  `json:"goods_list"        description:"商品列表"` // 用于测试代码
}

type OrderAddRes struct {
	Id uint `json:"id"`
}

type OrderAddGoodsInfo struct {
	GoodsId        int    `json:"goods_id"`
	GoodsOptionsId int    `json:"goods_options_id"`
	Count          int    `json:"count"`
	Remark         string `json:"remark"`
	Price          int    `json:"price"`
	CouponPrice    int    `json:"coupon_price"`
	ActualPrice    int    `json:"actual_price"`
}

type OrderGoodsInfo struct {
	Id             int `json:"id,omitempty"`
	OrderId        int `json:"order_id"`
	GoodsId        int `json:"goods_id"`
	GoodsOptionsId int `json:"goods_options_id"`
	//商品详情
	GoodsInfo *BaseGoodsColumns
	//注意：api层不需要做orm关联 关联了也没有意义
	//GoodsInfo   *BaseGoodsColumns `orm:"with:id=goods_id" json:"goods_info"`
	Count       int    `json:"count"`
	Remark      string `json:"remark"`
	Status      int    `json:"status"`
	Price       int    `json:"price"`
	CouponPrice int    `json:"coupon_price"`
	ActualPrice int    `json:"actual_price"`
}

type OrderListReq struct {
	g.Meta `path:"/order/list" method:"get" tags:"前台订单" summary:"获取订单列表"`
	CommonPaginationReq
	Status   int    `json:"status" description:"订单状态"`
	DateGte  string `json:"date_gte" description:"开始时间"`
	DateLte  string `json:"date_lte" description:"结束时间"`
	Page     int    `json:"page" description:"页码"`
	PageSize int    `json:"page_size" description:"每页数量"` // 用于测试代码
}

type OrderListRes struct {
	CommonPaginationRes
	List []OrderInfo `json:"list" description:"订单列表"`
}

type OrderDetailReq struct {
	g.Meta `path:"/order/detail" method:"get" tags:"前台订单" summary:"获取订单详情"`
	Id     uint `json:"id" v:"required" description:"订单ID"`
}

type OrderDetailRes struct {
	OrderInfo OrderInfo        `json:"order_info" description:"订单信息"`
	GoodsInfo []OrderGoodsInfo `json:"goods_info" description:"订单商品信息"`
}

type OrderCancelReq struct {
	g.Meta `path:"/order/cancel" method:"post" tags:"前台订单" summary:"取消订单"`
	Id     uint   `json:"id" v:"required" description:"订单ID"`
	Reason string `json:"reason" description:"取消原因"`
}

type OrderCancelRes struct {
	Id uint `json:"id" description:"订单ID"`
}

type OrderPayReq struct {
	g.Meta  `path:"/order/pay" method:"post" tags:"前台订单" summary:"订单支付"`
	Id      uint `json:"id" v:"required" description:"订单ID"`
	PayType int  `json:"pay_type" v:"required" description:"支付方式：1微信 2支付宝 3云闪付"`
}

type OrderPayRes struct {
	PayUrl string `json:"pay_url" description:"支付URL"`
}

type OrderInfo struct {
	Id               uint        `json:"id" description:"订单ID"`
	Number           string      `json:"number" description:"订单编号"`
	UserId           uint        `json:"user_id" description:"用户ID"`
	PayType          int         `json:"pay_type" description:"支付方式"`
	Remark           string      `json:"remark" description:"备注"`
	PayAt            *gtime.Time `json:"pay_at" description:"支付时间"`
	Status           int         `json:"status" description:"订单状态: 1待支付 2已支付待发货 3已发货 4已收货待评价 5已评价 6已取消"`
	ConsigneeName    string      `json:"consignee_name" description:"收货人姓名"`
	ConsigneePhone   string      `json:"consignee_phone" description:"收货人手机号"`
	ConsigneeAddress string      `json:"consignee_address" description:"收货人详细地址"`
	Price            int         `json:"price" description:"订单金额 单位分"`
	CouponPrice      int         `json:"coupon_price" description:"优惠券金额 单位分"`
	ActualPrice      int         `json:"actual_price" description:"实际支付金额 单位分"`
	CreatedAt        *gtime.Time `json:"created_at" description:"创建时间"`
	UpdatedAt        *gtime.Time `json:"updated_at" description:"更新时间"`
}

// 添加缺失的OrderConfirmReq和OrderConfirmRes结构体
type OrderConfirmReq struct {
	g.Meta `path:"/order/confirm" method:"post" tags:"前台订单" summary:"确认收货"`
	Id     uint `json:"id" v:"required" description:"订单ID"`
}

type OrderConfirmRes struct {
	Id uint `json:"id" description:"订单ID"`
}
