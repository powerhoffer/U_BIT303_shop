package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Employee = cEmployee{}

type cEmployee struct{}

func (c *cEmployee) Register(ctx context.Context, req *backendApi.EmployeeRegisterReq) (res *backendApi.EmployeeRegisterRes, err error) {
	var input model.EmployeeRegisterInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().Register(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeRegisterRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) Login(ctx context.Context, req *backendApi.EmployeeLoginReq) (res *backendApi.EmployeeLoginRes, err error) {
	var input model.EmployeeLoginInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().Login(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeLoginRes{
		Token:    out.Token,
		ExpireAt: out.ExpireAt,
		Employee: toApiEmployee(out.Employee),
	}, nil
}

func (c *cEmployee) Logout(ctx context.Context, req *backendApi.EmployeeLogoutReq) (res *backendApi.EmployeeLogoutRes, err error) {
	return &backendApi.EmployeeLogoutRes{Message: "Logged out successfully"}, nil
}

func (c *cEmployee) Info(ctx context.Context, req *backendApi.EmployeeInfoReq) (res *backendApi.EmployeeInfoRes, err error) {
	out, err := service.Employee().Info(ctx, currentEmployeeId(ctx))
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeInfoRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) UpdatePassword(ctx context.Context, req *backendApi.EmployeeUpdatePasswordReq) (res *backendApi.EmployeeUpdatePasswordRes, err error) {
	err = service.Employee().UpdatePassword(ctx, model.EmployeeUpdatePasswordInput{
		EmployeeId:  currentEmployeeId(ctx),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeUpdatePasswordRes{Message: "Password changed successfully"}, nil
}

func (c *cEmployee) ManageCreate(ctx context.Context, req *backendApi.EmployeeManageCreateReq) (res *backendApi.EmployeeManageCreateRes, err error) {
	var input model.EmployeeManageCreateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().ManageCreate(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeManageCreateRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) ManageList(ctx context.Context, req *backendApi.EmployeeManageListReq) (res *backendApi.EmployeeManageListRes, err error) {
	var input model.EmployeeManageListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Employee().ManageList(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeManageListRes{
		List:  toApiEmployees(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cEmployee) ManageDetail(ctx context.Context, req *backendApi.EmployeeManageDetailReq) (res *backendApi.EmployeeManageDetailRes, err error) {
	out, err := service.Employee().ManageDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeManageDetailRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) ManageUpdate(ctx context.Context, req *backendApi.EmployeeManageUpdateReq) (res *backendApi.EmployeeManageUpdateRes, err error) {
	out, err := service.Employee().ManageUpdate(ctx, model.EmployeeManageUpdateInput{
		Id:       req.Id,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeManageUpdateRes{Employee: toApiEmployee(out.Employee)}, nil
}

func (c *cEmployee) ManageStatus(ctx context.Context, req *backendApi.EmployeeManageStatusReq) (res *backendApi.EmployeeManageStatusRes, err error) {
	err = service.Employee().ManageUpdateStatus(ctx, model.EmployeeManageStatusInput{
		Id:     req.Id,
		Status: gconv.Int(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeManageStatusRes{Message: "Status updated successfully"}, nil
}

func (c *cEmployee) ManageResetPassword(ctx context.Context, req *backendApi.EmployeeManageResetPasswordReq) (res *backendApi.EmployeeManageResetPasswordRes, err error) {
	err = service.Employee().ManageResetPassword(ctx, model.EmployeeManageResetPasswordInput{
		Id:       req.Id,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.EmployeeManageResetPasswordRes{Message: "Password reset successfully"}, nil
}

func toApiEmployees(in []model.EmployeeBase) []backendApi.EmployeeBase {
	list := make([]backendApi.EmployeeBase, 0, len(in))
	for _, item := range in {
		list = append(list, toApiEmployee(item))
	}
	return list
}

func toApiEmployee(in model.EmployeeBase) backendApi.EmployeeBase {
	return backendApi.EmployeeBase{
		Id:       in.Id,
		Username: in.Username,
		RealName: in.RealName,
		Phone:    in.Phone,
		Email:    in.Email,
		Status:   in.Status,
	}
}
