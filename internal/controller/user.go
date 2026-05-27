package controller

import (
	"context"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var User = cUser{}

type cUser struct{}

func (c *cUser) Register(ctx context.Context, req *frontend.RegisterReq) (res *frontend.RegisterRes, err error) {
	data := model.RegisterInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	out, err := service.User().Register(ctx, data)
	if err != nil {
		return nil, err
	}
	return &frontend.RegisterRes{Id: out.Id}, nil
}

// 前台用户登录（JWT）
func (c *cUser) Login(ctx context.Context, req *frontend.LoginReq) (res *frontend.LoginRes, err error) {
	res = &frontend.LoginRes{}
	res.Token, res.Expire = service.UserAuth().LoginHandler(ctx)
	return
}

func (c *cUser) Info(ctx context.Context, req *frontend.UserInfoReq) (res *frontend.UserInfoRes, err error) {
	res = &frontend.UserInfoRes{}
	// 从 JWT 中获取当前用户 ID
	g.Log().Info(ctx, "Info 当前用户ID:", service.UserAuth().GetIdentity(ctx))
	userId := gconv.Uint(service.UserAuth().GetIdentity(ctx))
	if userId == 0 {
		return res, nil
	}
	var user entity.UserInfo
	if err = dao.UserInfo.Ctx(ctx).WherePri(userId).Scan(&user); err != nil {
		return nil, err
	}
	res.Id = uint(user.Id)
	res.Name = user.Name
	res.Avatar = user.Avatar
	res.Sex = uint8(user.Sex)
	res.Sign = user.Sign
	res.Status = uint8(user.Status)
	return res, nil
}

func (*cUser) UpdatePassword(ctx context.Context, req *frontend.UpdatePasswordReq) (res *frontend.UpdatePasswordRes, err error) {
	data := model.UpdatePasswordInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	out, err := service.User().UpdatePassword(ctx, data)
	if err != nil {
		return nil, err
	}
	return &frontend.UpdatePasswordRes{Id: out.Id}, nil
}

// 前台用户登出
func (c *cUser) Logout(ctx context.Context, req *frontend.LogoutReq) (res *frontend.LogoutRes, err error) {
	// 简单返回登出成功消息，依赖客户端删除token
	return &frontend.LogoutRes{
		Message: "登出成功",
	}, nil
}
