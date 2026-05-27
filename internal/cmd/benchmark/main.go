package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"bit303_shop/internal/consts"
// 	"bit303_shop/internal/logic/order"
// 	logicseckill "bit303_shop/internal/logic/seckill"
// 	"bit303_shop/internal/service"
// 	"bit303_shop/internal/test"
// 	"bit303_shop/utility/seckill"
// 	"math/rand"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"sync/atomic"
// 	"syscall"
// 	"time"

// 	// 导入MySQL驱动
// 	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
// 	// 导入Redis驱动
// 	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
// 	// 导入Kafka相关包
// 	"github.com/IBM/sarama"

// 	"github.com/gogf/gf/v2/database/gdb"
// 	"github.com/gogf/gf/v2/frame/g"
// 	"github.com/gogf/gf/v2/os/gcfg"
// 	"github.com/gogf/gf/v2/util/gconv"
// )

// func main() {
// 	// 加载测试配置文件
// 	err := g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetPath("internal/cmd/benchmark")
// 	if err != nil {
// 		fmt.Printf("加载配置文件路径失败: %v\n", err)
// 	}

// 	// 命令行参数
// 	testType := flag.String("type", "order", "测试类型: order、order-api、seckill、seckill-api、direct-seckill或all")
// 	users := flag.Int("c", 100, "并发用户数")
// 	duration := flag.Int("d", 60, "测试持续时间(秒)")
// 	rampUp := flag.Int("r", 10, "预热时间(秒)")
// 	rampUsers := flag.Int("ru", 10, "预热阶段用户数")
// 	goodsId := flag.Int64("gid", 1, "商品ID")
// 	optionsId := flag.Int64("oid", 1, "商品规格ID")
// 	thinkTime := flag.Int("t", 100, "思考时间(毫秒)")
// 	apiURL := flag.String("url", "http://localhost:8000", "API基础URL")
// 	dbUser := flag.String("dbuser", "root", "数据库用户名")
// 	dbPass := flag.String("dbpass", "123456", "数据库密码")
// 	dbHost := flag.String("dbhost", "127.0.0.1:3306", "数据库主机")
// 	dbName := flag.String("dbname", "gf_shop", "数据库名称")

// 	flag.Parse()

// 	// 根据命令行参数设置数据库配置
// 	gdb.SetConfig(gdb.Config{
// 		"default": gdb.ConfigGroup{
// 			gdb.ConfigNode{
// 				Link:  fmt.Sprintf("mysql:%s:%s@tcp(%s)/%s", *dbUser, *dbPass, *dbHost, *dbName),
// 				Debug: true,
// 			},
// 		},
// 	})

// 	// 创建上下文并处理中断信号
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// 处理信号
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
// 	go func() {
// 		<-sigChan
// 		fmt.Println("\n接收到中断信号，正在清理资源并退出...")
// 		cancel()
// 		// 给一点时间清理
// 		time.Sleep(2 * time.Second)
// 		os.Exit(1)
// 	}()

// 	// 根据测试类型执行不同的测试
// 	switch *testType {
// 	case "all":
// 		// 设置数据库配置
// 		gdb.SetConfig(gdb.Config{
// 			"default": gdb.ConfigGroup{
// 				gdb.ConfigNode{
// 					Type:  "mysql",
// 					Host:  "127.0.0.1",
// 					Port:  "3306",
// 					User:  *dbUser,
// 					Pass:  *dbPass,
// 					Name:  *dbName,
// 					Debug: true,
// 				},
// 			},
// 		})

// 		// 验证数据库连接
// 		db := g.DB()
// 		_, err := db.Query(ctx, "SELECT 1")
// 		if err != nil {
// 			fmt.Printf("数据库连接测试失败: %v\n", err)
// 			fmt.Println("请检查数据库配置后重试")
// 			os.Exit(1)
// 		}
// 		fmt.Println("数据库连接测试成功")

// 		// 注册服务
// 		service.RegisterOrder(order.New())
// 		service.RegisterSeckill(logicseckill.New())

// 		// 准备测试数据
// 		prepareTestData(ctx)

// 		// 运行订单性能测试
// 		fmt.Println("\n========== 开始订单性能测试 ==========")
// 		orderResult, err := test.RunOrderBenchmark(ctx, *users, time.Duration(*duration)*time.Second)
// 		if err != nil {
// 			fmt.Printf("订单测试执行失败: %v\n", err)
// 		} else {
// 			test.PrintOrderBenchmarkResult(orderResult)
// 		}

// 		// 运行秒杀性能测试
// 		fmt.Println("\n========== 开始秒杀性能测试 ==========")
// 		seckillConfig := &test.BenchmarkConfig{
// 			ConcurrentUsers: *users,
// 			Duration:        time.Duration(*duration) * time.Second,
// 			RampUpTime:      time.Duration(*rampUp) * time.Second,
// 			RampUpUsers:     *rampUsers,
// 			ThinkTime:       time.Duration(*thinkTime) * time.Millisecond,
// 			GoodsId:         *goodsId,
// 			GoodsOptionsId:  *optionsId,
// 		}
// 		seckillResult, err := test.RunSeckillBenchmark(ctx, seckillConfig)
// 		if err != nil {
// 			fmt.Printf("秒杀测试执行失败: %v\n", err)
// 		} else {
// 			test.PrintBenchmarkResult(seckillResult)
// 		}

// 		// 运行直接秒杀测试
// 		fmt.Println("\n========== 开始直接秒杀测试 ==========")
// 		TestDirectSeckill(*goodsId, *optionsId, *users)

// 		fmt.Println("\n========== 所有性能测试完成 ==========")
// 	case "order":
// 		runOrderBenchmark(ctx, *users, time.Duration(*duration)*time.Second)
// 	case "order-api":
// 		runOrderAPIBenchmark(ctx, *apiURL, *users, time.Duration(*duration)*time.Second, time.Duration(*thinkTime)*time.Millisecond)
// 	case "seckill":
// 		// 注册秒杀服务
// 		service.RegisterSeckill(logicseckill.New())

// 		// 在运行测试前确保数据库配置正确
// 		gdb.SetConfig(gdb.Config{
// 			"default": gdb.ConfigGroup{
// 				gdb.ConfigNode{
// 					Type:  "mysql",
// 					Host:  "127.0.0.1",
// 					Port:  "3306",
// 					User:  *dbUser,
// 					Pass:  *dbPass,
// 					Name:  *dbName,
// 					Debug: true,
// 				},
// 			},
// 		})

// 		// 验证数据库连接
// 		db := g.DB()
// 		_, err := db.Query(ctx, "SELECT 1")
// 		if err != nil {
// 			fmt.Printf("数据库连接测试失败: %v\n", err)
// 			fmt.Println("请检查数据库配置后重试")
// 			os.Exit(1)
// 		}
// 		fmt.Println("数据库连接测试成功")

// 		// 创建测试配置
// 		config := &test.BenchmarkConfig{
// 			ConcurrentUsers: *users,
// 			Duration:        time.Duration(*duration) * time.Second,
// 			RampUpTime:      time.Duration(*rampUp) * time.Second,
// 			RampUpUsers:     *rampUsers,
// 			ThinkTime:       time.Duration(*thinkTime) * time.Millisecond,
// 			GoodsId:         *goodsId,
// 			GoodsOptionsId:  *optionsId,
// 		}

// 		// 运行秒杀性能测试
// 		result, err := test.RunSeckillBenchmark(ctx, config)
// 		if err != nil {
// 			fmt.Printf("秒杀测试执行失败: %v\n", err)
// 			os.Exit(1)
// 		}

// 		// 打印测试结果
// 		test.PrintBenchmarkResult(result)
// 	case "seckill-api":
// 		// 确保数据库配置正确
// 		gdb.SetConfig(gdb.Config{
// 			"default": gdb.ConfigGroup{
// 				gdb.ConfigNode{
// 					Type:  "mysql",
// 					Host:  "127.0.0.1",
// 					Port:  "3306",
// 					User:  *dbUser,
// 					Pass:  *dbPass,
// 					Name:  *dbName,
// 					Debug: true,
// 				},
// 			},
// 		})

// 		// 验证数据库连接
// 		db := g.DB()
// 		_, err := db.Query(ctx, "SELECT 1")
// 		if err != nil {
// 			fmt.Printf("数据库连接测试失败: %v\n", err)
// 			fmt.Println("请检查数据库配置后重试")
// 			os.Exit(1)
// 		}
// 		fmt.Println("数据库连接测试成功")

// 		runSeckillAPIBenchmark(ctx, *apiURL, *users, time.Duration(*duration)*time.Second, time.Duration(*thinkTime)*time.Millisecond, *goodsId, *optionsId)
// 	case "direct-seckill":
// 		// 注册秒杀服务
// 		service.RegisterSeckill(logicseckill.New())

// 		// 确保数据库配置正确
// 		gdb.SetConfig(gdb.Config{
// 			"default": gdb.ConfigGroup{
// 				gdb.ConfigNode{
// 					Type:  "mysql",
// 					Host:  "127.0.0.1",
// 					Port:  "3306",
// 					User:  *dbUser,
// 					Pass:  *dbPass,
// 					Name:  *dbName,
// 					Debug: true,
// 				},
// 			},
// 		})

// 		// 验证数据库连接
// 		db := g.DB()
// 		_, err := db.Query(ctx, "SELECT 1")
// 		if err != nil {
// 			fmt.Printf("数据库连接测试失败: %v\n", err)
// 			fmt.Println("请检查数据库配置后重试")
// 			os.Exit(1)
// 		}
// 		fmt.Println("数据库连接测试成功")

// 		// 运行直接秒杀测试
// 		TestDirectSeckill(*goodsId, *optionsId, *users)
// 	default:
// 		fmt.Printf("不支持的测试类型: %s\n", *testType)
// 		os.Exit(1)
// 	}

// 	fmt.Println("性能测试完成")
// }

// // 运行订单性能测试
// func runOrderBenchmark(ctx context.Context, concurrentUsers int, duration time.Duration) {
// 	fmt.Println("开始执行订单性能测试...")
// 	fmt.Printf("并发用户数: %d, 测试时间: %v\n", concurrentUsers, duration)

// 	// 注册订单服务
// 	service.RegisterOrder(order.New())

// 	// 在调用测试前确保数据库配置正确
// 	// 这里强制设置数据库配置，避免使用错误的配置
// 	gdb.SetConfig(gdb.Config{
// 		"default": gdb.ConfigGroup{
// 			gdb.ConfigNode{
// 				Type:  "mysql",
// 				Host:  "127.0.0.1",
// 				Port:  "3306",
// 				User:  "root",
// 				Pass:  "123456",
// 				Name:  "gf_shop",
// 				Debug: true,
// 			},
// 		},
// 	})

// 	// 验证数据库连接
// 	db := g.DB()
// 	_, err := db.Query(ctx, "SELECT 1")
// 	if err != nil {
// 		fmt.Printf("数据库连接测试失败: %v\n", err)
// 		fmt.Println("请检查数据库配置后重试")
// 		os.Exit(1)
// 	}
// 	fmt.Println("数据库连接测试成功")

// 	// 在测试开始前确保数据库中有必要的测试数据
// 	prepareTestData(ctx)

// 	// 运行性能测试
// 	result, err := test.RunOrderBenchmark(ctx, concurrentUsers, duration)
// 	if err != nil {
// 		fmt.Printf("订单测试执行失败: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// 打印测试结果
// 	test.PrintOrderBenchmarkResult(result)
// }

// // 运行订单API性能测试
// func runOrderAPIBenchmark(ctx context.Context, baseURL string, concurrentUsers int, duration, thinkTime time.Duration) {
// 	fmt.Println("开始执行订单API性能测试...")
// 	fmt.Printf("API基础URL: %s\n", baseURL)
// 	fmt.Printf("并发用户数: %d, 测试时间: %v, 思考时间: %v\n", concurrentUsers, duration, thinkTime)

// 	// 确保数据库配置正确
// 	gdb.SetConfig(gdb.Config{
// 		"default": gdb.ConfigGroup{
// 			gdb.ConfigNode{
// 				Type:  "mysql",
// 				Host:  "127.0.0.1",
// 				Port:  "3306",
// 				User:  "root",
// 				Pass:  "123456",
// 				Name:  "gf_shop",
// 				Debug: true,
// 			},
// 		},
// 	})

// 	// 验证数据库连接
// 	db := g.DB()
// 	_, err := db.Query(ctx, "SELECT 1")
// 	if err != nil {
// 		fmt.Printf("数据库连接测试失败: %v\n", err)
// 		fmt.Println("请检查数据库配置后重试")
// 		os.Exit(1)
// 	}
// 	fmt.Println("数据库连接测试成功")

// 	// 在测试开始前确保数据库中有必要的测试数据
// 	prepareTestData(ctx)

// 	// 运行API性能测试
// 	result, err := test.RunOrderAPIBenchmark(ctx, baseURL, concurrentUsers, duration, thinkTime)
// 	if err != nil {
// 		fmt.Printf("订单API测试执行失败: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// 打印测试结果
// 	test.PrintAPIBenchmarkResult(result)
// }

// // 运行秒杀API性能测试
// func runSeckillAPIBenchmark(ctx context.Context, baseURL string, concurrentUsers int, duration, thinkTime time.Duration, goodsId, optionsId int64) {
// 	fmt.Println("开始执行秒杀API性能测试...")
// 	fmt.Printf("API基础URL: %s\n", baseURL)
// 	fmt.Printf("并发用户数: %d, 测试时间: %v, 思考时间: %v\n", concurrentUsers, duration, thinkTime)
// 	fmt.Printf("商品ID: %d, 规格ID: %d\n", goodsId, optionsId)

// 	// 确保数据库配置正确
// 	gdb.SetConfig(gdb.Config{
// 		"default": gdb.ConfigGroup{
// 			gdb.ConfigNode{
// 				Type:  "mysql",
// 				Host:  "127.0.0.1",
// 				Port:  "3306",
// 				User:  "root",
// 				Pass:  "123456",
// 				Name:  "gf_shop",
// 				Debug: true,
// 			},
// 		},
// 	})

// 	// 验证数据库连接
// 	db := g.DB()
// 	_, err := db.Query(ctx, "SELECT 1")
// 	if err != nil {
// 		fmt.Printf("数据库连接测试失败: %v\n", err)
// 		fmt.Println("请检查数据库配置后重试")
// 		os.Exit(1)
// 	}
// 	fmt.Println("数据库连接测试成功")

// 	// 在测试开始前确保数据库中有必要的测试数据
// 	prepareTestData(ctx)

// 	// 创建API测试配置
// 	config := &test.APITestConfig{
// 		BaseURL:         baseURL,
// 		ConcurrentUsers: concurrentUsers,
// 		Duration:        duration,
// 		ThinkTime:       thinkTime,
// 		RequestTimeout:  30 * time.Second,
// 	}

// 	// 运行API性能测试
// 	result, err := test.RunSeckillAPIBenchmark(ctx, config, goodsId, optionsId)
// 	if err != nil {
// 		fmt.Printf("秒杀API测试执行失败: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// 打印测试结果
// 	test.PrintAPIBenchmarkResult(result)
// }

// // 准备测试所需的数据
// func prepareTestData(ctx context.Context) {
// 	fmt.Println("正在准备测试数据...")

// 	// 检查并创建测试用户
// 	createTestUser(ctx)

// 	// 检查并创建测试商品
// 	createTestGoods(ctx)

// 	// 检查并创建测试地址
// 	createTestAddress(ctx)

// 	fmt.Println("测试数据准备完成")
// }

// // 创建测试用户
// func createTestUser(ctx context.Context) {
// 	// 查询是否已存在测试用户
// 	userCount, err := g.DB().Ctx(ctx).Model("user").Where("id", 1).Count()
// 	if err != nil {
// 		fmt.Printf("查询测试用户失败: %v\n", err)
// 		return
// 	}

// 	if userCount == 0 {
// 		// 创建测试用户
// 		_, err = g.DB().Ctx(ctx).Model("user").Data(g.Map{
// 			"id":         1,
// 			"name":       "测试用户",
// 			"password":   "e10adc3949ba59abbe56e057f20f883e", // 123456的MD5
// 			"email":      "test@example.com",
// 			"phone":      "13800138000",
// 			"status":     1,
// 			"avatar":     "https://example.com/avatar.png",
// 			"created_at": time.Now(),
// 			"updated_at": time.Now(),
// 		}).Insert()

// 		if err != nil {
// 			fmt.Printf("创建测试用户失败: %v\n", err)
// 		} else {
// 			fmt.Println("成功创建测试用户")
// 		}
// 	} else {
// 		fmt.Println("测试用户已存在")
// 	}

// 	// 为API测试创建多个测试用户
// 	for i := 1; i <= 50; i++ {
// 		username := fmt.Sprintf("test_user_%d", i)

// 		// 查询是否已存在此测试用户
// 		exists, err := g.DB().Ctx(ctx).Model("user").Where("name", username).Count()
// 		if err != nil {
// 			fmt.Printf("查询用户 %s 失败: %v\n", username, err)
// 			continue
// 		}

// 		if exists == 0 {
// 			// 创建API测试用户
// 			_, err = g.DB().Ctx(ctx).Model("user").Data(g.Map{
// 				"name":       username,
// 				"password":   "e10adc3949ba59abbe56e057f20f883e", // 123456的MD5
// 				"email":      fmt.Sprintf("test%d@example.com", i),
// 				"phone":      fmt.Sprintf("138%08d", i),
// 				"status":     1,
// 				"avatar":     "https://example.com/avatar.png",
// 				"created_at": time.Now(),
// 				"updated_at": time.Now(),
// 			}).Insert()

// 			if err != nil {
// 				fmt.Printf("创建用户 %s 失败: %v\n", username, err)
// 			} else {
// 				fmt.Printf("成功创建用户 %s\n", username)
// 			}
// 		}
// 	}
// }

// // 创建测试商品
// func createTestGoods(ctx context.Context) {
// 	// 查询是否已存在测试商品
// 	goodsCount, err := g.DB().Ctx(ctx).Model("goods_info").Where("id", 1).Count()
// 	if err != nil {
// 		fmt.Printf("查询测试商品失败: %v\n", err)
// 		return
// 	}

// 	if goodsCount == 0 {
// 		// 创建测试商品
// 		_, err = g.DB().Ctx(ctx).Model("goods_info").Data(g.Map{
// 			"id":                 1,
// 			"pic_url":            "https://example.com/goods1.png",
// 			"name":               "测试商品1",
// 			"price":              1000, // 10元，单位分
// 			"level1_category_id": 1,
// 			"level2_category_id": 1,
// 			"level3_category_id": 1,
// 			"brand":              "测试品牌",
// 			"stock":              100,
// 			"sale":               0,
// 			"tags":               "测试,性能测试",
// 			"detail_info":        "这是一个测试商品详情",
// 			"created_at":         time.Now(),
// 			"updated_at":         time.Now(),
// 		}).Insert()

// 		if err != nil {
// 			fmt.Printf("创建测试商品失败: %v\n", err)
// 		} else {
// 			fmt.Println("成功创建测试商品")
// 		}

// 		// 创建测试商品规格
// 		_, err = g.DB().Ctx(ctx).Model("goods_options_info").Data(g.Map{
// 			"id":         1,
// 			"goods_id":   1,
// 			"pic_url":    "https://example.com/goods1_option1.png",
// 			"name":       "默认规格",
// 			"price":      1000, // 10元，单位分
// 			"stock":      100,
// 			"created_at": time.Now(),
// 			"updated_at": time.Now(),
// 		}).Insert()

// 		if err != nil {
// 			fmt.Printf("创建测试商品规格失败: %v\n", err)
// 		} else {
// 			fmt.Println("成功创建测试商品规格")
// 		}

// 		// 创建秒杀商品
// 		_, err = g.DB().Ctx(ctx).Model("seckill_goods").Data(g.Map{
// 			"goods_id":         1,
// 			"goods_options_id": 1,
// 			"original_price":   1000,
// 			"seckill_price":    500, // 5元，单位分
// 			"seckill_stock":    1000,
// 			"status":           1, // 进行中
// 			"start_time":       time.Now(),
// 			"end_time":         time.Now().Add(24 * time.Hour),
// 			"created_at":       time.Now(),
// 			"updated_at":       time.Now(),
// 		}).Insert()

// 		if err != nil {
// 			fmt.Printf("创建测试秒杀商品失败: %v\n", err)
// 		} else {
// 			fmt.Println("成功创建测试秒杀商品")
// 		}
// 	} else {
// 		fmt.Println("测试商品已存在")
// 	}
// }

// // 创建测试地址
// func createTestAddress(ctx context.Context) {
// 	// 查询是否已存在测试地址
// 	addressCount, err := g.DB().Ctx(ctx).Model("address").Where("id", 1).Count()
// 	if err != nil {
// 		fmt.Printf("查询测试地址失败: %v\n", err)
// 		return
// 	}

// 	if addressCount == 0 {
// 		// 创建测试地址
// 		_, err = g.DB().Ctx(ctx).Model("address").Data(g.Map{
// 			"id":         1,
// 			"user_id":    1,
// 			"link_man":   "测试收货人",
// 			"link_phone": "13800138000",
// 			"province":   "北京市",
// 			"city":       "北京市",
// 			"district":   "海淀区",
// 			"address":    "测试详细地址",
// 			"is_default": 1,
// 			"created_at": time.Now(),
// 			"updated_at": time.Now(),
// 		}).Insert()

// 		if err != nil {
// 			fmt.Printf("创建测试地址失败: %v\n", err)
// 		} else {
// 			fmt.Println("成功创建测试地址")
// 		}

// 		// 为API测试用户创建地址
// 		for i := 1; i <= 50; i++ {
// 			// 获取用户ID
// 			var userId int64
// 			err := g.DB().Ctx(ctx).Model("user").
// 				Where("name", fmt.Sprintf("test_user_%d", i)).
// 				Fields("id").
// 				Scan(&userId)

// 			if err != nil || userId == 0 {
// 				continue
// 			}

// 			// 为该用户创建地址
// 			_, err = g.DB().Ctx(ctx).Model("address").Data(g.Map{
// 				"user_id":    userId,
// 				"link_man":   fmt.Sprintf("测试用户%d", i),
// 				"link_phone": fmt.Sprintf("138%08d", i),
// 				"province":   "北京市",
// 				"city":       "北京市",
// 				"district":   "海淀区",
// 				"address":    fmt.Sprintf("测试地址%d号", i),
// 				"is_default": 1,
// 				"created_at": time.Now(),
// 				"updated_at": time.Now(),
// 			}).Insert()

// 			if err != nil {
// 				fmt.Printf("为用户%d创建地址失败: %v\n", i, err)
// 			}
// 		}
// 	} else {
// 		fmt.Println("测试地址已存在")
// 	}
// }

// // TestDirectSeckill 执行直接秒杀测试
// func TestDirectSeckill(goodsId, goodsOptionsId int64, concurrentUsers int) {
// 	ctx := context.Background()
// 	fmt.Println("=== 开始直接秒杀功能测试 ===")
// 	fmt.Printf("商品ID: %d, 规格ID: %d, 并发用户数: %d\n", goodsId, goodsOptionsId, concurrentUsers)
// 	fmt.Println("正在准备测试数据...")

// 	// 初始化测试数据
// 	initTestData(ctx, goodsId, goodsOptionsId)

// 	// 测试直接秒杀，使用纯Redis操作
// 	testDirectSeckillByRedis(ctx, goodsId, goodsOptionsId, concurrentUsers)

// 	// 确保在测试结束时同步库存到数据库
// 	fmt.Println("\n执行最终库存同步...")
// 	stockManager := seckill.NewStockManager(1000)
// 	err := stockManager.SyncStockToDatabase(ctx, int32(goodsId), int32(goodsOptionsId))
// 	if err != nil {
// 		fmt.Printf("最终同步库存失败: %v\n", err)
// 	} else {
// 		fmt.Println("成功将最终库存同步到数据库")
// 	}
// }

// // 初始化测试数据
// func initTestData(ctx context.Context, goodsId, goodsOptionsId int64) {
// 	fmt.Println("正在准备测试数据...")

// 	// 检查秒杀商品是否存在
// 	var seckillCount int
// 	seckillCount, err := g.DB().Ctx(ctx).Model("seckill_goods").
// 		Where("goods_id", goodsId).
// 		Where("goods_options_id", goodsOptionsId).
// 		Count()
// 	if err != nil {
// 		fmt.Printf("检查秒杀商品失败: %v\n", err)
// 		return
// 	}

// 	if seckillCount == 0 {
// 		fmt.Println("秒杀商品不存在，请先确保秒杀商品已创建")
// 		return
// 	}

// 	// 更新秒杀商品的库存为一个足够大的值
// 	initialStock := 1000 // 设置一个足够大的初始库存
// 	_, err = g.DB().Ctx(ctx).Model("seckill_goods").
// 		Data(g.Map{"seckill_stock": initialStock}).
// 		Where("goods_id", goodsId).
// 		Where("goods_options_id", goodsOptionsId).
// 		Update()

// 	if err != nil {
// 		fmt.Printf("更新秒杀商品库存失败: %v\n", err)
// 	} else {
// 		fmt.Printf("成功将秒杀商品库存更新为: %d\n", initialStock)
// 	}

// 	// 同时更新商品规格表的库存
// 	if goodsOptionsId > 0 {
// 		_, err = g.DB().Ctx(ctx).Model("goods_options_info").
// 			Data(g.Map{"stock": initialStock}).
// 			Where("id", goodsOptionsId).
// 			Update()

// 		if err != nil {
// 			fmt.Printf("更新商品规格库存失败: %v\n", err)
// 		} else {
// 			fmt.Printf("成功将商品规格库存更新为: %d\n", initialStock)
// 		}
// 	}

// 	// 更新商品主表的库存
// 	_, err = g.DB().Ctx(ctx).Model("goods_info").
// 		Data(g.Map{"stock": initialStock}).
// 		Where("id", goodsId).
// 		Update()

// 	if err != nil {
// 		fmt.Printf("更新商品主表库存失败: %v\n", err)
// 	} else {
// 		fmt.Printf("成功将商品主表库存更新为: %d\n", initialStock)
// 	}

// 	// 清理Redis缓存
// 	patterns := []string{
// 		fmt.Sprintf("%s%d:%d*", consts.SeckillStockPrefix, goodsId, goodsOptionsId),
// 		fmt.Sprintf("%s%d:%d*", consts.SeckillResultPrefix, goodsId, goodsOptionsId),
// 	}

// 	for _, pattern := range patterns {
// 		keys, err := g.Redis().Do(ctx, "KEYS", pattern)
// 		if err != nil {
// 			fmt.Printf("获取Redis键失败: %v\n", err)
// 			continue
// 		}

// 		if len(keys.Slice()) > 0 {
// 			_, err = g.Redis().Do(ctx, "DEL", keys.Slice()...)
// 			if err != nil {
// 				fmt.Printf("删除Redis键失败: %v\n", err)
// 			}
// 		}
// 	}

// 	// 初始化Redis库存
// 	stockKey := fmt.Sprintf("%s%d:%d", consts.SeckillStockPrefix, goodsId, goodsOptionsId)
// 	_, err = g.Redis().Do(ctx, "SET", stockKey, initialStock)
// 	if err != nil {
// 		fmt.Printf("初始化Redis库存失败: %v\n", err)
// 	} else {
// 		fmt.Printf("成功初始化Redis库存为: %d\n", initialStock)
// 	}

// 	// 预热缓存
// 	err = service.Seckill().WarmUpCache(ctx, goodsId, goodsOptionsId)
// 	if err != nil {
// 		fmt.Printf("预热缓存失败: %v\n", err)
// 	} else {
// 		fmt.Println("成功预热缓存")
// 	}
// }

// // 测试直接秒杀，使用纯Redis操作
// func testDirectSeckillByRedis(ctx context.Context, goodsId, goodsOptionsId int64, concurrentUsers int) {
// 	fmt.Println("\n=== 开始直接秒杀功能测试 ===")
// 	fmt.Printf("商品ID: %d, 规格ID: %d, 并发用户数: %d\n", goodsId, goodsOptionsId, concurrentUsers)
// 	fmt.Println("正在准备测试数据...")

// 	// 初始化测试数据
// 	initTestData(ctx, goodsId, goodsOptionsId)

// 	// 创建库存管理器
// 	stockManager := seckill.NewStockManager(1000)

// 	// 统计数据
// 	var successCount, failCount int32
// 	startTime := time.Now()

// 	// 获取初始库存
// 	stockKey := fmt.Sprintf("%s%d:%d", consts.SeckillStockPrefix, goodsId, goodsOptionsId)
// 	stockBefore, _ := g.Redis().Do(ctx, "GET", stockKey)
// 	initialStock := stockBefore.Int()

// 	fmt.Printf("秒杀前库存: %d\n", initialStock)

// 	// 预先查询商品信息和商品规格信息，避免重复查询
// 	_, err := g.DB().Ctx(ctx).Model("goods_info").Where("id", goodsId).One()
// 	if err != nil {
// 		fmt.Printf("获取商品信息失败: %v\n", err)
// 		return
// 	}

// 	goodsOptions, err := g.DB().Ctx(ctx).Model("goods_options_info").Where("id", goodsOptionsId).One()
// 	if err != nil {
// 		fmt.Printf("获取商品规格信息失败: %v\n", err)
// 		return
// 	}

// 	// 计算价格
// 	price := gconv.Float64(goodsOptions["price"])

// 	// 创建等待组来同步goroutine
// 	var wg sync.WaitGroup

// 	// 创建工作池
// 	workerCount := 20
// 	if concurrentUsers < workerCount {
// 		workerCount = concurrentUsers
// 	}

// 	// 创建错误通道，用于接收错误
// 	errCh := make(chan error, concurrentUsers)
// 	defer close(errCh)

// 	// 并发处理秒杀请求
// 	for w := 0; w < workerCount; w++ {
// 		wg.Add(1)
// 		go func(workerId int) {
// 			defer wg.Done()

// 			// 使用单独的上下文，避免互相影响
// 			workerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 			defer cancel()

// 			// 计算每个worker处理的用户数量范围
// 			userPerWorker := concurrentUsers / workerCount
// 			start := workerId*userPerWorker + 1
// 			end := (workerId + 1) * userPerWorker

// 			// 最后一个worker处理剩余的用户
// 			if workerId == workerCount-1 {
// 				end = concurrentUsers
// 			}

// 			// 处理分配给该worker的用户
// 			for i := start; i <= end; i++ {
// 				userId := uint(i)

// 				// 使用StockManager扣减库存，确保库存操作一致性
// 				remain, err := stockManager.DeductStock(workerCtx, int64(userId), int32(goodsId), int32(goodsOptionsId), 1)
// 				if err != nil {
// 					fmt.Printf("用户%d秒杀失败，库存不足\n", userId)
// 					atomic.AddInt32(&failCount, 1)
// 					continue
// 				}

// 				// 库存扣减成功，创建订单
// 				// 生成订单号
// 				orderNo := generateOrderNo(userId)

// 				// 插入秒杀订单记录
// 				seckillOrderData := g.Map{
// 					"user_id":          userId,
// 					"order_no":         orderNo,
// 					"original_price":   int(price * 100), // 单位为分
// 					"seckill_price":    int(price * 80),  // 秒杀价格为原价的8折，单位为分
// 					"goods_id":         goodsId,
// 					"goods_options_id": goodsOptionsId,
// 					"status":           1, // 已支付状态
// 					"pay_time":         time.Now(),
// 					"created_at":       time.Now(),
// 					"updated_at":       time.Now(),
// 				}

// 				r, err := g.DB().Ctx(workerCtx).Model("seckill_order").Insert(seckillOrderData)
// 				if err != nil {
// 					fmt.Printf("用户%d创建秒杀订单失败: %v\n", userId, err)
// 					// 回滚库存
// 					_, _ = stockManager.AddStock(workerCtx, int32(goodsId), int32(goodsOptionsId), 1)
// 					atomic.AddInt32(&failCount, 1)
// 					continue
// 				}

// 				orderId, err := r.LastInsertId()
// 				if err != nil {
// 					fmt.Printf("用户%d获取秒杀订单ID失败: %v\n", userId, err)
// 					// 订单已创建，不回滚
// 					atomic.AddInt32(&successCount, 1)
// 					fmt.Printf("用户%d秒杀成功，订单号: %s\n", userId, orderNo)
// 					continue
// 				}

// 				atomic.AddInt32(&successCount, 1)
// 				fmt.Printf("用户%d秒杀成功，订单号: %s, 订单ID: %d, 剩余库存: %d\n", userId, orderNo, orderId, remain)
// 			}
// 		}(w)
// 	}

// 	// 等待所有goroutine完成
// 	wg.Wait()

// 	// 获取最终库存
// 	stockAfter, _ := g.Redis().Do(ctx, "GET", stockKey)
// 	finalStock := stockAfter.Int()

// 	// 计算统计结果
// 	duration := time.Since(startTime)
// 	successRate := float64(successCount) / float64(concurrentUsers) * 100

// 	fmt.Printf("\n=== 直接秒杀测试结果 ===\n")
// 	fmt.Printf("测试用户数: %d\n", concurrentUsers)
// 	fmt.Printf("成功用户数: %d\n", successCount)
// 	fmt.Printf("失败用户数: %d\n", failCount)
// 	fmt.Printf("成功率: %.2f%%\n", successRate)
// 	fmt.Printf("测试耗时: %v\n", duration)
// 	fmt.Printf("每秒处理请求: %.2f\n", float64(concurrentUsers)/duration.Seconds())
// 	fmt.Printf("初始库存: %d\n", initialStock)
// 	fmt.Printf("剩余库存: %d\n", finalStock)
// 	fmt.Printf("消耗库存: %d\n", initialStock-finalStock)

// 	// 同步Redis库存到seckill_goods表
// 	fmt.Println("\n开始同步Redis库存到数据库...")
// 	err = stockManager.SyncStockToDatabase(ctx, int32(goodsId), int32(goodsOptionsId))
// 	if err != nil {
// 		fmt.Printf("同步库存到数据库失败: %v\n", err)
// 	} else {
// 		fmt.Println("成功同步库存到数据库")
// 	}
// }

// // generateOrderNo 生成订单号
// func generateOrderNo(userId uint) string {
// 	// 订单号格式: 时间戳 + 用户ID后4位 + 4位随机数
// 	timestamp := time.Now().Format("20060102150405")
// 	userIdStr := fmt.Sprintf("%04d", userId%10000)

// 	// 生成4位随机数
// 	rand.Seed(time.Now().UnixNano())
// 	random := fmt.Sprintf("%04d", rand.Intn(10000))

// 	return fmt.Sprintf("%s%s%s", timestamp, userIdStr, random)
// }

// // 初始化Kafka配置
// func initKafka(ctx context.Context) {
// 	fmt.Println("正在初始化Kafka...")

// 	// Kafka配置
// 	kafkaConfig := sarama.NewConfig()
// 	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
// 	kafkaConfig.Producer.Retry.Max = 3
// 	kafkaConfig.Producer.Return.Successes = true

// 	// 设置Kafka主题
// 	topic := consts.KafkaTopicSeckill

// 	// 检查主题是否存在，不存在则创建
// 	brokers := []string{"localhost:9092"}

// 	// 尝试连接Kafka
// 	admin, err := sarama.NewClusterAdmin(brokers, kafkaConfig)
// 	if err != nil {
// 		fmt.Printf("连接Kafka失败: %v\n", err)
// 		fmt.Println("无法创建Kafka管理客户端，但测试会继续进行")
// 		return
// 	}
// 	defer admin.Close()

// 	// 检查主题是否存在
// 	topics, err := admin.ListTopics()
// 	if err != nil {
// 		fmt.Printf("列出Kafka主题失败: %v\n", err)
// 		return
// 	}

// 	// 如果主题不存在，则创建
// 	if _, exists := topics[topic]; !exists {
// 		topicDetail := &sarama.TopicDetail{
// 			NumPartitions:     1,
// 			ReplicationFactor: 1,
// 		}
// 		err = admin.CreateTopic(topic, topicDetail, false)
// 		if err != nil {
// 			fmt.Printf("创建Kafka主题失败: %v\n", err)
// 		} else {
// 			fmt.Printf("成功创建Kafka主题: %s\n", topic)
// 		}
// 	} else {
// 		fmt.Printf("Kafka主题已存在: %s\n", topic)
// 	}

// 	// 直接设置环境变量来配置Kafka
// 	// 这将直接影响运行时系统
// 	os.Setenv("GF_MQ_KAFKA_BROKERS", "localhost:9092")
// 	os.Setenv("GF_MQ_KAFKA_PRODUCER_RETRYMAX", "3")
// 	os.Setenv("GF_MQ_KAFKA_PRODUCER_TIMEOUT", "5")

// 	fmt.Println("成功设置Kafka环境变量")
// 	fmt.Println("Kafka初始化完成")
// }
