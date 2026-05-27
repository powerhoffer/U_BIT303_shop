package model

import (
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
)

type OrderListInput struct {
	Page           int    // 分页号码
	Size           int    // 分页数量，最大50
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

type OrderListOutput struct {
	List  []OrderListOutputItem
	Page  int `json:"page" description:"分页码"`
	Size  int `json:"size" description:"分页数量"`
	Total int `json:"total" description:"数据总数"`
}

type OrderListOutputItem struct {
	entity.OrderInfo
}

type OrderDetailInput struct {
	Id uint
}

type OrderDetailOutput struct {
	do.OrderInfo
	GoodsInfo []*do.OrderGoodsInfo `orm:"with:order_id=id"`
}

type OrderAddInput struct {
	UserId           uint
	Number           string
	Remark           string `description:"备注"`
	Price            int    `description:"订单金额 单位分"`
	CouponPrice      int    `description:"优惠券金额 单位分"`
	ActualPrice      int    `description:"实际支付金额 单位分"`
	ConsigneeName    string `description:"收货人姓名"`
	ConsigneePhone   string `description:"收货人手机号"`
	ConsigneeAddress string `description:"收货人详细地址"`
	//商品订单维度
	OrderAddGoodsInfos []*OrderAddGoodsInfo
}

type OrderAddGoodsInfo struct {
	Id             int
	OrderId        int
	GoodsId        int
	GoodsOptionsId int
	Count          int
	Remark         string
	Price          int
	CouponPrice    int
	ActualPrice    int
}

type OrderAddOutput struct {
	Id uint `json:"id"`
}

// SeckillOrderInput 秒杀订单输入
type SeckillOrderInput struct {
	UserId         uint    `json:"userId"`         // 用户ID
	GoodsId        uint    `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint    `json:"goodsOptionsId"` // 商品规格ID
	Count          uint    `json:"count"`          // 购买数量
	Price          float64 `json:"price"`          // 价格
	OrderNo        string  `json:"orderNo"`        // 订单号
}

// OrderCreateInput 订单创建输入
type OrderCreateInput struct {
	UserId         uint    `json:"userId"`         // 用户ID
	GoodsId        uint    `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint    `json:"goodsOptionsId"` // 商品选项ID
	Count          uint    `json:"count"`          // 购买数量
	Price          float64 `json:"price"`          // 总价
	Status         int     `json:"status"`         // 订单状态
	OrderNo        string  `json:"orderNo"`        // 订单号
	Address        string  `json:"address"`        // 收货地址
	Phone          string  `json:"phone"`          // 手机号码
	Remark         string  `json:"remark"`         // 备注
}
