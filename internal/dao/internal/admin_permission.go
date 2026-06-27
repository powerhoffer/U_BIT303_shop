// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminPermissionDao is the data access object for the table admin_permission.
type AdminPermissionDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  AdminPermissionColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// AdminPermissionColumns defines and stores column names for the table admin_permission.
type AdminPermissionColumns struct {
	Id        string // Permission ID
	Name      string // Permission name
	GroupName string // Permission group
	Method    string // HTTP method
	Path      string // API path
	Status    string // Status: 1 enabled 0 disabled
	CreatedAt string // Created time
	UpdatedAt string // Updated time
	DeletedAt string // Deleted time
}

// adminPermissionColumns holds the columns for the table admin_permission.
var adminPermissionColumns = AdminPermissionColumns{
	Id:        "id",
	Name:      "name",
	GroupName: "group_name",
	Method:    "method",
	Path:      "path",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewAdminPermissionDao creates and returns a new DAO object for table data access.
func NewAdminPermissionDao(handlers ...gdb.ModelHandler) *AdminPermissionDao {
	return &AdminPermissionDao{
		group:    "default",
		table:    "admin_permission",
		columns:  adminPermissionColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AdminPermissionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AdminPermissionDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AdminPermissionDao) Columns() AdminPermissionColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AdminPermissionDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AdminPermissionDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AdminPermissionDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
