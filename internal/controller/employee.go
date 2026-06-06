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

func (c *cEmployee) ManageCreate(ctx context.Context, req *employee.ManageCreateReq) (res *employee.ManageCreateRes, err error) {
	var input model.EmployeeManageCreateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().ManageCreate(ctx, input)
	if err != nil {
		return nil, err
	}
	return &employee.ManageCreateRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) ManageList(ctx context.Context, req *employee.ManageListReq) (res *employee.ManageListRes, err error) {
	var input model.EmployeeManageListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().ManageList(ctx, input)
	if err != nil {
		return nil, err
	}
	return &employee.ManageListRes{
		List:  toApiEmployees(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cEmployee) ManageDetail(ctx context.Context, req *employee.ManageDetailReq) (res *employee.ManageDetailRes, err error) {
	out, err := service.Employee().ManageDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &employee.ManageDetailRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) ManageUpdate(ctx context.Context, req *employee.ManageUpdateReq) (res *employee.ManageUpdateRes, err error) {
	out, err := service.Employee().ManageUpdate(ctx, model.EmployeeManageUpdateInput{
		Id:       req.Id,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &employee.ManageUpdateRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) ManageStatus(ctx context.Context, req *employee.ManageStatusReq) (res *employee.ManageStatusRes, err error) {
	err = service.Employee().ManageUpdateStatus(ctx, model.EmployeeManageStatusInput{
		Id:     req.Id,
		Status: gconv.Int(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return &employee.ManageStatusRes{Message: "状态更新成功"}, nil
}

func (c *cEmployee) ManageResetPassword(ctx context.Context, req *employee.ManageResetPasswordReq) (res *employee.ManageResetPasswordRes, err error) {
	err = service.Employee().ManageResetPassword(ctx, model.EmployeeManageResetPasswordInput{
		Id:       req.Id,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &employee.ManageResetPasswordRes{Message: "密码重置成功"}, nil
}

func currentEmployeeId(ctx context.Context) uint {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar(consts.CtxEmployeeId).Uint()
}

func toApiEmployees(in []model.EmployeeBase) []employee.EmployeeBase {
	list := make([]employee.EmployeeBase, 0, len(in))
	for _, item := range in {
		list = append(list, toApiEmployee(item))
	}
	return list
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
