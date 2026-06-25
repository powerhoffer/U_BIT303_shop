// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderItem is the golang structure for table order_item.
type OrderItem struct {
	Id            uint        `json:"id"            orm:"id"              ` // Order item ID
	OrderId       uint        `json:"orderId"       orm:"order_id"        ` // Order ID
	EmployeeId    uint        `json:"employeeId"    orm:"employee_id"     ` // Employee ID
	GoodsId       uint        `json:"goodsId"       orm:"goods_id"        ` // Goods ID
	GoodsName     string      `json:"goodsName"     orm:"goods_name"      ` // Goods name snapshot
	GoodsImageUrl string      `json:"goodsImageUrl" orm:"goods_image_url" ` // Goods image snapshot
	PointsPrice   uint        `json:"pointsPrice"   orm:"points_price"    ` // Points price snapshot
	Count         uint        `json:"count"         orm:"count"           ` // Goods count
	TotalPoints   uint        `json:"totalPoints"   orm:"total_points"    ` // Item total points
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      ` // Created time
}
