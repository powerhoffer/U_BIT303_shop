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
	IPoints interface {
		Balance(ctx context.Context, employeeId uint) (out model.PointsBalanceOutput, err error)
		Records(ctx context.Context, in model.PointsRecordsInput) (out model.PointsRecordsOutput, err error)
		ManageRecords(ctx context.Context, in model.PointsRecordsInput) (out model.PointsRecordsOutput, err error)
		ManageAdd(ctx context.Context, in model.PointsChangeInput) (out model.PointsChangeOutput, err error)
		ManageDeduct(ctx context.Context, in model.PointsChangeInput) (out model.PointsChangeOutput, err error)
	}
)

var (
	localPoints IPoints
)

func Points() IPoints {
	if localPoints == nil {
		panic("implement not found for interface IPoints, forgot register?")
	}
	return localPoints
}

func RegisterPoints(i IPoints) {
	localPoints = i
}
