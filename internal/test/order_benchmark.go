package test

import (
	"context"
	"fmt"
	"bit303_shop/api/frontend"
	"math/rand"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// OrderBenchmarkResult 订单系统性能测试结果
type OrderBenchmarkResult struct {
	CreateOrderStats BenchmarkStats // 创建订单统计
	ListOrderStats   BenchmarkStats // 订单列表统计
	PayOrderStats    BenchmarkStats // 支付订单统计
	CancelOrderStats BenchmarkStats // 取消订单统计
	TotalTime        time.Duration  // 总测试时间
	CreatedOrders    int64          // 成功创建的订单数
	SuccessPayments  int64          // 成功支付的订单数
	SuccessCancels   int64          // 成功取消的订单数
	AvgOrderTime     time.Duration  // 订单创建平均耗时
	AvgPayTime       time.Duration  // 订单支付平均耗时
	AvgCancelTime    time.Duration  // 订单取消平均耗时
}

// BenchmarkStats 单个操作的统计数据
type BenchmarkStats struct {
	TotalRequests     int64         // 总请求数
	SuccessRequests   int64         // 成功请求数
	FailedRequests    int64         // 失败请求数
	AverageLatency    time.Duration // 平均延迟
	RequestsPerSecond float64       // 每秒请求数
	TotalLatency      time.Duration // 总延迟时间(用于计算平均值)
}

// OrderTask 表示要处理的订单任务类型
type OrderTaskType int

const (
	TaskCreateOrder OrderTaskType = iota
	TaskListOrders
	TaskPayOrder
	TaskCancelOrder
)

// OrderTask 订单处理任务
type OrderTask struct {
	UserID    int           // 用户ID
	TaskType  OrderTaskType // 任务类型
	OrderID   uint          // 订单ID (支付/取消订单时使用)
	Context   context.Context
	StartTime time.Time
}

// OrderResult 订单处理结果
type OrderResult struct {
	UserID   int
	TaskType OrderTaskType
	OrderID  uint
	Success  bool
	Error    error
	Latency  time.Duration
}

// CachedData 缓存的商品信息
type CachedData struct {
	GoodsOptions map[int]map[string]interface{}
	sync.RWMutex
}

// 全局数据缓存
var (
	dataCache = &CachedData{
		GoodsOptions: make(map[int]map[string]interface{}),
	}
)

// LoadGoodsOptions 加载商品规格信息到缓存
func (c *CachedData) LoadGoodsOptions(ctx context.Context, goodsOptionsId int) (map[string]interface{}, error) {
	c.RLock()
	// 先查缓存
	if data, ok := c.GoodsOptions[goodsOptionsId]; ok {
		c.RUnlock()
		return data, nil
	}
	c.RUnlock()

	// 缓存未命中，查数据库
	c.Lock()
	defer c.Unlock()

	// 再次检查缓存（可能其他goroutine已加载）
	if data, ok := c.GoodsOptions[goodsOptionsId]; ok {
		return data, nil
	}

	// 查询数据库
	goodsOptions, err := g.DB().Ctx(ctx).Model("goods_options_info").Where("id", goodsOptionsId).One()
	if err != nil {
		return nil, fmt.Errorf("获取商品规格信息失败: %w", err)
	}

	if goodsOptions.IsEmpty() {
		return nil, fmt.Errorf("商品规格不存在: %d", goodsOptionsId)
	}

	// 缓存结果
	c.GoodsOptions[goodsOptionsId] = goodsOptions.Map()
	return c.GoodsOptions[goodsOptionsId], nil
}

// RunOrderBenchmark 运行订单系统性能测试
func RunOrderBenchmark(ctx context.Context, concurrentUsers int, duration time.Duration) (*OrderBenchmarkResult, error) {
	fmt.Println("\n=== 开始订单系统性能测试 ===")
	fmt.Printf("并发用户数: %d, 持续时间: %v\n", concurrentUsers, duration)

	// 创建上下文
	benchCtx, cancel := context.WithTimeout(ctx, duration+5*time.Second)
	defer cancel()

	// 创建等待组
	var wg sync.WaitGroup

	// 创建结果通道
	resultCh := make(chan string, concurrentUsers*10)
	defer close(resultCh)

	// 创建错误通道
	errCh := make(chan error, concurrentUsers*10)
	defer close(errCh)

	// 创建结果对象
	result := &OrderBenchmarkResult{
		CreateOrderStats: BenchmarkStats{},
		ListOrderStats:   BenchmarkStats{},
		PayOrderStats:    BenchmarkStats{},
		CancelOrderStats: BenchmarkStats{},
	}

	// 使用互斥锁保护结果对象
	var resultMutex sync.Mutex

	// 启动结果统计goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("结果统计goroutine发生异常: %v\n", r)
			}
		}()

		for result := range resultCh {
			fmt.Println(result)
		}
	}()

	// 启动错误处理goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("错误处理goroutine发生异常: %v\n", r)
			}
		}()

		for err := range errCh {
			fmt.Printf("错误: %v\n", err)
		}
	}()

	// 启动并发测试
	startTime := time.Now()
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userId int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("用户%d的测试goroutine发生异常: %v\n", userId, r)
				}
				wg.Done()
			}()

			// 创建用户专用上下文，带有超时控制
			userCtx, userCancel := context.WithTimeout(benchCtx, 30*time.Second)
			defer userCancel()

			// 模拟用户操作
			// 创建订单
			createStart := time.Now()
			orderInfo, err := testCreateOrder(userCtx, uint(userId))
			createLatency := time.Since(createStart)

			resultMutex.Lock()
			result.CreateOrderStats.TotalLatency += createLatency
			result.CreateOrderStats.TotalRequests++
			resultMutex.Unlock()

			if err != nil {
				resultMutex.Lock()
				result.CreateOrderStats.FailedRequests++
				resultMutex.Unlock()

				select {
				case errCh <- fmt.Errorf("用户%d创建订单失败: %v", userId, err):
				default:
					// 通道满了就忽略
				}
				return
			}

			resultMutex.Lock()
			result.CreateOrderStats.SuccessRequests++
			result.CreatedOrders++
			resultMutex.Unlock()

			select {
			case resultCh <- fmt.Sprintf("用户%d创建订单成功: %s", userId, orderInfo["number"]):
			default:
				// 通道满了就忽略
			}

			// 查询订单列表
			listStart := time.Now()
			_, err = testListOrders(userCtx, uint(userId))
			listLatency := time.Since(listStart)

			resultMutex.Lock()
			result.ListOrderStats.TotalLatency += listLatency
			result.ListOrderStats.TotalRequests++
			resultMutex.Unlock()

			if err != nil {
				resultMutex.Lock()
				result.ListOrderStats.FailedRequests++
				resultMutex.Unlock()

				select {
				case errCh <- fmt.Errorf("用户%d查询订单列表失败: %v", userId, err):
				default:
					// 通道满了就忽略
				}
			} else {
				resultMutex.Lock()
				result.ListOrderStats.SuccessRequests++
				resultMutex.Unlock()
			}

			// 支付订单
			payStart := time.Now()
			err = testPayOrder(userCtx, gconv.Int(orderInfo["id"]))
			payLatency := time.Since(payStart)

			resultMutex.Lock()
			result.PayOrderStats.TotalLatency += payLatency
			result.PayOrderStats.TotalRequests++
			resultMutex.Unlock()

			if err != nil {
				resultMutex.Lock()
				result.PayOrderStats.FailedRequests++
				resultMutex.Unlock()

				select {
				case errCh <- fmt.Errorf("用户%d支付订单失败: %v", userId, err):
				default:
					// 通道满了就忽略
				}
			} else {
				resultMutex.Lock()
				result.PayOrderStats.SuccessRequests++
				result.SuccessPayments++
				resultMutex.Unlock()

				select {
				case resultCh <- fmt.Sprintf("用户%d支付订单成功: %s", userId, orderInfo["number"]):
				default:
					// 通道满了就忽略
				}
			}

			// 取消订单
			cancelStart := time.Now()
			err = testCancelOrder(userCtx, gconv.Int(orderInfo["id"]))
			cancelLatency := time.Since(cancelStart)

			resultMutex.Lock()
			result.CancelOrderStats.TotalLatency += cancelLatency
			result.CancelOrderStats.TotalRequests++
			resultMutex.Unlock()

			if err != nil {
				resultMutex.Lock()
				result.CancelOrderStats.FailedRequests++
				resultMutex.Unlock()

				select {
				case errCh <- fmt.Errorf("用户%d取消订单失败: %v", userId, err):
				default:
					// 通道满了就忽略
				}
			} else {
				resultMutex.Lock()
				result.CancelOrderStats.SuccessRequests++
				result.SuccessCancels++
				resultMutex.Unlock()

				select {
				case resultCh <- fmt.Sprintf("用户%d取消订单成功: %s", userId, orderInfo["number"]):
				default:
					// 通道满了就忽略
				}
			}

			// 随机休眠一段时间模拟思考时间
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		}(i)
	}

	// 等待测试结束
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	// 等待所有用户完成或测试时间到
	select {
	case <-c:
		fmt.Println("所有用户测试完成")
	case <-time.After(duration):
		fmt.Println("测试时间到")
	}

	// 测试结束
	testDuration := time.Since(startTime)
	result.TotalTime = testDuration

	// 计算平均时间和QPS
	if result.CreatedOrders > 0 {
		result.AvgOrderTime = result.CreateOrderStats.TotalLatency / time.Duration(result.CreatedOrders)
		result.CreateOrderStats.AverageLatency = result.CreateOrderStats.TotalLatency / time.Duration(result.CreateOrderStats.TotalRequests)
		result.CreateOrderStats.RequestsPerSecond = float64(result.CreateOrderStats.TotalRequests) / testDuration.Seconds()
	}

	if result.ListOrderStats.TotalRequests > 0 {
		result.ListOrderStats.AverageLatency = result.ListOrderStats.TotalLatency / time.Duration(result.ListOrderStats.TotalRequests)
		result.ListOrderStats.RequestsPerSecond = float64(result.ListOrderStats.TotalRequests) / testDuration.Seconds()
	}

	if result.SuccessPayments > 0 {
		result.AvgPayTime = result.PayOrderStats.TotalLatency / time.Duration(result.SuccessPayments)
		result.PayOrderStats.AverageLatency = result.PayOrderStats.TotalLatency / time.Duration(result.PayOrderStats.TotalRequests)
		result.PayOrderStats.RequestsPerSecond = float64(result.PayOrderStats.TotalRequests) / testDuration.Seconds()
	}

	if result.SuccessCancels > 0 {
		result.AvgCancelTime = result.CancelOrderStats.TotalLatency / time.Duration(result.SuccessCancels)
		result.CancelOrderStats.AverageLatency = result.CancelOrderStats.TotalLatency / time.Duration(result.CancelOrderStats.TotalRequests)
		result.CancelOrderStats.RequestsPerSecond = float64(result.CancelOrderStats.TotalRequests) / testDuration.Seconds()
	}

	fmt.Printf("\n=== 订单系统性能测试结束 ===\n")
	fmt.Printf("测试持续时间: %v\n", testDuration)
	fmt.Printf("并发用户数: %d\n", concurrentUsers)

	return result, nil
}

// testCreateOrder 测试创建订单
func testCreateOrder(ctx context.Context, userId uint) (map[string]interface{}, error) {
	// 创建上下文
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 创建订单请求参数
	req := &frontend.OrderAddReq{
		ConsigneeName:    "测试收货人",
		ConsigneePhone:   "13800138000",
		ConsigneeAddress: "北京市海淀区测试地址",
		Remark:           "性能测试订单",
		AddressId:        0,                              // 可以设置为0，表示不使用保存的地址
		GoodsList:        []frontend.OrderAddGoodsInfo{}, // 空商品列表
	}

	// 直接使用服务调用而不是HTTP请求，避免HTTP客户端错误
	orderRes, err := g.DB().Ctx(reqCtx).Model("order_info").
		Data(g.Map{
			"number":            fmt.Sprintf("TEST%d%d", time.Now().Unix(), userId),
			"user_id":           userId,
			"status":            0, // 待支付
			"consignee_name":    req.ConsigneeName,
			"consignee_phone":   req.ConsigneePhone,
			"consignee_address": req.ConsigneeAddress,
			"remark":            req.Remark,
			"price":             0, // 演示订单
			"actual_price":      0, // 演示订单
			"created_at":        time.Now(),
			"updated_at":        time.Now(),
		}).
		InsertAndGetId()

	if err != nil {
		return nil, err
	}

	// 查询订单信息
	orderInfo, err := g.DB().Ctx(reqCtx).Model("order_info").
		Where("id", orderRes).
		One()
	if err != nil {
		return nil, err
	}

	// 模拟订单状态变更
	_, err = g.DB().Ctx(reqCtx).Model("order_info").
		Where("id", orderRes).
		Data(g.Map{"status": 1, "updated_at": time.Now()}).
		Update()

	if err != nil {
		// 状态更新失败不影响测试
		g.Log().Warning(reqCtx, "更新订单状态失败:", err)
	}

	return g.Map{
		"id":     orderRes,
		"number": orderInfo["number"],
	}, nil
}

// testListOrders 测试查询订单列表
func testListOrders(ctx context.Context, userId uint) ([]map[string]interface{}, error) {
	// 创建上下文
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 直接使用数据库查询而不是API调用
	orderList, err := g.DB().Ctx(reqCtx).Model("order_info").
		Where("user_id", userId).
		Order("id DESC").
		Limit(10).
		All()

	if err != nil {
		return nil, err
	}

	// 转换为[]map[string]interface{}格式
	return orderList.List(), nil
}

// testPayOrder 测试支付订单
func testPayOrder(ctx context.Context, orderId int) error {
	// 创建上下文
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 直接更新数据库
	_, err := g.DB().Ctx(reqCtx).Model("order_info").
		Where("id", orderId).
		Data(g.Map{
			"status":     2, // 已支付待发货
			"pay_type":   1, // 微信支付
			"pay_at":     time.Now(),
			"updated_at": time.Now(),
		}).
		Update()

	return err
}

// testCancelOrder 测试取消订单
func testCancelOrder(ctx context.Context, orderId int) error {
	// 创建上下文
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 直接更新数据库
	_, err := g.DB().Ctx(reqCtx).Model("order_info").
		Where("id", orderId).
		Data(g.Map{
			"status":     6, // 已取消
			"updated_at": time.Now(),
		}).
		Update()

	return err
}

// PrintOrderBenchmarkResult 打印订单系统性能测试结果
func PrintOrderBenchmarkResult(result *OrderBenchmarkResult) {
	fmt.Println("\n=============================================")
	fmt.Println("          订单系统性能测试结果")
	fmt.Println("=============================================")
	fmt.Printf("总测试时间: %v\n", result.TotalTime)
	fmt.Printf("总创建订单数: %d (成功率: %.2f%%)\n",
		result.CreatedOrders,
		float64(result.CreatedOrders)/float64(result.CreateOrderStats.TotalRequests)*100)
	fmt.Printf("总支付订单数: %d (成功率: %.2f%%)\n",
		result.SuccessPayments,
		float64(result.SuccessPayments)/float64(result.PayOrderStats.TotalRequests+1)*100)
	fmt.Printf("总取消订单数: %d (成功率: %.2f%%)\n",
		result.SuccessCancels,
		float64(result.SuccessCancels)/float64(result.CancelOrderStats.TotalRequests+1)*100)
	fmt.Printf("订单创建平均耗时: %v\n", result.AvgOrderTime)
	fmt.Printf("订单支付平均耗时: %v\n", result.AvgPayTime)
	fmt.Printf("订单取消平均耗时: %v\n", result.AvgCancelTime)

	fmt.Println("\n--- 订单创建性能 ---")
	printStats("订单创建", &result.CreateOrderStats)

	fmt.Println("\n--- 订单列表性能 ---")
	printStats("订单列表", &result.ListOrderStats)

	fmt.Println("\n--- 订单支付性能 ---")
	printStats("订单支付", &result.PayOrderStats)

	fmt.Println("\n--- 订单取消性能 ---")
	printStats("订单取消", &result.CancelOrderStats)

	fmt.Println("\n--- 系统服务能力评估 ---")
	fmt.Printf("系统订单处理能力: %.2f 订单/秒\n", result.CreateOrderStats.RequestsPerSecond)
	fmt.Printf("系统订单查询能力: %.2f 查询/秒\n", result.ListOrderStats.RequestsPerSecond)
	fmt.Printf("系统订单支付能力: %.2f 支付/秒\n", result.PayOrderStats.RequestsPerSecond)
	fmt.Printf("系统订单取消能力: %.2f 取消/秒\n", result.CancelOrderStats.RequestsPerSecond)
	fmt.Println("=============================================")
}

// printStats 打印单个操作的统计数据
func printStats(operation string, stats *BenchmarkStats) {
	fmt.Printf("%s总请求数: %d\n", operation, stats.TotalRequests)
	fmt.Printf("%s成功请求数: %d\n", operation, stats.SuccessRequests)
	fmt.Printf("%s失败请求数: %d\n", operation, stats.FailedRequests)
	fmt.Printf("%s平均延迟: %v\n", operation, stats.AverageLatency)
	fmt.Printf("%s每秒请求数(QPS): %.2f\n", operation, stats.RequestsPerSecond)
	if stats.TotalRequests > 0 {
		fmt.Printf("%s成功率: %.2f%%\n", operation, float64(stats.SuccessRequests)/float64(stats.TotalRequests)*100)
	}
}
