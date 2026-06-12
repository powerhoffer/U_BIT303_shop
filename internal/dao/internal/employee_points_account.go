// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EmployeePointsAccountDao is the data access object for the table employee_points_account.
type EmployeePointsAccountDao struct {
	table    string                       // table is the underlying table name of the DAO.
	group    string                       // group is the database configuration group name of the current DAO.
	columns  EmployeePointsAccountColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler           // handlers for customized model modification.
}

// EmployeePointsAccountColumns defines and stores column names for the table employee_points_account.
type EmployeePointsAccountColumns struct {
	Id         string // Credit account ID
	EmployeeId string // Employee ID
	Balance    string // Available credit balance
	Status     string // Status: 1 active, 0 disabled
	CreatedAt  string // Created time
	UpdatedAt  string // Updated time
	DeletedAt  string // Deleted time
}

// employeePointsAccountColumns holds the columns for the table employee_points_account.
var employeePointsAccountColumns = EmployeePointsAccountColumns{
	Id:         "id",
	EmployeeId: "employee_id",
	Balance:    "balance",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewEmployeePointsAccountDao creates and returns a new DAO object for table data access.
func NewEmployeePointsAccountDao(handlers ...gdb.ModelHandler) *EmployeePointsAccountDao {
	return &EmployeePointsAccountDao{
		group:    "default",
		table:    "employee_points_account",
		columns:  employeePointsAccountColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EmployeePointsAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EmployeePointsAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EmployeePointsAccountDao) Columns() EmployeePointsAccountColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EmployeePointsAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EmployeePointsAccountDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *EmployeePointsAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
