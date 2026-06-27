// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminInfo is the golang structure for table admin_info.
type AdminInfo struct {
	Id           uint        `json:"id"           orm:"id"            ` // Admin ID
	Username     string      `json:"username"     orm:"username"      ` // Login username
	PasswordHash string      `json:"passwordHash" orm:"password_hash" ` // bcrypt password hash
	RealName     string      `json:"realName"     orm:"real_name"     ` // Admin name
	Phone        string      `json:"phone"        orm:"phone"         ` // Phone
	Email        string      `json:"email"        orm:"email"         ` // Email
	Status       int         `json:"status"       orm:"status"        ` // Status: 1 normal 0 disabled
	IsSuper      int         `json:"isSuper"      orm:"is_super"      ` // Is super admin: 1 yes 0 no
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" ` // Last login time
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // Created time
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // Updated time
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    ` // Deleted time
}
