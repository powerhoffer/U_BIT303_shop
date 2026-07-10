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
	IStock interface {
		Adjust(ctx context.Context, in model.StockAdjustInput) (out model.StockAdjustOutput, err error)
		RecordChange(ctx context.Context, in model.StockRecordInput) error
		Records(ctx context.Context, in model.StockRecordsInput) (out model.StockRecordsOutput, err error)
	}
)

var (
	localStock IStock
)

func Stock() IStock {
	if localStock == nil {
		panic("implement not found for interface IStock, forgot register?")
	}
	return localStock
}

func RegisterStock(i IStock) {
	localStock = i
}
