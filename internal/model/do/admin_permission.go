// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminPermission is the golang structure of table admin_permission for DAO operations like Where/Data.
type AdminPermission struct {
	g.Meta    `orm:"table:admin_permission, do:true"`
	Id        any         // Permission ID
	Name      any         // Permission name
	GroupName any         // Permission group
	Method    any         // HTTP method
	Path      any         // API path
	Status    any         // Status: 1 enabled 0 disabled
	CreatedAt *gtime.Time // Created time
	UpdatedAt *gtime.Time // Updated time
	DeletedAt *gtime.Time // Deleted time
}
