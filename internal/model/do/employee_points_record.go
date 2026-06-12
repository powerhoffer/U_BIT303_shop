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
	Id                 any         // Credit record ID
	EmployeeId         any         // Employee ID
	ChangeType         any         // Change type: 1 add, 2 deduct
	Points             any         // Changed credits
	BeforeBalance      any         // Balance before change
	AfterBalance       any         // Balance after change
	OperatorEmployeeId any         // Operator employee ID
	Remark             any         // Remark
	CreatedAt          *gtime.Time // Created time
}
