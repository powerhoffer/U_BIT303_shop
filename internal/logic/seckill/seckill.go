package seckill

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"math/rand"
// 	"runtime"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"bit303_shop/internal/consts"
// 	"bit303_shop/internal/dao"
// 	"bit303_shop/internal/model"
// 	"bit303_shop/internal/model/entity"
// 	"bit303_shop/internal/service"
// 	"bit303_shop/utility/seckill"

// 	"github.com/IBM/sarama"
// 	"github.com/gogf/gf/v2/database/gdb"
// 	"github.com/gogf/gf/v2/errors/gerror"
// 	"github.com/gogf/gf/v2/frame/g"
// 	"github.com/gogf/gf/v2/os/gctx"
// 	"github.com/gogf/gf/v2/os/gtime"
// 	"github.com/gogf/gf/v2/util/gconv"
// 	"golang.org/x/sync/errgroup"
// )

// // 秒杀请求结构
// type seckillRequest struct {
// 	ctx        context.Context             // 请求上下文
// 	input      *model.SeckillDoInput       // 请求输入
// 	resultCh   chan *model.SeckillDoOutput // 结果通道
// 	createTime time.Time                   // 请求创建时间
// }

// // 秒杀系统指标
// type seckillMetrics struct {
// 	duration int64 // 平均处理时间(毫秒)
// 	errors   int64 // 错误数
// 	success  int64 // 成功数
// 	failures int64 // 失败数
// }

// // 秒杀业务实现
// type sSeckill struct {
// 	// 64位字段放在前面，确保对齐
// 	metrics *seckillMetrics

// 	// 指针和接口类型
// 	ntpTime       *gtime.Time             // 时间同步
// 	tokenBucket   *seckill.TokenBucket    // 令牌桶限流器
// 	leakyBucket   *seckill.LeakyBucket    // 漏桶限流器
// 	stockManager  *seckill.StockManager   // 库存管理器
// 	breaker       *seckill.CircuitBreaker // 系统熔断器
// 	kafkaProducer sarama.SyncProducer     // 消息队列生产者
// 	bloomFilter   *seckill.BloomFilter    // 布隆过滤器
// 	// cacheManager  *seckill.CacheManager   // 缓存管理器

// 	// 通道和整型
// 	requestChan  chan *seckillRequest // 请求队列
// 	responseChan chan *seckillRequest // 响应队列
// 	stopCh       chan struct{}        // 关闭信号
// 	workerCount  int32                // 工作线程数

// 	// 布尔值和配置
// 	kafkaEnabled bool                // Kafka是否启用
// 	config       model.SeckillConfig // 配置信息

// 	// 同步等待组
// 	wg sync.WaitGroup // 等待所有goroutine结束的等待组

// 	// 原子值
// 	isRunning atomic.Bool // 是否运行中
// }

// // 单例
// var (
// 	instance *sSeckill
// 	once     sync.Once
// )

// // initDatabaseConfig 初始化数据库和Redis配置
// func initDatabaseConfig(ctx context.Context) {
// 	// 1. 检查数据库配置
// 	dbConfig := g.Cfg().MustGet(ctx, "database.default.link").String()
// 	if dbConfig == "" || strings.Contains(dbConfig, "mysql:mysql:") {
// 		// 如果连接配置错误，提供修复建议
// 		g.Log().Warning(ctx, "数据库连接配置有误，当前配置:", dbConfig)

// 		// 推荐正确配置
// 		correctLink := "mysql:root:111111@tcp(127.0.0.1:3306)/shop"
// 		g.Log().Warning(ctx, "请在manifest/config/config.yaml中修改数据库配置为:")
// 		g.Log().Warning(ctx, "database:\n  default:\n    link: \""+correctLink+"\"")

// 		// 尝试强制使用正确的配置连接数据库（不修改配置文件）
// 		// 注意：这只是一个临时解决方案，不会持久化到配置文件中
// 		g.Log().Info(ctx, "将尝试使用默认连接串进行数据库操作")
// 	} else {
// 		g.Log().Debug(ctx, "数据库配置有效")
// 	}

// 	// 2. 检查Redis配置
// 	redisAddress := g.Cfg().MustGet(ctx, "redis.default.address").String()
// 	if redisAddress == "" {
// 		g.Log().Warning(ctx, "Redis配置缺失，请在配置文件中设置")
// 		g.Log().Warning(ctx, "redis:\n  default:\n    address: 127.0.0.1:6379\n    db: 1")

// 		// 尝试使用默认配置
// 		g.Log().Info(ctx, "将尝试使用默认地址127.0.0.1:6379连接Redis")
// 	} else {
// 		g.Log().Debug(ctx, "当前Redis配置:", redisAddress)
// 	}
// }

// // 创建秒杀系统实例
// func New() *sSeckill {
// 	once.Do(func() {
// 		// 计算工作池大小（CPU核心数的2倍，至少4个，最多32个）
// 		poolSize := runtime.NumCPU() * 2
// 		if poolSize < consts.SeckillMinWorkers {
// 			poolSize = consts.SeckillMinWorkers
// 		}
// 		if poolSize > consts.SeckillMaxWorkers {
// 			poolSize = consts.SeckillMaxWorkers
// 		}

// 		// 初始化数据库和Redis配置
// 		ctx := context.Background()
// 		initDatabaseConfig(ctx)

// 		// 创建布隆过滤器（位图大小100万，5个哈希函数）
// 		bloomFilter := seckill.NewBloomFilter("seckill:bloom", 1000000, 5)

// 		instance = &sSeckill{
// 			metrics: &seckillMetrics{
// 				duration: 0,
// 				errors:   0,
// 				success:  0,
// 				failures: 0,
// 			},
// 			ntpTime:      gtime.Now(),
// 			tokenBucket:  seckill.NewTokenBucket(consts.SeckillTokenBucketSize, consts.SeckillTokenRate),
// 			leakyBucket:  seckill.NewLeakyBucket(consts.SeckillLeakyBucketSize, consts.SeckillLeakyBucketRate),
// 			stockManager: seckill.NewStockManager(consts.SeckillDefaultStock),
// 			breaker:      seckill.NewCircuitBreaker("seckill", 100, consts.SeckillBreakerTimeout),
// 			bloomFilter:  bloomFilter,
// 			// cacheManager: seckill.NewCacheManager("seckill:", 10*time.Minute),
// 			requestChan:  make(chan *seckillRequest, 10000),
// 			responseChan: make(chan *seckillRequest, 10000),
// 			workerCount:  int32(poolSize),
// 			kafkaEnabled: false,
// 			stopCh:       make(chan struct{}),
// 			config: model.SeckillConfig{
// 				EnableKafka:             false,
// 				EnableTokenBucket:       true,
// 				EnableLeakyBucket:       true,
// 				EnableCircuitBreaker:    true,
// 				WorkerCount:             poolSize,
// 				QueueSize:               10000,
// 				TokenBucketSize:         consts.SeckillTokenBucketSize,
// 				TokenRate:               consts.SeckillTokenRate,
// 				LeakyBucketSize:         consts.SeckillLeakyBucketSize,
// 				LeakyRate:               consts.SeckillLeakyBucketRate,
// 				CircuitBreakerThreshold: 100,
// 			},
// 		}

// 		// 启动工作池
// 		instance.startWorkers()

// 		// 标记为正在运行
// 		instance.isRunning.Store(true)

// 		// 尝试初始化Kafka
// 		instance.initKafka()

// 		// 启动库存同步工作协程
// 		instance.startStockSyncWorker()

// 		// 预热缓存 - 使用默认值0作为goodsId和optionsId
// 		instance.WarmUpCache(ctx, 0, 0)
// 	})
// 	return instance
// }

// // initKafka 初始化Kafka生产者
// func (s *sSeckill) initKafka() {
// 	// 创建Kafka生产者配置
// 	config := sarama.NewConfig()
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Return.Successes = true
// 	config.Producer.Timeout = 5 * time.Second

// 	// 获取Kafka配置路径
// 	ctx := context.Background()

// 	// 尝试从配置中读取Kafka生产者配置
// 	if requiredAcks := g.Cfg().MustGet(ctx, "mq.kafka.producer.requiredAcks").String(); requiredAcks != "" {
// 		if requiredAcks == "all" {
// 			config.Producer.RequiredAcks = sarama.WaitForAll
// 		} else if requiredAcks == "one" {
// 			config.Producer.RequiredAcks = sarama.RequiredAcks(1) // 等待leader确认
// 		} else if requiredAcks == "none" {
// 			config.Producer.RequiredAcks = sarama.NoResponse
// 		}
// 	}

// 	if timeout := g.Cfg().MustGet(ctx, "mq.kafka.producer.timeout").Int(); timeout > 0 {
// 		config.Producer.Timeout = time.Duration(timeout) * time.Second
// 	}

// 	// 从配置文件读取Kafka地址
// 	kafkaAddrs := g.Cfg().MustGet(ctx, "mq.kafka.brokers").Strings()

// 	// 如果配置文件中没有，给出提示
// 	if len(kafkaAddrs) == 0 {
// 		g.Log().Warning(ctx, "Kafka地址未配置，将不使用Kafka")
// 		g.Log().Info(ctx, "=== Kafka配置诊断 ===")
// 		g.Log().Info(ctx, "如需使用Kafka，请在配置文件(manifest/config/config.yaml)中添加以下配置:")
// 		g.Log().Info(ctx, "mq:\n  kafka:\n    brokers: [\"localhost:9092\"]\n    topics:\n      seckill_orders: \"seckill_orders\"\n    producer:\n      requiredAcks: \"all\"\n      timeout: 5")
// 		return
// 	}

// 	// 记录配置的Kafka地址
// 	g.Log().Info(ctx, "Kafka地址配置:", kafkaAddrs)

// 	// 尝试连接Kafka
// 	producer, err := sarama.NewSyncProducer(kafkaAddrs, config)
// 	if err != nil {
// 		g.Log().Error(ctx, "连接Kafka失败:", err)
// 		g.Log().Info(ctx, "=== Kafka连接诊断 ===")
// 		g.Log().Info(ctx, "1. 请确保Kafka服务已启动")
// 		g.Log().Info(ctx, "2. 检查配置的地址是否正确")
// 		g.Log().Info(ctx, "3. 检查网络连接")
// 		g.Log().Info(ctx, "4. 检查防火墙设置")
// 		return
// 	}

// 	// 检查秒杀订单主题是否已存在
// 	adminClient, err := sarama.NewClusterAdmin(kafkaAddrs, config)
// 	if err == nil {
// 		defer adminClient.Close()

// 		seckillOrderTopic := "seckill_orders"
// 		configTopic := g.Cfg().MustGet(ctx, "mq.kafka.topics.seckill_orders").String()
// 		if configTopic != "" {
// 			seckillOrderTopic = configTopic
// 		}

// 		topics, err := adminClient.ListTopics()
// 		if err == nil {
// 			if _, exists := topics[seckillOrderTopic]; !exists {
// 				g.Log().Infof(ctx, "未找到秒杀订单主题[%s]，尝试创建...", seckillOrderTopic)
// 				// 尝试创建主题
// 				topicDetail := &sarama.TopicDetail{
// 					NumPartitions:     3, // 默认3个分区
// 					ReplicationFactor: 1, // 单节点测试环境使用1
// 				}
// 				err = adminClient.CreateTopic(seckillOrderTopic, topicDetail, false)
// 				if err != nil {
// 					g.Log().Warning(ctx, "创建主题失败:", err)
// 				} else {
// 					g.Log().Info(ctx, "成功创建秒杀订单主题:", seckillOrderTopic)
// 				}
// 			} else {
// 				g.Log().Info(ctx, "秒杀订单主题已存在:", seckillOrderTopic)
// 			}
// 		}
// 	}

// 	s.kafkaProducer = producer
// 	s.kafkaEnabled = true
// 	s.config.EnableKafka = true

// 	// 输出秒杀订单主题名称
// 	seckillOrderTopic := g.Cfg().MustGet(ctx, "mq.kafka.topics.seckill_orders").String()
// 	if seckillOrderTopic == "" {
// 		seckillOrderTopic = "seckill_orders" // 默认主题名
// 	}
// 	g.Log().Info(ctx, "Kafka连接成功，秒杀订单主题:", seckillOrderTopic)
// }

// // StartWorkers 启动工作线程池
// func (s *sSeckill) startWorkers() {
// 	workerCount := int(s.workerCount)
// 	s.wg.Add(workerCount)

// 	// 启动工作线程
// 	for i := 0; i < workerCount; i++ {
// 		go s.startWorker(i)
// 	}

// 	// 启动结果处理线程
// 	s.wg.Add(1)
// 	go s.processResponses()

// 	g.Log().Infof(context.Background(), "秒杀系统已启动，工作线程数: %d", workerCount)
// }

// // 工作线程启动
// func (s *sSeckill) startWorker(workerId int) {
// 	defer s.wg.Done()

// 	g.Log().Debug(context.Background(), "工作线程启动: ", workerId)

// 	for {
// 		select {
// 		case req := <-s.requestChan:
// 			// 处理请求
// 			s.processSeckillRequest(req)
// 		case <-s.stopCh:
// 			// 收到停止信号，退出
// 			g.Log().Debug(context.Background(), "工作线程退出: ", workerId)
// 			return
// 		}
// 	}
// }

// // 处理结果线程
// func (s *sSeckill) processResponses() {
// 	defer s.wg.Done()

// 	g.Log().Debug(context.Background(), "结果处理线程启动")

// 	for {
// 		select {
// 		case resp := <-s.responseChan:
// 			// 发送结果到客户端
// 			select {
// 			case resp.resultCh <- nil: // 假设结果已通过其他方式发送
// 			default:
// 				// 客户端可能已经超时或不再等待
// 			}
// 		case <-s.stopCh:
// 			// 收到停止信号，退出
// 			g.Log().Debug(context.Background(), "结果处理线程退出")
// 			return
// 		}
// 	}
// }

// // 停止秒杀系统
// func (s *sSeckill) Stop() {
// 	if !s.isRunning.Load() {
// 		return
// 	}

// 	// 标记为停止状态
// 	s.isRunning.Store(false)

// 	// 关闭资源
// 	close(s.stopCh)

// 	// 等待所有工作线程结束
// 	s.wg.Wait()

// 	// 关闭通道
// 	close(s.requestChan)
// 	close(s.responseChan)

// 	// 关闭限流器
// 	s.tokenBucket.Close()
// 	s.leakyBucket.Close()

// 	// 关闭Kafka生产者
// 	if s.kafkaEnabled && s.kafkaProducer != nil {
// 		_ = s.kafkaProducer.Close()
// 	}

// 	g.Log().Info(context.Background(), "秒杀系统已停止")
// }

// // 添加秒杀请求到处理队列
// func (s *sSeckill) addRequest(req *seckillRequest) bool {
// 	if !s.isRunning.Load() {
// 		return false
// 	}

// 	// 检查熔断器状态
// 	if s.config.EnableCircuitBreaker && s.breaker.GetState() == seckill.StateOpen {
// 		return false
// 	}

// 	// 令牌桶限流
// 	if s.config.EnableTokenBucket && !s.tokenBucket.Take() {
// 		return false
// 	}

// 	// 漏桶限流
// 	if s.config.EnableLeakyBucket && !s.leakyBucket.Take() {
// 		return false
// 	}

// 	// 尝试放入请求队列
// 	select {
// 	case s.requestChan <- req:
// 		return true
// 	default:
// 		// 队列已满，拒绝请求
// 		return false
// 	}
// }

// // 处理秒杀请求(worker调用)
// func (s *sSeckill) processSeckillRequest(req *seckillRequest) {
// 	// 使用布隆过滤器检查请求是否有效
// 	requestKey := fmt.Sprintf("%d:%d:%d", req.input.UserId, req.input.GoodsId, req.input.GoodsOptionsId)
// 	exists, err := s.bloomFilter.Exists(req.ctx, requestKey)
// 	if err != nil {
// 		g.Log().Error(req.ctx, "布隆过滤器检查失败:", err)
// 		// 布隆过滤器出错时，继续处理请求
// 	} else if !exists {
// 		// 请求不在布隆过滤器中，可能是无效请求
// 		g.Log().Warning(req.ctx, "请求可能无效:", requestKey)
// 		// 将请求添加到布隆过滤器
// 		_ = s.bloomFilter.Add(req.ctx, requestKey)
// 	}

// 	// 直接执行秒杀
// 	resp := s.doSeckillDirectly(req)

// 	// 尝试通过结果通道返回
// 	select {
// 	case req.resultCh <- resp:
// 		// 成功发送结果
// 	default:
// 		// 结果通道已满或已关闭，记录警告
// 		g.Log().Warning(req.ctx, "无法发送秒杀结果到客户端通道")
// 	}
// }

// // DoSeckill 执行秒杀
// func (s *sSeckill) DoSeckill(ctx context.Context, input *model.SeckillDoInput) (output *model.SeckillDoOutput, err error) {
// 	// 1. 请求预检查
// 	// 检查熔断器状态
// 	if s.config.EnableCircuitBreaker && s.breaker.GetState() == seckill.StateOpen {
// 		return s.createErrorResponse(input, consts.CodeSeckillCircuitOpen, "系统熔断，请稍后重试")
// 	}

// 	// 2. 幂等性检查 - 检查是否已经处理过相同的请求ID
// 	cacheKey := fmt.Sprintf("%s%s", consts.SeckillResultPrefix, input.RequestId)
// 	var cachedResult model.SeckillDoOutput

// 	// 尝试从缓存获取结果
// 	err = seckill.GetCache(ctx, cacheKey, &cachedResult)
// 	if err == nil {
// 		// 缓存命中，直接返回缓存的结果
// 		return &cachedResult, nil
// 	}

// 	// 3. 创建结果通道和请求
// 	resultCh := make(chan *model.SeckillDoOutput, 1)
// 	seckillReq := &seckillRequest{
// 		ctx:        ctx,
// 		input:      input,
// 		resultCh:   resultCh,
// 		createTime: time.Now(),
// 	}

// 	// 4. 添加到处理队列
// 	if !s.addRequest(seckillReq) {
// 		// 添加失败，返回系统繁忙
// 		return s.createErrorResponse(input, consts.CodeSeckillRateLimited, "系统繁忙，请稍后重试")
// 	}

// 	// 5. 等待处理结果，带超时
// 	select {
// 	case result := <-resultCh:
// 		// 接收到处理结果
// 		return result, nil
// 	case <-time.After(consts.SeckillProcessTimeout):
// 		// 处理超时
// 		return s.createErrorResponse(input, consts.CodeSeckillTimeout, "处理超时，请稍后重试")
// 	}
// }

// // doSeckillDirectly 直接执行秒杀逻辑 (同步模式)
// func (s *sSeckill) doSeckillDirectly(req *seckillRequest) *model.SeckillDoOutput {
// 	startTime := time.Now()
// 	ctx := req.ctx
// 	input := req.input

// 	// 结果初始化
// 	result := &model.SeckillDoOutput{
// 		RequestId:    input.RequestId,
// 		UserId:       input.UserId,
// 		GoodsId:      input.GoodsId,
// 		Count:        input.Count,
// 		Status:       consts.CodeSeckillFailed, // 默认失败
// 		CreatedAt:    time.Now(),
// 		IsProcessing: false,
// 	}

// 	// 请求处理
// 	defer func() {
// 		// 计算处理时间
// 		processTime := time.Since(startTime).Milliseconds()
// 		result.ProcessTime = processTime

// 		// 更新统计指标
// 		if result.Status == consts.CodeSeckillSuccess {
// 			atomic.AddInt64(&s.metrics.success, 1)
// 		} else {
// 			atomic.AddInt64(&s.metrics.failures, 1)
// 		}

// 		// 计算平均处理时间 (使用简单EMA算法)
// 		currentAvg := atomic.LoadInt64(&s.metrics.duration)
// 		newAvg := int64(float64(currentAvg)*0.8 + float64(processTime)*0.2)
// 		atomic.StoreInt64(&s.metrics.duration, newAvg)

// 		// 缓存结果用于幂等性检查
// 		cacheKey := fmt.Sprintf("%s%s", consts.SeckillResultPrefix, input.RequestId)
// 		_ = seckill.SetCache(ctx, cacheKey, result, 1800*time.Second)

// 		// 记录秒杀记录到数据库
// 		_ = s.recordSeckillAttempt(ctx, input, result)
// 	}()

// 	// 1. 检查熔断器状态
// 	if s.config.EnableCircuitBreaker && s.breaker.GetState() == seckill.StateOpen {
// 		result.Status = consts.CodeSeckillCircuitOpen
// 		result.Message = "系统熔断，请稍后重试"
// 		s.recordBreakerFailure()
// 		return result
// 	}

// 	// 2. 库存检查和扣减
// 	_, err := s.stockManager.DeductStock(ctx, int64(input.UserId), int32(input.GoodsId), int32(input.GoodsOptionsId), int32(input.Count))
// 	if err != nil {
// 		// 库存不足
// 		result.Status = consts.CodeSeckillNoStock
// 		result.Message = "商品库存不足"
// 		s.recordBreakerFailure()
// 		return result
// 	}

// 	// 使用errgroup进行并发控制
// 	eg, egCtx := errgroup.WithContext(ctx)

// 	// 设置订单号
// 	orderNo := s.generateOrderNo(input.UserId)
// 	result.OrderNo = orderNo

// 	// 3. 创建订单消息
// 	orderMsg := &model.SeckillOrderMsg{
// 		OrderNo:        orderNo,
// 		RequestId:      input.RequestId,
// 		UserId:         input.UserId,
// 		GoodsId:        input.GoodsId,
// 		GoodsOptionsId: input.GoodsOptionsId,
// 		Count:          input.Count,
// 		TotalPrice:     0, // 需要查询数据库获取价格
// 		Status:         1, // 待支付
// 		CreatedAt:      time.Now(),
// 		UserAddress:    input.UserAddress,
// 		UserPhone:      input.UserPhone,
// 		Remark:         input.Remark,
// 	}

// 	// 4. 查询商品价格
// 	eg.Go(func() error {
// 		goods, err := dao.GoodsInfo.Ctx(egCtx).Where("id", input.GoodsId).One()
// 		if err != nil {
// 			return err
// 		}
// 		if goods == nil {
// 			return errors.New("商品不存在")
// 		}

// 		// 计算总价
// 		goodsOptions, err := dao.GoodsOptionsInfo.Ctx(egCtx).Where("id", input.GoodsOptionsId).One()
// 		if err != nil {
// 			return err
// 		}
// 		if goodsOptions == nil {
// 			return errors.New("商品选项不存在")
// 		}

// 		// 使用选项价格，从map中读取price字段
// 		price := gconv.Float64(goodsOptions["price"])
// 		orderMsg.TotalPrice = price * gconv.Float64(input.Count)
// 		orderMsg.Price = price
// 		return nil
// 	})

// 	// 等待所有goroutine完成
// 	if err := eg.Wait(); err != nil {
// 		// 价格查询失败，回滚库存
// 		_, _ = s.stockManager.AddStock(ctx, int32(input.GoodsId), int32(input.GoodsOptionsId), int32(input.Count))

// 		result.Status = consts.CodeSeckillSystemError
// 		result.Message = "系统错误: " + err.Error()
// 		s.recordBreakerFailure()
// 		return result
// 	}

// 	// 5. 创建订单
// 	if s.kafkaEnabled && s.kafkaProducer != nil {
// 		// 异步创建订单 (通过Kafka)
// 		err = s.sendOrderMessage(ctx, orderMsg)
// 		if err != nil {
// 			// 发送失败，回滚库存
// 			_, _ = s.stockManager.AddStock(ctx, int32(input.GoodsId), int32(input.GoodsOptionsId), int32(input.Count))

// 			result.Status = consts.CodeSeckillSystemError
// 			result.Message = "消息发送失败: " + err.Error()
// 			s.recordBreakerFailure()
// 			return result
// 		}

// 		// Kafka方式下，设置处理中状态
// 		result.Status = consts.CodeSeckillSuccess
// 		result.Message = "秒杀请求已接受，正在处理订单"
// 		result.IsProcessing = true
// 		s.resetBreakerFailure()
// 	} else {
// 		// 直接创建订单
// 		err = s.createOrderDirectly(ctx, orderMsg)
// 		if err != nil {
// 			// 创建失败，回滚库存
// 			_, _ = s.stockManager.AddStock(ctx, int32(input.GoodsId), int32(input.GoodsOptionsId), int32(input.Count))

// 			result.Status = consts.CodeSeckillSystemError
// 			result.Message = "订单创建失败: " + err.Error()
// 			s.recordBreakerFailure()
// 			return result
// 		}

// 		// 秒杀成功
// 		result.Status = consts.CodeSeckillSuccess
// 		result.Message = "秒杀成功，秒杀订单已创建"
// 		s.resetBreakerFailure()
// 	}

// 	// 同步库存到数据库
// 	go func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				fmt.Printf("同步商品[%d:%d]库存时panic: %v\n", input.GoodsId, input.GoodsOptionsId, r)
// 			}
// 		}()

// 		// 使用新的上下文，因为原上下文可能已关闭
// 		syncCtx := context.Background()

// 		// 延迟几秒再执行同步，确保其他并发请求先完成
// 		time.Sleep(2 * time.Second)

// 		err := s.stockManager.SyncStockToDatabase(syncCtx, int32(input.GoodsId), int32(input.GoodsOptionsId))
// 		if err != nil {
// 			g.Log().Error(syncCtx, "同步商品库存失败:", input.GoodsId, input.GoodsOptionsId, err)
// 		} else {
// 			g.Log().Info(syncCtx, "成功同步商品库存:", input.GoodsId, input.GoodsOptionsId)
// 		}
// 	}()

// 	return result
// }

// // sendOrderMessage 发送订单消息到Kafka
// func (s *sSeckill) sendOrderMessage(ctx context.Context, orderMsg *model.SeckillOrderMsg) error {
// 	// 检查是否已发送过该订单消息(幂等性处理)
// 	sentKey := fmt.Sprintf("%s%s", consts.SeckillOrderSentPrefix, orderMsg.OrderNo)
// 	exists, err := seckill.Exists(ctx, sentKey)
// 	if err != nil {
// 		return err
// 	}

// 	if exists {
// 		g.Log().Info(ctx, "订单消息已发送，跳过", orderMsg.OrderNo)
// 		return nil
// 	}

// 	// 序列化消息
// 	msgData, err := json.Marshal(orderMsg)
// 	if err != nil {
// 		return err
// 	}

// 	// 使用专用的秒杀订单主题
// 	seckillOrderTopic := "seckill_orders" // 独立的秒杀订单主题

// 	// 尝试从配置中获取秒杀订单主题
// 	configTopic := g.Cfg().MustGet(ctx, "mq.kafka.topics.seckill_orders").String()
// 	if configTopic != "" {
// 		seckillOrderTopic = configTopic
// 		g.Log().Debug(ctx, "使用配置中的秒杀订单主题:", seckillOrderTopic)
// 	}

// 	// 准备Kafka消息 - 使用新的topic名称以区分普通订单
// 	kafkaMsg := &sarama.ProducerMessage{
// 		Topic: seckillOrderTopic,
// 		Key:   sarama.StringEncoder(orderMsg.OrderNo),
// 		Value: sarama.ByteEncoder(msgData),
// 	}

// 	// 发送消息前记录发送尝试
// 	g.Log().Infof(ctx, "发送秒杀订单消息: 订单号=%s, 用户ID=%d, 商品ID=%d",
// 		orderMsg.OrderNo, orderMsg.UserId, orderMsg.GoodsId)

// 	// 发送消息
// 	partition, offset, err := s.kafkaProducer.SendMessage(kafkaMsg)
// 	if err != nil {
// 		g.Log().Error(ctx, "发送秒杀订单消息失败:", err)
// 		return err
// 	}

// 	g.Log().Infof(ctx, "秒杀订单消息发送成功: 订单号=%s, 分区=%d, 偏移量=%d",
// 		orderMsg.OrderNo, partition, offset)

// 	// 标记为已发送
// 	err = seckill.SetCache(ctx, sentKey, "1", 24*time.Hour)
// 	if err != nil {
// 		g.Log().Warning(ctx, "标记订单消息发送状态失败", err)
// 	}

// 	return nil
// }

// // createOrderDirectly 直接创建订单 (不通过Kafka)
// func (s *sSeckill) createOrderDirectly(ctx context.Context, orderMsg *model.SeckillOrderMsg) error {
// 	g.Log().Info(ctx, "开始创建秒杀订单, 订单号:", orderMsg.OrderNo)

// 	// 检查订单是否已存在，避免重复创建
// 	orderExists, err := dao.SeckillOrder.Ctx(ctx).Where("order_no", orderMsg.OrderNo).Count()
// 	if err != nil {
// 		g.Log().Error(ctx, "查询秒杀订单失败:", err)
// 		return err
// 	}

// 	if orderExists > 0 {
// 		g.Log().Info(ctx, "秒杀订单已存在，跳过创建:", orderMsg.OrderNo)
// 		return nil
// 	}

// 	// 确保商品价格正确
// 	if orderMsg.Price <= 0 {
// 		// 如果价格未设置，尝试从数据库获取
// 		goods, err := dao.SeckillGoods.Ctx(ctx).
// 			Where("goods_id", orderMsg.GoodsId).
// 			Where("goods_options_id", orderMsg.GoodsOptionsId).
// 			One()
// 		if err != nil {
// 			g.Log().Error(ctx, "查询秒杀商品价格失败:", err)
// 			return err
// 		}
// 		if goods != nil {
// 			// 使用秒杀价格
// 			price := gconv.Float64(goods["seckill_price"])
// 			orderMsg.Price = price
// 			orderMsg.TotalPrice = price * gconv.Float64(orderMsg.Count)
// 		}
// 	}

// 	// 直接创建秒杀订单记录
// 	seckillOrderInput := orderMsg.ToSeckillOrderAddInput()

// 	seckillOrderOutput, err := s.AddSeckillOrder(ctx, seckillOrderInput)
// 	if err != nil {
// 		g.Log().Error(ctx, "创建秒杀订单记录失败:", err)
// 		return err
// 	}

// 	g.Log().Info(ctx, "创建秒杀订单记录成功, ID:", seckillOrderOutput.Id)

// 	// 同步到统计信息
// 	atomic.AddInt64(&s.metrics.success, 1)

// 	return nil
// }

// // createErrorResponse 创建错误响应
// func (s *sSeckill) createErrorResponse(req *model.SeckillDoInput, code int, message string) (*model.SeckillDoOutput, error) {
// 	result := &model.SeckillDoOutput{
// 		RequestId:    req.RequestId,
// 		UserId:       req.UserId,
// 		GoodsId:      req.GoodsId,
// 		Count:        req.Count,
// 		Status:       code,
// 		Message:      message,
// 		CreatedAt:    time.Now(),
// 		IsProcessing: false,
// 	}

// 	// 记录指标
// 	atomic.AddInt64(&s.metrics.failures, 1)

// 	// 缓存结果
// 	cacheKey := fmt.Sprintf("%s%s", consts.SeckillResultPrefix, req.RequestId)
// 	_ = seckill.SetCache(context.Background(), cacheKey, result, 1800*time.Second)

// 	return result, nil
// }

// // 生成订单号
// func (s *sSeckill) generateOrderNo(userId uint) string {
// 	// 订单号格式: 时间戳 + 用户ID后4位 + 4位随机数
// 	timestamp := time.Now().Format("20060102150405")
// 	userIdStr := fmt.Sprintf("%04d", userId%10000)
// 	random := fmt.Sprintf("%04d", rand.Intn(10000))
// 	return fmt.Sprintf("%s%s%s", timestamp, userIdStr, random)
// }

// // 检查熔断器状态
// func (s *sSeckill) checkCircuitBreaker() bool {
// 	return s.config.EnableCircuitBreaker && s.breaker.GetState() == seckill.StateOpen
// }

// // 记录熔断器失败
// func (s *sSeckill) recordBreakerFailure() {
// 	if s.config.EnableCircuitBreaker && s.breaker != nil {
// 		s.breaker.RecordFailure()
// 	}
// }

// // 重置熔断器失败计数
// func (s *sSeckill) resetBreakerFailure() {
// 	if s.config.EnableCircuitBreaker && s.breaker != nil {
// 		s.breaker.RecordSuccess()
// 	}
// }

// // CheckStock 检查商品库存
// func (s *sSeckill) CheckStock(ctx context.Context, in *model.SeckillCheckStockInput) (*model.SeckillCheckStockOutput, error) {
// 	current, err := s.stockManager.CheckStock(ctx, int32(in.GoodsId), int32(in.GoodsOptionsId))
// 	if err != nil {
// 		return nil, err
// 	}

// 	required := int32(in.Count)
// 	available := current >= required

// 	return &model.SeckillCheckStockOutput{
// 		Available: available,
// 		Current:   current,
// 		Required:  required,
// 	}, nil
// }

// // InitStock 初始化秒杀商品库存
// func (s *sSeckill) InitStock(ctx context.Context, req model.SeckillInitStockInput) (res *model.SeckillInitStockOutput, err error) {
// 	// 初始化库存
// 	err = s.stockManager.InitStock(ctx, int32(req.GoodsId), int32(req.GoodsOptionsId), int32(req.Stock))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &model.SeckillInitStockOutput{Success: true}, nil
// }

// // GetStats 获取秒杀系统统计信息
// func (s *sSeckill) GetStats(ctx context.Context, goodsId, optionsId int64) (stats *model.SeckillStatsOutput, err error) {
// 	stockStats := s.stockManager.GetStats()

// 	// 获取秒杀订单统计
// 	var orderCount int = 0
// 	var successOrderCount int = 0
// 	var todayOrderCount int = 0

// 	// 计算今天的开始时间
// 	now := time.Now()
// 	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

// 	// 查询总订单数
// 	orderCount, err = dao.SeckillOrder.Ctx(ctx).Count()
// 	if err != nil {
// 		g.Log().Warning(ctx, "查询秒杀订单总数失败:", err)
// 	}

// 	// 查询成功订单数
// 	successOrderCount, err = dao.SeckillOrder.Ctx(ctx).Where("status", 1).Count()
// 	if err != nil {
// 		g.Log().Warning(ctx, "查询秒杀成功订单数失败:", err)
// 	}

// 	// 查询今日订单数
// 	todayOrderCount, err = dao.SeckillOrder.Ctx(ctx).
// 		Where("created_at >= ?", todayStart).
// 		Count()
// 	if err != nil {
// 		g.Log().Warning(ctx, "查询今日秒杀订单数失败:", err)
// 	}

// 	// 计算平均处理时间和成功率
// 	avgProcessTime := float64(atomic.LoadInt64(&s.metrics.duration))
// 	successRate := 0.0
// 	totalRequests := atomic.LoadInt64(&s.metrics.success) + atomic.LoadInt64(&s.metrics.failures)
// 	if totalRequests > 0 {
// 		successRate = float64(atomic.LoadInt64(&s.metrics.success)) / float64(totalRequests) * 100
// 	}

// 	// 获取内存使用情况
// 	var memStats runtime.MemStats
// 	runtime.ReadMemStats(&memStats)

// 	// 输出额外的统计信息到日志
// 	g.Log().Info(ctx, "=== 秒杀系统性能统计 ===")
// 	g.Log().Infof(ctx, "订单总数: %d, 成功订单: %d, 今日订单: %d",
// 		orderCount, successOrderCount, todayOrderCount)
// 	g.Log().Infof(ctx, "成功率: %.2f%%, 平均处理时间: %.2fms",
// 		successRate, avgProcessTime)
// 	g.Log().Infof(ctx, "内存使用: %.2fMB, Goroutines: %d",
// 		float64(memStats.Alloc)/1024/1024, runtime.NumGoroutine())

// 	return &model.SeckillStatsOutput{
// 		Workers:        int(s.workerCount),
// 		QueueSize:      cap(s.requestChan),
// 		QueueCurrent:   len(s.requestChan),
// 		Successes:      atomic.LoadInt64(&s.metrics.success),
// 		Failures:       atomic.LoadInt64(&s.metrics.failures),
// 		Errors:         atomic.LoadInt64(&s.metrics.errors),
// 		AvgTime:        avgProcessTime,
// 		TokenBucket:    s.tokenBucket.GetMetrics(),
// 		LeakyBucket:    s.leakyBucket.GetMetrics(),
// 		StockManager:   stockStats,
// 		CircuitBreaker: s.breaker.GetMetrics(),
// 	}, nil
// }

// // Reset 重置秒杀系统
// func (s *sSeckill) Reset() {
// 	s.tokenBucket.Reset()
// 	s.leakyBucket.Reset()
// 	s.stockManager.Reset()
// 	s.breaker.Reset()

// 	atomic.StoreInt64(&s.metrics.success, 0)
// 	atomic.StoreInt64(&s.metrics.failures, 0)
// 	atomic.StoreInt64(&s.metrics.errors, 0)
// }

// // SetConfig 设置秒杀系统配置
// func (s *sSeckill) SetConfig(config *model.SeckillConfig) {
// 	// 更新配置
// 	s.config = *config

// 	// 根据需要重新初始化组件
// 	if s.tokenBucket != nil {
// 		s.tokenBucket.Close()
// 		s.tokenBucket = seckill.NewTokenBucket(config.TokenBucketSize, config.TokenRate)
// 	}

// 	if s.leakyBucket != nil {
// 		s.leakyBucket.Close()
// 		s.leakyBucket = seckill.NewLeakyBucket(config.LeakyBucketSize, config.LeakyRate)
// 	}

// 	// 如果工作线程数变更，需要重启工作池
// 	if int(s.workerCount) != config.WorkerCount {
// 		// 停止当前工作池
// 		s.Stop()

// 		// 更新工作线程数
// 		s.workerCount = int32(config.WorkerCount)

// 		// 重新创建通道
// 		s.requestChan = make(chan *seckillRequest, config.QueueSize)
// 		s.responseChan = make(chan *seckillRequest, config.QueueSize)

// 		// 创建停止信号
// 		s.stopCh = make(chan struct{})

// 		// 重启工作池
// 		s.startWorkers()

// 		// 标记为正在运行
// 		s.isRunning.Store(true)
// 	}
// }

// // 实现service.ISeckill接口
// func init() {
// 	service.RegisterSeckill(New())
// }

// // Close 关闭秒杀系统资源
// func (s *sSeckill) Close() error {
// 	s.Stop()
// 	return nil
// }

// // GetResult 获取秒杀结果
// func (s *sSeckill) GetResult(ctx context.Context, req model.SeckillResultInput) (res *model.SeckillResultOutput, err error) {
// 	// 查询订单状态
// 	orderInfo, err := dao.OrderInfo.Ctx(ctx).Where("number", req.OrderNo).One()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if orderInfo == nil {
// 		return &model.SeckillResultOutput{
// 			OrderNo: req.OrderNo,
// 			Status:  consts.CodeSeckillProcessing,
// 		}, nil
// 	}

// 	return &model.SeckillResultOutput{
// 		OrderNo: req.OrderNo,
// 		Status:  gconv.Int(orderInfo["status"]),
// 	}, nil
// }

// // Detail 获取秒杀商品详情
// func (s *sSeckill) Detail(ctx context.Context, req *model.SeckillDetailReq) (res *model.SeckillDetailRes, err error) {
// 	// 查询秒杀商品详情
// 	seckillGoods, err := dao.SeckillGoods.Ctx(ctx).WherePri(req.Id).One()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if seckillGoods == nil {
// 		return nil, errors.New("秒杀商品不存在")
// 	}

// 	// 转换为响应
// 	res = &model.SeckillDetailRes{
// 		Id:             gconv.Uint(seckillGoods["id"]),
// 		GoodsId:        gconv.Uint(seckillGoods["goods_id"]),
// 		GoodsOptionsId: gconv.Uint(seckillGoods["goods_options_id"]),
// 		OriginalPrice:  gconv.Float64(seckillGoods["original_price"]),
// 		SeckillPrice:   gconv.Float64(seckillGoods["seckill_price"]),
// 		SeckillStock:   gconv.Int(seckillGoods["seckill_stock"]),
// 		StartTime:      gconv.String(seckillGoods["start_time"]),
// 		EndTime:        gconv.String(seckillGoods["end_time"]),
// 		Status:         gconv.Int(seckillGoods["status"]),
// 	}

// 	return res, nil
// }

// // List 获取秒杀商品列表
// func (s *sSeckill) List(ctx context.Context, req *model.SeckillListInput) (res *model.SeckillListOutput, err error) {
// 	// 查询秒杀商品列表
// 	m := dao.SeckillGoods.Ctx(ctx)

// 	// 分页查询
// 	total, err := m.Count()
// 	if err != nil {
// 		return nil, err
// 	}

// 	list, err := m.Page(req.Page, req.Size).All()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.SeckillListOutput{
// 		List:  list,
// 		Page:  req.Page,
// 		Size:  req.Size,
// 		Total: total,
// 	}, nil
// }

// // WarmUpCache 预热秒杀系统缓存
// func (s *sSeckill) WarmUpCache(ctx context.Context, goodsId int64, optionsId int64) error {
// 	g.Log().Info(ctx, "正在预热秒杀系统缓存...")

// 	// 创建后台上下文
// 	if ctx == nil {
// 		ctx = gctx.New()
// 	}

// 	// 预热过程
// 	go func() {
// 		// 加载活跃商品库存到Redis
// 		if err := s.warmUpStockCache(ctx, goodsId, optionsId); err != nil {
// 			g.Log().Error(ctx, "预热库存缓存失败:", err)
// 		}

// 		// 预热布隆过滤器
// 		if err := s.warmUpBloomFilter(ctx); err != nil {
// 			g.Log().Error(ctx, "预热布隆过滤器失败:", err)
// 		}

// 		// 预热结束
// 		g.Log().Info(ctx, "秒杀系统缓存预热完成")
// 	}()

// 	return nil
// }

// // warmUpStockCache 预热库存缓存
// func (s *sSeckill) warmUpStockCache(ctx context.Context, goodsId int64, optionsId int64) error {
// 	// 查询秒杀商品
// 	seckillGoods, err := dao.SeckillGoods.Ctx(ctx).
// 		Where("goods_id", goodsId).
// 		Where("goods_options_id", optionsId).
// 		One()

// 	if err != nil {
// 		return err
// 	}

// 	if seckillGoods == nil {
// 		return errors.New("秒杀商品不存在")
// 	}

// 	// 初始化Redis库存
// 	stock := gconv.Int32(seckillGoods["seckill_stock"])
// 	return s.stockManager.InitStock(ctx, int32(goodsId), int32(optionsId), stock)
// }

// // warmUpBloomFilter 预热布隆过滤器
// func (s *sSeckill) warmUpBloomFilter(ctx context.Context) error {
// 	// 查询所有秒杀商品
// 	var goods []map[string]interface{}
// 	err := g.DB().Model("seckill_goods").
// 		Where("status = ?", 1).
// 		Scan(&goods)
// 	if err != nil {
// 		return err
// 	}

// 	// 将商品ID添加到布隆过滤器
// 	for _, item := range goods {
// 		goodsId := gconv.Int64(item["goods_id"])
// 		goodsOptionsId := gconv.Int64(item["goods_options_id"])
// 		key := fmt.Sprintf("%d:%d", goodsId, goodsOptionsId)
// 		if err := s.bloomFilter.Add(ctx, key); err != nil {
// 			g.Log().Error(ctx, fmt.Sprintf("添加商品[%s]到布隆过滤器失败:", key), err)
// 			continue
// 		}
// 		g.Log().Info(ctx, fmt.Sprintf("商品[%s]已添加到布隆过滤器", key))
// 	}

// 	return nil
// }

// // GetSeckillInfo 获取秒杀商品信息
// func (s *sSeckill) GetSeckillInfo(ctx context.Context, goodsId int64, optionsId int64) (*entity.SeckillGoods, error) {
// 	// 使用缓存管理器生成缓存键
// 	cacheKey := fmt.Sprintf("goods:%d:%d", goodsId, optionsId)

// 	// 懒初始化缓存管理器
// 	//if s.cacheManager == nil {
// 	//	s.cacheManager = seckill.NewCacheManager("seckill:", 10*time.Minute)
// 	//}

// 	// 首先尝试从缓存获取
// 	var seckillGoods entity.SeckillGoods
// 	err := seckill.GetCache(ctx, cacheKey, &seckillGoods)
// 	if err == nil {
// 		return &seckillGoods, nil
// 	}

// 	// 缓存未命中，从数据库查询
// 	// 使用索引优化查询：使用已创建的联合索引(goods_id, goods_options_id)
// 	err = dao.SeckillGoods.Ctx(ctx).
// 		Where("goods_id", goodsId).
// 		Where("goods_options_id", optionsId).
// 		OrderDesc("id"). // 按ID倒序，获取最新的配置
// 		Scan(&seckillGoods)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if seckillGoods.Id == 0 {
// 		return nil, gerror.Newf("秒杀商品不存在: goods_id=%d, options_id=%d", goodsId, optionsId)
// 	}

// 	// 将查询结果存入缓存，缓存1小时
// 	_ = seckill.SetCache(ctx, cacheKey, seckillGoods, time.Hour)

// 	return &seckillGoods, nil
// }

// // UpdateOrderStatus 更新订单状态
// func (s *sSeckill) UpdateOrderStatus(ctx context.Context, req model.SeckillUpdateOrderStatusInput) (res *model.SeckillUpdateOrderStatusOutput, err error) {
// 	// 更新订单状态
// 	_, err = dao.OrderInfo.Ctx(ctx).
// 		Data(g.Map{"status": req.Status}).
// 		Where("number", req.OrderNo).
// 		Update()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.SeckillUpdateOrderStatusOutput{Success: true}, nil
// }

// // GetConfig 获取当前秒杀配置
// func (s *sSeckill) GetConfig() *model.SeckillConfig {
// 	config := s.config
// 	return &config
// }

// // AddSeckillGoods 添加秒杀商品
// func (s *sSeckill) AddSeckillGoods(ctx context.Context, req *model.SeckillGoodsAddInput) (res *model.SeckillGoodsAddOutput, err error) {
// 	// 检查商品是否存在
// 	goods, err := dao.GoodsInfo.Ctx(ctx).WherePri(req.GoodsId).One()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if goods == nil {
// 		return nil, gerror.New("商品不存在")
// 	}

// 	// 检查商品规格是否存在
// 	if req.GoodsOptionsId > 0 {
// 		options, err := dao.GoodsOptionsInfo.Ctx(ctx).WherePri(req.GoodsOptionsId).One()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if options == nil {
// 			return nil, gerror.New("商品规格不存在")
// 		}
// 	}

// 	// 检查时间
// 	startTime := req.StartTime.Time
// 	endTime := req.EndTime.Time
// 	if startTime.After(endTime) {
// 		return nil, gerror.New("开始时间不能晚于结束时间")
// 	}

// 	// 添加秒杀商品记录
// 	data := &entity.SeckillGoods{
// 		GoodsId:        req.GoodsId,
// 		GoodsOptionsId: req.GoodsOptionsId,
// 		OriginalPrice:  req.OriginalPrice,
// 		SeckillPrice:   req.SeckillPrice,
// 		SeckillStock:   req.SeckillStock,
// 		StartTime:      req.StartTime,
// 		EndTime:        req.EndTime,
// 		Status:         req.Status,
// 		CreatedAt:      gtime.Now(),
// 		UpdatedAt:      gtime.Now(),
// 	}

// 	id, err := dao.SeckillGoods.Ctx(ctx).Data(data).InsertAndGetId()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 初始化秒杀商品的库存到缓存
// 	stockKey := fmt.Sprintf("%d-%d", req.GoodsId, req.GoodsOptionsId)
// 	err = s.stockManager.InitStock(ctx, int32(req.GoodsId), int32(req.GoodsOptionsId), int32(req.SeckillStock))
// 	if err != nil {
// 		g.Log().Error(ctx, "初始化秒杀库存失败:", err, stockKey)
// 	}

// 	return &model.SeckillGoodsAddOutput{
// 		Id: id,
// 	}, nil
// }

// // UpdateSeckillGoods 更新秒杀商品
// func (s *sSeckill) UpdateSeckillGoods(ctx context.Context, req *model.SeckillGoodsUpdateInput) (res *model.SeckillGoodsUpdateOutput, err error) {
// 	// 检查秒杀商品是否存在
// 	seckillGoods, err := dao.SeckillGoods.Ctx(ctx).WherePri(req.Id).One()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if seckillGoods == nil {
// 		return nil, gerror.New("秒杀商品不存在")
// 	}

// 	// 准备更新数据
// 	data := g.Map{
// 		"updated_at": gtime.Now(),
// 	}

// 	// 更新价格
// 	if req.OriginalPrice > 0 {
// 		data["original_price"] = req.OriginalPrice
// 	}
// 	if req.SeckillPrice > 0 {
// 		data["seckill_price"] = req.SeckillPrice
// 	}

// 	// 更新库存
// 	if req.SeckillStock > 0 {
// 		data["seckill_stock"] = req.SeckillStock
// 		// 同步更新缓存中的库存
// 		goodsId := gconv.Int64(seckillGoods["goods_id"])
// 		optionsId := gconv.Int64(seckillGoods["goods_options_id"])
// 		err = s.stockManager.InitStock(ctx, int32(goodsId), int32(optionsId), int32(req.SeckillStock))
// 		if err != nil {
// 			g.Log().Error(ctx, "更新秒杀库存失败:", err)
// 		}
// 	}

// 	// 更新时间
// 	if req.StartTime != nil {
// 		data["start_time"] = req.StartTime
// 	}
// 	if req.EndTime != nil {
// 		data["end_time"] = req.EndTime
// 	}

// 	// 更新状态
// 	if req.Status >= 0 {
// 		data["status"] = req.Status
// 	}

// 	// 执行更新
// 	_, err = dao.SeckillGoods.Ctx(ctx).WherePri(req.Id).Data(data).Update()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.SeckillGoodsUpdateOutput{
// 		Id: req.Id,
// 	}, nil
// }

// // AddSeckillOrder 添加秒杀订单
// func (s *sSeckill) AddSeckillOrder(ctx context.Context, req *model.SeckillOrderAddInput) (res *model.SeckillOrderAddOutput, err error) {
// 	// 不再检查原订单表，直接添加秒杀订单记录
// 	data := &entity.SeckillOrder{
// 		UserId:         req.UserId,
// 		GoodsId:        req.GoodsId,
// 		GoodsOptionsId: req.GoodsOptionsId,
// 		OriginalPrice:  req.OriginalPrice,
// 		SeckillPrice:   req.SeckillPrice,
// 		Status:         req.Status,
// 		OrderNo:        req.OrderNo,        // 订单编号
// 		Count:          req.Count,          // 数量
// 		PayTime:        req.PayTime,        // 支付时间
// 		CancelTime:     req.CancelTime,     // 取消时间
// 		ConsigneeName:  req.ConsigneeName,  // 收货人
// 		ConsigneePhone: req.ConsigneePhone, // 联系电话
// 		Address:        req.Address,        // 地址
// 		Remark:         req.Remark,         // 备注
// 		CreatedAt:      gtime.Now(),
// 		UpdatedAt:      gtime.Now(),
// 	}

// 	id, err := dao.SeckillOrder.Ctx(ctx).Data(data).InsertAndGetId()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.SeckillOrderAddOutput{
// 		Id: id,
// 	}, nil
// }

// // BatchAddSeckillOrders 批量添加秒杀订单
// func (s *sSeckill) BatchAddSeckillOrders(ctx context.Context, orders []*model.SeckillOrderAddInput) ([]int64, error) {
// 	if len(orders) == 0 {
// 		return nil, nil
// 	}

// 	// 准备数据
// 	batchData := make([]g.Map, 0, len(orders))
// 	now := gtime.Now()

// 	for _, order := range orders {
// 		batchData = append(batchData, g.Map{
// 			"user_id":          order.UserId,
// 			"goods_id":         order.GoodsId,
// 			"goods_options_id": order.GoodsOptionsId,
// 			"original_price":   order.OriginalPrice,
// 			"seckill_price":    order.SeckillPrice,
// 			"status":           order.Status,
// 			"order_no":         order.OrderNo,
// 			"count":            order.Count,
// 			"pay_time":         order.PayTime,
// 			"cancel_time":      order.CancelTime,
// 			"consignee_name":   order.ConsigneeName,
// 			"consignee_phone":  order.ConsigneePhone,
// 			"address":          order.Address,
// 			"remark":           order.Remark,
// 			"created_at":       now,
// 			"updated_at":       now,
// 		})
// 	}

// 	// 批量插入
// 	result, err := dao.SeckillOrder.Ctx(ctx).
// 		Data(batchData).
// 		Batch(len(batchData)).
// 		Insert()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 获取插入ID
// 	lastInsertId, err := result.LastInsertId()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 生成ID列表
// 	ids := make([]int64, len(batchData))
// 	for i := 0; i < len(batchData); i++ {
// 		ids[i] = lastInsertId + int64(i)
// 	}

// 	return ids, nil
// }

// // UpdateSeckillOrder 更新秒杀订单
// func (s *sSeckill) UpdateSeckillOrder(ctx context.Context, req *model.SeckillOrderUpdateInput) (res *model.SeckillOrderUpdateOutput, err error) {
// 	// 检查秒杀订单是否存在
// 	seckillOrder, err := dao.SeckillOrder.Ctx(ctx).WherePri(req.Id).One()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if seckillOrder == nil {
// 		return nil, gerror.New("秒杀订单不存在")
// 	}

// 	// 准备更新数据
// 	data := g.Map{
// 		"updated_at": gtime.Now(),
// 	}

// 	// 更新状态
// 	if req.Status >= 0 {
// 		data["status"] = req.Status
// 	}

// 	// 执行更新
// 	_, err = dao.SeckillOrder.Ctx(ctx).WherePri(req.Id).Data(data).Update()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &model.SeckillOrderUpdateOutput{
// 		Id: req.Id,
// 	}, nil
// }

// // GetSeckillOrderByOrderId 根据订单ID获取秒杀订单
// func (s *sSeckill) GetSeckillOrderByOrderId(ctx context.Context, req *model.SeckillOrderByOrderIdInput) (res *model.SeckillOrderByOrderIdOutput, err error) {
// 	// 查询秒杀订单记录
// 	seckillOrder, err := dao.SeckillOrder.Ctx(ctx).Where("order_no", req.OrderId).One()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if seckillOrder == nil {
// 		return nil, gerror.New("秒杀订单不存在")
// 	}

// 	res = &model.SeckillOrderByOrderIdOutput{
// 		Id:             gconv.Int64(seckillOrder["id"]),
// 		OrderId:        gconv.Int64(seckillOrder["order_id"]),
// 		UserId:         gconv.Int64(seckillOrder["user_id"]),
// 		GoodsId:        gconv.Int64(seckillOrder["goods_id"]),
// 		GoodsOptionsId: gconv.Int64(seckillOrder["goods_options_id"]),
// 		SeckillPrice:   gconv.Int(seckillOrder["seckill_price"]),
// 		Status:         gconv.Int(seckillOrder["status"]),
// 		CreatedAt:      gtime.NewFromStr(gconv.String(seckillOrder["created_at"])),
// 		UpdatedAt:      gtime.NewFromStr(gconv.String(seckillOrder["updated_at"])),
// 	}

// 	return res, nil
// }

// // BatchGetSeckillOrders 批量获取秒杀订单
// func (s *sSeckill) BatchGetSeckillOrders(ctx context.Context, orderNos []string) (map[string]*entity.SeckillOrder, error) {
// 	if len(orderNos) == 0 {
// 		return nil, nil
// 	}

// 	var orders []*entity.SeckillOrder
// 	err := dao.SeckillOrder.Ctx(ctx).
// 		WhereIn("order_no", orderNos).
// 		Scan(&orders)

// 	if err != nil {
// 		return nil, err
// 	}

// 	result := make(map[string]*entity.SeckillOrder, len(orders))
// 	for _, order := range orders {
// 		result[order.OrderNo] = order
// 	}

// 	return result, nil
// }

// // BatchProcessSeckill 批量处理秒杀请求
// func (s *sSeckill) BatchProcessSeckill(ctx context.Context, inputs []*model.SeckillDoInput) ([]*model.SeckillDoOutput, error) {
// 	if len(inputs) == 0 {
// 		return nil, nil
// 	}

// 	startTime := time.Now()
// 	g.Log().Infof(ctx, "开始批量处理%d个秒杀请求", len(inputs))

// 	// 结果集
// 	results := make([]*model.SeckillDoOutput, len(inputs))

// 	// 先进行幂等性检查，查找已处理的请求
// 	requestIds := make([]string, len(inputs))
// 	for i, input := range inputs {
// 		requestIds[i] = input.RequestId
// 		// 初始化结果
// 		results[i] = &model.SeckillDoOutput{
// 			RequestId:    input.RequestId,
// 			UserId:       input.UserId,
// 			GoodsId:      input.GoodsId,
// 			Count:        input.Count,
// 			Status:       consts.CodeSeckillFailed, // 默认失败
// 			CreatedAt:    time.Now(),
// 			IsProcessing: false,
// 		}
// 	}

// 	// 从缓存批量获取已处理的结果
// 	cacheKeys := make([]string, len(requestIds))
// 	for i, reqId := range requestIds {
// 		cacheKeys[i] = fmt.Sprintf("%s%s", consts.SeckillResultPrefix, reqId)
// 	}

// 	// 简化处理：单个获取
// 	cachedResults := make(map[string]*model.SeckillDoOutput)
// 	for i, reqId := range requestIds {
// 		cacheKey := cacheKeys[i]
// 		var cachedResult model.SeckillDoOutput
// 		err := seckill.GetCache(ctx, cacheKey, &cachedResult)
// 		if err == nil {
// 			// 缓存命中
// 			cachedResults[reqId] = &cachedResult
// 			results[i] = &cachedResult
// 			g.Log().Debugf(ctx, "请求[%s]缓存命中，状态: %d", reqId, cachedResult.Status)
// 		}
// 	}

// 	// 筛选出未处理的请求
// 	remainingInputs := make([]*model.SeckillDoInput, 0)
// 	remainingIndexes := make([]int, 0)

// 	for i, input := range inputs {
// 		if _, ok := cachedResults[input.RequestId]; !ok {
// 			// 未找到缓存结果，需要处理
// 			remainingInputs = append(remainingInputs, input)
// 			remainingIndexes = append(remainingIndexes, i)
// 		}
// 	}

// 	if len(remainingInputs) == 0 {
// 		// 所有请求都已处理过，直接返回缓存结果
// 		g.Log().Info(ctx, "所有请求都已处理过，直接返回缓存结果")
// 		return results, nil
// 	}

// 	g.Log().Infof(ctx, "需要处理%d个新请求", len(remainingInputs))

// 	// 准备批量扣减库存的数据
// 	stockItems := make([]struct {
// 		GoodsId  int32
// 		OptionId int32
// 		Quantity int32
// 	}, len(remainingInputs))

// 	for i, input := range remainingInputs {
// 		stockItems[i] = struct {
// 			GoodsId  int32
// 			OptionId int32
// 			Quantity int32
// 		}{
// 			GoodsId:  int32(input.GoodsId),
// 			OptionId: int32(input.GoodsOptionsId),
// 			Quantity: int32(input.Count),
// 		}
// 	}

// 	// 使用任意一个用户ID作为批量扣减的标识 (可以进一步优化为分组)
// 	userId := int64(remainingInputs[0].UserId)

// 	// 批量扣减库存
// 	// 记录库存扣减结果但不使用
// 	_, err := s.stockManager.BatchDeductStock(ctx, userId, stockItems)
// 	if err != nil {
// 		// 库存扣减失败，设置所有结果为失败
// 		g.Log().Error(ctx, "批量扣减库存失败:", err)
// 		for _, i := range remainingIndexes {
// 			results[i].Status = consts.CodeSeckillNoStock
// 			results[i].Message = "商品库存不足: " + err.Error()
// 			results[i].ProcessTime = time.Since(startTime).Milliseconds()

// 			// 记录秒杀尝试
// 			_ = s.recordSeckillAttempt(ctx, inputs[i], results[i])
// 		}
// 		return results, nil
// 	}

// 	// 准备批量查询商品价格
// 	goodsOptionIds := make(map[string]bool)
// 	goodsPrices := make(map[string]float64)

// 	// 收集所有需要查询价格的商品ID和选项ID
// 	for _, input := range remainingInputs {
// 		key := fmt.Sprintf("%d:%d", input.GoodsId, input.GoodsOptionsId)
// 		goodsOptionIds[key] = true
// 	}

// 	// 查询所有相关商品的价格
// 	for key := range goodsOptionIds {
// 		parts := strings.Split(key, ":")
// 		if len(parts) != 2 {
// 			continue
// 		}

// 		goodsId := gconv.Int(parts[0])
// 		optionId := gconv.Int(parts[1])

// 		// 优先查询秒杀价格
// 		goods, err := dao.SeckillGoods.Ctx(ctx).
// 			Where("goods_id", goodsId).
// 			Where("goods_options_id", optionId).
// 			One()

// 		if err == nil && goods != nil {
// 			// 使用秒杀价格
// 			goodsPrices[key] = gconv.Float64(goods["seckill_price"])
// 		} else {
// 			// 回退到普通商品价格
// 			option, err := dao.GoodsOptionsInfo.Ctx(ctx).Where("id", optionId).One()
// 			if err == nil && option != nil {
// 				goodsPrices[key] = gconv.Float64(option["price"])
// 			} else {
// 				// 默认价格
// 				goodsPrices[key] = 0
// 			}
// 		}
// 	}

// 	// 准备批量创建订单的数据
// 	orderInputs := make([]*model.SeckillOrderAddInput, len(remainingInputs))
// 	for i, input := range remainingInputs {
// 		// 设置订单号
// 		orderNo := s.generateOrderNo(input.UserId)

// 		// 查询价格
// 		priceKey := fmt.Sprintf("%d:%d", input.GoodsId, input.GoodsOptionsId)
// 		price := goodsPrices[priceKey]

// 		// 当找不到价格时设置默认值
// 		if price <= 0 {
// 			price = 999 // 默认价格，单位元
// 		}

// 		// 转换为分
// 		priceInCents := int(price * 100)

// 		orderInputs[i] = &model.SeckillOrderAddInput{
// 			UserId:         int64(input.UserId),
// 			GoodsId:        int64(input.GoodsId),
// 			GoodsOptionsId: int64(input.GoodsOptionsId),
// 			SeckillPrice:   priceInCents,
// 			OriginalPrice:  priceInCents,
// 			Status:         1, // 待支付
// 			OrderNo:        orderNo,
// 			Count:          int(input.Count),
// 			ConsigneeName:  fmt.Sprintf("用户%d", input.UserId),
// 			ConsigneePhone: input.UserPhone,
// 			Address:        input.UserAddress,
// 			Remark:         input.Remark,
// 			PayTime:        gtime.Now(),
// 		}

// 		// 在结果中记录订单号
// 		idx := remainingIndexes[i]
// 		results[idx].OrderNo = orderNo
// 	}

// 	// 批量创建订单
// 	orderIds, err := s.BatchAddSeckillOrders(ctx, orderInputs)
// 	if err != nil {
// 		// 订单创建失败，回滚库存
// 		g.Log().Error(ctx, "批量创建订单失败:", err)
// 		// 注意：真实环境应添加补偿逻辑确保回滚成功
// 		for i, item := range stockItems {
// 			s.stockManager.AddStock(ctx, item.GoodsId, item.OptionId, item.Quantity)
// 			idx := remainingIndexes[i]
// 			results[idx].Status = consts.CodeSeckillSystemError
// 			results[idx].Message = "订单创建失败: " + err.Error()
// 			results[idx].ProcessTime = time.Since(startTime).Milliseconds()

// 			// 记录秒杀尝试
// 			_ = s.recordSeckillAttempt(ctx, inputs[idx], results[idx])
// 		}
// 		return results, nil
// 	}

// 	g.Log().Infof(ctx, "批量创建订单成功，订单IDs: %v", orderIds)

// 	// 更新结果状态为成功
// 	for j, idx := range remainingIndexes {
// 		results[idx].Status = consts.CodeSeckillSuccess
// 		results[idx].Message = "秒杀成功，秒杀订单已创建"
// 		results[idx].ProcessTime = time.Since(startTime).Milliseconds()

// 		// 缓存结果用于幂等性检查
// 		cacheKey := fmt.Sprintf("%s%s", consts.SeckillResultPrefix, results[idx].RequestId)
// 		_ = seckill.SetCache(ctx, cacheKey, results[idx], 1800*time.Second)

// 		// 记录秒杀尝试
// 		_ = s.recordSeckillAttempt(ctx, inputs[idx], results[idx])

// 		// 记录订单创建成功（如果有订单ID列表）
// 		if len(orderIds) > j {
// 			g.Log().Debugf(ctx, "订单[%s]创建成功，ID: %d", results[idx].OrderNo, orderIds[j])
// 		}
// 	}

// 	// 异步同步库存到数据库
// 	go func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				g.Log().Error(context.Background(), "同步库存时发生panic:", r)
// 			}
// 		}()

// 		// 使用新的上下文，因为原上下文可能已关闭
// 		syncCtx := context.Background()

// 		// 延迟几秒再执行同步，确保其他并发请求先完成
// 		time.Sleep(2 * time.Second)

// 		// 准备同步项
// 		syncItems := make([]struct {
// 			GoodsId  int32
// 			OptionId int32
// 		}, len(stockItems))

// 		for i, item := range stockItems {
// 			syncItems[i] = struct {
// 				GoodsId  int32
// 				OptionId int32
// 			}{
// 				GoodsId:  item.GoodsId,
// 				OptionId: item.OptionId,
// 			}
// 		}

// 		err := s.stockManager.BatchSyncStockToDatabase(syncCtx, syncItems)
// 		if err != nil {
// 			g.Log().Error(syncCtx, "批量同步库存失败:", err)
// 		} else {
// 			g.Log().Info(syncCtx, "批量同步库存成功")
// 		}
// 	}()

// 	return results, nil
// }

// // 添加定期同步库存的方法
// func (s *sSeckill) startStockSyncWorker() {
// 	g.Log().Info(context.Background(), "启动秒杀库存同步任务...")

// 	go func() {
// 		// 定义指数退避重试策略
// 		var syncInterval time.Duration = 5 * time.Second
// 		var maxSyncInterval time.Duration = 60 * time.Second
// 		var consecutiveErrors int = 0

// 		ticker := time.NewTicker(syncInterval)
// 		defer ticker.Stop()

// 		// 立即执行一次同步，不必等待第一个ticker触发
// 		ctx := context.Background()
// 		if err := s.syncAllSeckillGoods(ctx); err != nil {
// 			consecutiveErrors++
// 			// 首次同步失败，增加日志但不调整间隔
// 			g.Log().Warningf(ctx, "首次同步失败: %v，将在%s后重试", err, syncInterval)
// 		} else {
// 			consecutiveErrors = 0
// 		}

// 		for {
// 			select {
// 			case <-s.stopCh:
// 				g.Log().Info(context.Background(), "停止秒杀库存同步任务")
// 				return
// 			case <-ticker.C:
// 				ctx := context.Background()
// 				if err := s.syncAllSeckillGoods(ctx); err != nil {
// 					// 同步失败，增加连续错误计数
// 					consecutiveErrors++

// 					// 如果连续错误超过阈值，指数级增加同步间隔
// 					if consecutiveErrors > 3 {
// 						// 计算新的间隔时间(指数退避策略)
// 						syncInterval = time.Duration(float64(syncInterval) * 1.5)
// 						if syncInterval > maxSyncInterval {
// 							syncInterval = maxSyncInterval
// 						}

// 						ticker.Reset(syncInterval)
// 						g.Log().Warningf(ctx, "连续%d次同步失败，调整同步间隔为%s",
// 							consecutiveErrors, syncInterval)
// 					}
// 				} else {
// 					// 同步成功，检查是否需要恢复正常间隔
// 					consecutiveErrors = 0

// 					// 如果当前间隔不是默认间隔，则恢复
// 					if syncInterval != 5*time.Second {
// 						syncInterval = 5 * time.Second
// 						ticker.Reset(syncInterval)
// 						g.Log().Info(ctx, "同步恢复正常，重置同步间隔为5秒")
// 					}
// 				}
// 			}
// 		}
// 	}()
// }

// // syncAllSeckillGoods 同步所有秒杀商品库存
// func (s *sSeckill) syncAllSeckillGoods(ctx context.Context) error {
// 	g.Log().Info(ctx, "执行定期秒杀库存同步...")

// 	// 创建错误返回
// 	var syncError error = nil

// 	// 检查数据库连接
// 	if _, err := g.DB().Ctx(ctx).GetOne(ctx, "SELECT 1"); err != nil {
// 		g.Log().Errorf(ctx, "数据库连接测试失败: %v", err)
// 		return fmt.Errorf("数据库连接测试失败: %v", err)
// 	}

// 	// 获取当前活跃的数据库配置
// 	dbConfig := g.DB().GetConfig()
// 	if dbConfig != nil {
// 		// 记录当前配置用于调试
// 		g.Log().Infof(ctx, "当前数据库配置: %v", dbConfig)
// 	}

// 	// 获取所有活跃的秒杀商品
// 	var seckillGoods []entity.SeckillGoods
// 	err := dao.SeckillGoods.Ctx(ctx).
// 		WhereIn("status", g.Slice{0, 1}). // 未开始或进行中
// 		Scan(&seckillGoods)

// 	if err != nil {
// 		g.Log().Error(ctx, "查询秒杀商品失败:", err)
// 		return err
// 	}

// 	if len(seckillGoods) == 0 {
// 		g.Log().Info(ctx, "没有需要同步的秒杀商品")
// 		return nil
// 	}

// 	g.Log().Infof(ctx, "开始同步%d个秒杀商品的库存", len(seckillGoods))

// 	// 准备批量同步的数据
// 	syncItems := make([]struct {
// 		GoodsId  int32
// 		OptionId int32
// 	}, 0, len(seckillGoods))

// 	for _, goods := range seckillGoods {
// 		syncItems = append(syncItems, struct {
// 			GoodsId  int32
// 			OptionId int32
// 		}{
// 			GoodsId:  int32(goods.GoodsId),
// 			OptionId: int32(goods.GoodsOptionsId),
// 		})
// 	}

// 	// 使用批量同步方法
// 	if err := s.stockManager.BatchSyncStockToDatabase(ctx, syncItems); err != nil {
// 		g.Log().Error(ctx, "批量同步商品库存失败:", err)
// 		syncError = err

// 		// 如果批量同步失败，尝试对每个商品单独同步
// 		g.Log().Info(ctx, "开始尝试单个同步每个商品...")
// 		successCount := 0

// 		for _, item := range syncItems {
// 			if err := s.stockManager.SyncStockToDatabase(ctx, item.GoodsId, item.OptionId); err != nil {
// 				g.Log().Warningf(ctx, "单个同步商品[%d:%d]库存失败: %v",
// 					item.GoodsId, item.OptionId, err)
// 			} else {
// 				successCount++
// 			}
// 		}

// 		if successCount > 0 {
// 			// 部分成功，降级错误等级
// 			g.Log().Infof(ctx, "单个同步完成，成功: %d/%d", successCount, len(syncItems))
// 			syncError = fmt.Errorf("部分同步成功 (%d/%d)", successCount, len(syncItems))
// 		}
// 	} else {
// 		g.Log().Info(ctx, "批量同步商品库存成功")
// 	}

// 	return syncError
// }

// // ForceSyncStock 强制立即同步指定商品的库存
// func (s *sSeckill) ForceSyncStock(ctx context.Context, goodsId, optionsId int64) error {
// 	g.Log().Infof(ctx, "强制同步商品[%d:%d]库存", goodsId, optionsId)

// 	// 检查商品是否存在
// 	exists, err := dao.SeckillGoods.Ctx(ctx).
// 		Where("goods_id", goodsId).
// 		Where("goods_options_id", optionsId).
// 		Count()

// 	if err != nil {
// 		return err
// 	}

// 	if exists == 0 {
// 		return gerror.Newf("秒杀商品[%d:%d]不存在", goodsId, optionsId)
// 	}

// 	// 立即同步库存
// 	return s.stockManager.SyncStockToDatabase(ctx, int32(goodsId), int32(optionsId))
// }

// // ForceSyncAllStock 强制立即同步所有秒杀商品库存
// //func (s *sSeckill) ForceSyncAllStock(ctx context.Context) error {
// //	g.Log().Info(ctx, "强制同步所有秒杀商品库存")
// //	s.syncAllSeckillGoods(ctx)
// //	return nil
// //}

// // ensureCorrectDatabaseConnection 确保数据库连接正确
// func (s *sSeckill) ensureCorrectDatabaseConnection(ctx context.Context) error {
// 	// 首先检查配置
// 	initDatabaseConfig(ctx)

// 	// 获取当前数据库配置
// 	dbConfig := g.Cfg().MustGet(ctx, "database.default.link").String()
// 	if dbConfig == "" {
// 		g.Log().Error(ctx, "数据库连接字符串为空")
// 		g.Log().Info(ctx, "=== 数据库配置诊断 ===")
// 		g.Log().Info(ctx, "请确保在manifest/config/config.yaml中配置了有效的数据库连接字符串")
// 		g.Log().Info(ctx, "正确的配置格式示例:")
// 		g.Log().Info(ctx, "database:\n  default:\n    link: \"mysql:root:111111@tcp(127.0.0.1:3306)/gf_shop\"")

// 		return fmt.Errorf("请在配置文件中设置数据库连接")
// 	}

// 	// 解析连接字符串，判断格式是否正确
// 	if !strings.Contains(dbConfig, "@tcp(") || !strings.Contains(dbConfig, ")/") {
// 		g.Log().Error(ctx, "数据库连接字符串格式不正确:", dbConfig)
// 		g.Log().Info(ctx, "=== 数据库配置格式诊断 ===")
// 		g.Log().Info(ctx, "正确的格式应为: mysql:用户名:密码@tcp(主机:端口)/数据库名")
// 		g.Log().Info(ctx, "示例: mysql:root:111111@tcp(127.0.0.1:3306)/gf_shop")

// 		return fmt.Errorf("数据库连接字符串格式不正确")
// 	}

// 	// 隐藏密码后输出配置
// 	maskedConfig := dbConfig
// 	if idx := strings.Index(maskedConfig, ":"); idx > 0 {
// 		if idx2 := strings.Index(maskedConfig[idx+1:], "@"); idx2 > 0 {
// 			// 替换密码部分为******
// 			maskedConfig = maskedConfig[:idx+1] + "******" + maskedConfig[idx+1+idx2:]
// 		}
// 	}
// 	g.Log().Debug(ctx, "当前数据库连接配置:", maskedConfig)

// 	// 尝试连接数据库
// 	try := 0
// 	maxTries := 3
// 	for try < maxTries {
// 		// 使用配置的连接信息进行测试
// 		if _, err := g.DB().GetOne(ctx, "SELECT 1"); err != nil {
// 			try++
// 			errMsg := err.Error()

// 			// 更详细的错误分类和诊断
// 			if strings.Contains(errMsg, "Access denied") {
// 				// 用户名密码错误
// 				g.Log().Errorf(ctx, "数据库连接失败，用户名或密码错误: %v", err)
// 				g.Log().Info(ctx, "=== 数据库认证诊断 ===")
// 				g.Log().Info(ctx, "1. 请检查配置文件中的用户名和密码是否正确")
// 				g.Log().Info(ctx, "2. 请确保该用户有权限访问指定的数据库")
// 				g.Log().Info(ctx, "3. 检查MySQL的用户认证设置")

// 				// 建议修改配置文件
// 				correctLink := "mysql:root:111111@tcp(127.0.0.1:3306)/shop"
// 				g.Log().Warning(ctx, "请在manifest/config/config.yaml中修改数据库配置为:")
// 				g.Log().Warning(ctx, "database:\n  default:\n    link: \""+correctLink+"\"")

// 				return fmt.Errorf("数据库用户名或密码错误: %v", err)
// 			} else if strings.Contains(errMsg, "connection refused") ||
// 				strings.Contains(errMsg, "dial tcp") ||
// 				strings.Contains(errMsg, "i/o timeout") {
// 				// 连接问题 - 服务未启动或网络问题
// 				g.Log().Warningf(ctx, "数据库连接测试失败(尝试%d/%d): %v", try, maxTries, err)
// 				g.Log().Info(ctx, "=== 数据库连接诊断 ===")
// 				g.Log().Info(ctx, "1. 请确保MySQL服务已启动")
// 				g.Log().Info(ctx, "2. 检查主机名和端口是否正确")
// 				g.Log().Info(ctx, "3. 检查防火墙设置")
// 				g.Log().Info(ctx, "4. 如果使用远程数据库，请确保网络连接正常")

// 				if try >= maxTries {
// 					return fmt.Errorf("数据库连接测试最终失败: %v", err)
// 				}
// 				// 短暂等待后重试
// 				time.Sleep(500 * time.Millisecond)
// 			} else if strings.Contains(errMsg, "unknown database") {
// 				// 数据库不存在
// 				g.Log().Errorf(ctx, "数据库不存在: %v", err)
// 				g.Log().Info(ctx, "=== 数据库诊断 ===")
// 				g.Log().Info(ctx, "请确保指定的数据库已创建")
// 				g.Log().Info(ctx, "可以使用以下命令创建数据库:")
// 				g.Log().Info(ctx, "CREATE DATABASE gf_shop;")

// 				return fmt.Errorf("指定的数据库不存在: %v", err)
// 			} else {
// 				// 其他错误
// 				g.Log().Errorf(ctx, "数据库连接错误: %v", err)
// 				g.Log().Info(ctx, "=== 通用数据库错误诊断 ===")
// 				g.Log().Info(ctx, "1. 检查MySQL版本兼容性")
// 				g.Log().Info(ctx, "2. 检查MySQL配置")
// 				g.Log().Info(ctx, "3. 检查系统资源使用情况")

// 				return err
// 			}
// 		} else {
// 			// 连接成功
// 			g.Log().Info(ctx, "数据库连接测试成功")
// 			// 输出一些数据库信息，帮助诊断
// 			version, err := g.DB().GetValue(ctx, "SELECT VERSION()")
// 			if err == nil {
// 				g.Log().Info(ctx, "数据库版本:", version)
// 			}
// 			return nil
// 		}
// 	}

// 	return fmt.Errorf("数据库连接失败，尝试次数已用完")
// }

// // recordSeckillAttempt 记录秒杀尝试
// func (s *sSeckill) recordSeckillAttempt(ctx context.Context, input *model.SeckillDoInput, result *model.SeckillDoOutput) error {
// 	// 使用新的SeckillRecord表记录秒杀尝试
// 	now := gtime.Now()

// 	// 记录到数据库
// 	g.Log().Infof(ctx, "记录秒杀尝试: 用户ID=%d, 商品ID=%d, 状态=%d, 订单号=%s",
// 		input.UserId, input.GoodsId, result.Status, result.OrderNo)

// 	// 检查秒杀记录表中是否已存在该请求ID的记录
// 	recordExists := 0
// 	var err error

// 	// 使用事务防止并发问题
// 	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
// 		// 检查表是否存在
// 		tableExists, err := tx.GetOne("SHOW TABLES LIKE 'seckill_record'")
// 		if err != nil {
// 			g.Log().Error(ctx, "检查秒杀记录表是否存在失败:", err)
// 			return err
// 		}

// 		// 如果表不存在，创建表
// 		if tableExists == nil {
// 			g.Log().Warning(ctx, "秒杀记录表不存在，开始创建表...")
// 			createTableSQL := `
// 			CREATE TABLE IF NOT EXISTS seckill_record (
// 				id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
// 				request_id     VARCHAR(128) NOT NULL COMMENT '请求ID',
// 				user_id        BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
// 				goods_id       BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
// 				goods_options_id BIGINT UNSIGNED NOT NULL COMMENT '商品规格ID',
// 				count          INT NOT NULL DEFAULT 1 COMMENT '商品数量',
// 				status         INT NOT NULL COMMENT '秒杀状态：0-成功 其他-失败',
// 				order_no       VARCHAR(128) DEFAULT '' COMMENT '订单编号',
// 				message        VARCHAR(256) DEFAULT '' COMMENT '消息',
// 				process_time   BIGINT DEFAULT 0 COMMENT '处理时间(毫秒)',
// 				created_at     DATETIME NOT NULL COMMENT '创建时间',
// 				updated_at     DATETIME NOT NULL COMMENT '更新时间',
// 				PRIMARY KEY (id),
// 				UNIQUE KEY uk_request_id (request_id),
// 				KEY idx_user_id (user_id),
// 				KEY idx_goods_id (goods_id),
// 				KEY idx_order_no (order_no)
// 			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀记录表';
// 			`
// 			_, err := tx.Exec(createTableSQL)
// 			if err != nil {
// 				g.Log().Error(ctx, "创建秒杀记录表失败:", err)
// 				return err
// 			}
// 			g.Log().Info(ctx, "秒杀记录表创建成功")
// 		} else {
// 			// 检查是否已有记录
// 			result, err := tx.GetOne("SELECT COUNT(1) as count FROM seckill_record WHERE request_id=?", input.RequestId)
// 			if err != nil {
// 				g.Log().Error(ctx, "查询秒杀记录失败:", err)
// 				return err
// 			}

// 			if result != nil {
// 				recordExists = gconv.Int(result["count"])
// 			}
// 		}

// 		// 如果已有记录，不重复创建
// 		if recordExists > 0 {
// 			g.Log().Info(ctx, "已存在相同请求ID的记录，跳过:", input.RequestId)
// 			return nil
// 		}

// 		// 准备数据
// 		data := g.Map{
// 			"request_id":       input.RequestId,
// 			"user_id":          input.UserId,
// 			"goods_id":         input.GoodsId,
// 			"goods_options_id": input.GoodsOptionsId,
// 			"count":            input.Count,
// 			"status":           result.Status,
// 			"order_no":         result.OrderNo,
// 			"message":          result.Message,
// 			"process_time":     result.ProcessTime,
// 			"created_at":       now,
// 			"updated_at":       now,
// 		}

// 		// 插入记录到秒杀记录表
// 		_, err = tx.Model("seckill_record").Data(data).Insert()
// 		if err != nil {
// 			g.Log().Error(ctx, "插入秒杀记录失败:", err)
// 			return err
// 		}

// 		g.Log().Info(ctx, "秒杀记录保存成功, 请求ID:", input.RequestId)
// 		return nil
// 	})

// 	if err != nil {
// 		// 如果事务失败，使用SeckillOrder表或日志记录
// 		g.Log().Warning(ctx, "秒杀记录表存储失败，降级使用日志记录:", err)

// 		// 准备日志数据
// 		data := g.Map{
// 			"request_id":       input.RequestId,
// 			"user_id":          input.UserId,
// 			"goods_id":         input.GoodsId,
// 			"goods_options_id": input.GoodsOptionsId,
// 			"count":            input.Count,
// 			"status":           result.Status,
// 			"order_no":         result.OrderNo,
// 			"message":          result.Message,
// 			"process_time":     result.ProcessTime,
// 			"created_at":       now,
// 			"updated_at":       now,
// 		}

// 		// 将数据格式化为JSON并记录到日志
// 		jsonBytes, _ := json.Marshal(data)
// 		g.Log().Info(ctx, "秒杀记录降级存储:", string(jsonBytes))
// 	}

// 	return err
// }
