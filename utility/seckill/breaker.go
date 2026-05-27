package seckill

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// CircuitBreakerState 熔断器状态
type CircuitBreakerState int32

const (
	StateClosed   CircuitBreakerState = iota // 关闭状态 - 正常请求
	StateOpen                                // 开启状态 - 拒绝所有请求
	StateHalfOpen                            // 半开状态 - 允许部分请求通过
)

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	name            string           // 熔断器名称
	state           int32            // 当前状态
	failures        int32            // 失败计数
	successes       int32            // 成功计数
	threshold       int32            // 失败阈值
	successesNeeded int32            // 半开状态下需要的连续成功数
	timeout         time.Duration    // 熔断超时时间
	lastStateChange time.Time        // 最后状态变更时间
	mutex           sync.RWMutex     // 读写锁
	counts          map[string]int64 // 各类统计计数
	countsMutex     sync.Mutex       // 计数锁
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, threshold int32, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:            name,
		state:           int32(StateClosed),
		failures:        0,
		successes:       0,
		threshold:       threshold,
		successesNeeded: 5, // 默认半开状态需要5次连续成功
		timeout:         timeout,
		lastStateChange: time.Now(),
		counts:          make(map[string]int64),
	}
}

// Execute 执行请求
func (cb *CircuitBreaker) Execute(ctx context.Context, req func() (interface{}, error)) (interface{}, error) {
	// 检查熔断器状态
	if !cb.AllowRequest() {
		cb.incrementMetric("rejected")
		g.Log().Warning(ctx, "Circuit breaker is open, request rejected")
		return nil, ErrCircuitBreakerOpen
	}

	// 执行请求
	result, err := req()

	// 处理结果
	if err != nil {
		cb.RecordFailure()
		cb.incrementMetric("failures")
		return result, err
	}

	// 记录成功
	cb.RecordSuccess()
	cb.incrementMetric("successes")
	return result, nil
}

// AllowRequest 判断是否允许请求通过
func (cb *CircuitBreaker) AllowRequest() bool {
	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))
	switch state {
	case StateClosed:
		return true
	case StateOpen:
		// 检查是否超时，如果超时则切换到半开状态
		if time.Since(cb.lastStateChange) > cb.timeout {
			cb.moveToHalfOpen()
			return true
		}
		return false
	case StateHalfOpen:
		// 半开状态下只允许有限的请求通过
		cb.mutex.RLock()
		current := cb.successes + cb.failures
		cb.mutex.RUnlock()

		// 根据当前已处理的请求数决定是否允许新请求
		return current < cb.successesNeeded
	default:
		return true
	}
}

// RecordSuccess 记录成功
func (cb *CircuitBreaker) RecordSuccess() {
	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch state {
	case StateHalfOpen:
		// 在半开状态下累计成功计数
		cb.successes++
		// 如果达到连续成功阈值，关闭熔断器
		if cb.successes >= cb.successesNeeded {
			cb.moveToClosed()
		}
	case StateOpen:
		// 不应该在开启状态下有成功，但如果有，切换到半开
		cb.moveToHalfOpen()
	case StateClosed:
		// 重置失败计数
		cb.failures = 0
	}
}

// RecordFailure 记录失败
func (cb *CircuitBreaker) RecordFailure() {
	state := CircuitBreakerState(atomic.LoadInt32(&cb.state))

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	switch state {
	case StateClosed:
		// 关闭状态下累计失败次数
		cb.failures++
		// 如果失败次数达到阈值，开启熔断
		if cb.failures >= cb.threshold {
			cb.moveToOpen()
		}
	case StateHalfOpen:
		// 半开状态下任何失败都会重新开启熔断
		cb.moveToOpen()
	}
}

// Reset 重置熔断器状态
func (cb *CircuitBreaker) Reset() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	atomic.StoreInt32(&cb.state, int32(StateClosed))
	cb.failures = 0
	cb.successes = 0
	cb.lastStateChange = time.Now()

	cb.countsMutex.Lock()
	defer cb.countsMutex.Unlock()
	for k := range cb.counts {
		cb.counts[k] = 0
	}
}

// moveToOpen 切换到开启状态
func (cb *CircuitBreaker) moveToOpen() {
	atomic.StoreInt32(&cb.state, int32(StateOpen))
	cb.lastStateChange = time.Now()
	cb.failures = 0
	cb.successes = 0
	g.Log().Warning(context.Background(), "Circuit breaker moved to OPEN state: ", cb.name)
}

// moveToHalfOpen 切换到半开状态
func (cb *CircuitBreaker) moveToHalfOpen() {
	atomic.StoreInt32(&cb.state, int32(StateHalfOpen))
	cb.lastStateChange = time.Now()
	cb.failures = 0
	cb.successes = 0
	g.Log().Info(context.Background(), "Circuit breaker moved to HALF-OPEN state: ", cb.name)
}

// moveToClosed 切换到关闭状态
func (cb *CircuitBreaker) moveToClosed() {
	atomic.StoreInt32(&cb.state, int32(StateClosed))
	cb.lastStateChange = time.Now()
	cb.failures = 0
	cb.successes = 0
	g.Log().Info(context.Background(), "Circuit breaker moved to CLOSED state: ", cb.name)
}

// incrementMetric 增加计数指标
func (cb *CircuitBreaker) incrementMetric(name string) {
	cb.countsMutex.Lock()
	defer cb.countsMutex.Unlock()
	cb.counts[name]++
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	return CircuitBreakerState(atomic.LoadInt32(&cb.state))
}

// GetMetrics 获取指标
func (cb *CircuitBreaker) GetMetrics() map[string]int64 {
	cb.countsMutex.Lock()
	defer cb.countsMutex.Unlock()

	// 复制一份指标数据返回
	metrics := make(map[string]int64)
	for k, v := range cb.counts {
		metrics[k] = v
	}

	// 添加状态信息
	metrics["state"] = int64(cb.GetState())

	return metrics
}

// 错误类型
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open")
)
