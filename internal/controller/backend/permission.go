package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Permission = cPermission{}

type cPermission struct{}

func (c *cPermission) Create(ctx context.Context, req *backendApi.PermissionCreateReq) (res *backendApi.PermissionCreateRes, err error) {
	var input model.PermissionCreateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Permission().Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.PermissionCreateRes{Permission: toApiPermission(out.Permission)}, nil
}

func (c *cPermission) List(ctx context.Context, req *backendApi.PermissionListReq) (res *backendApi.PermissionListRes, err error) {
	var input model.PermissionListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Permission().List(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.PermissionListRes{List: toApiPermissions(out.List), Total: out.Total, Page: out.Page, Size: out.Size}, nil
}

func (c *cPermission) Detail(ctx context.Context, req *backendApi.PermissionDetailReq) (res *backendApi.PermissionDetailRes, err error) {
	out, err := service.Permission().Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &backendApi.PermissionDetailRes{Permission: toApiPermission(out.Permission)}, nil
}

func (c *cPermission) Update(ctx context.Context, req *backendApi.PermissionUpdateReq) (res *backendApi.PermissionUpdateRes, err error) {
	var input model.PermissionUpdateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Permission().Update(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.PermissionUpdateRes{Permission: toApiPermission(out.Permission)}, nil
}

func (c *cPermission) Status(ctx context.Context, req *backendApi.PermissionStatusReq) (res *backendApi.PermissionStatusRes, err error) {
	err = service.Permission().UpdateStatus(ctx, model.PermissionStatusInput{Id: req.Id, Status: req.Status})
	if err != nil {
		return nil, err
	}
	return &backendApi.PermissionStatusRes{Message: "Status updated successfully"}, nil
}

func toApiPermissions(in []model.PermissionBase) []backendApi.PermissionBase {
	list := make([]backendApi.PermissionBase, 0, len(in))
	for _, item := range in {
		list = append(list, toApiPermission(item))
	}
	return list
}

func toApiPermission(in model.PermissionBase) backendApi.PermissionBase {
	return backendApi.PermissionBase{Id: in.Id, Name: in.Name, GroupName: in.GroupName, Method: in.Method, Path: in.Path, Status: in.Status}
}
