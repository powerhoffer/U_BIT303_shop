// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SeckillOrder is the golang structure for table seckill_order.
type SeckillOrder struct {
	Id             int64       `json:"id"             description:"主键ID"`
	OrderNo        string      `json:"orderNo"        description:"订单编号"`
	UserId         int64       `json:"userId"         description:"用户ID"`
	GoodsId        int64       `json:"goodsId"        description:"商品ID"`
	GoodsOptionsId int64       `json:"goodsOptionsId" description:"商品规格ID"`
	OriginalPrice  int         `json:"originalPrice"  description:"原始价格 单位分"`
	SeckillPrice   int         `json:"seckillPrice"   description:"秒杀价格 单位分"`
	Status         int         `json:"status"         description:"订单状态：0-待支付 1-已支付 2-已取消 3-已退款"`
	PayTime        *gtime.Time `json:"payTime"        description:"支付时间"`
	CancelTime     *gtime.Time `json:"cancelTime"     description:"取消时间"`
	Count          int         `json:"count"          description:"商品数量"`
	ConsigneeName  string      `json:"consigneeName"  description:"收货人姓名"`
	ConsigneePhone string      `json:"consigneePhone" description:"收货人电话"`
	Address        string      `json:"address"        description:"收货地址"`
	Remark         string      `json:"remark"         description:"备注"`
	CreatedAt      *gtime.Time `json:"createdAt"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      description:"更新时间"`
}
