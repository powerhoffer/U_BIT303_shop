// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsAccount is the golang structure of table employee_points_account for DAO operations like Where/Data.
type EmployeePointsAccount struct {
	g.Meta     `orm:"table:employee_points_account, do:true"`
	Id         any         // Credit account ID
	EmployeeId any         // Employee ID
	Balance    any         // Available credit balance
	Status     any         // Status: 1 active, 0 disabled
	CreatedAt  *gtime.Time // Created time
	UpdatedAt  *gtime.Time // Updated time
	DeletedAt  *gtime.Time // Deleted time
}
