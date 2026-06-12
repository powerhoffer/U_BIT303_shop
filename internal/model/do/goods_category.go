// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsCategory is the golang structure of table goods_category for DAO operations like Where/Data.
type GoodsCategory struct {
	g.Meta    `orm:"table:goods_category, do:true"`
	Id        any         // Goods category ID
	Name      any         // Category name
	Sort      any         // Sort order
	Status    any         // Status: 1 enabled, 0 disabled
	CreatedAt *gtime.Time // Created time
	UpdatedAt *gtime.Time // Updated time
	DeletedAt *gtime.Time // Deleted time
}
