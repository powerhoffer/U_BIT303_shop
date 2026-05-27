package backend

import (
	"bit303_shop/internal/model/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DoReq 秒杀请求
type DoReq struct {
	g.Meta         `path:"/seckill/do" method:"post" tags:"秒杀管理" summary:"执行秒杀(调试用)"`
	UserId         uint   `json:"userId" v:"required#用户ID不能为空"`               // 用户ID
	GoodsId        uint   `json:"goodsId" v:"required#商品ID不能为空"`              // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId" v:"required#商品选项ID不能为空"`     // 商品选项ID
	Count          uint   `json:"count" v:"required|min:1#购买数量不能为空|购买数量最小为1"` // 购买数量
	RequestId      string `json:"requestId" v:"required#请求ID不能为空"`            // 请求ID(幂等性)
	UserAddress    string `json:"userAddress"`                                // 收货地址
	UserPhone      string `json:"userPhone"`                                  // 手机号码
	Remark         string `json:"remark"`                                     // 备注
}

// DoRes 秒杀响应
type DoRes struct {
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

// ToModel 将请求转换为模型
func (req *DoReq) ToModel() *SeckillDoInput {
	return &SeckillDoInput{
		UserId:         req.UserId,
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Count:          req.Count,
		RequestId:      req.RequestId,
		UserAddress:    req.UserAddress,
		UserPhone:      req.UserPhone,
		Remark:         req.Remark,
	}
}

// FromModel 从模型转换为响应
func (res *DoRes) FromModel(output *SeckillDoOutput) {
	if output == nil {
		return
	}
	res.RequestId = output.RequestId
	res.OrderNo = output.OrderNo
	res.UserId = output.UserId
	res.GoodsId = output.GoodsId
	res.Count = output.Count
	res.Status = output.Status
	res.Message = output.Message
	res.ProcessTime = output.ProcessTime
	res.IsProcessing = output.IsProcessing
}

// CheckStockReq 检查库存请求
type CheckStockReq struct {
	g.Meta         `path:"/seckill/check-stock" method:"post" tags:"秒杀管理" summary:"检查库存"`
	GoodsId        uint `json:"goodsId" v:"required#商品ID不能为空"`              // 商品ID
	GoodsOptionsId uint `json:"goodsOptionsId" v:"required#商品选项ID不能为空"`     // 商品选项ID
	Count          uint `json:"count" v:"required|min:1#购买数量不能为空|购买数量最小为1"` // 购买数量
}

// CheckStockRes 检查库存响应
type CheckStockRes struct {
	Available bool  `json:"available"` // 库存是否充足
	Current   int32 `json:"current"`   // 当前库存
	Required  int32 `json:"required"`  // 需要的库存
}

// ToModel 将请求转换为模型
func (req *CheckStockReq) ToModel() *SeckillCheckStockInput {
	return &SeckillCheckStockInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Count:          req.Count,
	}
}

// FromModel 从模型转换为响应
func (res *CheckStockRes) FromModel(output *SeckillCheckStockOutput) {
	if output == nil {
		return
	}
	res.Available = output.Available
	res.Current = output.Current
	res.Required = output.Required
}

// InitStockReq 初始化库存请求
type InitStockReq struct {
	g.Meta         `path:"/seckill/init-stock" method:"post" tags:"秒杀管理" summary:"初始化库存"`
	GoodsId        uint  `json:"goodsId" v:"required#商品ID不能为空"`          // 商品ID
	GoodsOptionsId uint  `json:"goodsOptionsId" v:"required#商品选项ID不能为空"` // 商品选项ID
	Stock          int32 `json:"stock" v:"required|min:1#库存不能为空|库存最小为1"` // 库存
}

// InitStockRes 初始化库存响应
type InitStockRes struct {
	Success bool `json:"success"` // 是否成功
}

// ToModel 将请求转换为模型
func (req *InitStockReq) ToModel() *SeckillInitStockInput {
	return &SeckillInitStockInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Stock:          int(req.Stock),
	}
}

// GetStatsReq 获取统计信息请求
type GetStatsReq struct {
	g.Meta         `path:"/seckill/stats" method:"get" tags:"秒杀管理" summary:"获取统计信息"`
	GoodsId        int64 `json:"goodsId" d:"0"`        // 商品ID
	GoodsOptionsId int64 `json:"goodsOptionsId" d:"0"` // 商品选项ID
}

// GetStatsRes 获取统计信息响应
type GetStatsRes struct {
	Workers        int              `json:"workers"`        // 工作线程数
	QueueSize      int              `json:"queueSize"`      // 队列大小
	QueueCurrent   int              `json:"queueCurrent"`   // 当前队列中的请求数
	Successes      int64            `json:"successes"`      // 成功请求数
	Failures       int64            `json:"failures"`       // 失败请求数
	Errors         int64            `json:"errors"`         // 错误请求数
	AvgTime        float64          `json:"avgTime"`        // 平均处理时间
	TokenBucket    map[string]int64 `json:"tokenBucket"`    // 令牌桶统计
	LeakyBucket    map[string]int64 `json:"leakyBucket"`    // 漏桶统计
	StockManager   map[string]int64 `json:"stockManager"`   // 库存管理器统计
	CircuitBreaker map[string]int64 `json:"circuitBreaker"` // 熔断器统计
}

// FromModel 从模型转换为响应
func (res *GetStatsRes) FromModel(stats *SeckillStatsOutput) {
	if stats == nil {
		return
	}
	res.Workers = stats.Workers
	res.QueueSize = stats.QueueSize
	res.QueueCurrent = stats.QueueCurrent
	res.Successes = stats.Successes
	res.Failures = stats.Failures
	res.Errors = stats.Errors
	res.AvgTime = stats.AvgTime
	res.TokenBucket = stats.TokenBucket
	res.LeakyBucket = stats.LeakyBucket
	res.StockManager = stats.StockManager
	res.CircuitBreaker = stats.CircuitBreaker
}

// ResetReq 重置系统请求
type ResetReq struct {
	g.Meta `path:"/seckill/reset" method:"post" tags:"秒杀管理" summary:"重置系统"`
}

// ResetRes 重置系统响应
type ResetRes struct {
	Success bool `json:"success"` // 是否成功
}

// SetConfigReq 设置配置请求
type SetConfigReq struct {
	g.Meta                  `path:"/seckill/config" method:"post" tags:"秒杀管理" summary:"设置配置"`
	EnableKafka             bool  `json:"enableKafka"`             // 是否启用Kafka
	EnableTokenBucket       bool  `json:"enableTokenBucket"`       // 是否启用令牌桶限流
	EnableLeakyBucket       bool  `json:"enableLeakyBucket"`       // 是否启用漏桶限流
	EnableCircuitBreaker    bool  `json:"enableCircuitBreaker"`    // 是否启用熔断器
	WorkerCount             int   `json:"workerCount"`             // 工作线程数
	QueueSize               int   `json:"queueSize"`               // 队列大小
	TokenBucketSize         int32 `json:"tokenBucketSize"`         // 令牌桶大小
	TokenRate               int32 `json:"tokenRate"`               // 令牌发放速率
	LeakyBucketSize         int32 `json:"leakyBucketSize"`         // 漏桶大小
	LeakyRate               int32 `json:"leakyRate"`               // 漏桶处理速率
	CircuitBreakerThreshold int32 `json:"circuitBreakerThreshold"` // 熔断器阈值
}

// SetConfigRes 设置配置响应
type SetConfigRes struct {
	Success bool `json:"success"` // 是否成功
}

// ToModel 将请求转换为模型
func (req *SetConfigReq) ToModel() *SeckillConfig {
	return &SeckillConfig{
		EnableKafka:             req.EnableKafka,
		EnableTokenBucket:       req.EnableTokenBucket,
		EnableLeakyBucket:       req.EnableLeakyBucket,
		EnableCircuitBreaker:    req.EnableCircuitBreaker,
		WorkerCount:             req.WorkerCount,
		QueueSize:               req.QueueSize,
		TokenBucketSize:         req.TokenBucketSize,
		TokenRate:               req.TokenRate,
		LeakyBucketSize:         req.LeakyBucketSize,
		LeakyRate:               req.LeakyRate,
		CircuitBreakerThreshold: req.CircuitBreakerThreshold,
	}
}

// SeckillCheckStockInput 检查库存输入
type SeckillCheckStockInput struct {
	GoodsId        uint `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint `json:"goodsOptionsId"` // 商品选项ID
	Count          uint `json:"count"`          // 购买数量
}

// SeckillCheckStockOutput 检查库存输出
type SeckillCheckStockOutput struct {
	Available bool  `json:"available"` // 库存是否充足
	Current   int32 `json:"current"`   // 当前库存
	Required  int32 `json:"required"`  // 需要的库存
}

// InitializeReq 初始化系统请求
type InitializeReq struct {
	g.Meta `path:"/seckill/initialize" method:"post" tags:"秒杀管理" summary:"初始化秒杀系统"`
}

// InitializeRes 初始化系统响应
type InitializeRes struct {
	Success bool `json:"success"` // 是否成功
}

// AddGoodsReq 添加秒杀商品请求
type AddGoodsReq struct {
	g.Meta         `path:"/seckill/goods/add" method:"post" tags:"秒杀商品" summary:"添加秒杀商品"`
	GoodsId        int64       `json:"goodsId"        v:"required#商品ID不能为空" description:"商品ID"`
	GoodsOptionsId int64       `json:"goodsOptionsId" v:"required#商品规格ID不能为空" description:"商品规格ID"`
	OriginalPrice  int         `json:"originalPrice"  v:"required#原始价格不能为空" description:"原始价格 单位分"`
	SeckillPrice   int         `json:"seckillPrice"   v:"required#秒杀价格不能为空" description:"秒杀价格 单位分"`
	SeckillStock   int         `json:"seckillStock"   v:"required#秒杀库存不能为空" description:"秒杀库存"`
	StartTime      *gtime.Time `json:"startTime"      v:"required#开始时间不能为空" description:"秒杀开始时间"`
	EndTime        *gtime.Time `json:"endTime"        v:"required#结束时间不能为空" description:"秒杀结束时间"`
}

// AddGoodsRes 添加秒杀商品响应
type AddGoodsRes struct {
	Id int64 `json:"id" description:"秒杀商品ID"`
}

// UpdateGoodsStatusReq 更新秒杀商品状态请求
type UpdateGoodsStatusReq struct {
	g.Meta `path:"/seckill/goods/status" method:"post" tags:"秒杀商品" summary:"更新秒杀商品状态"`
	Id     int64 `json:"id"     v:"required#秒杀商品ID不能为空" description:"秒杀商品ID"`
	Status int   `json:"status" v:"required#状态不能为空" description:"状态：0-未开始 1-进行中 2-已结束"`
}

// UpdateGoodsStatusRes 更新秒杀商品状态响应
type UpdateGoodsStatusRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// UpdateOrderStatusReq 更新秒杀订单状态请求
type UpdateOrderStatusReq struct {
	g.Meta  `path:"/seckill/order/status" method:"post" tags:"秒杀订单" summary:"更新秒杀订单状态"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空" description:"订单号"`
	Status  int    `json:"status"  v:"required#状态不能为空" description:"状态：0-待支付 1-已支付 2-已取消"`
}

// UpdateOrderStatusRes 更新秒杀订单状态响应
type UpdateOrderStatusRes struct {
	Success bool `json:"success" description:"是否成功"`
}

// ListGoodsReq 秒杀商品列表请求
type ListGoodsReq struct {
	g.Meta `path:"/seckill/goods/list" method:"get" tags:"秒杀商品" summary:"获取秒杀商品列表"`
	Page   int `json:"page" v:"required#页码不能为空" description:"页码"`
	Size   int `json:"size" v:"required#每页数量不能为空" description:"每页数量"`
	Status int `json:"status" description:"状态：0-未开始 1-进行中 2-已结束"`
}

// ListGoodsRes 秒杀商品列表响应
type ListGoodsRes struct {
	List  []*entity.SeckillGoods `json:"list" description:"秒杀商品列表"`
	Page  int                    `json:"page" description:"页码"`
	Size  int                    `json:"size" description:"每页数量"`
	Total int                    `json:"total" description:"总数量"`
}

// ListOrdersReq 秒杀订单列表请求
type ListOrdersReq struct {
	g.Meta  `path:"/seckill/order/list" method:"get" tags:"秒杀订单" summary:"获取秒杀订单列表"`
	Page    int   `json:"page" v:"required#页码不能为空" description:"页码"`
	Size    int   `json:"size" v:"required#每页数量不能为空" description:"每页数量"`
	UserId  int64 `json:"userId" description:"用户ID"`
	GoodsId int64 `json:"goodsId" description:"商品ID"`
	Status  int   `json:"status" description:"状态：0-待支付 1-已支付 2-已取消"`
}

// ListOrdersRes 秒杀订单列表响应
type ListOrdersRes struct {
	List  []*entity.SeckillOrder `json:"list" description:"秒杀订单列表"`
	Page  int                    `json:"page" description:"页码"`
	Size  int                    `json:"size" description:"每页数量"`
	Total int                    `json:"total" description:"总数量"`
}
