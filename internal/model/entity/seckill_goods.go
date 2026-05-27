// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SeckillGoods is the golang structure for table seckill_goods.
type SeckillGoods struct {
	Id             int64       `json:"id"             description:"主键ID"`
	GoodsId        int64       `json:"goodsId"        description:"商品ID"`
	GoodsOptionsId int64       `json:"goodsOptionsId" description:"商品规格ID"`
	OriginalPrice  int         `json:"originalPrice"  description:"原始价格 单位分"`
	SeckillPrice   int         `json:"seckillPrice"   description:"秒杀价格 单位分"`
	SeckillStock   int         `json:"seckillStock"   description:"秒杀库存"`
	StartTime      *gtime.Time `json:"startTime"      description:"秒杀开始时间"`
	EndTime        *gtime.Time `json:"endTime"        description:"秒杀结束时间"`
	Status         int         `json:"status"         description:"状态：0-未开始 1-进行中 2-已结束"`
	CreatedAt      *gtime.Time `json:"createdAt"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      description:"更新时间"`
}
