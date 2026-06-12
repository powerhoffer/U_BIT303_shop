// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// EmployeePointsRecord is the golang structure for table employee_points_record.
type EmployeePointsRecord struct {
	Id                 uint        `json:"id"                 orm:"id"                   ` // Credit record ID
	EmployeeId         uint        `json:"employeeId"         orm:"employee_id"          ` // Employee ID
	ChangeType         int         `json:"changeType"         orm:"change_type"          ` // Change type: 1 add, 2 deduct
	Points             uint        `json:"points"             orm:"points"               ` // Changed credits
	BeforeBalance      uint        `json:"beforeBalance"      orm:"before_balance"       ` // Balance before change
	AfterBalance       uint        `json:"afterBalance"       orm:"after_balance"        ` // Balance after change
	OperatorEmployeeId uint        `json:"operatorEmployeeId" orm:"operator_employee_id" ` // Operator employee ID
	Remark             string      `json:"remark"             orm:"remark"               ` // Remark
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           ` // Created time
}
