// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SeckillGoodsDao is the data access object for table seckill_goods.
type SeckillGoodsDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns SeckillGoodsColumns // columns contains all the column names of Table for convenient usage.
}

// SeckillGoodsColumns defines and stores column names for table seckill_goods.
type SeckillGoodsColumns struct {
	Id             string // 主键ID
	GoodsId        string // 商品ID
	GoodsOptionsId string // 商品规格ID
	OriginalPrice  string // 原始价格 单位分
	SeckillPrice   string // 秒杀价格 单位分
	SeckillStock   string // 秒杀库存
	StartTime      string // 秒杀开始时间
	EndTime        string // 秒杀结束时间
	Status         string // 状态：0-未开始 1-进行中 2-已结束
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// seckillGoodsColumns holds the columns for table seckill_goods.
var seckillGoodsColumns = SeckillGoodsColumns{
	Id:             "id",
	GoodsId:        "goods_id",
	GoodsOptionsId: "goods_options_id",
	OriginalPrice:  "original_price",
	SeckillPrice:   "seckill_price",
	SeckillStock:   "seckill_stock",
	StartTime:      "start_time",
	EndTime:        "end_time",
	Status:         "status",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewSeckillGoodsDao creates and returns a new DAO object for table data access.
func NewSeckillGoodsDao() *SeckillGoodsDao {
	return &SeckillGoodsDao{
		group:   "default",
		table:   "seckill_goods",
		columns: seckillGoodsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SeckillGoodsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SeckillGoodsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SeckillGoodsDao) Columns() SeckillGoodsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SeckillGoodsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SeckillGoodsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SeckillGoodsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
