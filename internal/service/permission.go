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
	IPermission interface {
		Create(ctx context.Context, in model.PermissionCreateInput) (out model.PermissionCreateOutput, err error)
		List(ctx context.Context, in model.PermissionListInput) (out model.PermissionListOutput, err error)
		Detail(ctx context.Context, id uint) (out model.PermissionDetailOutput, err error)
		Update(ctx context.Context, in model.PermissionUpdateInput) (out model.PermissionUpdateOutput, err error)
		UpdateStatus(ctx context.Context, in model.PermissionStatusInput) error
	}
)

var (
	localPermission IPermission
)

func Permission() IPermission {
	if localPermission == nil {
		panic("implement not found for interface IPermission, forgot register?")
	}
	return localPermission
}

func RegisterPermission(i IPermission) {
	localPermission = i
}
