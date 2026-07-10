package backend

import (
	"context"

	"bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Goods = cGoods{}

type cGoods struct{}

func (c *cGoods) Create(ctx context.Context, req *backend.GoodsCreateReq) (res *backend.GoodsCreateRes, err error) {
	var input model.GoodsCreateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	input.OperatorAdminId = currentAdminId(ctx)
	out, err := service.Goods().Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backend.GoodsCreateRes{Id: out.Id}, nil
}

func (c *cGoods) List(ctx context.Context, req *backend.GoodsListReq) (res *backend.GoodsListRes, err error) {
	var input model.GoodsListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Goods().List(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backend.GoodsListRes{
		List:  toApiGoodsItems(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cGoods) Detail(ctx context.Context, req *backend.GoodsDetailReq) (res *backend.GoodsDetailRes, err error) {
	out, err := service.Goods().Detail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &backend.GoodsDetailRes{Goods: toApiGoodsItem(out.Goods)}, nil
}

func (c *cGoods) Update(ctx context.Context, req *backend.GoodsUpdateReq) (res *backend.GoodsUpdateRes, err error) {
	var input model.GoodsUpdateInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	input.OperatorAdminId = currentAdminId(ctx)
	out, err := service.Goods().Update(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backend.GoodsUpdateRes{Goods: toApiGoodsItem(out.Goods)}, nil
}

func (c *cGoods) Status(ctx context.Context, req *backend.GoodsStatusReq) (res *backend.GoodsStatusRes, err error) {
	err = service.Goods().UpdateStatus(ctx, model.GoodsStatusInput{
		Id:     req.Id,
		Status: gconv.Int(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return &backend.GoodsStatusRes{Message: "Status updated successfully"}, nil
}

func toApiGoodsItems(in []model.GoodsItem) []backend.GoodsItem {
	list := make([]backend.GoodsItem, 0, len(in))
	for _, item := range in {
		list = append(list, toApiGoodsItem(item))
	}
	return list
}

func toApiGoodsItem(in model.GoodsItem) backend.GoodsItem {
	return backend.GoodsItem{
		Id:           in.Id,
		CategoryId:   in.CategoryId,
		CategoryName: in.CategoryName,
		Name:         in.Name,
		ImageUrl:     in.ImageUrl,
		PointsPrice:  in.PointsPrice,
		Stock:        in.Stock,
		Description:  in.Description,
		Status:       in.Status,
	}
}
