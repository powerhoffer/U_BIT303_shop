// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeeInfo is the golang structure of table employee_info for DAO operations like Where/Data.
type EmployeeInfo struct {
	g.Meta       `orm:"table:employee_info, do:true"`
	Id           any         // 员工ID
	Username     any         // 登录账号
	PasswordHash any         // bcrypt密码哈希
	RealName     any         // 员工姓名
	Phone        any         // 手机号
	Email        any         // 邮箱
	Status       any         // 状态：1正常 0禁用
	LastLoginAt  *gtime.Time // 最后登录时间
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
	DeletedAt    *gtime.Time // 删除时间
}
