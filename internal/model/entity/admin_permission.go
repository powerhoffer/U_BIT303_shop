// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminPermission is the golang structure for table admin_permission.
type AdminPermission struct {
	Id        uint        `json:"id"        orm:"id"         ` // Permission ID
	Name      string      `json:"name"      orm:"name"       ` // Permission name
	GroupName string      `json:"groupName" orm:"group_name" ` // Permission group
	Method    string      `json:"method"    orm:"method"     ` // HTTP method
	Path      string      `json:"path"      orm:"path"       ` // API path
	Status    int         `json:"status"    orm:"status"     ` // Status: 1 enabled 0 disabled
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` // Created time
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` // Updated time
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" ` // Deleted time
}
