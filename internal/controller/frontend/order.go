package frontend

import (
	"context"

	frontendApi "bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Order = cOrder{}

type cOrder struct{}

func (c *cOrder) Create(ctx context.Context, req *frontendApi.OrderCreateReq) (res *frontendApi.OrderCreateRes, err error) {
	out, err := service.Order().Create(ctx, model.OrderCreateInput{
		EmployeeId: currentEmployeeId(ctx),
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &frontendApi.OrderCreateRes{Order: toFrontendApiOrderItem(out.Order)}, nil
}

func (c *cOrder) List(ctx context.Context, req *frontendApi.OrderListReq) (res *frontendApi.OrderListRes, err error) {
	var input model.OrderListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	input.EmployeeId = currentEmployeeId(ctx)
	out, err := service.Order().List(ctx, input)
	if err != nil {
		return nil, err
	}
	return &frontendApi.OrderListRes{
		List:  toFrontendApiOrderItems(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cOrder) Detail(ctx context.Context, req *frontendApi.OrderDetailReq) (res *frontendApi.OrderDetailRes, err error) {
	out, err := service.Order().Detail(ctx, model.OrderDetailInput{
		EmployeeId: currentEmployeeId(ctx),
		Id:         req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &frontendApi.OrderDetailRes{Order: toFrontendApiOrderDetail(out.Order)}, nil
}

func (c *cOrder) Cancel(ctx context.Context, req *frontendApi.OrderCancelReq) (res *frontendApi.OrderCancelRes, err error) {
	out, err := service.Order().Cancel(ctx, model.OrderCancelInput{
		EmployeeId: currentEmployeeId(ctx),
		Id:         req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &frontendApi.OrderCancelRes{Order: toFrontendApiOrderItem(out.Order)}, nil
}

func toFrontendApiOrderItems(in []model.OrderBase) []frontendApi.OrderItem {
	list := make([]frontendApi.OrderItem, 0, len(in))
	for _, item := range in {
		list = append(list, toFrontendApiOrderItem(item))
	}
	return list
}

func toFrontendApiOrderItem(in model.OrderBase) frontendApi.OrderItem {
	return frontendApi.OrderItem{
		Id:          in.Id,
		OrderNo:     in.OrderNo,
		EmployeeId:  in.EmployeeId,
		TotalPoints: in.TotalPoints,
		Status:      in.Status,
		Remark:      in.Remark,
		CreatedAt:   in.CreatedAt,
	}
}

func toFrontendApiOrderDetail(in model.OrderDetail) frontendApi.OrderDetail {
	items := make([]frontendApi.OrderGoodsItem, 0, len(in.Items))
	for _, item := range in.Items {
		items = append(items, frontendApi.OrderGoodsItem{
			Id:            item.Id,
			OrderId:       item.OrderId,
			GoodsId:       item.GoodsId,
			GoodsName:     item.GoodsName,
			GoodsImageUrl: item.GoodsImageUrl,
			PointsPrice:   item.PointsPrice,
			Count:         item.Count,
			TotalPoints:   item.TotalPoints,
			CreatedAt:     item.CreatedAt,
		})
	}
	return frontendApi.OrderDetail{
		OrderItem: toFrontendApiOrderItem(in.OrderBase),
		Items:     items,
	}
}
