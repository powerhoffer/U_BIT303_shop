// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EmployeePointsRecordDao is the data access object for the table employee_points_record.
type EmployeePointsRecordDao struct {
	table    string                      // table is the underlying table name of the DAO.
	group    string                      // group is the database configuration group name of the current DAO.
	columns  EmployeePointsRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler          // handlers for customized model modification.
}

// EmployeePointsRecordColumns defines and stores column names for the table employee_points_record.
type EmployeePointsRecordColumns struct {
	Id                 string // Credit record ID
	EmployeeId         string // Employee ID
	ChangeType         string // Change type: 1 add, 2 deduct
	Points             string // Changed credits
	BeforeBalance      string // Balance before change
	AfterBalance       string // Balance after change
	OperatorEmployeeId string // Operator employee ID
	OperatorAdminId    string // Operator admin ID
	Remark             string // Remark
	CreatedAt          string // Created time
}

// employeePointsRecordColumns holds the columns for the table employee_points_record.
var employeePointsRecordColumns = EmployeePointsRecordColumns{
	Id:                 "id",
	EmployeeId:         "employee_id",
	ChangeType:         "change_type",
	Points:             "points",
	BeforeBalance:      "before_balance",
	AfterBalance:       "after_balance",
	OperatorEmployeeId: "operator_employee_id",
	OperatorAdminId:    "operator_admin_id",
	Remark:             "remark",
	CreatedAt:          "created_at",
}

// NewEmployeePointsRecordDao creates and returns a new DAO object for table data access.
func NewEmployeePointsRecordDao(handlers ...gdb.ModelHandler) *EmployeePointsRecordDao {
	return &EmployeePointsRecordDao{
		group:    "default",
		table:    "employee_points_record",
		columns:  employeePointsRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EmployeePointsRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EmployeePointsRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EmployeePointsRecordDao) Columns() EmployeePointsRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EmployeePointsRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EmployeePointsRecordDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *EmployeePointsRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
