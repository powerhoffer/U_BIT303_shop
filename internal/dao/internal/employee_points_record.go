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
	Id                 string // 积分流水ID
	EmployeeId         string // 员工ID
	ChangeType         string // 变动类型：1增加 2扣除
	Points             string // 变动积分
	BeforeBalance      string // 变动前积分
	AfterBalance       string // 变动后积分
	OperatorEmployeeId string // 操作员工ID
	Remark             string // 备注
	CreatedAt          string // 创建时间
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
