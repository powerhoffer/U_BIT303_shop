// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OrderInfo is the golang structure for table order_info.
type OrderInfo struct {
	Id          uint        `json:"id"          orm:"id"           ` // Order ID
	OrderNo     string      `json:"orderNo"     orm:"order_no"     ` // Order number
	EmployeeId  uint        `json:"employeeId"  orm:"employee_id"  ` // Employee ID
	TotalPoints uint        `json:"totalPoints" orm:"total_points" ` // Total points
	Status      int         `json:"status"      orm:"status"       ` // Status: 1 pending 2 completed 3 cancelled
	Remark      string      `json:"remark"      orm:"remark"       ` // Remark
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   ` // Created time
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   ` // Updated time
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   ` // Deleted time
}
