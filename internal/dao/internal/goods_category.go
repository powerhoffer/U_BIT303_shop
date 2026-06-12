// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GoodsCategoryDao is the data access object for the table goods_category.
type GoodsCategoryDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  GoodsCategoryColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// GoodsCategoryColumns defines and stores column names for the table goods_category.
type GoodsCategoryColumns struct {
	Id        string // Goods category ID
	Name      string // Category name
	Sort      string // Sort order
	Status    string // Status: 1 enabled, 0 disabled
	CreatedAt string // Created time
	UpdatedAt string // Updated time
	DeletedAt string // Deleted time
}

// goodsCategoryColumns holds the columns for the table goods_category.
var goodsCategoryColumns = GoodsCategoryColumns{
	Id:        "id",
	Name:      "name",
	Sort:      "sort",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewGoodsCategoryDao creates and returns a new DAO object for table data access.
func NewGoodsCategoryDao(handlers ...gdb.ModelHandler) *GoodsCategoryDao {
	return &GoodsCategoryDao{
		group:    "default",
		table:    "goods_category",
		columns:  goodsCategoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GoodsCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GoodsCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GoodsCategoryDao) Columns() GoodsCategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GoodsCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GoodsCategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GoodsCategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
