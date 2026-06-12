// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsCategory is the golang structure for table goods_category.
type GoodsCategory struct {
	Id        uint        `json:"id"        orm:"id"         ` // Goods category ID
	Name      string      `json:"name"      orm:"name"       ` // Category name
	Sort      uint        `json:"sort"      orm:"sort"       ` // Sort order
	Status    int         `json:"status"    orm:"status"     ` // Status: 1 enabled, 0 disabled
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // Created time
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // Updated time
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" ` // Deleted time
}
