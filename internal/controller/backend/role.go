package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Role = cRole{}

type cRole struct{}

func (c *cRole) Create(ctx context.Context, req *backendApi.RoleCreateReq) (res *backendApi.RoleCreateRes, err error) {
	out, err := service.Role().Create(ctx, model.RoleCreateInput{Name: req.Name, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &backendApi.RoleCreateRes{Role: toApiRole(out.Role)}, nil
}

func (c *cRole) List(ctx context.Context, req *backendApi.RoleListReq) (res *backendApi.RoleListRes, err error) {
	var input model.RoleListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Role().List(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.RoleListRes{List: toApiRoles(out.List), Total: out.Total, Page: out.Page, Size: out.Size}, nil
}

func (c *cRole) Detail(ctx context.Context, req *backendApi.RoleDetailReq) (res *backendApi.RoleDetailRes, err error) {
	out, err := service.Role().Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &backendApi.RoleDetailRes{Role: toApiRole(out.Role)}, nil
}

func (c *cRole) Update(ctx context.Context, req *backendApi.RoleUpdateReq) (res *backendApi.RoleUpdateRes, err error) {
	out, err := service.Role().Update(ctx, model.RoleUpdateInput{Id: req.Id, Name: req.Name, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &backendApi.RoleUpdateRes{Role: toApiRole(out.Role)}, nil
}

func (c *cRole) Status(ctx context.Context, req *backendApi.RoleStatusReq) (res *backendApi.RoleStatusRes, err error) {
	err = service.Role().UpdateStatus(ctx, model.RoleStatusInput{Id: req.Id, Status: req.Status})
	if err != nil {
		return nil, err
	}
	return &backendApi.RoleStatusRes{Message: "Status updated successfully"}, nil
}

func (c *cRole) Permissions(ctx context.Context, req *backendApi.RolePermissionsReq) (res *backendApi.RolePermissionsRes, err error) {
	err = service.Role().AssignPermissions(ctx, model.RolePermissionsInput{Id: req.Id, PermissionIds: req.PermissionIds})
	if err != nil {
		return nil, err
	}
	return &backendApi.RolePermissionsRes{Message: "Permissions updated successfully"}, nil
}

func toApiRoles(in []model.RoleBase) []backendApi.RoleBase {
	list := make([]backendApi.RoleBase, 0, len(in))
	for _, item := range in {
		list = append(list, toApiRole(item))
	}
	return list
}

func toApiRole(in model.RoleBase) backendApi.RoleBase {
	return backendApi.RoleBase{Id: in.Id, Name: in.Name, Description: in.Description, Status: in.Status, PermissionIds: in.PermissionIds}
}
