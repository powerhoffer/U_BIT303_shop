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
	Id        any         // 商品分类ID
	Name      any         // 分类名称
	Sort      any         // 排序值
	Status    any         // 状态：1启用 0停用
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
	DeletedAt *gtime.Time // 删除时间
}
