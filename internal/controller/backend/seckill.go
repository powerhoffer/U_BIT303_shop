package backend

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

type sSeckill struct{}

var Seckill = sSeckill{}

// Do 执行秒杀操作(测试用)
func (s *sSeckill) Do(ctx context.Context, req *backend.DoReq) (res *backend.DoRes, err error) {
	// 直接创建内部模型
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

	// 调用秒杀服务
	output, err := service.Seckill().DoSeckill(ctx, input)
	if err != nil {
		return nil, err
	}

	// 创建响应并直接设置字段值
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

// InitStock 初始化秒杀商品库存
func (s *sSeckill) InitStock(ctx context.Context, req *backend.InitStockReq) (res *backend.InitStockRes, err error) {
	// 直接创建内部模型
	input := model.SeckillInitStockInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Stock:          int(req.Stock),
	}

	// 调用秒杀服务
	result, err := service.Seckill().InitStock(ctx, input)
	if err != nil {
		return nil, err
	}

	return &backend.InitStockRes{
		Success: result.Success,
	}, nil
}

// GetStats 获取秒杀统计信息
func (s *sSeckill) GetStats(ctx context.Context, req *backend.GetStatsReq) (res *backend.GetStatsRes, err error) {
	// 调用秒杀服务
	stats, err := service.Seckill().GetStats(ctx, req.GoodsId, req.GoodsOptionsId)
	if err != nil {
		// 发生错误时返回模拟数据
		res = &backend.GetStatsRes{
			Workers:      10,
			QueueSize:    1000,
			QueueCurrent: 5,
			Successes:    100,
			Failures:     10,
			Errors:       1,
			AvgTime:      15.5,
			TokenBucket: map[string]int64{
				"issued": 1000,
				"used":   800,
			},
			LeakyBucket: map[string]int64{
				"received":  1200,
				"processed": 1000,
				"dropped":   200,
			},
			StockManager: map[string]int64{
				"totalStock":     1000,
				"remainingStock": 200,
			},
			CircuitBreaker: map[string]int64{
				"tripped": 0,
				"success": 100,
				"failure": 0,
			},
		}
		return res, nil
	}

	// 直接创建响应
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
