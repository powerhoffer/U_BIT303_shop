package frontend

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

var Consignee = cConsignee{}

type cConsignee struct{}

func (c *cConsignee) List(ctx context.Context, req *frontend.ConsigneeListReq) (res *frontend.ConsigneeListRes, err error) {
	getListRes, err := service.Consignee().GetList(ctx, model.ConsigneeGetListInput{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}

	return &frontend.ConsigneeListRes{
		List:  getListRes.List,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total}, nil
}

func (c *cConsignee) Add(ctx context.Context, req *frontend.AddConsigneeReq) (res *frontend.AddConsigneeRes, err error) {
	data := model.AddConsigneeInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	
	// 获取当前用户ID
	data.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))
	
	out, err := service.Consignee().Add(ctx, data)
	if err != nil {
		return nil, err
	}
	return &frontend.AddConsigneeRes{Id: out.Id}, nil
}

func (c *cConsignee) Update(ctx context.Context, req *frontend.UpdateConsigneeReq) (res *frontend.UpdateConsigneeRes, err error) {
	data := model.UpdateConsigneeInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	
	// 获取当前用户ID
	data.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))
	
	out, err := service.Consignee().Update(ctx, data)
	if err != nil {
		return nil, err
	}
	return &frontend.UpdateConsigneeRes{Id: out.Id}, nil
}

func (c *cConsignee) Delete(ctx context.Context, req *frontend.DeleteConsigneeReq) (res *frontend.DeleteConsigneeRes, err error) {
	out, err := service.Consignee().Delete(ctx, model.DeleteConsigneeInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &frontend.DeleteConsigneeRes{Id: out.Id}, nil
}