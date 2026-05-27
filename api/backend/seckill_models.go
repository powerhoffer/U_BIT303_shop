package backend

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SeckillInitStockInput 初始化库存输入
type SeckillInitStockInput struct {
	GoodsId        uint `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint `json:"goodsOptionsId"` // 商品选项ID
	Stock          int  `json:"stock"`          // 库存数量
}

// SeckillInitStockOutput 初始化库存输出
type SeckillInitStockOutput struct {
	Success bool `json:"success"` // 是否成功
}

// SeckillStatsInput 秒杀统计信息输入
type SeckillStatsInput struct {
	GoodsId        int64 `json:"goodsId"`        // 商品ID
	GoodsOptionsId int64 `json:"goodsOptionsId"` // 商品选项ID
}

// SeckillStatsOutput 秒杀统计信息输出
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

// SeckillDoInput 后台秒杀请求输入
type SeckillDoInput struct {
	UserId         uint   `json:"userId"`         // 用户ID
	GoodsId        uint   `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId"` // 商品选项ID
	Count          uint   `json:"count"`          // 购买数量
	RequestId      string `json:"requestId"`      // 请求ID(幂等性)
	UserAddress    string `json:"userAddress"`    // 收货地址
	UserPhone      string `json:"userPhone"`      // 手机号码
	Remark         string `json:"remark"`         // 备注
}

// SeckillDoOutput 后台秒杀请求输出
type SeckillDoOutput struct {
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

// SeckillGoodsInfo 秒杀商品信息
type SeckillGoodsInfo struct {
	Id             int64       `json:"id"`             // 秒杀商品ID
	GoodsId        int64       `json:"goodsId"`        // 商品ID
	GoodsOptionsId int64       `json:"goodsOptionsId"` // 商品选项ID
	Price          float64     `json:"price"`          // 价格
	Stock          int         `json:"stock"`          // 库存
	StartTime      *gtime.Time `json:"startTime"`      // 开始时间
	EndTime        *gtime.Time `json:"endTime"`        // 结束时间
	Status         int         `json:"status"`         // 状态
}
