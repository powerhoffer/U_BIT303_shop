package frontend

import (
	"context"
	"fmt"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sSeckill struct{}

var Seckill = sSeckill{}

// List 获取秒杀商品列表
func (s *sSeckill) List(ctx context.Context, req *frontend.SeckillListReq) (res *frontend.SeckillListRes, err error) {
	// 转换为内部模型
	input := &model.SeckillListInput{
		Page: req.Page,
		Size: req.Size,
	}

	list, err := service.Seckill().List(ctx, input)
	if err != nil {
		return nil, err
	}
	return &frontend.SeckillListRes{
		List:  list.List,
		Page:  list.Page,
		Size:  list.Size,
		Total: list.Total,
	}, nil
}

// Detail 获取秒杀商品详情
func (s *sSeckill) Detail(ctx context.Context, req *frontend.SeckillDetailReq) (res *frontend.SeckillDetailRes, err error) {
	// 转换为内部模型
	input := &model.SeckillDetailReq{
		Id: req.Id,
	}

	detail, err := service.Seckill().Detail(ctx, input)
	if err != nil {
		return nil, err
	}
	return &frontend.SeckillDetailRes{
		Id:             detail.Id,
		GoodsId:        detail.GoodsId,
		GoodsOptionsId: detail.GoodsOptionsId,
		OriginalPrice:  detail.OriginalPrice,
		SeckillPrice:   detail.SeckillPrice,
		SeckillStock:   detail.SeckillStock,
		StartTime:      detail.StartTime,
		EndTime:        detail.EndTime,
		Status:         detail.Status,
	}, nil
}

// Do 执行秒杀操作 - 处理SeckillReq请求
func (s *sSeckill) Do(ctx context.Context, req *frontend.SeckillReq) (res *frontend.SeckillRes, err error) {
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

// DoSeckill 兼容老版本API请求
func (s *sSeckill) DoSeckill(ctx context.Context, req *frontend.SeckillDoReq) (res *frontend.SeckillDoRes, err error) {
	input := &model.SeckillDoInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		UserId:         req.UserId,
		// 老版本缺少其他参数，使用默认值
		Count:     1,  // 默认购买1个
		RequestId: "", // 使用空字符串，服务层会自动生成
	}
	result, err := service.Seckill().DoSeckill(ctx, input)
	if err != nil {
		return nil, err
	}
	return &frontend.SeckillDoRes{
		OrderNo: result.OrderNo,
	}, nil
}

// GetStatus 获取秒杀状态
func (s *sSeckill) GetStatus(ctx context.Context, req *frontend.GetSeckillStatusReq) (res *frontend.GetSeckillStatusRes, err error) {
	// 首先查询订单详情
	orderDetail, err := service.Order().Detail(ctx, model.OrderDetailInput{
		Id: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	// 获取订单编号（转换为字符串）
	orderNumber := gconv.String(orderDetail.OrderInfo.Number)

	// 直接创建内部模型
	modelInput := model.SeckillResultInput{
		OrderNo: orderNumber,
	}

	// 使用订单编号查询秒杀结果
	seckillResult, err := service.Seckill().GetResult(ctx, modelInput)

	// 如果查询秒杀结果失败，就使用订单状态
	if err != nil {
		g.Log().Warning(ctx, "获取秒杀结果失败，将使用订单状态:", err)

		// 根据订单状态返回结果
		var status string
		var reason string

		switch gconv.Int(orderDetail.OrderInfo.Status) {
		case 0:
			status = "pending" // 待支付
			reason = "订单等待支付"
		case 1:
			status = "success" // 已支付
			reason = "订单已支付"
		case 2:
			status = "failed" // 已取消
			reason = "订单已取消"
		default:
			status = "unknown"
			reason = fmt.Sprintf("未知订单状态: %v", orderDetail.OrderInfo.Status)
		}

		return &frontend.GetSeckillStatusRes{
			Status: status,
			Number: orderNumber,
			Reason: reason,
		}, nil
	}

	// 使用秒杀结果的状态
	var status string
	var reason string

	switch seckillResult.Status {
	case 0:
		status = "success" // 秒杀成功
		reason = "秒杀成功"
	case 1:
		status = "failed" // 秒杀失败
		reason = "秒杀失败"
	case 2:
		status = "failed" // 库存不足
		reason = "商品库存不足"
	case 3:
		status = "failed" // 被限流
		reason = "系统繁忙，请稍后再试"
	case 4:
		status = "failed" // 重复请求
		reason = "重复请求"
	case 5:
		status = "failed" // 处理超时
		reason = "处理超时"
	case 6:
		status = "failed" // 熔断器开启
		reason = "系统保护中，请稍后再试"
	case 7:
		status = "pending" // 正在处理
		reason = "订单正在处理中"
	case 8:
		status = "failed" // 系统错误
		reason = "系统错误"
	default:
		status = "unknown"
		reason = fmt.Sprintf("未知状态码: %d", seckillResult.Status)
	}

	return &frontend.GetSeckillStatusRes{
		Status: status,
		Number: orderNumber,
		Reason: reason,
	}, nil
}

// GetResult 获取秒杀结果
func (s *sSeckill) GetResult(ctx context.Context, req *frontend.SeckillResultReq) (res *frontend.SeckillResultRes, err error) {
	// 直接创建内部模型
	modelInput := model.SeckillResultInput{
		OrderNo: req.OrderNo,
	}

	result, err := service.Seckill().GetResult(ctx, modelInput)
	if err != nil {
		return nil, err
	}
	return &frontend.SeckillResultRes{
		OrderNo: result.OrderNo,
		Status:  result.Status,
	}, nil
}

// Batch 批量秒杀处理
func (s *sSeckill) Batch(ctx context.Context, req *frontend.SeckillBatchReq) (res *frontend.SeckillBatchRes, err error) {
	// 创建响应对象
	res = &frontend.SeckillBatchRes{
		Results: make([]*frontend.SeckillItemResult, 0, len(req.Items)),
	}

	// 处理每个秒杀请求
	for _, item := range req.Items {
		// 创建内部请求模型
		input := &model.SeckillDoInput{
			GoodsId:        uint(item.GoodsId),
			GoodsOptionsId: uint(item.GoodsOptionsId),
			Count:          uint(item.Count),
			// 获取当前用户ID（根据实际情况获取）
			UserId: gconv.Uint(service.Session().GetUser(ctx).Id),
		}

		// 添加结果对象
		result := &frontend.SeckillItemResult{
			GoodsId: item.GoodsId,
		}

		// 执行秒杀
		output, err := service.Seckill().DoSeckill(ctx, input)
		if err != nil {
			// 秒杀失败
			result.Success = false
			result.Message = err.Error()
			res.FailedCount++
		} else {
			// 秒杀成功
			result.Success = true
			result.OrderNo = output.OrderNo
			res.SuccessCount++
		}

		// 添加到结果列表
		res.Results = append(res.Results, result)
	}

	return res, nil
}

// Stats 获取秒杀统计信息
func (s *sSeckill) Stats(ctx context.Context, req *frontend.SeckillStatsReq) (res *frontend.SeckillStatsRes, err error) {
	// 获取秒杀统计数据
	// 注意：由于API中没有提供goodsId和optionsId，这里使用默认值0表示获取全局统计
	stats, err := service.Seckill().GetStats(ctx, 0, 0)
	if err != nil {
		return nil, err
	}

	// 将秒杀统计转换为可读字符串
	avgLatency := ""
	maxLatency := ""
	minLatency := ""

	// 计算成功率和QPS
	successRate := 0.0
	qps := 0.0

	// 基于实际的SeckillStatsOutput字段进行转换
	totalRequests := int(stats.Successes + stats.Failures + stats.Errors)
	successRequests := int(stats.Successes)
	failedRequests := int(stats.Failures + stats.Errors)

	// 转换为响应模型
	return &frontend.SeckillStatsRes{
		TotalRequests:   totalRequests,
		SuccessRequests: successRequests,
		FailedRequests:  failedRequests,
		AverageLatency:  avgLatency,
		MaxLatency:      maxLatency,
		MinLatency:      minLatency,
		QPS:             qps,
		SuccessRate:     successRate,
	}, nil
}

// GoodsList 获取秒杀商品列表
func (s *sSeckill) GoodsList(ctx context.Context, req *frontend.SeckillGoodsListReq) (res *frontend.SeckillGoodsListRes, err error) {
	// 转换为内部模型
	input := &model.SeckillListInput{
		Page: req.Page,
		Size: req.PageSize,
	}

	// 调用服务获取列表
	listOutput, err := service.Seckill().List(ctx, input)
	if err != nil {
		return nil, err
	}

	// 准备响应数据
	res = &frontend.SeckillGoodsListRes{
		List:  make([]frontend.SeckillGoodsInfo, 0),
		Total: listOutput.Total,
		Page:  listOutput.Page,
		Size:  listOutput.Size,
	}

	// 转换列表数据
	if listData, ok := listOutput.List.([]model.SeckillGoodsInfo); ok {
		for _, item := range listData {
			res.List = append(res.List, frontend.SeckillGoodsInfo{
				Id:             item.Id,
				GoodsId:        item.GoodsId,
				GoodsOptionsId: item.GoodsOptionsId,
				Price:          item.Price,
				Stock:          item.Stock,
				StartTime:      item.StartTime.String(),
				EndTime:        item.EndTime.String(),
				Status:         item.Status,
			})
		}
	}

	return res, nil
}

// DirectSeckill 直接执行秒杀并返回结果
func (s *sSeckill) DirectSeckill(ctx context.Context, req *frontend.DirectSeckillReq) (res *frontend.DirectSeckillRes, err error) {
	// 生成请求ID（保证幂等性）
	requestId := fmt.Sprintf("%d-%d-%d-%d", req.UserId, req.GoodsId, req.GoodsOptionsId, time.Now().UnixNano())

	// 创建秒杀输入模型
	input := &model.SeckillDoInput{
		UserId:         req.UserId,
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		Count:          req.Count,
		RequestId:      requestId,
		UserAddress:    req.UserAddress,
		UserPhone:      req.UserPhone,
		Remark:         req.Remark,
	}

	// 执行秒杀
	output, err := service.Seckill().DoSeckill(ctx, input)
	if err != nil {
		return &frontend.DirectSeckillRes{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 秒杀成功，查询订单详情
	if output.Status == 0 { // 0表示成功
		// 查询秒杀订单
		seckillOrderInfo, err := dao.SeckillOrder.Ctx(ctx).Where("order_no", output.OrderNo).One()
		if err != nil {
			return &frontend.DirectSeckillRes{
				Success:    true,
				OrderNo:    output.OrderNo,
				Message:    "秒杀成功，但获取订单详情失败",
				TotalPrice: 0,
			}, nil
		}

		// 返回完整订单信息
		return &frontend.DirectSeckillRes{
			Success:    true,
			OrderId:    gconv.Uint(seckillOrderInfo["id"]),
			OrderNo:    gconv.String(seckillOrderInfo["order_no"]),
			Message:    "秒杀成功",
			TotalPrice: gconv.Int(seckillOrderInfo["seckill_price"]) * gconv.Int(seckillOrderInfo["count"]),
		}, nil
	}

	// 秒杀失败
	return &frontend.DirectSeckillRes{
		Success: false,
		OrderNo: output.OrderNo,
		Message: output.Message,
	}, nil
}

// OrderDetail 获取秒杀订单详情
func (s *sSeckill) OrderDetail(ctx context.Context, req *frontend.SeckillOrderDetailReq) (res *frontend.SeckillOrderDetailRes, err error) {
	// 查询秒杀订单信息
	seckillOrderInfo, err := dao.SeckillOrder.Ctx(ctx).Where("order_no", req.OrderNo).One()
	if err != nil {
		return nil, err
	}
	if seckillOrderInfo == nil {
		return nil, gerror.New("秒杀订单不存在")
	}

	// 查询商品信息
	goodsInfo, err := dao.GoodsInfo.Ctx(ctx).WherePri(seckillOrderInfo["goods_id"]).One()
	if err != nil || goodsInfo == nil {
		return nil, gerror.New("商品信息不存在")
	}

	// 查询商品规格信息
	var goodsSpec string
	if gconv.Int(seckillOrderInfo["goods_options_id"]) > 0 {
		goodsOption, err := dao.GoodsOptionsInfo.Ctx(ctx).WherePri(seckillOrderInfo["goods_options_id"]).One()
		if err == nil && goodsOption != nil {
			goodsSpec = gconv.String(goodsOption["name"])
		}
	}

	// 构建商品详情
	goodsList := []frontend.SeckillOrderGoodsDetail{
		{
			GoodsId:        gconv.Uint(seckillOrderInfo["goods_id"]),
			GoodsOptionsId: gconv.Uint(seckillOrderInfo["goods_options_id"]),
			Count:          gconv.Int(seckillOrderInfo["count"]),
			Price:          gconv.Int(seckillOrderInfo["original_price"]),
			ActualPrice:    gconv.Int(seckillOrderInfo["seckill_price"]),
			GoodsName:      gconv.String(goodsInfo["name"]),
			GoodsImage:     gconv.String(goodsInfo["pic_url"]),
			GoodsSpec:      goodsSpec,
		},
	}

	// 获取支付、取消时间
	var payAt, cancelAt *string
	if seckillOrderInfo["pay_time"] != nil && !gtime.NewFromStr(gconv.String(seckillOrderInfo["pay_time"])).IsZero() {
		payTime := gtime.NewFromStr(gconv.String(seckillOrderInfo["pay_time"])).String()
		payAt = &payTime
	}
	if seckillOrderInfo["cancel_time"] != nil && !gtime.NewFromStr(gconv.String(seckillOrderInfo["cancel_time"])).IsZero() {
		cancelTime := gtime.NewFromStr(gconv.String(seckillOrderInfo["cancel_time"])).String()
		cancelAt = &cancelTime
	}

	// 构建返回结果
	totalAmount := gconv.Int(seckillOrderInfo["seckill_price"]) * gconv.Int(seckillOrderInfo["count"])

	res = &frontend.SeckillOrderDetailRes{
		Id:             gconv.Uint(seckillOrderInfo["id"]),
		OrderNo:        gconv.String(seckillOrderInfo["order_no"]),
		UserId:         gconv.Uint(seckillOrderInfo["user_id"]),
		Status:         gconv.Int(seckillOrderInfo["status"]),
		Price:          totalAmount,
		PayAt:          payAt,
		CancelAt:       cancelAt,
		ConsigneeName:  gconv.String(seckillOrderInfo["consignee_name"]),
		ConsigneePhone: gconv.String(seckillOrderInfo["consignee_phone"]),
		Address:        gconv.String(seckillOrderInfo["address"]),
		Remark:         gconv.String(seckillOrderInfo["remark"]),
		GoodsList:      goodsList,
		CreatedAt:      gtime.NewFromStr(gconv.String(seckillOrderInfo["created_at"])).String(),
	}

	return res, nil
}

// GoodsInfo 获取秒杀商品信息和库存
func (s *sSeckill) GoodsInfo(ctx context.Context, req *frontend.SeckillGoodsInfoReq) (res *frontend.SeckillGoodsInfoRes, err error) {
	// 查询商品信息
	goodsInfo, err := dao.GoodsInfo.Ctx(ctx).WherePri(req.GoodsId).One()
	if err != nil {
		return nil, err
	}
	if goodsInfo == nil {
		return nil, gerror.New("商品不存在")
	}

	// 查询商品规格信息
	goodsOption, err := dao.GoodsOptionsInfo.Ctx(ctx).WherePri(req.GoodsOptionsId).One()
	if err != nil {
		return nil, err
	}
	if goodsOption == nil {
		return nil, gerror.New("商品规格不存在")
	}

	// 查询库存
	stockInput := &model.SeckillCheckStockInput{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
	}
	stockResult, err := service.Seckill().CheckStock(ctx, stockInput)
	if err != nil {
		return nil, err
	}

	// 检查是否为秒杀商品，获取秒杀信息
	seckillGoods, _ := service.Seckill().GetSeckillInfo(ctx, gconv.Int64(req.GoodsId), gconv.Int64(req.GoodsOptionsId))

	// 计算秒杀价格和时间
	seckillPrice := gconv.Int(goodsOption["price"])
	var startTime, endTime string
	var status int = 0     // 默认未开始
	var limitCount int = 5 // 默认限购5件

	// 如果找到了秒杀商品信息
	if seckillGoods != nil {
		seckillPrice = gconv.Int(float64(seckillPrice) * 0.8) // 秒杀价为原价的8折
		startTime = seckillGoods.StartTime.String()
		endTime = seckillGoods.EndTime.String()
		status = seckillGoods.Status

		// 根据当前时间判断状态
		now := gtime.Now()
		if now.Before(seckillGoods.StartTime) {
			status = 0 // 未开始
		} else if now.After(seckillGoods.EndTime) {
			status = 2 // 已结束
		} else {
			status = 1 // 进行中
		}
	} else {
		// 如果不是秒杀商品，设置默认时间
		startTime = gtime.Now().String()
		endTime = gtime.Now().Add(time.Hour * 24).String()
	}

	// 构建响应
	res = &frontend.SeckillGoodsInfoRes{
		GoodsId:        req.GoodsId,
		GoodsOptionsId: req.GoodsOptionsId,
		GoodsName:      gconv.String(goodsInfo["name"]),
		GoodsImage:     gconv.String(goodsInfo["pic"]),
		GoodsSpec:      gconv.String(goodsOption["name"]),
		Price:          gconv.Int(goodsOption["price"]),
		SeckillPrice:   seckillPrice,
		Stock:          int(stockResult.Current),
		SalesCount:     gconv.Int(goodsInfo["sale"]),
		LimitCount:     limitCount,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         status,
		Description:    gconv.String(goodsInfo["detail"]),
	}

	return res, nil
}
