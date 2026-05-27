package mq

// import (
// 	"context"
// 	"encoding/json"
// 	"bit303_shop/internal/consts"
// 	"bit303_shop/internal/dao"
// 	"bit303_shop/internal/model"
// 	"sync"
// 	"time"

// 	"github.com/IBM/sarama"
// 	"github.com/gogf/gf/v2/frame/g"
// 	"github.com/gogf/gf/v2/os/gctx"
// 	"github.com/gogf/gf/v2/os/gtime"
// 	"github.com/gogf/gf/v2/util/gconv"
// )

// type SeckillConsumer struct {
// 	consumer sarama.Consumer
// 	producer sarama.SyncProducer
// 	done     chan struct{}
// 	wg       sync.WaitGroup
// }

// // 初始化时自动注册
// func init() {
// 	// 异步启动消费者
// 	go func() {
// 		// 等待系统初始化完成
// 		time.Sleep(5 * time.Second)
// 		ctx := gctx.New()

// 		// 尝试创建消费者
// 		consumer, err := NewSeckillConsumer()
// 		if err != nil {
// 			g.Log().Warning(ctx, "秒杀Kafka消费者创建失败，将不会处理异步订单: ", err)
// 			return
// 		}

// 		// 启动消费者
// 		if err := consumer.Start(ctx); err != nil {
// 			g.Log().Error(ctx, "秒杀Kafka消费者启动失败: ", err)
// 			consumer.Stop()
// 			return
// 		}

// 		g.Log().Info(ctx, "秒杀Kafka消费者启动成功，开始处理订单")
// 	}()
// }

// func NewSeckillConsumer() (*SeckillConsumer, error) {
// 	// 从配置获取Kafka地址
// 	kafkaAddrs := g.Cfg().MustGet(gctx.New(), "kafka.addrs", []string{"localhost:9092"}).Strings()

// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true
// 	config.Producer.Return.Successes = true
// 	config.Producer.Timeout = 5 * time.Second

// 	consumer, err := sarama.NewConsumer(kafkaAddrs, config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	producer, err := sarama.NewSyncProducer(kafkaAddrs, config)
// 	if err != nil {
// 		consumer.Close()
// 		return nil, err
// 	}

// 	return &SeckillConsumer{
// 		consumer: consumer,
// 		producer: producer,
// 		done:     make(chan struct{}),
// 	}, nil
// }

// func (c *SeckillConsumer) Start(ctx context.Context) error {
// 	partitionConsumer, err := c.consumer.ConsumePartition(consts.SeckillKafkaTopic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		return err
// 	}

// 	c.wg.Add(1)
// 	go func() {
// 		defer c.wg.Done()
// 		defer partitionConsumer.Close()

// 		for {
// 			select {
// 			case msg := <-partitionConsumer.Messages():
// 				if err := c.processMessage(ctx, msg); err != nil {
// 					g.Log().Error(ctx, "处理秒杀订单消息失败:", err)
// 					c.sendToRetryQueue(ctx, msg.Value)
// 				}
// 			case err := <-partitionConsumer.Errors():
// 				g.Log().Error(ctx, "消费秒杀订单消息失败:", err)
// 			case <-c.done:
// 				return
// 			}
// 		}
// 	}()

// 	return nil
// }

// func (c *SeckillConsumer) Stop() {
// 	close(c.done)
// 	c.wg.Wait()
// 	if c.consumer != nil {
// 		if err := c.consumer.Close(); err != nil {
// 			g.Log().Error(gctx.New(), "Failed to close consumer:", err)
// 		}
// 	}
// 	if c.producer != nil {
// 		if err := c.producer.Close(); err != nil {
// 			g.Log().Error(gctx.New(), "Failed to close producer:", err)
// 		}
// 	}
// }

// func (c *SeckillConsumer) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
// 	var orderMsg model.SeckillOrderMsg
// 	if err := json.Unmarshal(msg.Value, &orderMsg); err != nil {
// 		return err
// 	}

// 	g.Log().Info(ctx, "收到秒杀订单处理请求:", orderMsg.OrderNo)

// 	// 开启事务
// 	tx, err := g.DB().Begin(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// 计算单价和原价
// 	var price float64
// 	if orderMsg.Count > 0 {
// 		price = orderMsg.TotalPrice / float64(orderMsg.Count)
// 	}
// 	// 设置单价
// 	orderMsg.Price = price

// 	// 转换为秒杀订单输入
// 	seckillOrderInput := orderMsg.ToSeckillOrderAddInput()

// 	// 创建秒杀订单
// 	seckillOrderData := g.Map{
// 		"user_id":          seckillOrderInput.UserId,
// 		"order_no":         seckillOrderInput.OrderNo,
// 		"original_price":   seckillOrderInput.OriginalPrice,
// 		"seckill_price":    seckillOrderInput.SeckillPrice,
// 		"goods_id":         seckillOrderInput.GoodsId,
// 		"goods_options_id": seckillOrderInput.GoodsOptionsId,
// 		"status":           seckillOrderInput.Status,
// 		"pay_time":         seckillOrderInput.PayTime,
// 		"count":            seckillOrderInput.Count,
// 		"consignee_name":   seckillOrderInput.ConsigneeName,
// 		"consignee_phone":  seckillOrderInput.ConsigneePhone,
// 		"address":          seckillOrderInput.Address,
// 		"remark":           seckillOrderInput.Remark,
// 		"created_at":       gtime.Now(),
// 		"updated_at":       gtime.Now(),
// 	}

// 	// 创建秒杀订单，不再创建普通订单
// 	orderId, err := dao.SeckillOrder.Ctx(ctx).TX(tx).InsertAndGetId(seckillOrderData)
// 	if err != nil {
// 		return err
// 	}

// 	// 提交事务
// 	if err = tx.Commit(); err != nil {
// 		return err
// 	}

// 	// 发送订单创建成功的通知
// 	go c.sendNotification(ctx, gconv.Int64(orderMsg.UserId), orderMsg.OrderNo)

// 	g.Log().Info(ctx, "秒杀订单创建成功，ID:", orderId, "订单号:", orderMsg.OrderNo)
// 	return nil
// }

// func (c *SeckillConsumer) sendToRetryQueue(ctx context.Context, msg []byte) {
// 	retryMsg := &sarama.ProducerMessage{
// 		Topic: consts.SeckillKafkaTopic + "_retry",
// 		Value: sarama.ByteEncoder(msg),
// 	}
// 	_, _, err := c.producer.SendMessage(retryMsg)
// 	if err != nil {
// 		g.Log().Error(ctx, "发送重试消息失败:", err)
// 	}
// }

// func (c *SeckillConsumer) sendNotification(ctx context.Context, userId int64, orderNo string) error {
// 	// 创建通知消息
// 	notification := &model.OrderNotification{
// 		OrderNo:   orderNo,
// 		UserId:    uint(userId),
// 		Status:    "success",
// 		Message:   "您的秒杀订单已创建成功，请尽快支付",
// 		Timestamp: time.Now().Unix(),
// 	}

// 	// 序列化消息
// 	msgData, err := json.Marshal(notification)
// 	if err != nil {
// 		return err
// 	}

// 	// 发送通知消息
// 	notifyMsg := &sarama.ProducerMessage{
// 		Topic: "order_notifications",
// 		Key:   sarama.StringEncoder(orderNo),
// 		Value: sarama.ByteEncoder(msgData),
// 	}

// 	// 异步发送，不关心结果
// 	_, _, _ = c.producer.SendMessage(notifyMsg)

// 	// 同时将通知保存到Redis，方便用户查询
// 	notifyKey := "notify:order:" + orderNo
// 	_, _ = g.Redis().Do(ctx, "SETEX", notifyKey, 3600*24, string(msgData))

// 	return nil
// }
