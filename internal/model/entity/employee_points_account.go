// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsAccount is the golang structure for table employee_points_account.
type EmployeePointsAccount struct {
	Id         uint        `json:"id"         orm:"id"          ` // Credit account ID
	EmployeeId uint        `json:"employeeId" orm:"employee_id" ` // Employee ID
	Balance    uint        `json:"balance"    orm:"balance"     ` // Available credit balance
	Status     int         `json:"status"     orm:"status"      ` // Status: 1 active, 0 disabled
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // Created time
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // Updated time
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  ` // Deleted time
}
