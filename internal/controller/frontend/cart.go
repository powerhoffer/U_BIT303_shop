package frontend

import (
	"context"

	frontendApi "bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Cart = cCart{}

type cCart struct{}

func (c *cCart) Add(ctx context.Context, req *frontendApi.CartAddReq) (res *frontendApi.CartAddRes, err error) {
	out, err := service.Cart().Add(ctx, model.CartAddInput{
		EmployeeId: currentEmployeeId(ctx),
		GoodsId:    req.GoodsId,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &frontendApi.CartAddRes{Id: out.Id}, nil
}

func (c *cCart) List(ctx context.Context, req *frontendApi.CartListReq) (res *frontendApi.CartListRes, err error) {
	var input model.CartListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	input.EmployeeId = currentEmployeeId(ctx)
	out, err := service.Cart().List(ctx, input)
	if err != nil {
		return nil, err
	}
	return &frontendApi.CartListRes{
		List:  toFrontendApiCartItems(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cCart) Update(ctx context.Context, req *frontendApi.CartUpdateReq) (res *frontendApi.CartUpdateRes, err error) {
	out, err := service.Cart().Update(ctx, model.CartUpdateInput{
		EmployeeId: currentEmployeeId(ctx),
		Id:         req.Id,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &frontendApi.CartUpdateRes{Id: out.Id}, nil
}

func (c *cCart) Remove(ctx context.Context, req *frontendApi.CartRemoveReq) (res *frontendApi.CartRemoveRes, err error) {
	out, err := service.Cart().Remove(ctx, model.CartRemoveInput{
		EmployeeId: currentEmployeeId(ctx),
		Id:         req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &frontendApi.CartRemoveRes{Id: out.Id}, nil
}

func toFrontendApiCartItems(in []model.CartItem) []frontendApi.CartItem {
	list := make([]frontendApi.CartItem, 0, len(in))
	for _, item := range in {
		list = append(list, frontendApi.CartItem{
			Id:           item.Id,
			GoodsId:      item.GoodsId,
			CategoryId:   item.CategoryId,
			CategoryName: item.CategoryName,
			GoodsName:    item.GoodsName,
			ImageUrl:     item.ImageUrl,
			PointsPrice:  item.PointsPrice,
			Stock:        item.Stock,
			Count:        item.Count,
			TotalPoints:  item.TotalPoints,
		})
	}
	return list
}
