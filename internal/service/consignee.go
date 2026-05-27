// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"bit303_shop/internal/model"

	"golang.org/x/net/context"
)

type (
	IConsignee interface {
		// GetList 查询内容列表
		GetList(ctx context.Context, in model.ConsigneeGetListInput) (out *model.ConsigneeGetListOutput, err error)
		// Add 添加收货地址
		Add(ctx context.Context, in model.AddConsigneeInput) (out *model.AddConsigneeOutput, err error)
		// Update 更新收货地址
		Update(ctx context.Context, in model.UpdateConsigneeInput) (out *model.UpdateConsigneeOutput, err error)
		// Delete 删除收货地址
		Delete(ctx context.Context, in model.DeleteConsigneeInput) (out *model.DeleteConsigneeOutput, err error)
		// AdminGetList 后台管理列表（包含用户信息）
		AdminGetList(ctx context.Context, page, size int) (out *model.ConsigneeAdminListOutput, err error)
		// AdminDelete 后台管理删除
		AdminDelete(ctx context.Context, id uint) error
	}
)

var (
	localConsignee IConsignee
)

func Consignee() IConsignee {
	if localConsignee == nil {
		panic("implement not found for interface IConsignee, forgot register?")
	}
	return localConsignee
}

func RegisterConsignee(i IConsignee) {
	localConsignee = i
}
