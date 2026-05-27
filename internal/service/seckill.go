// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
)

type (
	ISeckill interface {
		// 停止秒杀系统
		Stop()
		// DoSeckill 执行秒杀操作(外部调用入口)
		DoSeckill(ctx context.Context, req *model.SeckillDoInput) (res *model.SeckillDoOutput, err error)
		// CheckStock 检查商品库存
		CheckStock(ctx context.Context, in *model.SeckillCheckStockInput) (*model.SeckillCheckStockOutput, error)
		// InitStock 初始化秒杀商品库存
		InitStock(ctx context.Context, req model.SeckillInitStockInput) (res *model.SeckillInitStockOutput, err error)
		// GetStats 获取秒杀系统统计信息
		GetStats(ctx context.Context, goodsId, optionsId int64) (stats *model.SeckillStatsOutput, err error)
		// Reset 重置秒杀系统
		Reset()
		// SetConfig 设置秒杀系统配置
		SetConfig(config *model.SeckillConfig)
		// Close 关闭秒杀系统资源
		Close() error
		// GetResult 获取秒杀结果
		GetResult(ctx context.Context, req model.SeckillResultInput) (res *model.SeckillResultOutput, err error)
		// Detail 获取秒杀商品详情
		Detail(ctx context.Context, req *model.SeckillDetailReq) (res *model.SeckillDetailRes, err error)
		// List 获取秒杀商品列表
		List(ctx context.Context, req *model.SeckillListInput) (res *model.SeckillListOutput, err error)
		// WarmUpCache 缓存预热
		WarmUpCache(ctx context.Context, goodsId int64, optionsId int64) error
		// GetSeckillInfo 获取秒杀商品信息
		GetSeckillInfo(ctx context.Context, goodsId int64, optionsId int64) (*entity.SeckillGoods, error)
		// UpdateOrderStatus 更新订单状态
		UpdateOrderStatus(ctx context.Context, req model.SeckillUpdateOrderStatusInput) (res *model.SeckillUpdateOrderStatusOutput, err error)
		// GetConfig 获取当前秒杀配置
		GetConfig() *model.SeckillConfig
		// AddSeckillGoods 添加秒杀商品
		AddSeckillGoods(ctx context.Context, req *model.SeckillGoodsAddInput) (res *model.SeckillGoodsAddOutput, err error)
		// UpdateSeckillGoods 更新秒杀商品
		UpdateSeckillGoods(ctx context.Context, req *model.SeckillGoodsUpdateInput) (res *model.SeckillGoodsUpdateOutput, err error)
		// AddSeckillOrder 添加秒杀订单
		AddSeckillOrder(ctx context.Context, req *model.SeckillOrderAddInput) (res *model.SeckillOrderAddOutput, err error)
		// UpdateSeckillOrder 更新秒杀订单
		UpdateSeckillOrder(ctx context.Context, req *model.SeckillOrderUpdateInput) (res *model.SeckillOrderUpdateOutput, err error)
		// GetSeckillOrderByOrderId 根据订单ID获取秒杀订单
		GetSeckillOrderByOrderId(ctx context.Context, req *model.SeckillOrderByOrderIdInput) (res *model.SeckillOrderByOrderIdOutput, err error)
	}
)

var (
	localSeckill ISeckill
)

func Seckill() ISeckill {
	if localSeckill == nil {
		panic("implement not found for interface ISeckill, forgot register?")
	}
	return localSeckill
}

func RegisterSeckill(i ISeckill) {
	localSeckill = i
}
