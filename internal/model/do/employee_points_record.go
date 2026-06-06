// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsRecord is the golang structure of table employee_points_record for DAO operations like Where/Data.
type EmployeePointsRecord struct {
	g.Meta             `orm:"table:employee_points_record, do:true"`
	Id                 any         // 积分流水ID
	EmployeeId         any         // 员工ID
	ChangeType         any         // 变动类型：1增加 2扣除
	Points             any         // 变动积分
	BeforeBalance      any         // 变动前积分
	AfterBalance       any         // 变动后积分
	OperatorEmployeeId any         // 操作员工ID
	Remark             any         // 备注
	CreatedAt          *gtime.Time // 创建时间
}
