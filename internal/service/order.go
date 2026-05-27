// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/model"
)

type (
	IOrder interface {
		// 下单
		Add(ctx context.Context, in model.OrderAddInput) (out *model.OrderAddOutput, err error)
		List(ctx context.Context, in model.OrderListInput) (out *model.OrderListOutput, err error)
		Detail(ctx context.Context, in model.OrderDetailInput) (out *model.OrderDetailOutput, err error)
		// Pay 支付订单
		Pay(ctx context.Context, in *frontend.OrderPayReq) (out *frontend.OrderPayRes, err error)
		// Cancel 取消订单
		Cancel(ctx context.Context, in *frontend.OrderCancelReq) (out *frontend.OrderCancelRes, err error)
		// Confirm 确认收货
		Confirm(ctx context.Context, in *frontend.OrderConfirmReq) (out *frontend.OrderConfirmRes, err error)
		// UpdateOrderStatus 更新订单状态(后台)
		UpdateOrderStatus(ctx context.Context, in *backend.OrderUpdateStatusReq) (out *backend.OrderUpdateStatusRes, err error)
		// Delete 删除订单(后台)
		Delete(ctx context.Context, in *backend.OrderDeleteReq) (out *backend.OrderDeleteRes, err error)
		// Refund 订单退款(后台)
		Refund(ctx context.Context, in *backend.OrderRefundReq) (out *backend.OrderRefundRes, err error)
		// CreateSeckillOrder 创建秒杀订单
		CreateSeckillOrder(ctx context.Context, in model.SeckillOrderInput) (out *model.OrderAddOutput, err error)
		// ProcessOrderMessage 处理订单消息
		ProcessOrderMessage(ctx context.Context, message []byte) error
	}
)

var (
	localOrder IOrder
)

func Order() IOrder {
	if localOrder == nil {
		panic("implement not found for interface IOrder, forgot register?")
	}
	return localOrder
}

func RegisterOrder(i IOrder) {
	localOrder = i
}
