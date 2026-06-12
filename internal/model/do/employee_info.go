// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeeInfo is the golang structure of table employee_info for DAO operations like Where/Data.
type EmployeeInfo struct {
	g.Meta       `orm:"table:employee_info, do:true"`
	Id           any         // Employee ID
	Username     any         // Login username
	PasswordHash any         // BCrypt password hash
	RealName     any         // Employee name
	Phone        any         // Phone number
	Email        any         // Email address
	Status       any         // Status: 1 active, 0 disabled
	LastLoginAt  *gtime.Time // Last login time
	CreatedAt    *gtime.Time // Created time
	UpdatedAt    *gtime.Time // Updated time
	DeletedAt    *gtime.Time // Deleted time
}
