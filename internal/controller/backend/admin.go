package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Admin = cAdmin{}

type cAdmin struct{}

func (c *cAdmin) Login(ctx context.Context, req *backendApi.AdminLoginReq) (res *backendApi.AdminLoginRes, err error) {
	var input model.AdminLoginInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Admin().Login(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminLoginRes{Token: out.Token, ExpireAt: out.ExpireAt, Admin: toApiAdmin(out.Admin)}, nil
}

func (c *cAdmin) Logout(ctx context.Context, req *backendApi.AdminLogoutReq) (res *backendApi.AdminLogoutRes, err error) {
	return &backendApi.AdminLogoutRes{Message: "Logged out successfully"}, nil
}

func (c *cAdmin) Info(ctx context.Context, req *backendApi.AdminInfoReq) (res *backendApi.AdminInfoRes, err error) {
	out, err := service.Admin().Info(ctx, currentAdminId(ctx))
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminInfoRes{Admin: toApiAdmin(out.Admin)}, nil
}

func (c *cAdmin) UpdatePassword(ctx context.Context, req *backendApi.AdminUpdatePasswordReq) (res *backendApi.AdminUpdatePasswordRes, err error) {
	err = service.Admin().UpdatePassword(ctx, model.AdminUpdatePasswordInput{AdminId: currentAdminId(ctx), OldPassword: req.OldPassword, NewPassword: req.NewPassword})
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminUpdatePasswordRes{Message: "Password changed successfully"}, nil
}

func (c *cAdmin) ManageCreate(ctx context.Context, req *backendApi.AdminManageCreateReq) (res *backendApi.AdminManageCreateRes, err error) {
	var input model.AdminManageCreateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Admin().ManageCreate(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageCreateRes{Admin: toApiAdmin(out.Admin)}, nil
}

func (c *cAdmin) ManageList(ctx context.Context, req *backendApi.AdminManageListReq) (res *backendApi.AdminManageListRes, err error) {
	var input model.AdminManageListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Admin().ManageList(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageListRes{List: toApiAdmins(out.List), Total: out.Total, Page: out.Page, Size: out.Size}, nil
}

func (c *cAdmin) ManageDetail(ctx context.Context, req *backendApi.AdminManageDetailReq) (res *backendApi.AdminManageDetailRes, err error) {
	out, err := service.Admin().ManageDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageDetailRes{Admin: toApiAdmin(out.Admin)}, nil
}

func (c *cAdmin) ManageUpdate(ctx context.Context, req *backendApi.AdminManageUpdateReq) (res *backendApi.AdminManageUpdateRes, err error) {
	var input model.AdminManageUpdateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Admin().ManageUpdate(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageUpdateRes{Admin: toApiAdmin(out.Admin)}, nil
}

func (c *cAdmin) ManageStatus(ctx context.Context, req *backendApi.AdminManageStatusReq) (res *backendApi.AdminManageStatusRes, err error) {
	err = service.Admin().ManageUpdateStatus(ctx, model.AdminManageStatusInput{Id: req.Id, Status: req.Status})
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageStatusRes{Message: "Status updated successfully"}, nil
}

func (c *cAdmin) ManageResetPassword(ctx context.Context, req *backendApi.AdminManageResetPasswordReq) (res *backendApi.AdminManageResetPasswordRes, err error) {
	err = service.Admin().ManageResetPassword(ctx, model.AdminManageResetPasswordInput{Id: req.Id, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageResetPasswordRes{Message: "Password reset successfully"}, nil
}

func (c *cAdmin) ManageRoles(ctx context.Context, req *backendApi.AdminManageRolesReq) (res *backendApi.AdminManageRolesRes, err error) {
	err = service.Admin().ManageAssignRoles(ctx, model.AdminManageRolesInput{Id: req.Id, RoleIds: req.RoleIds})
	if err != nil {
		return nil, err
	}
	return &backendApi.AdminManageRolesRes{Message: "Roles updated successfully"}, nil
}

func toApiAdmins(in []model.AdminBase) []backendApi.AdminBase {
	list := make([]backendApi.AdminBase, 0, len(in))
	for _, item := range in {
		list = append(list, toApiAdmin(item))
	}
	return list
}

func toApiAdmin(in model.AdminBase) backendApi.AdminBase {
	return backendApi.AdminBase{Id: in.Id, Username: in.Username, RealName: in.RealName, Phone: in.Phone, Email: in.Email, Status: in.Status, IsSuper: in.IsSuper, RoleIds: in.RoleIds}
}
