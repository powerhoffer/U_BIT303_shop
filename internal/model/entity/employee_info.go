// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeeInfo is the golang structure for table employee_info.
type EmployeeInfo struct {
	Id           uint        `json:"id"           orm:"id"            ` // Employee ID
	Username     string      `json:"username"     orm:"username"      ` // Login username
	PasswordHash string      `json:"passwordHash" orm:"password_hash" ` // BCrypt password hash
	RealName     string      `json:"realName"     orm:"real_name"     ` // Employee name
	Phone        string      `json:"phone"        orm:"phone"         ` // Phone number
	Email        string      `json:"email"        orm:"email"         ` // Email address
	Status       int         `json:"status"       orm:"status"        ` // Status: 1 active, 0 disabled
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" ` // Last login time
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // Created time
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // Updated time
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    ` // Deleted time
}
