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
	ICart interface {
		Add(ctx context.Context, in model.CartAddInput) (out model.CartAddOutput, err error)
		List(ctx context.Context, in model.CartListInput) (out model.CartListOutput, err error)
		Update(ctx context.Context, in model.CartUpdateInput) (out model.CartUpdateOutput, err error)
		Remove(ctx context.Context, in model.CartRemoveInput) (out model.CartRemoveOutput, err error)
	}
)

var (
	localCart ICart
)

func Cart() ICart {
	if localCart == nil {
		panic("implement not found for interface ICart, forgot register?")
	}
	return localCart
}

func RegisterCart(i ICart) {
	localCart = i
}
