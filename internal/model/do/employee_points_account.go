// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsAccount is the golang structure of table employee_points_account for DAO operations like Where/Data.
type EmployeePointsAccount struct {
	g.Meta     `orm:"table:employee_points_account, do:true"`
	Id         any         // 积分账户ID
	EmployeeId any         // 员工ID
	Balance    any         // 当前可用积分
	Status     any         // 状态：1正常 0停用
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 删除时间
}
