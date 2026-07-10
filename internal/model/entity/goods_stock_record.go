// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GoodsStockRecord is the golang structure for table goods_stock_record.
type GoodsStockRecord struct {
	Id             uint        `json:"id"             orm:"id"              ` // 库存流水ID
	GoodsId        uint        `json:"goodsId"        orm:"goods_id"        ` // 商品ID
	GoodsName      string      `json:"goodsName"      orm:"goods_name"      ` // 商品名称快照
	ChangeType     int         `json:"changeType"     orm:"change_type"     ` // 变动类型：1初始库存 2后台增加 3后台扣减 4订单扣减 5取消订单恢复
	ChangeQuantity int         `json:"changeQuantity" orm:"change_quantity" ` // 库存变动数量，正数增加，负数扣减
	BeforeStock    uint        `json:"beforeStock"    orm:"before_stock"    ` // 变动前库存
	AfterStock     uint        `json:"afterStock"     orm:"after_stock"     ` // 变动后库存
	BizType        string      `json:"bizType"        orm:"biz_type"        ` // 业务来源
	BizId          uint        `json:"bizId"          orm:"biz_id"          ` // 关联业务ID
	OperatorType   int         `json:"operatorType"   orm:"operator_type"   ` // 操作者类型：0系统 1管理员 2员工
	OperatorId     uint        `json:"operatorId"     orm:"operator_id"     ` // 操作者ID
	Remark         string      `json:"remark"         orm:"remark"          ` // 备注
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      ` // 创建时间
}
