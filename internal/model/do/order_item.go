// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItem is the golang structure of table order_item for DAO operations like Where/Data.
type OrderItem struct {
	g.Meta        `orm:"table:order_item, do:true"`
	Id            any         // Order item ID
	OrderId       any         // Order ID
	EmployeeId    any         // Employee ID
	GoodsId       any         // Goods ID
	GoodsName     any         // Goods name snapshot
	GoodsImageUrl any         // Goods image snapshot
	PointsPrice   any         // Points price snapshot
	Count         any         // Goods count
	TotalPoints   any         // Item total points
	CreatedAt     *gtime.Time // Created time
}
