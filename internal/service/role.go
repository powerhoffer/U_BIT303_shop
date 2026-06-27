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
	IRole interface {
		Create(ctx context.Context, in model.RoleCreateInput) (out model.RoleCreateOutput, err error)
		List(ctx context.Context, in model.RoleListInput) (out model.RoleListOutput, err error)
		Detail(ctx context.Context, id uint) (out model.RoleDetailOutput, err error)
		Update(ctx context.Context, in model.RoleUpdateInput) (out model.RoleUpdateOutput, err error)
		UpdateStatus(ctx context.Context, in model.RoleStatusInput) error
		AssignPermissions(ctx context.Context, in model.RolePermissionsInput) error
	}
)

var (
	localRole IRole
)

func Role() IRole {
	if localRole == nil {
		panic("implement not found for interface IRole, forgot register?")
	}
	return localRole
}

func RegisterRole(i IRole) {
	localRole = i
}
