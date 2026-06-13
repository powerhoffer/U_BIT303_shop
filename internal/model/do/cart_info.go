// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CartInfo is the golang structure of table cart_info for DAO operations like Where/Data.
type CartInfo struct {
	g.Meta     `orm:"table:cart_info, do:true"`
	Id         any         // 购物车ID
	EmployeeId any         // 员工ID
	GoodsId    any         // 商品ID
	Count      any         // 商品数量
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 删除时间
}
