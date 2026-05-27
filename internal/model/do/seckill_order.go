// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SeckillOrder is the golang structure of table seckill_order for DAO operations like Where/Data.
type SeckillOrder struct {
	g.Meta         `orm:"table:seckill_order, do:true"`
	Id             interface{} // 主键ID
	OrderNo        interface{} // 订单编号
	UserId         interface{} // 用户ID
	GoodsId        interface{} // 商品ID
	GoodsOptionsId interface{} // 商品规格ID
	OriginalPrice  interface{} // 原始价格 单位分
	SeckillPrice   interface{} // 秒杀价格 单位分
	Status         interface{} // 订单状态：0-待支付 1-已支付 2-已取消 3-已退款
	PayTime        *gtime.Time // 支付时间
	CancelTime     *gtime.Time // 取消时间
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
