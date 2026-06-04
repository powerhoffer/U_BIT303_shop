package controller

import (
	"context"

	"bit303_shop/api/employee"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var Employee = cEmployee{}

type cEmployee struct{}

func (c *cEmployee) Register(ctx context.Context, req *employee.RegisterReq) (res *employee.RegisterRes, err error) {
	var input model.EmployeeRegisterInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().Register(ctx, input)
	if err != nil {
		return nil, err
	}
	return &employee.RegisterRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) Login(ctx context.Context, req *employee.LoginReq) (res *employee.LoginRes, err error) {
	var input model.EmployeeLoginInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().Login(ctx, input)
	if err != nil {
		return nil, err
	}
	return &employee.LoginRes{
		Token:    out.Token,
		ExpireAt: out.ExpireAt,
		Employee: toApiEmployee(out.Employee),
	}, nil
}

func (c *cEmployee) Logout(ctx context.Context, req *employee.LogoutReq) (res *employee.LogoutRes, err error) {
	return &employee.LogoutRes{Message: "登出成功"}, nil
}

func (c *cEmployee) Info(ctx context.Context, req *employee.InfoReq) (res *employee.InfoRes, err error) {
	out, err := service.Employee().Info(ctx, currentEmployeeId(ctx))
	if err != nil {
		return nil, err
	}
	return &employee.InfoRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) UpdatePassword(ctx context.Context, req *employee.UpdatePasswordReq) (res *employee.UpdatePasswordRes, err error) {
	err = service.Employee().UpdatePassword(ctx, model.EmployeeUpdatePasswordInput{
		EmployeeId:  currentEmployeeId(ctx),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &employee.UpdatePasswordRes{Message: "修改密码成功"}, nil
}

func currentEmployeeId(ctx context.Context) uint {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar(consts.CtxEmployeeId).Uint()
}

func toApiEmployee(in model.EmployeeBase) employee.EmployeeBase {
	return employee.EmployeeBase{
		Id:       in.Id,
		Username: in.Username,
		RealName: in.RealName,
		Phone:    in.Phone,
		Email:    in.Email,
		Status:   in.Status,
	}
}
