package frontend

import (
	"bit303_shop/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
)

type SeckillListReq struct {
	g.Meta `path:"/seckill/list" method:"post" tags:"前台秒杀" summary:"秒杀商品列表接口"`
	CommonPaginationReq
}

type SeckillListRes struct {
	List  interface{} `json:"list" description:"列表"`
	Page  int         `json:"page" description:"分页码"`
	Size  int         `json:"size" description:"分页数量"`
	Total int         `json:"total" description:"数据总数"`
}

type SeckillDetailReq struct {
	g.Meta `path:"/seckill/goods/detail" method:"get" tags:"秒杀" summary:"获取秒杀商品详情"`
	Id     uint `json:"id" v:"required#商品ID不能为空"` // 秒杀商品ID
}

type SeckillDetailRes struct {
	Id             uint    `json:"id"`             // 秒杀商品ID
	GoodsId        uint    `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint    `json:"goodsOptionsId"` // 商品规格ID
	OriginalPrice  float64 `json:"originalPrice"`  // 原价
	SeckillPrice   float64 `json:"seckillPrice"`   // 秒杀价
	SeckillStock   int     `json:"seckillStock"`   // 秒杀库存
	StartTime      string  `json:"startTime"`      // 开始时间
	EndTime        string  `json:"endTime"`        // 结束时间
	Status         int     `json:"status"`         // 状态
}

type SeckillDoReq struct {
	g.Meta         `path:"/seckill/do" method:"post" tags:"前台秒杀" summary:"执行秒杀"`
	GoodsId        uint `json:"goods_id"        description:"商品ID"`
	GoodsOptionsId uint `json:"goods_options_id" description:"商品规格ID"`
	UserId         uint `json:"user_id"         description:"用户ID"`
}

type SeckillDoRes struct {
	OrderNo string `json:"order_no" description:"订单号"`
}

type SeckillResultReq struct {
	g.Meta  `path:"/seckill/result" method:"post" tags:"前台秒杀" summary:"获取秒杀结果"`
	OrderNo string `json:"order_no" description:"订单号"`
}

type SeckillResultRes struct {
	OrderNo string `json:"order_no" description:"订单号"`
	Status  int    `json:"status"   description:"状态"`
}

type ListReq struct {
	g.Meta `path:"/seckill/list" method:"get" tags:"秒杀" summary:"秒杀商品列表"`
	Page   int `json:"page" v:"required#请输入页码" dc:"页码"`
	Size   int `json:"size" v:"required#请输入每页数量" dc:"每页数量"`
}

type ListRes struct {
	List  []*entity.SeckillGoods `json:"list" dc:"秒杀商品列表"`
	Total int                    `json:"total" dc:"总数"`
}

type DetailReq struct {
	g.Meta `path:"/seckill/detail" method:"get" tags:"秒杀" summary:"秒杀商品详情"`
	Id     int64 `json:"id" v:"required#请输入商品ID" dc:"商品ID"`
}

type DetailRes struct {
	Goods *entity.SeckillGoods `json:"goods" dc:"秒杀商品信息"`
}

type InitStockReq struct {
	g.Meta         `path:"/seckill/init-stock" method:"post" tags:"秒杀" summary:"初始化秒杀商品库存"`
	GoodsId        int64 `json:"goodsId" v:"required#请输入商品ID" dc:"商品ID"`
	GoodsOptionsId int64 `json:"goodsOptionsId" v:"required#请输入商品规格ID" dc:"商品规格ID"`
	Stock          int   `json:"stock" v:"required#请输入库存数量" dc:"库存数量"`
}

type InitStockRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

// UpdateOrderStatusReq 更新订单状态请求
type UpdateOrderStatusReq struct {
	g.Meta  `path:"/order/status" method:"post" tags:"订单" summary:"更新订单状态"`
	OrderNo string `json:"order_no" v:"required#请输入订单号" description:"订单号"`
	Status  int    `json:"status"  v:"required#请输入订单状态" description:"订单状态：0-待支付 1-已支付 2-已取消"`
}

// UpdateOrderStatusRes 更新订单状态响应
type UpdateOrderStatusRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// GetSeckillResultReq 获取秒杀结果请求
type GetSeckillResultReq struct {
	g.Meta  `path:"/seckill/result" method:"get" tags:"秒杀" summary:"获取秒杀结果"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空" dc:"订单号"`
}

// GetSeckillResultRes 获取秒杀结果响应
type GetSeckillResultRes struct {
	OrderNo string `json:"orderNo" dc:"订单号"`
	Status  int    `json:"status" dc:"订单状态：0-待支付 1-已支付 2-已取消"`
}

// GetSeckillListReq 获取秒杀商品列表请求
type GetSeckillListReq struct {
	g.Meta `path:"/seckill/list" method:"get" tags:"秒杀" summary:"获取秒杀商品列表"`
	Page   int `json:"page" v:"required#页码不能为空" dc:"页码"`
	Size   int `json:"size" v:"required#每页数量不能为空" dc:"每页数量"`
}

// GetSeckillListRes 获取秒杀商品列表响应
type GetSeckillListRes struct {
	List  []SeckillGoodsInfo `json:"list" dc:"秒杀商品列表"`
	Page  int                `json:"page" dc:"页码"`
	Size  int                `json:"size" dc:"每页数量"`
	Total int                `json:"total" dc:"总数"`
}

// SeckillGoodsInfo 秒杀商品信息
type SeckillGoodsInfo struct {
	Id             int64   `json:"id" dc:"秒杀商品ID"`
	GoodsId        int64   `json:"goodsId" dc:"商品ID"`
	GoodsOptionsId int64   `json:"goodsOptionsId" dc:"商品规格ID"`
	Price          float64 `json:"price" dc:"秒杀价格"`
	Stock          int     `json:"stock" dc:"库存"`
	StartTime      string  `json:"startTime" dc:"开始时间"`
	EndTime        string  `json:"endTime" dc:"结束时间"`
	Status         int     `json:"status" dc:"状态：0-未开始 1-进行中 2-已结束"`
}

// SeckillReq 秒杀请求
type SeckillReq struct {
	g.Meta         `path:"/seckill/do" method:"post" tags:"秒杀" summary:"执行秒杀"`
	UserId         uint   `json:"userId" v:"required#用户ID不能为空"`               // 用户ID
	GoodsId        uint   `json:"goodsId" v:"required#商品ID不能为空"`              // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId" v:"required#商品选项ID不能为空"`     // 商品选项ID
	Count          uint   `json:"count" v:"required|min:1#购买数量不能为空|购买数量最小为1"` // 购买数量
	RequestId      string `json:"requestId" v:"required#请求ID不能为空"`            // 请求ID(幂等性)
	UserAddress    string `json:"userAddress"`                                // 收货地址
	UserPhone      string `json:"userPhone"`                                  // 手机号码
	Remark         string `json:"remark"`                                     // 备注
}

// SeckillRes 秒杀响应
type SeckillRes struct {
	RequestId    string `json:"requestId"`    // 请求ID
	OrderNo      string `json:"orderNo"`      // 订单号
	UserId       uint   `json:"userId"`       // 用户ID
	GoodsId      uint   `json:"goodsId"`      // 商品ID
	Count        uint   `json:"count"`        // 购买数量
	Status       int    `json:"status"`       // 秒杀状态(0:成功 其他:失败)
	Message      string `json:"message"`      // 消息
	ProcessTime  int64  `json:"processTime"`  // 处理时间(毫秒)
	IsProcessing bool   `json:"isProcessing"` // 是否正在处理
}

// SeckillGoodsListReq 秒杀商品列表请求
type SeckillGoodsListReq struct {
	g.Meta   `path:"/seckill/goods/list" method:"get" tags:"秒杀" summary:"获取秒杀商品列表"`
	Page     int `json:"page" v:"min:1#页码最小为1" d:"1"`                  // 页码
	PageSize int `json:"pageSize" v:"min:1|max:50#每页数量在1-50之间" d:"10"` // 每页数量
}

// SeckillGoodsListRes 秒杀商品列表响应
type SeckillGoodsListRes struct {
	List  []SeckillGoodsInfo `json:"list"`  // 商品列表
	Total int                `json:"total"` // 总数
	Page  int                `json:"page"`  // 页码
	Size  int                `json:"size"`  // 每页数量
}

// SeckillStatsReq 获取秒杀统计请求
type SeckillStatsReq struct {
	g.Meta `path:"/seckill/stats" method:"post" tags:"秒杀" summary:"获取秒杀统计"`
}

// SeckillStatsRes 获取秒杀统计响应
type SeckillStatsRes struct {
	TotalRequests   int     `json:"totalRequests" dc:"总请求数"`
	SuccessRequests int     `json:"successRequests" dc:"成功请求数"`
	FailedRequests  int     `json:"failedRequests" dc:"失败请求数"`
	AverageLatency  string  `json:"averageLatency" dc:"平均延迟"`
	MaxLatency      string  `json:"maxLatency" dc:"最大延迟"`
	MinLatency      string  `json:"minLatency" dc:"最小延迟"`
	QPS             float64 `json:"qps" dc:"每秒请求数"`
	SuccessRate     float64 `json:"successRate" dc:"成功率"`
}

// SeckillBatchReq 批量秒杀请求
type SeckillBatchReq struct {
	g.Meta `path:"/seckill/batch" method:"post" tags:"秒杀" summary:"批量秒杀请求API"`
	Items  []*SeckillBatchItem `json:"items" v:"required|min-length:1#请至少提交一个秒杀商品" dc:"批量秒杀商品列表"`
}

// SeckillBatchItem 批量秒杀项
type SeckillBatchItem struct {
	GoodsId        int `json:"goodsId" v:"required#商品ID不能为空" dc:"商品ID"`
	GoodsOptionsId int `json:"goodsOptionsId" dc:"商品规格ID"`
	Count          int `json:"count" v:"required|min:1|max:5#购买数量必须在1-5之间" dc:"购买数量"`
}

// SeckillBatchRes 批量秒杀响应
type SeckillBatchRes struct {
	SuccessCount int                  `json:"successCount" dc:"成功数量"`
	FailedCount  int                  `json:"failedCount" dc:"失败数量"`
	Results      []*SeckillItemResult `json:"results" dc:"处理结果"`
}

// SeckillItemResult 单个秒杀项处理结果
type SeckillItemResult struct {
	GoodsId int    `json:"goodsId" dc:"商品ID"`
	Success bool   `json:"success" dc:"是否成功"`
	OrderNo string `json:"orderNo" dc:"订单号，成功时返回"`
	Message string `json:"message" dc:"失败原因，失败时返回"`
}

// GetSeckillStatusReq 获取秒杀状态请求
type GetSeckillStatusReq struct {
	g.Meta  `path:"/seckill/status" method:"get" tags:"秒杀" summary:"获取秒杀状态"`
	OrderId uint `json:"orderId" v:"required#订单ID不能为空"` // 订单ID
}

// GetSeckillStatusRes 获取秒杀状态响应
type GetSeckillStatusRes struct {
	Status string `json:"status"` // 状态: success, failed, pending
	Number string `json:"number"` // 订单编号
	Reason string `json:"reason"` // 原因
}

// DirectSeckillReq 直接执行秒杀请求
type DirectSeckillReq struct {
	g.Meta         `path:"/seckill/direct" method:"post" tags:"秒杀" summary:"直接执行秒杀并返回结果"`
	UserId         uint   `json:"userId" v:"required#用户ID不能为空"`               // 用户ID
	GoodsId        uint   `json:"goodsId" v:"required#商品ID不能为空"`              // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId" v:"required#商品选项ID不能为空"`     // 商品选项ID
	Count          uint   `json:"count" v:"required|min:1#购买数量不能为空|购买数量最小为1"` // 购买数量
	UserAddress    string `json:"userAddress"`                                // 收货地址
	UserPhone      string `json:"userPhone"`                                  // 手机号码
	Remark         string `json:"remark"`                                     // 备注
}

// DirectSeckillRes 直接执行秒杀响应
type DirectSeckillRes struct {
	Success    bool   `json:"success"`    // 是否成功
	OrderId    uint   `json:"orderId"`    // 订单ID
	OrderNo    string `json:"orderNo"`    // 订单编号
	Message    string `json:"message"`    // 消息
	TotalPrice int    `json:"totalPrice"` // 总价(分)
}

// SeckillOrderDetailReq 秒杀订单详情请求
type SeckillOrderDetailReq struct {
	g.Meta  `path:"/seckill/order/detail" method:"get" tags:"秒杀" summary:"获取秒杀订单详情"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空"` // 订单号
}

// SeckillOrderGoodsDetail 秒杀订单商品详情
type SeckillOrderGoodsDetail struct {
	GoodsId        uint   `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId"` // 商品规格ID
	Count          int    `json:"count"`          // 数量
	Price          int    `json:"price"`          // 价格(分)
	ActualPrice    int    `json:"actualPrice"`    // 实际价格(分)
	GoodsName      string `json:"goodsName"`      // 商品名称
	GoodsImage     string `json:"goodsImage"`     // 商品图片
	GoodsSpec      string `json:"goodsSpec"`      // 商品规格
}

// SeckillOrderDetailRes 秒杀订单详情响应
type SeckillOrderDetailRes struct {
	Id             uint                      `json:"id"`             // 订单ID
	OrderNo        string                    `json:"orderNo"`        // 订单号
	UserId         uint                      `json:"userId"`         // 用户ID
	Status         int                       `json:"status"`         // 订单状态
	Price          int                       `json:"price"`          // 总价(分)
	PayAt          *string                   `json:"payAt"`          // 支付时间
	CancelAt       *string                   `json:"cancelAt"`       // 取消时间
	ConsigneeName  string                    `json:"consigneeName"`  // 收货人姓名
	ConsigneePhone string                    `json:"consigneePhone"` // 收货人电话
	Address        string                    `json:"address"`        // 收货地址
	Remark         string                    `json:"remark"`         // 备注
	GoodsList      []SeckillOrderGoodsDetail `json:"goodsList"`      // 商品详情列表
	CreatedAt      string                    `json:"createdAt"`      // 创建时间
}

// SeckillGoodsInfoReq 获取秒杀商品信息和库存请求
type SeckillGoodsInfoReq struct {
	g.Meta         `path:"/seckill/goods/info" method:"get" tags:"秒杀" summary:"获取秒杀商品信息和库存"`
	GoodsId        uint `json:"goodsId" v:"required#商品ID不能为空"`          // 商品ID
	GoodsOptionsId uint `json:"goodsOptionsId" v:"required#商品规格ID不能为空"` // 商品规格ID
}

// SeckillGoodsInfoRes 获取秒杀商品信息和库存响应
type SeckillGoodsInfoRes struct {
	GoodsId        uint   `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId"` // 商品规格ID
	GoodsName      string `json:"goodsName"`      // 商品名称
	GoodsImage     string `json:"goodsImage"`     // 商品图片
	GoodsSpec      string `json:"goodsSpec"`      // 商品规格
	Price          int    `json:"price"`          // 原价(分)
	SeckillPrice   int    `json:"seckillPrice"`   // 秒杀价(分)
	Stock          int    `json:"stock"`          // 库存
	SalesCount     int    `json:"salesCount"`     // 销量
	LimitCount     int    `json:"limitCount"`     // 限购数量
	StartTime      string `json:"startTime"`      // 开始时间
	EndTime        string `json:"endTime"`        // 结束时间
	Status         int    `json:"status"`         // 状态：0-未开始 1-进行中 2-已结束
	Description    string `json:"description"`    // 商品描述
}
