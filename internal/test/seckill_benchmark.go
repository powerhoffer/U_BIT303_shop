package test

import (
	"context"
	"encoding/json"
	"fmt"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// BenchmarkResult 性能测试结果
type BenchmarkResult struct {
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

// LatencyRecord 延迟记录
type LatencyRecord struct {
	timestamp time.Time
	latency   time.Duration
	success   bool
	errorMsg  string
}

// BenchmarkConfig 性能测试配置
type BenchmarkConfig struct {
	ConcurrentUsers int           // 并发用户数
	Duration        time.Duration // 测试持续时间
	RampUpTime      time.Duration // 预热时间
	RampUpUsers     int           // 预热阶段用户数
	ThinkTime       time.Duration // 用户思考时间
	GoodsId         int64         // 测试商品ID
	GoodsOptionsId  int64         // 测试商品规格ID
}

// DefaultBenchmarkConfig 默认配置
func DefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		ConcurrentUsers: 100,
		Duration:        time.Minute * 5,
		RampUpTime:      time.Second * 30,
		RampUpUsers:     10,
		ThinkTime:       time.Millisecond * 100,
		GoodsId:         1,
		GoodsOptionsId:  1,
	}
}

// SeckillMessage Kafka消息结构
type SeckillMessage struct {
	UserId         uint      `json:"userId"`
	GoodsId        uint      `json:"goodsId"`
	GoodsOptionsId uint      `json:"goodsOptionsId"`
	CreateTime     time.Time `json:"createTime"`
}

// RunSeckillBenchmark 运行秒杀性能测试
func RunSeckillBenchmark(ctx context.Context, config *BenchmarkConfig) (*BenchmarkResult, error) {
	if config == nil {
		config = DefaultBenchmarkConfig()
	}

	// 确保服务已注册
	if service.Seckill() == nil {
		return nil, fmt.Errorf("seckill service not registered")
	}

	g.Log().Info(ctx, "开始秒杀性能测试...")
	g.Log().Info(ctx, fmt.Sprintf("配置信息: 并发用户数=%d, 持续时间=%v, 预热时间=%v",
		config.ConcurrentUsers, config.Duration, config.RampUpTime))

	// 初始化测试数据
	if err := prepareTestData(ctx, config); err != nil {
		return nil, fmt.Errorf("prepare test data failed: %v", err)
	}

	var (
		wg            sync.WaitGroup
		successCount  atomic.Int64
		failedCount   atomic.Int64
		totalRequests atomic.Int64
		startTime     = time.Now()
		latencyChan   = make(chan LatencyRecord, config.ConcurrentUsers*100)
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
		collectLatencyStats(latencyChan, &latencyRecords)
	}()

	// 分批启动用户（实现预热）
	batchSize := config.ConcurrentUsers / 10
	if batchSize < 1 {
		batchSize = 1
	}

	g.Log().Info(ctx, "开始预热阶段...")
	for batch := 0; batch < 10; batch++ {
		usersInBatch := batchSize
		if batch == 9 {
			usersInBatch = config.ConcurrentUsers - (9 * batchSize)
		}

		g.Log().Info(ctx, fmt.Sprintf("启动第 %d 批用户, 数量: %d", batch+1, usersInBatch))

		for i := 0; i < usersInBatch; i++ {
			userID := batch*batchSize + i + 1
			wg.Add(1)
			go func(userID int) {
				defer wg.Done()
				userCtx := context.WithValue(ctx, "userId", int64(userID))

				for {
					select {
					case <-stopChan:
						return
					default:
						requestStart := time.Now()
						req := &model.SeckillDoInput{
							UserId:         uint(userID),
							GoodsId:        uint(config.GoodsId),
							GoodsOptionsId: uint(config.GoodsOptionsId),
						}

						_, err := service.Seckill().DoSeckill(userCtx, req)
						latency := time.Since(requestStart)

						if err != nil {
							failedCount.Add(1)
							errorLock.Lock()
							errorStats[err.Error()]++
							errorLock.Unlock()
							select {
							case latencyChan <- LatencyRecord{
								timestamp: time.Now(),
								latency:   latency,
								success:   false,
								errorMsg:  err.Error(),
							}:
							case <-stopChan:
								return
							}
							if totalRequests.Load()%1000 == 0 {
								g.Log().Debug(userCtx, fmt.Sprintf("User %d failed: %v", userID, err))
							}
						} else {
							successCount.Add(1)
							select {
							case latencyChan <- LatencyRecord{
								timestamp: time.Now(),
								latency:   latency,
								success:   true,
							}:
							case <-stopChan:
								return
							}
							if totalRequests.Load()%1000 == 0 {
								g.Log().Debug(userCtx, fmt.Sprintf("User %d succeeded, latency: %v", userID, latency))
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

						time.Sleep(config.ThinkTime)
					}
				}
			}(userID)
		}
		if batch < 9 {
			time.Sleep(config.RampUpTime / 10)
		}
	}

	g.Log().Info(ctx, "预热完成，开始正式测试...")

	// 等待指定的测试时间
	time.Sleep(config.Duration)

	g.Log().Info(ctx, "测试时间到，开始清理...")

	// 先关闭停止信号，让所有测试 goroutine 结束
	close(stopChan)

	// 等待所有测试 goroutine 结束
	wg.Wait()

	// 关闭延迟统计通道
	close(latencyChan)

	// 等待统计 goroutine 结束
	statsWg.Wait()

	// 关闭Kafka连接和其他资源
	err := service.Seckill().Close()
	if err != nil {
		g.Log().Warning(ctx, "关闭Kafka连接时出错:", err)
	} else {
		g.Log().Info(ctx, "Kafka连接已关闭")
	}

	totalTime := time.Since(startTime)

	// 计算延迟统计
	p95 := calculatePercentileLatency(latencyRecords, 95)
	p99 := calculatePercentileLatency(latencyRecords, 99)
	minLatency, maxLatency := calculateMinMaxLatency(latencyRecords)
	avgLatency := calculateAverageLatency(latencyRecords)

	// 计算性能指标
	result := &BenchmarkResult{
		TotalRequests:     totalRequests.Load(),
		SuccessRequests:   successCount.Load(),
		FailedRequests:    failedCount.Load(),
		TotalTime:         totalTime,
		AverageLatency:    avgLatency,
		MaxLatency:        maxLatency,
		MinLatency:        minLatency,
		P95Latency:        p95,
		P99Latency:        p99,
		RequestsPerSecond: float64(totalRequests.Load()) / totalTime.Seconds(),
		ErrorDetails:      errorStats,
	}

	// 注意：我们不再清理测试数据，以便在测试后查看数据库中的结果
	// 如果需要清理数据，可以手动调用cleanupTestData函数

	g.Log().Info(ctx, fmt.Sprintf("秒杀性能测试完成: 总请求=%d, 成功=%d, 失败=%d, 成功率=%.2f%%, QPS=%.2f",
		result.TotalRequests, result.SuccessRequests, result.FailedRequests,
		float64(result.SuccessRequests)/float64(result.TotalRequests)*100,
		result.RequestsPerSecond))
	return result, nil
}

// PrintBenchmarkResult 打印性能测试结果
func PrintBenchmarkResult(result *BenchmarkResult) {
	fmt.Println("\n=== 秒杀系统性能测试结果 ===")
	fmt.Printf("总请求数: %d\n", result.TotalRequests)
	fmt.Printf("成功请求数: %d\n", result.SuccessRequests)
	fmt.Printf("失败请求数: %d\n", result.FailedRequests)
	fmt.Printf("总测试时间: %v\n", result.TotalTime)
	fmt.Printf("平均延迟: %v\n", result.AverageLatency)
	fmt.Printf("最大延迟: %v\n", result.MaxLatency)
	fmt.Printf("最小延迟: %v\n", result.MinLatency)
	fmt.Printf("P95延迟: %v\n", result.P95Latency)
	fmt.Printf("P99延迟: %v\n", result.P99Latency)
	fmt.Printf("每秒请求数(QPS): %.2f\n", result.RequestsPerSecond)
	fmt.Printf("成功率: %.2f%%\n", float64(result.SuccessRequests)/float64(result.TotalRequests)*100)

	if len(result.ErrorDetails) > 0 {
		fmt.Println("\n错误统计:")
		for errMsg, count := range result.ErrorDetails {
			fmt.Printf("%s: %d\n", errMsg, count)
		}
	}
}

// 辅助函数

func prepareTestData(ctx context.Context, config *BenchmarkConfig) error {
	g.Log().Info(ctx, "开始准备测试数据...")

	// 0. 检查 Kafka 连接和主题
	if err := verifyKafkaConnection(ctx); err != nil {
		return fmt.Errorf("Kafka验证失败: %v", err)
	}

	// 注释掉清理测试数据的步骤，以保留之前的测试数据
	// 1. 清理数据库中的测试数据
	// if err := cleanupTestData(ctx, config); err != nil {
	// 	return fmt.Errorf("清理测试数据失败: %v", err)
	// }

	// 2. 检查是否已存在测试商品，如果存在则跳过创建
	var goodsCount int
	goodsCount, err := g.DB().Ctx(ctx).Model("goods_info").Where("id", config.GoodsId).Count()
	if err != nil {
		return fmt.Errorf("检查商品是否存在失败: %v", err)
	}

	if goodsCount == 0 {
		// 创建测试商品
		if err := createTestGoods(ctx, config); err != nil {
			return fmt.Errorf("创建测试商品失败: %v", err)
		}
	} else {
		g.Log().Info(ctx, "测试商品已存在，跳过创建")
	}

	// 3. 检查是否已存在商品规格，如果存在则跳过创建
	var optionsCount int
	optionsCount, err = g.DB().Ctx(ctx).Model("goods_options_info").Where("id", config.GoodsOptionsId).Count()
	if err != nil {
		return fmt.Errorf("检查商品规格是否存在失败: %v", err)
	}

	if optionsCount == 0 {
		// 创建商品规格
		if err := createTestGoodsOptions(ctx, config); err != nil {
			return fmt.Errorf("创建商品规格失败: %v", err)
		}
	} else {
		g.Log().Info(ctx, "商品规格已存在，跳过创建")
	}

	// 4. 检查是否已存在秒杀商品，如果存在则跳过创建
	var seckillCount int
	seckillCount, err = g.DB().Ctx(ctx).Model("seckill_goods").Where("goods_id", config.GoodsId).Where("goods_options_id", config.GoodsOptionsId).Count()
	if err != nil {
		return fmt.Errorf("检查秒杀商品是否存在失败: %v", err)
	}

	if seckillCount == 0 {
		// 创建秒杀商品
		if err := createSeckillGoods(ctx, config); err != nil {
			return fmt.Errorf("创建秒杀商品失败: %v", err)
		}
	} else {
		g.Log().Info(ctx, "秒杀商品已存在，跳过创建")
	}

	// 5. 清理Redis缓存
	if err := cleanupRedisCache(ctx); err != nil {
		return fmt.Errorf("清理Redis缓存失败: %v", err)
	}

	// 6. 执行缓存预热
	g.Log().Info(ctx, "开始执行缓存预热...")
	err = service.Seckill().WarmUpCache(ctx, int64(config.GoodsId), int64(config.GoodsOptionsId))
	if err != nil {
		g.Log().Warning(ctx, fmt.Sprintf("缓存预热失败: %v", err))
		// 即使预热失败，也继续进行数据验证
	}

	// 7. 验证数据初始化
	for i := 0; i < 5; i++ { // 最多重试5次
		if err := verifyTestData(ctx, config); err != nil {
			if i < 4 {
				g.Log().Warning(ctx, fmt.Sprintf("数据验证失败，等待重试(%d/5): %v", i+1, err))
				time.Sleep(time.Second * 2)
				continue
			}
			return fmt.Errorf("数据验证失败: %v", err)
		}
		break
	}

	// 8. 等待确保缓存就绪
	time.Sleep(2 * time.Second)
	g.Log().Info(ctx, "测试数据准备完成")
	return nil
}

func verifyKafkaConnection(ctx context.Context) error {
	g.Log().Info(ctx, "正在验证Kafka连接...")

	// 从配置获取Kafka配置
	kafkaConfig := g.Cfg().MustGet(ctx, "mq.kafka")
	if kafkaConfig.IsNil() {
		return fmt.Errorf("Kafka配置不存在，请检查config.yaml中的mq.kafka配置")
	}

	// 获取broker列表
	configMap := kafkaConfig.Map()
	var brokers []string

	if brokersVal, ok := configMap["brokers"]; ok && brokersVal != nil {
		// 尝试将值转换为字符串切片
		switch v := brokersVal.(type) {
		case []string:
			brokers = v
		case []interface{}:
			for _, b := range v {
				brokers = append(brokers, fmt.Sprintf("%v", b))
			}
		case string:
			brokers = []string{v}
		default:
			// 尝试最后的方式获取值
			brokerStr := fmt.Sprintf("%v", brokersVal)
			if brokerStr != "" && brokerStr != "<nil>" {
				brokers = []string{brokerStr}
			}
		}
	}

	// 如果没有找到，使用默认值
	if len(brokers) == 0 {
		g.Log().Warning(ctx, "未找到配置的Kafka brokers，将使用默认值: localhost:9092")
		brokers = []string{"localhost:9092"}
	}

	g.Log().Info(ctx, fmt.Sprintf("尝试连接Kafka brokers: %v", brokers))

	// 创建Kafka配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Timeout = time.Second * 5

	// 尝试创建生产者
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("连接Kafka失败: %v，请确保Kafka服务正在运行", err)
	}
	defer producer.Close()

	// 尝试创建管理客户端
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("创建Kafka管理员客户端失败: %v", err)
	}
	defer admin.Close()

	// 获取主题列表
	topics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("获取Kafka主题列表失败: %v", err)
	}

	// 检查并创建必要的主题
	requiredTopics := []string{consts.KafkaTopicSeckill, consts.KafkaTopicSeckillComplete}
	for _, topic := range requiredTopics {
		if _, exists := topics[topic]; !exists {
			g.Log().Info(ctx, fmt.Sprintf("主题 %s 不存在，正在创建...", topic))

			// 设置一天的保留期 (86400000 ms)
			retentionStr := "86400000"

			topicDetail := &sarama.TopicDetail{
				NumPartitions:     1,
				ReplicationFactor: 1,
				ConfigEntries: map[string]*string{
					"retention.ms": &retentionStr,
				},
			}
			if err := admin.CreateTopic(topic, topicDetail, false); err != nil {
				if err != sarama.ErrTopicAlreadyExists {
					return fmt.Errorf("创建Kafka主题 %s 失败: %v", topic, err)
				}
				g.Log().Info(ctx, fmt.Sprintf("主题 %s 已存在", topic))
			} else {
				g.Log().Info(ctx, fmt.Sprintf("主题 %s 创建成功", topic))
			}
		} else {
			g.Log().Info(ctx, fmt.Sprintf("主题 %s 已存在", topic))
		}
	}

	// 发送测试消息以验证连接
	testMsg := &sarama.ProducerMessage{
		Topic: consts.KafkaTopicSeckill,
		Key:   sarama.StringEncoder("test"),
		Value: sarama.StringEncoder("test-message"),
	}

	partition, offset, err := producer.SendMessage(testMsg)
	if err != nil {
		return fmt.Errorf("发送测试消息失败: %v", err)
	}

	g.Log().Info(ctx, fmt.Sprintf("测试消息发送成功: topic=%s, partition=%d, offset=%d",
		consts.KafkaTopicSeckill, partition, offset))

	g.Log().Info(ctx, "Kafka连接验证成功")
	return nil
}

func cleanupTestData(ctx context.Context, config *BenchmarkConfig) error {
	tables := []string{
		"seckill_goods",
		"goods_info",
		"goods_options_info",
	}

	for _, table := range tables {
		_, err := g.DB().Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), config.GoodsId)
		if err != nil {
			return fmt.Errorf("清理表%s失败: %v", table, err)
		}
	}

	return nil
}

func createTestGoods(ctx context.Context, config *BenchmarkConfig) error {
	goodsInfo := g.Map{
		"id":                 int(config.GoodsId),
		"pic_url":            "test.jpg",
		"name":               "测试商品",
		"price":              10000,
		"level1_category_id": 1,
		"level2_category_id": 1,
		"level3_category_id": 1,
		"brand":              "测试品牌",
		"stock":              10000,
		"sale":               0,
		"tags":               "测试标签",
		"detail_info":        "测试商品详情",
		"created_at":         gtime.Now(),
		"updated_at":         gtime.Now(),
	}
	_, err := dao.GoodsInfo.Ctx(ctx).Data(goodsInfo).Save()
	return err
}

func createTestGoodsOptions(ctx context.Context, config *BenchmarkConfig) error {
	goodsOptions := g.Map{
		"id":         int(config.GoodsOptionsId),
		"goods_id":   int(config.GoodsId),
		"pic_url":    "test.jpg",
		"name":       "默认规格",
		"price":      10000,
		"stock":      10000,
		"created_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}
	_, err := dao.GoodsOptionsInfo.Ctx(ctx).Data(goodsOptions).Save()
	return err
}

func createSeckillGoods(ctx context.Context, config *BenchmarkConfig) error {
	now := gtime.Now()
	seckillGoodsData := g.Map{
		"goods_id":         int(config.GoodsId),
		"goods_options_id": int(config.GoodsOptionsId),
		"seckill_price":    8000,
		"status":           1,
		"seckill_stock":    1000,
		"start_time":       now,
		"end_time":         now.Add(time.Hour * 24),
		"original_price":   10000,
		"created_at":       now,
		"updated_at":       now,
	}
	_, err := dao.SeckillGoods.Ctx(ctx).Data(seckillGoodsData).Save()
	return err
}

func cleanupRedisCache(ctx context.Context) error {
	patterns := []string{
		fmt.Sprintf("%s*", consts.SeckillTokenBucketPrefix),
		fmt.Sprintf("%s*", consts.SeckillLeakyBucketPrefix),
		fmt.Sprintf("%s*", consts.SeckillStockPrefix),
		fmt.Sprintf("%s*", consts.SeckillLockPrefix),
		fmt.Sprintf("%s*", consts.SeckillResultPrefix),
		fmt.Sprintf("%s*", consts.SeckillMetricsPrefix),
		fmt.Sprintf("%s*", consts.SeckillGoodsPrefix),
		fmt.Sprintf("%s*", consts.SeckillQueuePrefix),
		fmt.Sprintf("%s*", consts.SeckillSuccessPrefix),
		fmt.Sprintf("%s*", consts.SeckillUserBoughtPrefix),
	}

	for _, pattern := range patterns {
		keys, err := g.Redis().Do(ctx, "KEYS", pattern)
		if err != nil {
			return fmt.Errorf("获取Redis键失败: %v", err)
		}
		if !keys.IsNil() && len(keys.Strings()) > 0 {
			_, err = g.Redis().Do(ctx, "DEL", keys.Interfaces()...)
			if err != nil {
				return fmt.Errorf("删除Redis键失败: %v", err)
			}
		}
	}
	return nil
}

func verifyTestData(ctx context.Context, config *BenchmarkConfig) error {
	// 1. 验证数据库记录
	var seckillGoods *entity.SeckillGoods
	err := dao.SeckillGoods.Ctx(ctx).
		Where(dao.SeckillGoods.Columns().GoodsId, config.GoodsId).
		Where(dao.SeckillGoods.Columns().GoodsOptionsId, config.GoodsOptionsId).
		Scan(&seckillGoods)
	if err != nil {
		return fmt.Errorf("验证秒杀商品失败: %v", err)
	}
	if seckillGoods == nil {
		return fmt.Errorf("秒杀商品不存在")
	}

	// 2. 验证Redis缓存
	// 2.1 验证商品信息
	goodsKey := fmt.Sprintf("%s%d:%d", consts.SeckillGoodsPrefix, config.GoodsId, config.GoodsOptionsId)
	goodsData, err := g.Redis().Do(ctx, "GET", goodsKey)
	if err != nil {
		return fmt.Errorf("验证Redis商品信息失败: %v", err)
	}
	if goodsData.IsNil() {
		// 如果缓存不存在，尝试重新预热
		g.Log().Warning(ctx, "商品信息未缓存，尝试重新预热")
		if err := service.Seckill().WarmUpCache(ctx, int64(config.GoodsId), int64(config.GoodsOptionsId)); err != nil {
			return fmt.Errorf("重新预热失败: %v", err)
		}
		// 重新检查缓存
		goodsData, err = g.Redis().Do(ctx, "GET", goodsKey)
		if err != nil || goodsData.IsNil() {
			return fmt.Errorf("重新预热后商品信息仍未缓存")
		}
	}

	// 2.2 验证库存信息
	stockKey := fmt.Sprintf("%s%d:%d", consts.SeckillStockPrefix, config.GoodsId, config.GoodsOptionsId)
	stockData, err := g.Redis().Do(ctx, "GET", stockKey)
	if err != nil {
		return fmt.Errorf("验证Redis库存失败: %v", err)
	}
	if stockData.IsNil() {
		return fmt.Errorf("Redis库存未初始化")
	}

	// 2.3 验证库存数量
	redisStock := stockData.Int64()
	dbStock := int64(seckillGoods.SeckillStock)
	if redisStock != dbStock {
		g.Log().Warning(ctx, fmt.Sprintf("Redis库存数量不匹配: Redis=%d, DB=%d", redisStock, dbStock))
		// 修正库存数量
		_, err = g.Redis().Do(ctx, "SET", stockKey, seckillGoods.SeckillStock)
		if err != nil {
			return fmt.Errorf("修正Redis库存失败: %v", err)
		}
	}

	g.Log().Debug(ctx, fmt.Sprintf("数据验证成功: 商品ID=%d, 规格ID=%d, 库存=%d",
		config.GoodsId, config.GoodsOptionsId, redisStock))
	return nil
}

func collectLatencyStats(latencyChan chan LatencyRecord, records *[]time.Duration) {
	for record := range latencyChan {
		*records = append(*records, record.latency)
	}
}

func calculatePercentileLatency(latencies []time.Duration, percentile float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	// 对延迟进行排序
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 计算百分位数
	index := int(float64(len(sorted)) * percentile / 100)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func calculateMinMaxLatency(latencies []time.Duration) (time.Duration, time.Duration) {
	if len(latencies) == 0 {
		return 0, 0
	}

	min := latencies[0]
	max := latencies[0]
	for _, latency := range latencies {
		if latency < min {
			min = latency
		}
		if latency > max {
			max = latency
		}
	}
	return min, max
}

func calculateAverageLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}
	return total / time.Duration(len(latencies))
}

// Kafka相关函数
func initKafkaProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %v", err)
	}
	return producer, nil
}

func sendToKafka(producer sarama.SyncProducer, msg *SeckillMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: "seckill_orders",
		Value: sarama.StringEncoder(msgBytes),
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", msg.UserId)),
	}

	_, _, err = producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

// benchmarkCleanupTestData 清理基准测试数据
func benchmarkCleanupTestData(goodsId, goodsOptionsId int) {
	ctx := context.Background()

	// 获取所有相关的键
	keys := []string{
		fmt.Sprintf("%s*", consts.SeckillTokenBucketPrefix),
		fmt.Sprintf("%s*", consts.SeckillLeakyBucketPrefix),
		fmt.Sprintf("%s*", consts.SeckillStockPrefix),
		fmt.Sprintf("%s*", consts.SeckillLockPrefix),
		fmt.Sprintf("%s*", consts.SeckillResultPrefix),
		fmt.Sprintf("%s*", consts.SeckillMetricsPrefix),
		fmt.Sprintf("%s*", consts.SeckillGoodsPrefix),
		fmt.Sprintf("%s*", consts.SeckillQueuePrefix),
		fmt.Sprintf("%s*", consts.SeckillSuccessPrefix),
		fmt.Sprintf("%s*", consts.SeckillUserBoughtPrefix),
	}

	// 删除测试数据，但保留商品信息
	for _, keyPattern := range keys {
		delKeys, err := g.Redis().Do(ctx, "KEYS", keyPattern)
		if err != nil {
			g.Log().Error(ctx, "failed to get keys for pattern", keyPattern, ":", err)
			continue
		}

		keyArray := delKeys.Array()
		if len(keyArray) == 0 {
			continue
		}

		// 准备删除参数
		delArgs := []interface{}{"DEL"}
		for _, k := range keyArray {
			key := gconv.String(k)
			// 跳过商品信息和库存键，避免影响后续测试
			if key == fmt.Sprintf("%s%d:%d", consts.SeckillGoodsPrefix, goodsId, goodsOptionsId) ||
				key == fmt.Sprintf("%s%d:%d", consts.SeckillStockPrefix, goodsId, goodsOptionsId) {
				continue
			}
			delArgs = append(delArgs, key)
		}

		// 如果只有DEL命令没有键，跳过
		if len(delArgs) <= 1 {
			continue
		}

		// 批量删除键
		_, err = g.Redis().Do(ctx, "DEL", delArgs[1:]...)
		if err != nil {
			g.Log().Error(ctx, "failed to delete keys for pattern", keyPattern, ":", err)
		}
	}

	// 重置库存和计数器
	goodsKey := fmt.Sprintf("%s%d:%d", consts.SeckillGoodsPrefix, goodsId, goodsOptionsId)
	stockKey := fmt.Sprintf("%s%d:%d", consts.SeckillStockPrefix, goodsId, goodsOptionsId)
	successKey := fmt.Sprintf("%s%d:%d", consts.SeckillSuccessPrefix, goodsId, goodsOptionsId)

	// 获取商品信息
	goodsInfo, err := g.Redis().Do(ctx, "GET", goodsKey)
	if err != nil || goodsInfo.IsNil() {
		g.Log().Warning(ctx, "商品信息不存在，无法恢复库存")
		return
	}

	// 重置库存和成功计数
	_, err = g.Redis().Do(ctx, "SET", stockKey, consts.SeckillDefaultStock)
	if err != nil {
		g.Log().Error(ctx, "重置库存失败:", err)
	}

	_, err = g.Redis().Do(ctx, "SET", successKey, 0)
	if err != nil {
		g.Log().Error(ctx, "重置成功计数失败:", err)
	}

	fmt.Println("测试数据清理完成")
}

// BenchmarkSeckill 秒杀系统性能测试
func BenchmarkSeckill(b *testing.B) {
	// 初始化测试参数
	goodsId := 1
	goodsOptionsId := 1
	concurrency := 100        // 并发协程数
	requestsPerRoutine := 100 // 每个协程发送的请求数

	// 测试开始前重置计数器
	var (
		successCount int64 = 0
		failCount    int64 = 0
		errorStats   sync.Map
		minLatency   int64 = int64(time.Hour)
		maxLatency   int64 = 0
		totalLatency int64 = 0
		requestCount int64 = 0
	)

	// 并发执行秒杀请求
	start := time.Now()
	var wg sync.WaitGroup

	// 创建多个协程并发测试
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(routineId int) {
			defer wg.Done()

			// 每个协程发送多个请求
			for j := 0; j < requestsPerRoutine; j++ {
				// 随机生成用户ID (1-10000)
				userId := rand.Intn(10000) + 1

				// 构建请求
				req := &model.SeckillDoInput{
					GoodsId:        uint(goodsId),
					GoodsOptionsId: uint(goodsOptionsId),
					UserId:         uint(userId),
				}

				// 记录请求开始时间
				requestStart := time.Now()

				// 发送秒杀请求
				_, err := service.Seckill().DoSeckill(context.Background(), req)

				// 计算请求延迟
				latency := time.Since(requestStart)
				atomic.AddInt64(&requestCount, 1)
				atomic.AddInt64(&totalLatency, int64(latency))

				// 更新最大/最小延迟
				for {
					currentMin := atomic.LoadInt64(&minLatency)
					if int64(latency) >= currentMin || atomic.CompareAndSwapInt64(&minLatency, currentMin, int64(latency)) {
						break
					}
				}

				for {
					currentMax := atomic.LoadInt64(&maxLatency)
					if int64(latency) <= currentMax || atomic.CompareAndSwapInt64(&maxLatency, currentMax, int64(latency)) {
						break
					}
				}

				// 统计成功/失败数
				if err != nil {
					atomic.AddInt64(&failCount, 1)
					// 统计错误类型
					errMsg := err.Error()
					if _, ok := errorStats.Load(errMsg); !ok {
						errorStats.Store(errMsg, int64(1))
					} else {
						errorStats.Store(errMsg, func() int64 {
							if val, ok := errorStats.Load(errMsg); ok {
								return val.(int64) + 1
							}
							return int64(1)
						}())
					}
				} else {
					atomic.AddInt64(&successCount, 1)
				}

				// 随机休眠一小段时间，模拟真实用户行为
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			}
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()

	// 计算耗时和QPS
	duration := time.Since(start)
	qps := float64(requestCount) / duration.Seconds()

	// 计算成功率
	successRate := float64(successCount) / float64(requestCount) * 100

	// 输出测试结果
	fmt.Printf("\n=== 秒杀系统性能测试结果 ===\n")
	fmt.Printf("总请求数: %d\n", requestCount)
	fmt.Printf("成功请求数: %d\n", successCount)
	fmt.Printf("失败请求数: %d\n", failCount)
	fmt.Printf("总测试时间: %v\n", duration)
	fmt.Printf("平均延迟: %v\n", time.Duration(totalLatency/requestCount))
	fmt.Printf("最大延迟: %v\n", time.Duration(maxLatency))
	fmt.Printf("最小延迟: %v\n", time.Duration(minLatency))
	fmt.Printf("每秒请求数(QPS): %.2f\n", qps)
	fmt.Printf("成功率: %.2f%%\n", successRate)

	fmt.Printf("\n错误统计:\n")
	errorStats.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d\n", key, value)
		return true
	})

	// 输出Redis统计信息
	fmt.Println("\n正在清理测试数据...")
	benchmarkCleanupTestData(goodsId, goodsOptionsId)
}
