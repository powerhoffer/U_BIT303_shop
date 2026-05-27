// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SeckillOrderDao is the data access object for table seckill_order.
type SeckillOrderDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns SeckillOrderColumns // columns contains all the column names of Table for convenient usage.
}

// SeckillOrderColumns defines and stores column names for table seckill_order.
type SeckillOrderColumns struct {
	Id             string // 主键ID
	OrderNo        string // 订单编号
	UserId         string // 用户ID
	GoodsId        string // 商品ID
	GoodsOptionsId string // 商品规格ID
	OriginalPrice  string // 原始价格 单位分
	SeckillPrice   string // 秒杀价格 单位分
	Status         string // 订单状态：0-待支付 1-已支付 2-已取消 3-已退款
	PayTime        string // 支付时间
	CancelTime     string // 取消时间
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// seckillOrderColumns holds the columns for table seckill_order.
var seckillOrderColumns = SeckillOrderColumns{
	Id:             "id",
	OrderNo:        "order_no",
	UserId:         "user_id",
	GoodsId:        "goods_id",
	GoodsOptionsId: "goods_options_id",
	OriginalPrice:  "original_price",
	SeckillPrice:   "seckill_price",
	Status:         "status",
	PayTime:        "pay_time",
	CancelTime:     "cancel_time",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewSeckillOrderDao creates and returns a new DAO object for table data access.
func NewSeckillOrderDao() *SeckillOrderDao {
	return &SeckillOrderDao{
		group:   "default",
		table:   "seckill_order",
		columns: seckillOrderColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SeckillOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SeckillOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SeckillOrderDao) Columns() SeckillOrderColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SeckillOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SeckillOrderDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SeckillOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
