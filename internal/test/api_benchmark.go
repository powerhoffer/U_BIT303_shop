package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// APITestConfig API测试配置
type APITestConfig struct {
	BaseURL         string        // API基础URL
	ConcurrentUsers int           // 并发用户数
	Duration        time.Duration // 测试持续时间
	ThinkTime       time.Duration // 思考时间
	RequestTimeout  time.Duration // 请求超时时间
}

// APIBenchmarkResult API性能测试结果
type APIBenchmarkResult struct {
	TotalRequests     int64            // 总请求数
	SuccessRequests   int64            // 成功请求数
	FailedRequests    int64            // 失败请求数
	TotalTime         time.Duration    // 总测试时间
	AverageLatency    time.Duration    // 平均延迟
	MaxLatency        time.Duration    // 最大延迟
	MinLatency        time.Duration    // 最小延迟
	P95Latency        time.Duration    // 95分位延迟
	P99Latency        time.Duration    // 99分位延迟
	RequestsPerSecond float64          // 每秒请求数
	ErrorDetails      map[string]int64 // 错误详情统计
}

// APILatencyRecord API延迟记录
type APILatencyRecord struct {
	timestamp time.Time
	latency   time.Duration
	success   bool
	errorMsg  string
}

// SeckillRequest 秒杀请求
type SeckillRequest struct {
	UserId         uint   `json:"userId"`         // 用户ID
	GoodsId        uint   `json:"goodsId"`        // 商品ID
	GoodsOptionsId uint   `json:"goodsOptionsId"` // 商品选项ID
	Count          uint   `json:"count"`          // 购买数量
	RequestId      string `json:"requestId"`      // 请求ID
	UserAddress    string `json:"userAddress"`    // 收货地址
	UserPhone      string `json:"userPhone"`      // 手机号码
	Remark         string `json:"remark"`         // 备注
}

// SeckillResponse 秒杀响应
type SeckillResponse struct {
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

// RunOrderAPIBenchmark 运行订单API性能测试
func RunOrderAPIBenchmark(ctx context.Context, baseURL string, concurrentUsers int, duration, thinkTime time.Duration) (*APIBenchmarkResult, error) {
	config := &APITestConfig{
		BaseURL:         baseURL,
		ConcurrentUsers: concurrentUsers,
		Duration:        duration,
		ThinkTime:       thinkTime,
		RequestTimeout:  30 * time.Second,
	}
	return runAPIBenchmark(ctx, config, "order", 1, 1)
}

// RunSeckillAPIBenchmark 运行秒杀API性能测试
func RunSeckillAPIBenchmark(ctx context.Context, config *APITestConfig, goodsId, optionsId int64) (*APIBenchmarkResult, error) {
	return runAPIBenchmark(ctx, config, "seckill", goodsId, optionsId)
}

// runAPIBenchmark 运行API性能测试的通用函数
func runAPIBenchmark(ctx context.Context, config *APITestConfig, apiType string, goodsId, optionsId int64) (*APIBenchmarkResult, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	g.Log().Info(ctx, "开始API性能测试...")
	g.Log().Info(ctx, fmt.Sprintf("配置信息: API类型=%s, 并发用户数=%d, 持续时间=%v, 思考时间=%v",
		apiType, config.ConcurrentUsers, config.Duration, config.ThinkTime))

	var (
		wg            sync.WaitGroup
		successCount  atomic.Int64
		failedCount   atomic.Int64
		totalRequests atomic.Int64
		startTime     = time.Now()
		latencyChan   = make(chan APILatencyRecord, config.ConcurrentUsers*100)
		errorStats    = make(map[string]int64)
		errorLock     sync.Mutex
	)

	// 创建停止信号
	stopChan := make(chan struct{})

	// 启动延迟统计协程
	var latencyRecords []time.Duration
	var statsWg sync.WaitGroup
	statsWg.Add(1)
	go func() {
		defer statsWg.Done()
		collectAPILatencyStats(latencyChan, &latencyRecords)
	}()

	// 设置结束时间
	endTime := startTime.Add(config.Duration)

	// 启动所有用户协程
	for i := 0; i < config.ConcurrentUsers; i++ {
		userID := i + 1
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// 创建HTTP客户端
			client := &http.Client{
				Timeout: config.RequestTimeout,
			}

			// 生成随机请求ID前缀
			randomPrefix := fmt.Sprintf("%d-%d-", time.Now().UnixNano(), userID)

			// 循环发送请求，直到测试结束
			for time.Now().Before(endTime) {
				select {
				case <-stopChan:
					return
				default:
					// 生成唯一的请求ID
					requestID := fmt.Sprintf("%s%d", randomPrefix, rand.Int63())

					// 生成随机手机号
					phone := fmt.Sprintf("1%d%09d", rand.Intn(7)+3, rand.Intn(1000000000))

					// 执行API请求
					requestStart := time.Now()
					var statusCode int
					var respBody []byte
					var err error

					if apiType == "seckill" {
						// 构建秒杀请求
						requestBody := SeckillRequest{
							UserId:         uint(userID),
							GoodsId:        uint(goodsId),
							GoodsOptionsId: uint(optionsId),
							Count:          1,
							RequestId:      requestID,
							UserAddress:    "测试地址",
							UserPhone:      phone,
							Remark:         "API测试",
						}

						// 调用秒杀API
						statusCode, respBody, err = doRequest(client, fmt.Sprintf("%s/seckill/do", config.BaseURL), requestBody)
					} else {
						// 构建订单请求 (示例)
						orderRequest := map[string]interface{}{
							"userId": userID,
							"goods": []map[string]interface{}{
								{
									"goodsId": goodsId,
									"count":   1,
								},
							},
						}

						// 调用订单API
						statusCode, respBody, err = doRequest(client, fmt.Sprintf("%s/order/create", config.BaseURL), orderRequest)
					}

					// 计算延迟
					latency := time.Since(requestStart)

					// 处理响应
					if err != nil || statusCode != http.StatusOK {
						failedCount.Add(1)
						errorLock.Lock()
						errMsg := ""
						if err != nil {
							errMsg = err.Error()
						} else {
							errMsg = fmt.Sprintf("HTTP Status: %d", statusCode)
						}
						errorStats[errMsg]++
						errorLock.Unlock()

						latencyChan <- APILatencyRecord{
							timestamp: time.Now(),
							latency:   latency,
							success:   false,
							errorMsg:  errMsg,
						}
					} else {
						// 解析响应
						if apiType == "seckill" {
							var resp SeckillResponse
							if jsonErr := json.Unmarshal(respBody, &resp); jsonErr == nil {
								if resp.Status == 0 {
									successCount.Add(1)
									latencyChan <- APILatencyRecord{
										timestamp: time.Now(),
										latency:   latency,
										success:   true,
									}
								} else {
									failedCount.Add(1)
									errMsg := fmt.Sprintf("业务错误: %s", resp.Message)
									errorLock.Lock()
									errorStats[errMsg]++
									errorLock.Unlock()
									latencyChan <- APILatencyRecord{
										timestamp: time.Now(),
										latency:   latency,
										success:   false,
										errorMsg:  errMsg,
									}
								}
							} else {
								failedCount.Add(1)
								errMsg := fmt.Sprintf("JSON解析错误: %v", jsonErr)
								errorLock.Lock()
								errorStats[errMsg]++
								errorLock.Unlock()
								latencyChan <- APILatencyRecord{
									timestamp: time.Now(),
									latency:   latency,
									success:   false,
									errorMsg:  errMsg,
								}
							}
						} else {
							// 简单处理订单响应
							successCount.Add(1)
							latencyChan <- APILatencyRecord{
								timestamp: time.Now(),
								latency:   latency,
								success:   true,
							}
						}
					}

					totalRequests.Add(1)

					// 每1000个请求输出一次统计信息
					if totalRequests.Load()%1000 == 0 {
						success := successCount.Load()
						total := totalRequests.Load()
						g.Log().Info(ctx, fmt.Sprintf("进度: 总请求=%d, 成功=%d, 失败=%d, 成功率=%.2f%%",
							total, success, failedCount.Load(), float64(success)/float64(total)*100))
					}

					// 思考时间
					select {
					case <-stopChan:
						return
					case <-time.After(config.ThinkTime):
						// 继续下一次请求
					}
				}
			}
		}(userID)
	}

	// 等待测试结束
	time.Sleep(config.Duration)
	close(stopChan)
	wg.Wait()

	// 关闭延迟统计通道
	close(latencyChan)
	statsWg.Wait()

	// 计算统计指标
	totalTime := time.Since(startTime)
	requestsPerSecond := float64(totalRequests.Load()) / totalTime.Seconds()

	// 计算延迟指标
	minLatency, maxLatency := calculateMinMaxLatency(latencyRecords)
	avgLatency := calculateAverageLatency(latencyRecords)
	p95Latency := calculatePercentileLatency(latencyRecords, 95)
	p99Latency := calculatePercentileLatency(latencyRecords, 99)

	// 返回结果
	return &APIBenchmarkResult{
		TotalRequests:     totalRequests.Load(),
		SuccessRequests:   successCount.Load(),
		FailedRequests:    failedCount.Load(),
		TotalTime:         totalTime,
		AverageLatency:    avgLatency,
		MaxLatency:        maxLatency,
		MinLatency:        minLatency,
		P95Latency:        p95Latency,
		P99Latency:        p99Latency,
		RequestsPerSecond: requestsPerSecond,
		ErrorDetails:      errorStats,
	}, nil
}

// doRequest 执行HTTP请求
func doRequest(client *http.Client, url string, body interface{}) (int, []byte, error) {
	// 将请求体序列化为JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	return resp.StatusCode, respBody, nil
}

// collectAPILatencyStats 收集延迟统计数据
func collectAPILatencyStats(latencyChan chan APILatencyRecord, records *[]time.Duration) {
	for record := range latencyChan {
		if record.success {
			*records = append(*records, record.latency)
		}
	}
}

// PrintAPIBenchmarkResult 打印API性能测试结果
func PrintAPIBenchmarkResult(result *APIBenchmarkResult) {
	fmt.Println("\n================ API性能测试结果 ================")
	fmt.Printf("总请求数: %d\n", result.TotalRequests)
	fmt.Printf("成功请求数: %d\n", result.SuccessRequests)
	fmt.Printf("失败请求数: %d\n", result.FailedRequests)

	if result.TotalRequests > 0 {
		successRate := float64(result.SuccessRequests) / float64(result.TotalRequests) * 100
		fmt.Printf("成功率: %.2f%%\n", successRate)
	}

	fmt.Printf("总测试时间: %v\n", result.TotalTime)
	fmt.Printf("每秒请求数 (QPS): %.2f\n", result.RequestsPerSecond)
	fmt.Printf("平均延迟: %v\n", result.AverageLatency)
	fmt.Printf("最小延迟: %v\n", result.MinLatency)
	fmt.Printf("最大延迟: %v\n", result.MaxLatency)
	fmt.Printf("95分位延迟: %v\n", result.P95Latency)
	fmt.Printf("99分位延迟: %v\n", result.P99Latency)

	if len(result.ErrorDetails) > 0 {
		fmt.Println("\n错误详情:")
		for errMsg, count := range result.ErrorDetails {
			fmt.Printf("- %s: %d\n", errMsg, count)
		}
	}

	fmt.Println("================================================")
}
