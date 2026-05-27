package consts

import "time"

// 秒杀系统常量 - 与基础常量不冲突的部分
const (
	// 秒杀系统配置
	SeckillWorkerPoolSize   = 10000              // 秒杀请求队列容量
	SeckillMinWorkers       = 4                  // 最小工作协程数
	SeckillMaxWorkers       = 32                 // 最大工作协程数
	SeckillProcessTimeout   = 3 * time.Second    // 秒杀处理超时时间
	SeckillResultTTLSeconds = 1800 * time.Second // 结果缓存过期时间(30分钟)
	SeckillOrderTTLHours    = 24 * time.Hour     // 订单过期时间
	SeckillBreakerTimeout   = 5 * time.Second    // 熔断器超时时间
)

// 秒杀结果状态码
const (
	CodeSeckillSuccess     = 0 // 秒杀成功
	CodeSeckillFailed      = 1 // 秒杀失败(一般性失败)
	CodeSeckillNoStock     = 2 // 库存不足
	CodeSeckillRateLimited = 3 // 被限流
	CodeSeckillDuplicate   = 4 // 重复请求
	CodeSeckillTimeout     = 5 // 处理超时
	CodeSeckillCircuitOpen = 6 // 熔断器开启
	CodeSeckillProcessing  = 7 // 正在处理
	CodeSeckillSystemError = 8 // 系统错误
)
