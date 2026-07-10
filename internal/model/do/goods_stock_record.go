// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsStockRecord is the golang structure of table goods_stock_record for DAO operations like Where/Data.
type GoodsStockRecord struct {
	g.Meta         `orm:"table:goods_stock_record, do:true"`
	Id             any         // 库存流水ID
	GoodsId        any         // 商品ID
	GoodsName      any         // 商品名称快照
	ChangeType     any         // 变动类型：1初始库存 2后台增加 3后台扣减 4订单扣减 5取消订单恢复
	ChangeQuantity any         // 库存变动数量，正数增加，负数扣减
	BeforeStock    any         // 变动前库存
	AfterStock     any         // 变动后库存
	BizType        any         // 业务来源
	BizId          any         // 关联业务ID
	OperatorType   any         // 操作者类型：0系统 1管理员 2员工
	OperatorId     any         // 操作者ID
	Remark         any         // 备注
	CreatedAt      *gtime.Time // 创建时间
}
