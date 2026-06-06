// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsRecord is the golang structure for table employee_points_record.
type EmployeePointsRecord struct {
	Id                 uint        `json:"id"                 orm:"id"                   ` // 积分流水ID
	EmployeeId         uint        `json:"employeeId"         orm:"employee_id"          ` // 员工ID
	ChangeType         int         `json:"changeType"         orm:"change_type"          ` // 变动类型：1增加 2扣除
	Points             uint        `json:"points"             orm:"points"               ` // 变动积分
	BeforeBalance      uint        `json:"beforeBalance"      orm:"before_balance"       ` // 变动前积分
	AfterBalance       uint        `json:"afterBalance"       orm:"after_balance"        ` // 变动后积分
	OperatorEmployeeId uint        `json:"operatorEmployeeId" orm:"operator_employee_id" ` // 操作员工ID
	Remark             string      `json:"remark"             orm:"remark"               ` // 备注
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           ` // 创建时间
}
