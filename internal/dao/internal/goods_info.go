// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GoodsInfoDao is the data access object for the table goods_info.
type GoodsInfoDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  GoodsInfoColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// GoodsInfoColumns defines and stores column names for the table goods_info.
type GoodsInfoColumns struct {
	Id          string // 商品ID
	CategoryId  string // 商品分类ID
	Name        string // 商品名称
	ImageUrl    string // 商品图片
	PointsPrice string // 兑换所需积分
	Stock       string // 库存
	Description string // 商品简介
	Status      string // 状态：1上架 0下架
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// goodsInfoColumns holds the columns for the table goods_info.
var goodsInfoColumns = GoodsInfoColumns{
	Id:          "id",
	CategoryId:  "category_id",
	Name:        "name",
	ImageUrl:    "image_url",
	PointsPrice: "points_price",
	Stock:       "stock",
	Description: "description",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewGoodsInfoDao creates and returns a new DAO object for table data access.
func NewGoodsInfoDao(handlers ...gdb.ModelHandler) *GoodsInfoDao {
	return &GoodsInfoDao{
		group:    "default",
		table:    "goods_info",
		columns:  goodsInfoColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GoodsInfoDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GoodsInfoDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GoodsInfoDao) Columns() GoodsInfoColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GoodsInfoDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GoodsInfoDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GoodsInfoDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
