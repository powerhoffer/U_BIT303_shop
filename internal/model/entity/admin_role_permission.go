// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRolePermission is the golang structure for table admin_role_permission.
type AdminRolePermission struct {
	Id           uint        `json:"id"           orm:"id"            ` // Relation ID
	RoleId       uint        `json:"roleId"       orm:"role_id"       ` // Role ID
	PermissionId uint        `json:"permissionId" orm:"permission_id" ` // Permission ID
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // Created time
}
