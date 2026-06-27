// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminInfo is the golang structure of table admin_info for DAO operations like Where/Data.
type AdminInfo struct {
	g.Meta       `orm:"table:admin_info, do:true"`
	Id           any         // Admin ID
	Username     any         // Login username
	PasswordHash any         // bcrypt password hash
	RealName     any         // Admin name
	Phone        any         // Phone
	Email        any         // Email
	Status       any         // Status: 1 normal 0 disabled
	IsSuper      any         // Is super admin: 1 yes 0 no
	LastLoginAt  *gtime.Time // Last login time
	CreatedAt    *gtime.Time // Created time
	UpdatedAt    *gtime.Time // Updated time
	DeletedAt    *gtime.Time // Deleted time
}
