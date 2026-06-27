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
	IAdmin interface {
		Login(ctx context.Context, in model.AdminLoginInput) (out model.AdminLoginOutput, err error)
		Info(ctx context.Context, adminId uint) (out model.AdminInfoOutput, err error)
		UpdatePassword(ctx context.Context, in model.AdminUpdatePasswordInput) error
		ManageCreate(ctx context.Context, in model.AdminManageCreateInput) (out model.AdminManageCreateOutput, err error)
		ManageList(ctx context.Context, in model.AdminManageListInput) (out model.AdminManageListOutput, err error)
		ManageDetail(ctx context.Context, id uint) (out model.AdminManageDetailOutput, err error)
		ManageUpdate(ctx context.Context, in model.AdminManageUpdateInput) (out model.AdminManageUpdateOutput, err error)
		ManageUpdateStatus(ctx context.Context, in model.AdminManageStatusInput) error
		ManageResetPassword(ctx context.Context, in model.AdminManageResetPasswordInput) error
		ManageAssignRoles(ctx context.Context, in model.AdminManageRolesInput) error
	}
)

var (
	localAdmin IAdmin
)

func Admin() IAdmin {
	if localAdmin == nil {
		panic("implement not found for interface IAdmin, forgot register?")
	}
	return localAdmin
}

func RegisterAdmin(i IAdmin) {
	localAdmin = i
}
