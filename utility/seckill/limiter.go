package seckill

import (
	"sync"
	"sync/atomic"
	"time"
)

// TokenBucket 令牌桶限流器
// 以固定速率向桶中添加令牌，请求到来时必须获取令牌才能继续执行
type TokenBucket struct {
	capacity    int32         // 桶容量
	rate        int32         // 令牌发放速率 (个/秒)
	tokens      int32         // 当前令牌数
	lastRefresh time.Time     // 上次刷新时间
	mu          sync.Mutex    // 互斥锁
	stop        chan struct{} // 停止信号
	metricHits  atomic.Int64  // 计数器：请求次数
	metricMiss  atomic.Int64  // 计数器：限流次数
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(capacity, rate int32) *TokenBucket {
	tb := &TokenBucket{
		capacity:    capacity,
		rate:        rate,
		tokens:      capacity, // 初始时桶是满的
		lastRefresh: time.Now(),
		stop:        make(chan struct{}),
	}

	// 启动后台令牌填充
	go tb.scheduleTokens()

	return tb
}

// scheduleTokens 后台定时填充令牌
func (tb *TokenBucket) scheduleTokens() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tb.refreshTokens()
		case <-tb.stop:
			return
		}
	}
}

// refreshTokens 刷新令牌数量
func (tb *TokenBucket) refreshTokens() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefresh).Seconds()

	// 计算新增的令牌数
	newTokens := int32(elapsed * float64(tb.rate))
	if newTokens > 0 {
		// 添加令牌，但不超过容量
		currentTokens := atomic.LoadInt32(&tb.tokens)
		newAmount := currentTokens + newTokens
		if newAmount > tb.capacity {
			newAmount = tb.capacity
		}
		atomic.StoreInt32(&tb.tokens, newAmount)
		tb.lastRefresh = now
	}
}

// Take 获取令牌
func (tb *TokenBucket) Take() bool {
	tb.metricHits.Add(1)
	if tb.tokens <= 0 {
		tb.metricMiss.Add(1)
		return false
	}
	result := atomic.AddInt32(&tb.tokens, -1) >= 0
	if !result {
		tb.metricMiss.Add(1)
	}
	return result
}

// GetMetrics 获取指标数据
func (tb *TokenBucket) GetMetrics() map[string]int64 {
	return map[string]int64{
		"capacity":    int64(tb.capacity),
		"rate":        int64(tb.rate),
		"tokens":      int64(tb.tokens),
		"hits":        tb.metricHits.Load(),
		"misses":      tb.metricMiss.Load(),
		"reject_rate": tb.metricMiss.Load() * 100 / max(tb.metricHits.Load(), 1),
	}
}

// Reset 重置令牌桶
func (tb *TokenBucket) Reset() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	atomic.StoreInt32(&tb.tokens, tb.capacity)
	tb.metricHits.Store(0)
	tb.metricMiss.Store(0)
}

// Close 关闭令牌桶
func (tb *TokenBucket) Close() {
	close(tb.stop)
}

// LeakyBucket 漏桶限流器
// 请求以任意速率流入桶中，但从桶中流出的速率是固定的
type LeakyBucket struct {
	capacity   int32         // 桶容量
	rate       int32         // 处理速率（每秒）
	water      int32         // 当前水量
	mu         sync.Mutex    // 互斥锁
	inChan     chan struct{} // 请求输入通道
	outChan    chan struct{} // 请求输出通道
	stop       chan struct{} // 停止信号
	metricHits atomic.Int64  // 计数器：请求次数
	metricMiss atomic.Int64  // 计数器：限流次数
}

// NewLeakyBucket 创建漏桶
func NewLeakyBucket(capacity, rate int32) *LeakyBucket {
	lb := &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		water:    0,
		inChan:   make(chan struct{}, capacity),
		outChan:  make(chan struct{}, capacity),
		stop:     make(chan struct{}),
	}

	// 启动固定速率消费协程
	go lb.consume()

	return lb
}

// consume 以固定速率消费请求
func (lb *LeakyBucket) consume() {
	// 计算每个请求的处理间隔
	interval := time.Second / time.Duration(lb.rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case <-lb.inChan:
				// 消费一个请求
				atomic.AddInt32(&lb.water, -1)
				// 通知请求已处理完成
				select {
				case lb.outChan <- struct{}{}:
				default:
					// 如果outChan已满，丢弃
				}
			default:
				// 没有请求，跳过
			}
		case <-lb.stop:
			return
		}
	}
}

// Take 获取处理机会
func (lb *LeakyBucket) Take() bool {
	lb.metricHits.Add(1)

	lb.mu.Lock()
	if lb.water >= lb.capacity {
		lb.mu.Unlock()
		lb.metricMiss.Add(1)
		return false
	}
	atomic.AddInt32(&lb.water, 1)
	lb.mu.Unlock()

	select {
	case <-lb.outChan:
		atomic.AddInt32(&lb.water, -1)
		return true
	case <-time.After(time.Second * 3): // 3秒超时
		atomic.AddInt32(&lb.water, -1)
		lb.metricMiss.Add(1)
		return false
	}
}

// GetMetrics 获取指标数据
func (lb *LeakyBucket) GetMetrics() map[string]int64 {
	return map[string]int64{
		"capacity":    int64(lb.capacity),
		"rate":        int64(lb.rate),
		"water":       int64(lb.water),
		"hits":        lb.metricHits.Load(),
		"misses":      lb.metricMiss.Load(),
		"reject_rate": lb.metricMiss.Load() * 100 / max(lb.metricHits.Load(), 1),
	}
}

// Reset 重置漏桶
func (lb *LeakyBucket) Reset() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	atomic.StoreInt32(&lb.water, 0)
	lb.metricHits.Store(0)
	lb.metricMiss.Store(0)

	// 清空队列
	for {
		select {
		case <-lb.outChan:
			// 取出队列中的元素
		default:
			return
		}
	}
}

// Close 关闭漏桶
func (lb *LeakyBucket) Close() {
	close(lb.stop)
}

// 辅助函数 - 返回两个数中较大的数
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
