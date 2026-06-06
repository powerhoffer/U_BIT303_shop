// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsAccount is the golang structure for table employee_points_account.
type EmployeePointsAccount struct {
	Id         uint        `json:"id"         orm:"id"          ` // 积分账户ID
	EmployeeId uint        `json:"employeeId" orm:"employee_id" ` // 员工ID
	Balance    uint        `json:"balance"    orm:"balance"     ` // 当前可用积分
	Status     int         `json:"status"     orm:"status"      ` // 状态：1正常 0停用
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  ` // 创建时间
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  ` // 更新时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  ` // 删除时间
}
