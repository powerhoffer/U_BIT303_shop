// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRoleRelation is the golang structure of table admin_role_relation for DAO operations like Where/Data.
type AdminRoleRelation struct {
	g.Meta    `orm:"table:admin_role_relation, do:true"`
	Id        any         // Relation ID
	AdminId   any         // Admin ID
	RoleId    any         // Role ID
	CreatedAt *gtime.Time // Created time
}
