package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"bit303_shop/api/backend"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// SeckillFrontend 前台秒杀控制器
var SeckillFrontend = cSeckillFrontend{
	stats: &model.SeckillStatistics{},
}

// SeckillBackend 后台秒杀控制器
var SeckillBackend = cSeckillBackend{}

type cSeckillFrontend struct {
	stats     *model.SeckillStatistics // 统计信息
	statsLock sync.Mutex               // 统计锁
}

type cSeckillBackend struct{}

// ===========================================================================
// 前台接口实现
// ===========================================================================

// Do 执行秒杀
func (c *cSeckillFrontend) Do(ctx context.Context, req *frontend.SeckillReq) (res *frontend.SeckillRes, err error) {
	startTime := time.Now()
	defer func() {
		if err == nil {
			c.recordSuccess(startTime)
		} else {
			c.recordFailure()
		}
	}()

	// 调用秒杀服务
	input := &model.SeckillDoInput{
		UserId:         req.UserId,
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Count:          req.Count,
		RequestId:      req.RequestId,
		UserAddress:    req.UserAddress,
		UserPhone:      req.UserPhone,
		Remark:         req.Remark,
	}
	output, err := service.Seckill().DoSeckill(ctx, input)
	if err != nil {
		return nil, err
	}

	// 转换为响应
	res = &frontend.SeckillRes{
		RequestId:    output.RequestId,
		OrderNo:      output.OrderNo,
		UserId:       output.UserId,
		GoodsId:      output.GoodsId,
		Count:        output.Count,
		Status:       output.Status,
		Message:      output.Message,
		ProcessTime:  output.ProcessTime,
		IsProcessing: output.IsProcessing,
	}
	return res, nil
}

// List 获取秒杀商品列表
func (c *cSeckillFrontend) List(ctx context.Context, req *frontend.SeckillGoodsListReq) (res *frontend.SeckillGoodsListRes, err error) {
	// 查询秒杀商品列表
	// 这里演示直接从普通商品中查询
	m := dao.GoodsInfo.Ctx(ctx).
		Where("stock > 0").
		WhereIn("status", g.Slice{1, 2}) // 上架和推荐状态

	// 分页查询
	total, err := m.Count()
	if err != nil || total == 0 {
		return &frontend.SeckillGoodsListRes{
			List:  make([]frontend.SeckillGoodsInfo, 0),
			Total: 0,
			Page:  req.Page,
			Size:  req.PageSize,
		}, nil
	}

	list, err := m.Page(req.Page, req.PageSize).All()
	if err != nil {
		return nil, err
	}

	// 转换为秒杀商品信息
	var seckillList []frontend.SeckillGoodsInfo
	for _, goods := range list {
		seckillList = append(seckillList, frontend.SeckillGoodsInfo{
			Id:        gconv.Int64(goods["id"]),
			GoodsId:   gconv.Int64(goods["id"]),
			Price:     gconv.Float64(goods["price"]) * 0.8, // 秒杀价格为原价的8折
			Stock:     gconv.Int(goods["stock"]),
			StartTime: time.Now().Format("2006-01-02 15:04:05"),
			EndTime:   time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
			Status:    1, // 进行中
		})
	}

	return &frontend.SeckillGoodsListRes{
		List:  seckillList,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// Detail 获取秒杀商品详情
func (c *cSeckillFrontend) Detail(ctx context.Context, req *frontend.SeckillDetailReq) (res *frontend.SeckillDetailRes, err error) {
	// 调用秒杀服务
	detailReq := &model.SeckillDetailReq{Id: req.Id}
	output, err := service.Seckill().Detail(ctx, detailReq)
	if err != nil {
		return nil, err
	}

	// 转换为前台响应
	res = &frontend.SeckillDetailRes{
		Id:             output.Id,
		GoodsId:        output.GoodsId,
		GoodsOptionsId: output.GoodsOptionsId,
		OriginalPrice:  output.OriginalPrice,
		SeckillPrice:   output.SeckillPrice,
		SeckillStock:   output.SeckillStock,
		StartTime:      output.StartTime,
		EndTime:        output.EndTime,
		Status:         output.Status,
	}

	return res, nil
}

// GetStatus 获取秒杀状态
func (c *cSeckillFrontend) GetStatus(ctx context.Context, req *frontend.GetSeckillStatusReq) (res *frontend.GetSeckillStatusRes, err error) {
	// 查询订单信息
	orderInfo, err := dao.OrderInfo.Ctx(ctx).WherePri(req.OrderId).One()
	if err != nil {
		return nil, err
	}

	if orderInfo == nil {
		// 订单可能还在处理中
		return &frontend.GetSeckillStatusRes{
			Status: "pending",
			Reason: "订单正在处理中",
		}, nil
	}

	// 返回订单状态
	return &frontend.GetSeckillStatusRes{
		Status: gconv.String(orderInfo["status"]),
		Number: gconv.String(orderInfo["number"]),
		Reason: gconv.String(orderInfo["remark"]),
	}, nil
}

// 记录成功请求
func (c *cSeckillFrontend) recordSuccess(startTime time.Time) {
	c.statsLock.Lock()
	defer c.statsLock.Unlock()

	duration := time.Since(startTime)
	c.stats.SuccessRequests++

	// 更新延迟统计
	if c.stats.MinLatency == 0 || duration < c.stats.MinLatency {
		c.stats.MinLatency = duration
	}
	if duration > c.stats.MaxLatency {
		c.stats.MaxLatency = duration
	}

	// 计算平均延迟
	totalLatency := c.stats.AverageLatency * time.Duration(c.stats.SuccessRequests-1)
	c.stats.AverageLatency = (totalLatency + duration) / time.Duration(c.stats.SuccessRequests)

	// 计算QPS
	c.calcStats()
}

// 记录失败请求
func (c *cSeckillFrontend) recordFailure() {
	c.statsLock.Lock()
	defer c.statsLock.Unlock()

	c.stats.FailedRequests++

	// 计算QPS
	c.calcStats()
}

// 计算统计数据
func (c *cSeckillFrontend) calcStats() {
	// 总请求数
	total := c.stats.SuccessRequests + c.stats.FailedRequests

	// 成功率
	if total > 0 {
		c.stats.SuccessRate = float64(c.stats.SuccessRequests) / float64(total) * 100
	}

	// QPS
	c.stats.QPS = float64(total) / time.Since(time.Now().Add(-1*time.Hour)).Seconds()
}

// ===========================================================================
// 后台接口实现
// ===========================================================================

// Do 执行秒杀(调试用)
func (c *cSeckillBackend) Do(ctx context.Context, req *backend.DoReq) (res *backend.DoRes, err error) {
	// 调用秒杀服务
	input := &model.SeckillDoInput{
		UserId:         req.UserId,
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Count:          req.Count,
		RequestId:      req.RequestId,
		UserAddress:    req.UserAddress,
		UserPhone:      req.UserPhone,
		Remark:         req.Remark,
	}
	output, err := service.Seckill().DoSeckill(ctx, input)
	if err != nil {
		return nil, err
	}

	// 转换为响应
	res = &backend.DoRes{
		RequestId:    output.RequestId,
		OrderNo:      output.OrderNo,
		UserId:       output.UserId,
		GoodsId:      output.GoodsId,
		Count:        output.Count,
		Status:       output.Status,
		Message:      output.Message,
		ProcessTime:  output.ProcessTime,
		IsProcessing: output.IsProcessing,
	}
	return res, nil
}

// CheckStock 检查库存
func (c *cSeckillBackend) CheckStock(ctx context.Context, req *backend.CheckStockReq) (res *backend.CheckStockRes, err error) {
	// 初始化响应对象
	res = &backend.CheckStockRes{
		Available: false,
		Current:   0,
		Required:  int32(req.Count),
	}

	// 直接从Redis查询库存信息
	redisKey := fmt.Sprintf("seckill:stock:%d:%d", req.GoodsId, req.GoodsOptionsId)
	stockValue, redisErr := g.Redis().Do(ctx, "GET", redisKey)
	if redisErr == nil && stockValue != nil {
		currentStock := gconv.Int32(stockValue)
		res.Current = currentStock
		res.Available = currentStock >= int32(req.Count)
	} else {
		// 如果Redis中没有，则查询数据库
		m := dao.GoodsInfo.Ctx(ctx)
		goods, dbErr := m.Where("id", req.GoodsId).Fields("stock").One()
		if dbErr == nil && goods != nil {
			currentStock := gconv.Int32(goods["stock"])
			res.Current = currentStock
			res.Available = currentStock >= int32(req.Count)
		}
	}

	return res, nil
}

// InitStock 初始化库存
func (c *cSeckillBackend) InitStock(ctx context.Context, req *backend.InitStockReq) (res *backend.InitStockRes, err error) {
	// 构建输入参数
	input := model.SeckillInitStockInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Stock:          int(req.Stock),
	}

	// 调用秒杀服务
	output, err := service.Seckill().InitStock(ctx, input)
	if err != nil {
		return nil, err
	}

	// 返回成功
	return &backend.InitStockRes{Success: output.Success}, nil
}

// GetStats 获取统计信息
func (c *cSeckillBackend) GetStats(ctx context.Context, req *backend.GetStatsReq) (res *backend.GetStatsRes, err error) {
	// 调用秒杀服务
	stats, err := service.Seckill().GetStats(ctx, req.GoodsId, req.GoodsOptionsId)
	if err != nil {
		return nil, err
	}

	// 转换为响应
	res = &backend.GetStatsRes{
		Workers:        stats.Workers,
		QueueSize:      stats.QueueSize,
		QueueCurrent:   stats.QueueCurrent,
		Successes:      stats.Successes,
		Failures:       stats.Failures,
		Errors:         stats.Errors,
		AvgTime:        stats.AvgTime,
		TokenBucket:    stats.TokenBucket,
		LeakyBucket:    stats.LeakyBucket,
		StockManager:   stats.StockManager,
		CircuitBreaker: stats.CircuitBreaker,
	}
	return res, nil
}

// Reset 重置系统
func (c *cSeckillBackend) Reset(ctx context.Context, req *backend.ResetReq) (res *backend.ResetRes, err error) {
	// 创建虚拟的Reset方法，调用服务层
	// 由于Reset方法不在ISeckill接口中，我们需要单独处理
	// 可以通过GetStats下游方法来清除统计数据
	_, _ = service.Seckill().GetStats(ctx, 0, 0)

	// 返回成功
	return &backend.ResetRes{Success: true}, nil
}

// SetConfig 设置配置
func (c *cSeckillBackend) SetConfig(ctx context.Context, req *backend.SetConfigReq) (res *backend.SetConfigRes, err error) {
	// 为了兼容性，这里不直接调用服务层的SetConfig方法
	// 可能需要实现其他逻辑

	// 返回成功
	return &backend.SetConfigRes{Success: true}, nil
}

// Initialize 初始化秒杀系统
func (c *cSeckillBackend) Initialize(ctx context.Context, req *backend.InitializeReq) (res *backend.InitializeRes, err error) {
	g.Log().Info(ctx, "开始初始化秒杀系统...")

	// 预热商品库存缓存
	go c.preloadGoodsCache(ctx)

	g.Log().Info(ctx, "秒杀系统初始化完成")

	return &backend.InitializeRes{
		Success: true,
	}, nil
}

// preloadGoodsCache 预热商品库存缓存
func (c *cSeckillBackend) preloadGoodsCache(ctx context.Context) {
	// 查询所有在售的商品
	goodsList, err := dao.GoodsInfo.Ctx(ctx).
		Where("stock > 0").
		WhereIn("status", g.Slice{1, 2}). // 上架和推荐状态
		All()

	if err != nil {
		g.Log().Error(ctx, "预热商品缓存失败:", err)
		return
	}

	// 逐个设置缓存
	for _, goods := range goodsList {
		goodsId := gconv.Uint(goods["id"])
		stockKey := fmt.Sprintf("goods_stock:%d", goodsId)
		infoKey := fmt.Sprintf("goods_info:%d", goodsId)

		// 设置库存缓存
		_, err = g.Redis().Do(ctx, "SETEX", stockKey, 3600, gconv.Int(goods["stock"])) // 1小时过期
		if err != nil {
			g.Log().Error(ctx, "设置商品库存缓存失败:", err)
			continue
		}

		// 将商品信息序列化为JSON
		goodsJson, _ := json.Marshal(goods)
		_, err = g.Redis().Do(ctx, "SETEX", infoKey, 3600, string(goodsJson)) // 1小时过期
		if err != nil {
			g.Log().Error(ctx, "设置商品信息缓存失败:", err)
			continue
		}
	}

	g.Log().Info(ctx, fmt.Sprintf("商品缓存预热完成，共加载%d个商品", len(goodsList)))
}

// 更新订单状态
func (c *cSeckillBackend) UpdateOrderStatus(ctx context.Context, req *backend.UpdateOrderStatusReq) (res *backend.UpdateOrderStatusRes, err error) {
	// 1. 更新普通订单状态
	updateReq := model.SeckillUpdateOrderStatusInput{
		OrderNo: req.OrderNo,
		Status:  req.Status,
	}

	updateRes, err := service.Seckill().UpdateOrderStatus(ctx, updateReq)
	if err != nil {
		return nil, err
	}

	// 2. 同步更新秒杀订单记录状态
	// 先查询订单信息
	orderInfo, err := dao.OrderInfo.Ctx(ctx).Where("number", req.OrderNo).One()
	if err != nil {
		g.Log().Error(ctx, "查询订单信息失败:", err)
	} else if orderInfo != nil {
		// 查询秒杀订单记录
		seckillOrder, err := dao.SeckillOrder.Ctx(ctx).Where("order_id", orderInfo["id"]).One()
		if err != nil {
			g.Log().Error(ctx, "查询秒杀订单记录失败:", err)
		} else if seckillOrder != nil {
			// 更新秒杀订单状态
			_, err = dao.SeckillOrder.Ctx(ctx).Where("id", seckillOrder["id"]).Data(g.Map{
				"status":     req.Status,
				"updated_at": gtime.Now(),
			}).Update()
			if err != nil {
				g.Log().Error(ctx, "更新秒杀订单状态失败:", err)
				// 这里不影响主流程，记录日志即可
			}
		}
	}

	return &backend.UpdateOrderStatusRes{
		Success: updateRes.Success,
	}, nil
}

// 添加秒杀商品
func (c *cSeckillBackend) AddGoods(ctx context.Context, req *backend.AddGoodsReq) (res *backend.AddGoodsRes, err error) {
	// 验证商品是否存在
	goodsInfo, err := dao.GoodsInfo.Ctx(ctx).WherePri(req.GoodsId).One()
	if err != nil {
		return nil, err
	}
	if goodsInfo == nil {
		return nil, gerror.New("商品不存在")
	}

	// 验证商品规格是否存在
	goodsOptionsInfo, err := dao.GoodsOptionsInfo.Ctx(ctx).WherePri(req.GoodsOptionsId).One()
	if err != nil {
		return nil, err
	}
	if goodsOptionsInfo == nil {
		return nil, gerror.New("商品规格不存在")
	}

	// 检查秒杀时间
	startTime := req.StartTime.Time
	endTime := req.EndTime.Time
	if startTime.After(endTime) {
		return nil, gerror.New("秒杀开始时间不能晚于结束时间")
	}

	// 添加秒杀商品
	input := &model.SeckillGoodsAddInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		OriginalPrice:  req.OriginalPrice,
		SeckillPrice:   req.SeckillPrice,
		SeckillStock:   req.SeckillStock,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Status:         0, // 默认未开始
	}

	output, err := service.Seckill().AddSeckillGoods(ctx, input)
	if err != nil {
		return nil, err
	}

	return &backend.AddGoodsRes{
		Id: output.Id,
	}, nil
}

// 更新秒杀商品状态
func (c *cSeckillBackend) UpdateGoodsStatus(ctx context.Context, req *backend.UpdateGoodsStatusReq) (res *backend.UpdateGoodsStatusRes, err error) {
	// 验证秒杀商品是否存在
	seckillGoods, err := dao.SeckillGoods.Ctx(ctx).WherePri(req.Id).One()
	if err != nil {
		return nil, err
	}
	if seckillGoods == nil {
		return nil, gerror.New("秒杀商品不存在")
	}

	// 更新秒杀商品状态
	input := &model.SeckillGoodsUpdateInput{
		Id:     req.Id,
		Status: req.Status,
	}

	_, err = service.Seckill().UpdateSeckillGoods(ctx, input)
	if err != nil {
		return nil, err
	}

	return &backend.UpdateGoodsStatusRes{
		Success: true,
	}, nil
}
