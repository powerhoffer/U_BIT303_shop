// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"bit303_shop/internal/model"
)

type (
	IAdmin interface {
		Create(ctx context.Context, in *model.AdminCreateInput) (out model.AdminCreateOutput, err error)
		GetUserByUserNamePassword(ctx context.Context, in model.UserLoginInput) map[string]interface{}
		// Delete 删除
		Delete(ctx context.Context, id uint) error
		// Update 修改
		Update(ctx context.Context, in model.AdminUpdateInput) error
		// GetList 查询内容列表
		GetList(ctx context.Context, in model.AdminGetListInput) (out *model.AdminGetListOutput, err error)
		GetAdminByNamePassword(ctx context.Context, in model.UserLoginInput) map[string]interface{}
		// GetAdminByNamePasswordWithRoles 登录验证并返回包含角色信息的数据
		GetAdminByNamePasswordWithRoles(ctx context.Context, in model.UserLoginInput) map[string]interface{}
		// GetById 根据ID获取管理员信息
		GetById(ctx context.Context, id int) (*model.AdminInfo, error)
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
