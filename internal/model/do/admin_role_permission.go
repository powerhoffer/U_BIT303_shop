// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRolePermission is the golang structure of table admin_role_permission for DAO operations like Where/Data.
type AdminRolePermission struct {
	g.Meta       `orm:"table:admin_role_permission, do:true"`
	Id           any         // Relation ID
	RoleId       any         // Role ID
	PermissionId any         // Permission ID
	CreatedAt    *gtime.Time // Created time
}
