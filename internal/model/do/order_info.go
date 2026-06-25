// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderInfo is the golang structure of table order_info for DAO operations like Where/Data.
type OrderInfo struct {
	g.Meta      `orm:"table:order_info, do:true"`
	Id          any         // Order ID
	OrderNo     any         // Order number
	EmployeeId  any         // Employee ID
	TotalPoints any         // Total points
	Status      any         // Status: 1 pending 2 completed 3 cancelled
	Remark      any         // Remark
	CreatedAt   *gtime.Time // Created time
	UpdatedAt   *gtime.Time // Updated time
	DeletedAt   *gtime.Time // Deleted time
}
