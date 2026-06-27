// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRole is the golang structure of table admin_role for DAO operations like Where/Data.
type AdminRole struct {
	g.Meta      `orm:"table:admin_role, do:true"`
	Id          any         // Role ID
	Name        any         // Role name
	Description any         // Role description
	Status      any         // Status: 1 enabled 0 disabled
	CreatedAt   *gtime.Time // Created time
	UpdatedAt   *gtime.Time // Updated time
	DeletedAt   *gtime.Time // Deleted time
}
