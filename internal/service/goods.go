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
	IGoods interface {
		Create(ctx context.Context, in model.GoodsCreateInput) (out model.GoodsCreateOutput, err error)
		List(ctx context.Context, in model.GoodsListInput) (out model.GoodsListOutput, err error)
		Detail(ctx context.Context, id uint) (out model.GoodsDetailOutput, err error)
		Update(ctx context.Context, in model.GoodsUpdateInput) (out model.GoodsUpdateOutput, err error)
		UpdateStatus(ctx context.Context, in model.GoodsStatusInput) error
	}
)

var (
	localGoods IGoods
)

func Goods() IGoods {
	if localGoods == nil {
		panic("implement not found for interface IGoods, forgot register?")
	}
	return localGoods
}

func RegisterGoods(i IGoods) {
	localGoods = i
}
