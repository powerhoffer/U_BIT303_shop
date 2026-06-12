// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsInfo is the golang structure for table goods_info.
type GoodsInfo struct {
	Id          uint        `json:"id"          orm:"id"           ` // 商品ID
	CategoryId  uint        `json:"categoryId"  orm:"category_id"  ` // 商品分类ID
	Name        string      `json:"name"        orm:"name"         ` // 商品名称
	ImageUrl    string      `json:"imageUrl"    orm:"image_url"    ` // 商品图片
	PointsPrice uint        `json:"pointsPrice" orm:"points_price" ` // 兑换所需积分
	Stock       uint        `json:"stock"       orm:"stock"        ` // 库存
	Description string      `json:"description" orm:"description"  ` // 商品简介
	Status      int         `json:"status"      orm:"status"       ` // 状态：1上架 0下架
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` // 更新时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   ` // 删除时间
}
