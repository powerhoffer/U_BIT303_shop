package backend

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// Consignee 后台收货地址管理
var Consignee = cConsignee{}

type cConsignee struct{}

// List 收货地址列表
func (c *cConsignee) List(ctx context.Context, req *backend.ConsigneeListReq) (res *backend.ConsigneeListRes, err error) {
	out, err := service.Consignee().AdminGetList(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &backend.ConsigneeListRes{
		List:  out.List,
		Page:  out.Page,
		Size:  out.Size,
		Total: out.Total,
	}, nil
}

// Add 添加收货地址
func (c *cConsignee) Add(ctx context.Context, req *backend.ConsigneeAddReq) (res *backend.ConsigneeAddRes, err error) {
	data := model.AddConsigneeInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	data.UserId = req.UserId

	out, err := service.Consignee().Add(ctx, data)
	if err != nil {
		return nil, err
	}
	return &backend.ConsigneeAddRes{Id: out.Id}, nil
}

// Update 更新收货地址
func (c *cConsignee) Update(ctx context.Context, req *backend.ConsigneeUpdateReq) (res *backend.ConsigneeUpdateRes, err error) {
	// 直接从请求中获取 consignee_id，避免被 JWT 的 id 覆盖
	r := g.RequestFromCtx(ctx)
	consigneeId := r.Get("consignee_id").Uint()

	data := model.UpdateConsigneeInput{
		Id:        consigneeId,
		UserId:    req.UserId,
		IsDefault: req.IsDefault,
		Name:      req.Name,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		Town:      req.Town,
		Street:    req.Street,
		Detail:    req.Detail,
	}

	out, err := service.Consignee().Update(ctx, data)
	if err != nil {
		return nil, err
	}
	return &backend.ConsigneeUpdateRes{Id: out.Id}, nil
}

// Delete 删除收货地址
func (c *cConsignee) Delete(ctx context.Context, req *backend.ConsigneeDeleteReq) (res *backend.ConsigneeDeleteRes, err error) {
	// 直接从请求中获取 consignee_id，避免被 JWT 的 id 覆盖
	r := g.RequestFromCtx(ctx)
	consigneeId := r.Get("consignee_id").Uint()

	err = service.Consignee().AdminDelete(ctx, consigneeId)
	return &backend.ConsigneeDeleteRes{}, err
}
