// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SeckillGoods is the golang structure of table seckill_goods for DAO operations like Where/Data.
type SeckillGoods struct {
	g.Meta         `orm:"table:seckill_goods, do:true"`
	Id             interface{} // 主键ID
	GoodsId        interface{} // 商品ID
	GoodsOptionsId interface{} // 商品规格ID
	OriginalPrice  interface{} // 原始价格 单位分
	SeckillPrice   interface{} // 秒杀价格 单位分
	SeckillStock   interface{} // 秒杀库存
	StartTime      *gtime.Time // 秒杀开始时间
	EndTime        *gtime.Time // 秒杀结束时间
	Status         interface{} // 状态：0-未开始 1-进行中 2-已结束
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}
