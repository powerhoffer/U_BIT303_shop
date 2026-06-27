// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminRole is the golang structure for table admin_role.
type AdminRole struct {
	Id          uint        `json:"id"          orm:"id"          ` // Role ID
	Name        string      `json:"name"        orm:"name"        ` // Role name
	Description string      `json:"description" orm:"description" ` // Role description
	Status      int         `json:"status"      orm:"status"      ` // Status: 1 enabled 0 disabled
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  ` // Created time
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  ` // Updated time
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"  ` // Deleted time
}
