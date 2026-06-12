// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsInfo is the golang structure of table goods_info for DAO operations like Where/Data.
type GoodsInfo struct {
	g.Meta      `orm:"table:goods_info, do:true"`
	Id          any         // 商品ID
	CategoryId  any         // 商品分类ID
	Name        any         // 商品名称
	ImageUrl    any         // 商品图片
	PointsPrice any         // 兑换所需积分
	Stock       any         // 库存
	Description any         // 商品简介
	Status      any         // 状态：1上架 0下架
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}
