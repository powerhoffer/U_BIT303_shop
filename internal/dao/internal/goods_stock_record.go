// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// GoodsStockRecordDao is the data access object for the table goods_stock_record.
type GoodsStockRecordDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  GoodsStockRecordColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// GoodsStockRecordColumns defines and stores column names for the table goods_stock_record.
type GoodsStockRecordColumns struct {
	Id             string // 库存流水ID
	GoodsId        string // 商品ID
	GoodsName      string // 商品名称快照
	ChangeType     string // 变动类型：1初始库存 2后台增加 3后台扣减 4订单扣减 5取消订单恢复
	ChangeQuantity string // 库存变动数量，正数增加，负数扣减
	BeforeStock    string // 变动前库存
	AfterStock     string // 变动后库存
	BizType        string // 业务来源
	BizId          string // 关联业务ID
	OperatorType   string // 操作者类型：0系统 1管理员 2员工
	OperatorId     string // 操作者ID
	Remark         string // 备注
	CreatedAt      string // 创建时间
}

// goodsStockRecordColumns holds the columns for the table goods_stock_record.
var goodsStockRecordColumns = GoodsStockRecordColumns{
	Id:             "id",
	GoodsId:        "goods_id",
	GoodsName:      "goods_name",
	ChangeType:     "change_type",
	ChangeQuantity: "change_quantity",
	BeforeStock:    "before_stock",
	AfterStock:     "after_stock",
	BizType:        "biz_type",
	BizId:          "biz_id",
	OperatorType:   "operator_type",
	OperatorId:     "operator_id",
	Remark:         "remark",
	CreatedAt:      "created_at",
}

// NewGoodsStockRecordDao creates and returns a new DAO object for table data access.
func NewGoodsStockRecordDao(handlers ...gdb.ModelHandler) *GoodsStockRecordDao {
	return &GoodsStockRecordDao{
		group:    "default",
		table:    "goods_stock_record",
		columns:  goodsStockRecordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *GoodsStockRecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *GoodsStockRecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *GoodsStockRecordDao) Columns() GoodsStockRecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *GoodsStockRecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *GoodsStockRecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *GoodsStockRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
