package model

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// SeckillListInput 秒杀商品列表输入
type SeckillListInput struct {
	Page int `json:"page" description:"分页码"`
	Size int `json:"size" description:"分页数量"`
}

// SeckillListOutput 秒杀商品列表输出
type SeckillListOutput struct {
	List  interface{} `json:"list" description:"列表"`
	Page  int         `json:"page" description:"分页码"`
	Size  int         `json:"size" description:"分页数量"`
	Total int         `json:"total" description:"数据总数"`
}

// SeckillDetailInput 秒杀商品详情输入
type SeckillDetailInput struct {
	Id uint `json:"id" description:"秒杀商品ID"`
}

// SeckillDetailOutput 秒杀商品详情输出
type SeckillDetailOutput struct {
	Id             uint    `json:"id"            description:"秒杀商品ID"`
	GoodsId        uint    `json:"goods_id"       description:"商品ID"`
	GoodsOptionsId uint    `json:"goods_options_id" description:"商品规格ID"`
	OriginalPrice  float64 `json:"original_price" description:"原价"`
	SeckillPrice   float64 `json:"seckill_price"  description:"秒杀价"`
	SeckillStock   int     `json:"seckill_stock"  description:"秒杀库存"`
	StartTime      string  `json:"start_time"     description:"开始时间"`
	EndTime        string  `json:"end_time"       description:"结束时间"`
	Status         int     `json:"status"         description:"状态"`
}

// SeckillDoInput 秒杀请求输入
type SeckillDoInput struct {
	UserId         uint   `json:"userId" v:"required#用户ID不能为空"`               // 用户ID
	GoodsId        uint   `json:"goodsId" v:"required#商品ID不能为空"`              // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId" v:"required#商品选项ID不能为空"`     // 商品选项ID
	Count          uint   `json:"count" v:"required|min:1#购买数量不能为空|购买数量最小为1"` // 购买数量
	RequestId      string `json:"requestId" v:"required#请求ID不能为空"`            // 请求ID(幂等性)
	UserAddress    string `json:"userAddress"`                                // 收货地址
	UserPhone      string `json:"userPhone"`                                  // 手机号码
	Remark         string `json:"remark"`                                     // 备注
}

// SeckillDoOutput 秒杀请求输出
type SeckillDoOutput struct {
	RequestId    string    `json:"requestId"`    // 请求ID
	OrderNo      string    `json:"orderNo"`      // 订单号
	UserId       uint      `json:"userId"`       // 用户ID
	GoodsId      uint      `json:"goodsId"`      // 商品ID
	Count        uint      `json:"count"`        // 购买数量
	Status       int       `json:"status"`       // 秒杀状态(0:成功 其他:失败)
	Message      string    `json:"message"`      // 消息
	ProcessTime  int64     `json:"processTime"`  // 处理时间(毫秒)
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	IsProcessing bool      `json:"isProcessing"` // 是否正在处理
}

// SeckillResultInput 获取秒杀结果输入
type SeckillResultInput struct {
	OrderNo string `json:"order_no" description:"订单号"`
}

// SeckillResultOutput 获取秒杀结果输出
type SeckillResultOutput struct {
	OrderNo string `json:"order_no" description:"订单号"`
	Status  int    `json:"status"   description:"状态"`
}

// SeckillInitStockInput 初始化秒杀商品库存输入
type SeckillInitStockInput struct {
	GoodsId        uint `json:"goods_id"        description:"商品ID"`
	GoodsOptionsId uint `json:"goods_options_id" description:"商品规格ID"`
	Stock          int  `json:"stock"          description:"库存数量"`
}

// SeckillInitStockOutput 初始化秒杀商品库存输出
type SeckillInitStockOutput struct {
	Success bool `json:"success" description:"是否成功"`
}

// SeckillUpdateOrderStatusInput 更新订单状态输入
type SeckillUpdateOrderStatusInput struct {
	OrderNo string `json:"order_no"  description:"订单号"`
	Status  int    `json:"status"    description:"状态"`
}

// SeckillUpdateOrderStatusOutput 更新订单状态输出
type SeckillUpdateOrderStatusOutput struct {
	Success bool `json:"success" description:"是否成功"`
}

// SeckillStatsOutput 秒杀系统统计信息
type SeckillStatsOutput struct {
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

// SeckillOrderMsg 秒杀订单消息
type SeckillOrderMsg struct {
	OrderNo        string    `json:"orderNo"`        // 订单号
	RequestId      string    `json:"requestId"`      // 请求ID
	UserId         uint      `json:"userId"`         // 用户ID
	GoodsId        uint      `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint      `json:"goodsOptionsId"` // 商品选项ID
	Count          uint      `json:"count"`          // 购买数量
	TotalPrice     float64   `json:"totalPrice"`     // 总价
	Price          float64   `json:"price"`          // 单价
	Status         int       `json:"status"`         // 状态
	CreatedAt      time.Time `json:"createdAt"`      // 创建时间
	UserAddress    string    `json:"userAddress"`    // 收货地址
	UserPhone      string    `json:"userPhone"`      // 手机号码
	Remark         string    `json:"remark"`         // 备注
}

// SeckillGoodsInfo 秒杀商品信息
type SeckillGoodsInfo struct {
	Id             int64       `json:"id"`
	GoodsId        int64       `json:"goodsId"`
	GoodsOptionsId int64       `json:"goodsOptionsId"`
	Price          float64     `json:"price"`
	Stock          int         `json:"stock"`
	StartTime      *gtime.Time `json:"startTime"`
	EndTime        *gtime.Time `json:"endTime"`
	Status         int         `json:"status"`
}

// SeckillOrderCompleteMsg 秒杀订单完成消息
type SeckillOrderCompleteMsg struct {
	OrderNo        string      `json:"order_no"`         // 订单号
	UserId         int64       `json:"user_id"`          // 用户ID
	GoodsId        int64       `json:"goods_id"`         // 商品ID
	GoodsOptionsId int64       `json:"goods_options_id"` // 商品规格ID
	Status         int         `json:"status"`           // 订单状态
	CompletedAt    *gtime.Time `json:"completed_at"`     // 完成时间
}

// SeckillDetailReq 秒杀商品详情请求
type SeckillDetailReq struct {
	Id uint `json:"id" description:"秒杀商品ID"`
}

// SeckillDetailRes 秒杀商品详情响应
type SeckillDetailRes struct {
	Id             uint    `json:"id"            description:"秒杀商品ID"`
	GoodsId        uint    `json:"goods_id"       description:"商品ID"`
	GoodsOptionsId uint    `json:"goods_options_id" description:"商品规格ID"`
	OriginalPrice  float64 `json:"original_price" description:"原价"`
	SeckillPrice   float64 `json:"seckill_price"  description:"秒杀价"`
	SeckillStock   int     `json:"seckill_stock"  description:"秒杀库存"`
	StartTime      string  `json:"start_time"     description:"开始时间"`
	EndTime        string  `json:"end_time"       description:"结束时间"`
	Status         int     `json:"status"         description:"状态"`
}

// OrderNotification 订单通知
type OrderNotification struct {
	OrderId   uint   `json:"order_id"`  // 订单ID
	OrderNo   string `json:"order_no"`  // 订单编号
	UserId    uint   `json:"user_id"`   // 用户ID
	Status    string `json:"status"`    // 状态：success, failed
	Message   string `json:"message"`   // 消息内容
	Timestamp int64  `json:"timestamp"` // 时间戳
}

// SeckillRequest 秒杀请求
type SeckillRequest struct {
	GoodsId        uint `json:"goods_id" v:"required#商品ID不能为空"`           // 商品ID
	GoodsOptionsId uint `json:"goods_options_id" v:"required#商品规格ID不能为空"` // 商品规格ID
	Count          int  `json:"count" v:"required#商品数量不能为空"`              // 商品数量
}

// SeckillResponse 秒杀响应
type SeckillResponse struct {
	Success       bool   `json:"success"`        // 是否成功
	OrderId       uint   `json:"order_id"`       // 订单ID
	OrderNo       string `json:"order_no"`       // 订单编号
	Message       string `json:"message"`        // 消息
	QueuePosition int    `json:"queue_position"` // 队列位置（可选）
}

// StockCache 库存缓存
type StockCache struct {
	GoodsId uint `json:"goods_id"` // 商品ID
	Stock   int  `json:"stock"`    // 库存
	Version int  `json:"version"`  // 版本号，用于乐观锁
}

// SeckillStatistics 秒杀统计
type SeckillStatistics struct {
	TotalRequests   int           `json:"total_requests"`   // 总请求数
	SuccessRequests int           `json:"success_requests"` // 成功请求数
	FailedRequests  int           `json:"failed_requests"`  // 失败请求数
	TotalTime       time.Duration `json:"total_time"`       // 总时间
	AverageLatency  time.Duration `json:"average_latency"`  // 平均延迟
	MaxLatency      time.Duration `json:"max_latency"`      // 最大延迟
	MinLatency      time.Duration `json:"min_latency"`      // 最小延迟
	QPS             float64       `json:"qps"`              // 每秒请求数
	SuccessRate     float64       `json:"success_rate"`     // 成功率
}

// SeckillOrderAddInput 添加秒杀订单输入
type SeckillOrderAddInput struct {
	UserId         int64       `json:"userId"         description:"用户ID"`
	GoodsId        int64       `json:"goodsId"        description:"商品ID"`
	GoodsOptionsId int64       `json:"goodsOptionsId" description:"商品规格ID"`
	OriginalPrice  int         `json:"originalPrice"  description:"原始价格 单位分"`
	SeckillPrice   int         `json:"seckillPrice"   description:"秒杀价格 单位分"`
	Status         int         `json:"status"         description:"订单状态：0-待支付 1-已支付 2-已取消 3-已退款"`
	OrderNo        string      `json:"orderNo"        description:"订单编号"`
	Count          int         `json:"count"          description:"商品数量"`
	PayTime        *gtime.Time `json:"payTime"        description:"支付时间"`
	CancelTime     *gtime.Time `json:"cancelTime"     description:"取消时间"`
	ConsigneeName  string      `json:"consigneeName"  description:"收货人姓名"`
	ConsigneePhone string      `json:"consigneePhone" description:"收货人电话"`
	Address        string      `json:"address"        description:"收货地址"`
	Remark         string      `json:"remark"         description:"备注"`
}

// SeckillOrderAddOutput 添加秒杀订单输出
type SeckillOrderAddOutput struct {
	Id int64 `json:"id" description:"主键ID"`
}

// SeckillOrderUpdateInput 更新秒杀订单输入
type SeckillOrderUpdateInput struct {
	Id     int64 `json:"id"     v:"required#秒杀订单ID不能为空" description:"秒杀订单ID"`
	Status int   `json:"status" description:"状态：0-待支付 1-已支付 2-已取消"`
}

// SeckillOrderUpdateOutput 更新秒杀订单输出
type SeckillOrderUpdateOutput struct {
	Id int64 `json:"id" description:"秒杀订单ID"`
}

// SeckillOrderByOrderIdInput 根据订单ID查询秒杀订单输入
type SeckillOrderByOrderIdInput struct {
	OrderId int64 `json:"orderId" v:"required#订单ID不能为空" description:"订单ID"`
}

// SeckillOrderByOrderIdOutput 根据订单ID查询秒杀订单输出
type SeckillOrderByOrderIdOutput struct {
	Id             int64       `json:"id"             description:"主键ID"`
	OrderId        int64       `json:"orderId"        description:"订单ID"`
	UserId         int64       `json:"userId"         description:"用户ID"`
	GoodsId        int64       `json:"goodsId"        description:"商品ID"`
	GoodsOptionsId int64       `json:"goodsOptionsId" description:"商品规格ID"`
	SeckillPrice   int         `json:"seckillPrice"   description:"秒杀价格 单位分"`
	Status         int         `json:"status"         description:"状态：0-待支付 1-已支付 2-已取消"`
	CreatedAt      *gtime.Time `json:"createdAt"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      description:"更新时间"`
}

// SeckillOrderMsg 将秒杀订单消息转换为创建秒杀订单输入
func (msg *SeckillOrderMsg) ToSeckillOrderAddInput() *SeckillOrderAddInput {
	return &SeckillOrderAddInput{
		UserId:         int64(msg.UserId),
		GoodsId:        int64(msg.GoodsId),
		GoodsOptionsId: int64(msg.GoodsOptionsId),
		SeckillPrice:   int(msg.Price * 100),                           // 转换为分
		OriginalPrice:  int(msg.TotalPrice * 100 / float64(msg.Count)), // 转换为分
		Status:         1,                                              // 默认已支付状态
		OrderNo:        msg.OrderNo,
		Count:          int(msg.Count),
		PayTime:        gtime.Now(),
		ConsigneeName:  fmt.Sprintf("秒杀用户%d", msg.UserId),
		ConsigneePhone: msg.UserPhone,
		Address:        msg.UserAddress,
		Remark:         msg.Remark,
	}
}

// SeckillCheckStockInput 检查库存输入
type SeckillCheckStockInput struct {
	GoodsId        uint `json:"goodsId" v:"required#商品ID不能为空"`              // 商品ID
	GoodsOptionsId uint `json:"goodsOptionsId" v:"required#商品选项ID不能为空"`     // 商品选项ID
	Count          uint `json:"count" v:"required|min:1#购买数量不能为空|购买数量最小为1"` // 购买数量
}

// SeckillCheckStockOutput 检查库存输出
type SeckillCheckStockOutput struct {
	Available bool  `json:"available"` // 库存是否充足
	Current   int32 `json:"current"`   // 当前库存
	Required  int32 `json:"required"`  // 需要的库存
}

// SeckillConfig 秒杀配置
type SeckillConfig struct {
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

// SeckillGoodsAddInput 添加秒杀商品输入
type SeckillGoodsAddInput struct {
	GoodsId        int64       `json:"goodsId"        v:"required#商品ID不能为空" description:"商品ID"`
	GoodsOptionsId int64       `json:"goodsOptionsId" v:"required#商品规格ID不能为空" description:"商品规格ID"`
	OriginalPrice  int         `json:"originalPrice"  v:"required#原始价格不能为空" description:"原始价格 单位分"`
	SeckillPrice   int         `json:"seckillPrice"   v:"required#秒杀价格不能为空" description:"秒杀价格 单位分"`
	SeckillStock   int         `json:"seckillStock"   v:"required#秒杀库存不能为空" description:"秒杀库存"`
	StartTime      *gtime.Time `json:"startTime"      v:"required#开始时间不能为空" description:"秒杀开始时间"`
	EndTime        *gtime.Time `json:"endTime"        v:"required#结束时间不能为空" description:"秒杀结束时间"`
	Status         int         `json:"status"         description:"状态：0-未开始 1-进行中 2-已结束"`
}

// SeckillGoodsAddOutput 添加秒杀商品输出
type SeckillGoodsAddOutput struct {
	Id int64 `json:"id" description:"秒杀商品ID"`
}

// SeckillGoodsUpdateInput 更新秒杀商品输入
type SeckillGoodsUpdateInput struct {
	Id            int64       `json:"id"             v:"required#秒杀商品ID不能为空" description:"秒杀商品ID"`
	OriginalPrice int         `json:"originalPrice"  description:"原始价格 单位分"`
	SeckillPrice  int         `json:"seckillPrice"   description:"秒杀价格 单位分"`
	SeckillStock  int         `json:"seckillStock"   description:"秒杀库存"`
	StartTime     *gtime.Time `json:"startTime"      description:"秒杀开始时间"`
	EndTime       *gtime.Time `json:"endTime"        description:"秒杀结束时间"`
	Status        int         `json:"status"         description:"状态：0-未开始 1-进行中 2-已结束"`
}

// SeckillGoodsUpdateOutput 更新秒杀商品输出
type SeckillGoodsUpdateOutput struct {
	Id int64 `json:"id" description:"秒杀商品ID"`
}
