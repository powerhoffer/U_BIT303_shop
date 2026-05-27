package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

var User = cUser{}

type cUser struct{}

// List 用户列表
func (c *cUser) List(ctx context.Context, req *backend.UserGetListReq) (res *backend.UserGetListRes, err error) {
	out, err := service.User().GetList(ctx, model.UserGetListInput{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}

	// 转换输出
	list := make([]backend.UserListItem, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, backend.UserListItem{
			Id:     item.Id,
			Name:   item.Name,
			Avatar: item.Avatar,
			Sex:    item.Sex,
			Status: item.Status,
			Sign:   item.Sign,
		})
	}

	return &backend.UserGetListRes{
		List:  list,
		Page:  out.Page,
		Size:  out.Size,
		Total: out.Total,
	}, nil
}

// UpdateStatus 冻结/解冻用户
func (c *cUser) UpdateStatus(ctx context.Context, req *backend.UserUpdateStatusReq) (res *backend.UserUpdateStatusRes, err error) {
	err = service.User().UpdateStatus(ctx, req.Id, req.Status)
	if err != nil {
		return nil, err
	}
	return &backend.UserUpdateStatusRes{}, nil
}

// Delete 删除用户
func (c *cUser) Delete(ctx context.Context, req *backend.UserDeleteReq) (res *backend.UserDeleteRes, err error) {
	err = service.User().Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &backend.UserDeleteRes{}, nil
}

// PromoteToAdmin 提升用户为管理员
func (c *cUser) PromoteToAdmin(ctx context.Context, req *backend.UserPromoteToAdminReq) (res *backend.UserPromoteToAdminRes, err error) {
	// 获取用户信息
	userInfo, err := service.User().GetById(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, gerror.New("用户不存在")
	}

	// 创建管理员账号（密码由 Admin.Create 内部加密）
	out, err := service.Admin().Create(ctx, &model.AdminCreateInput{
		AdminCreateUpdateBase: model.AdminCreateUpdateBase{
			Name:     userInfo.Name,
			Password: req.Password,
			RoleIds:  req.RoleIds,
			IsAdmin:  0, // 非超管
		},
	})
	if err != nil {
		return nil, err
	}

	return &backend.UserPromoteToAdminRes{
		AdminId: out.AdminId,
	}, nil
}
