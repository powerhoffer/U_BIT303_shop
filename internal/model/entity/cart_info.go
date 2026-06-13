// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CartInfo is the golang structure for table cart_info.
type CartInfo struct {
	Id         uint        `json:"id"         orm:"id"          ` // 购物车ID
	EmployeeId uint        `json:"employeeId" orm:"employee_id" ` // 员工ID
	GoodsId    uint        `json:"goodsId"    orm:"goods_id"    ` // 商品ID
	Count      uint        `json:"count"      orm:"count"       ` // 商品数量
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // 更新时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  ` // 删除时间
}
