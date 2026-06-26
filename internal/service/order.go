// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"bit303_shop/internal/model"
	"context"
)

type (
	IOrder interface {
		Create(ctx context.Context, in model.OrderCreateInput) (out model.OrderCreateOutput, err error)
		List(ctx context.Context, in model.OrderListInput) (out model.OrderListOutput, err error)
		Detail(ctx context.Context, in model.OrderDetailInput) (out model.OrderDetailOutput, err error)
		Cancel(ctx context.Context, in model.OrderCancelInput) (out model.OrderCancelOutput, err error)
		ManageList(ctx context.Context, in model.BackendOrderListInput) (out model.BackendOrderListOutput, err error)
		ManageDetail(ctx context.Context, in model.BackendOrderDetailInput) (out model.BackendOrderDetailOutput, err error)
		ManageComplete(ctx context.Context, in model.BackendOrderCompleteInput) (out model.BackendOrderCompleteOutput, err error)
		ManageCancel(ctx context.Context, in model.BackendOrderCancelInput) (out model.BackendOrderCancelOutput, err error)
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
