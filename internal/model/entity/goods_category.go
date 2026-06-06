// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsCategory is the golang structure for table goods_category.
type GoodsCategory struct {
	Id        uint        `json:"id"        orm:"id"         ` // 商品分类ID
	Name      string      `json:"name"      orm:"name"       ` // 分类名称
	Sort      uint        `json:"sort"      orm:"sort"       ` // 排序值
	Status    int         `json:"status"    orm:"status"     ` // 状态：1启用 0停用
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // 更新时间
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" ` // 删除时间
}
