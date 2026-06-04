// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeeInfo is the golang structure for table employee_info.
type EmployeeInfo struct {
	Id           uint        `json:"id"           orm:"id"            ` // 员工ID
	Username     string      `json:"username"     orm:"username"      ` // 登录账号
	PasswordHash string      `json:"passwordHash" orm:"password_hash" ` // bcrypt密码哈希
	RealName     string      `json:"realName"     orm:"real_name"     ` // 员工姓名
	Phone        string      `json:"phone"        orm:"phone"         ` // 手机号
	Email        string      `json:"email"        orm:"email"         ` // 邮箱
	Status       int         `json:"status"       orm:"status"        ` // 状态：1正常 0禁用
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" ` // 最后登录时间
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    ` // 更新时间
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    ` // 删除时间
}
