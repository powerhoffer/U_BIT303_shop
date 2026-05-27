package backend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type OrderListReq struct {
	g.Meta `path:"/order/list" tags:"订单列表" method:"get" summary:"订单列表"`
	CommonPaginationReq
	Number         string `json:"number"           dc:"订单编号"`
	UserId         int    `json:"userId"           dc:"用户id"`
	PayType        int    `json:"payType"          dc:"支付方式 1微信 2支付宝 3云闪付"`
	PayAtGte       string `json:"payAtGte"         dc:"支付时间>="`
	PayAtLte       string `json:"payAtLte"         dc:"支付时间<="`
	Status         int    `json:"status"           dc:"订单状态： 1待支付 2已支付待发货 3已发货 4已收货待评价"`
	ConsigneePhone string `json:"consigneePhone"   dc:"收货人手机号"`
	PriceGte       int    `json:"priceGte"         dc:"订单金额>= 单位分"`
	PriceLte       int    `json:"priceLte"         dc:"订单金额<= 单位分"`
	DateGte        string `json:"dateGte"          dc:"創建时间>="`
	DateLte        string `json:"dateLte"          dc:"創建时间<="`
}

type OrderListRes struct {
	CommonPaginationRes
}

type OrderDetailReq struct {
	g.Meta `path:"/order/detail" tags:"订单详情" method:"get" summary:"订单详情"`
	Id     uint `json:"id"`
}

type OrderDetailRes struct {
	OrderInfoBase
	GoodsInfo []OrderGoodsInfoBase `json:"goods_info" dc:"订单商品列表"`
}

type OrderInfoBase struct {
	Id               int         `json:"id"               description:""`
	Number           string      `json:"number"           description:"订单编号"`
	UserId           int         `json:"userId"           description:"用户id"`
	PayType          int         `json:"payType"          description:"支付方式 1微信 2支付宝 3云闪付"`
	Remark           string      `json:"remark"           description:"备注"`
	PayAt            *gtime.Time `json:"payAt"            description:"支付时间"`
	Status           int         `json:"status"           description:"订单状态： 1待支付 2已支付待发货 3已发货 4已收货待评价"`
	ConsigneeName    string      `json:"consigneeName"    description:"收货人姓名"`
	ConsigneePhone   string      `json:"consigneePhone"   description:"收货人手机号"`
	ConsigneeAddress string      `json:"consigneeAddress" description:"收货人详细地址"`
	Price            int         `json:"price"            description:"订单金额 单位分"`
	CouponPrice      int         `json:"couponPrice"      description:"优惠券金额 单位分"`
	ActualPrice      int         `json:"actualPrice"      description:"实际支付金额 单位分"`
	CreatedAt        *gtime.Time `json:"createdAt"        description:""`
	UpdatedAt        *gtime.Time `json:"updatedAt"        description:""`
}

type OrderGoodsInfoBase struct {
	Id          int         `json:"id"          description:"商品维度的订单表"`
	OrderId     int         `json:"orderId"     description:"关联的主订单表"`
	GoodsId     int         `json:"goodsId"     description:"商品id"`
	Count       int         `json:"count"       description:"商品数量"`
	Remark      string      `json:"remark"      description:"备注"`
	Price       int         `json:"price"       description:"订单金额 单位分"`
	CouponPrice int         `json:"couponPrice" description:"优惠券金额 单位分"`
	ActualPrice int         `json:"actualPrice" description:"实际支付金额 单位分"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:""`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:""`
}

// 添加缺失的结构体定义
type OrderAddReq struct {
	g.Meta    `path:"/order/add" method:"post" tags:"订单管理" summary:"创建订单"`
	GoodsList []OrderGoodsInfo `json:"goods_list" dc:"商品列表"`
}

type OrderAddRes struct {
	Id uint `json:"id" dc:"订单ID"`
}

type OrderUpdateStatusReq struct {
	g.Meta `path:"/order/update/status" method:"post" tags:"订单管理" summary:"更新订单状态"`
	Id     uint `json:"id" v:"required" dc:"订单ID"`
	Status int  `json:"status" v:"required" dc:"订单状态"`
}

type OrderUpdateStatusRes struct {
	Id uint `json:"id" dc:"订单ID"`
}

type OrderDeleteReq struct {
	g.Meta `path:"/order/delete" method:"delete" tags:"订单管理" summary:"删除订单"`
	Id     uint `json:"id" v:"required" dc:"订单ID"`
}

type OrderDeleteRes struct {
	Id uint `json:"id" dc:"订单ID"`
}

type OrderRefundReq struct {
	g.Meta `path:"/order/refund" method:"post" tags:"订单管理" summary:"订单退款"`
	Id     uint   `json:"id" v:"required" dc:"订单ID"`
	Reason string `json:"reason" dc:"退款原因"`
}

type OrderRefundRes struct {
	Id uint `json:"id" dc:"订单ID"`
}

type OrderGoodsInfo struct {
	GoodsId        int    `json:"goods_id" dc:"商品ID"`
	GoodsOptionsId int    `json:"goods_options_id" dc:"商品选项ID"`
	Count          int    `json:"count" dc:"商品数量"`
	Price          int    `json:"price" dc:"价格"`
	CouponPrice    int    `json:"coupon_price" dc:"优惠券金额"`
	ActualPrice    int    `json:"actual_price" dc:"实际支付金额"`
	Remark         string `json:"remark" dc:"备注"`
}

type OrderInfo struct {
	OrderBase
	Number         string      `json:"number" dc:"订单编号"`
	UserId         int         `json:"user_id" dc:"用户ID"`
	PayType        int         `json:"pay_type" dc:"支付方式"`
	PayAt          *gtime.Time `json:"pay_at" dc:"支付时间"`
	ShipAt         *gtime.Time `json:"ship_at" dc:"发货时间"`
	FinishAt       *gtime.Time `json:"finish_at" dc:"完成时间"`
	Status         int         `json:"status" dc:"订单状态"`
	ConsigneeName  string      `json:"consignee_name" dc:"收货人姓名"`
	ConsigneePhone string      `json:"consignee_phone" dc:"收货人手机号"`
	Address        string      `json:"address" dc:"地址"`
}

type OrderBase struct {
	Id          int         `json:"id" dc:"订单ID"`
	Price       int         `json:"price" dc:"订单金额"`
	CouponPrice int         `json:"coupon_price" dc:"优惠券金额"`
	ActualPrice int         `json:"actual_price" dc:"实际支付金额"`
	Remark      string      `json:"remark" dc:"备注"`
	CreatedAt   *gtime.Time `json:"created_at" dc:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updated_at" dc:"更新时间"`
}
