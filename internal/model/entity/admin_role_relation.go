// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRoleRelation is the golang structure for table admin_role_relation.
type AdminRoleRelation struct {
	Id        uint        `json:"id"        orm:"id"         ` // Relation ID
	AdminId   uint        `json:"adminId"   orm:"admin_id"   ` // Admin ID
	RoleId    uint        `json:"roleId"    orm:"role_id"    ` // Role ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // Created time
}
