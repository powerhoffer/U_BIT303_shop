package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Order = cOrder{}

type cOrder struct{}

func (c *cOrder) List(ctx context.Context, req *backendApi.OrderListReq) (res *backendApi.OrderListRes, err error) {
	var input model.BackendOrderListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Order().ManageList(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.OrderListRes{
		List:  toBackendApiOrderItems(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cOrder) Detail(ctx context.Context, req *backendApi.OrderDetailReq) (res *backendApi.OrderDetailRes, err error) {
	out, err := service.Order().ManageDetail(ctx, model.BackendOrderDetailInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &backendApi.OrderDetailRes{Order: toBackendApiOrderDetail(out.Order)}, nil
}

func (c *cOrder) Complete(ctx context.Context, req *backendApi.OrderCompleteReq) (res *backendApi.OrderCompleteRes, err error) {
	out, err := service.Order().ManageComplete(ctx, model.BackendOrderCompleteInput{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &backendApi.OrderCompleteRes{Order: toBackendApiOrderItem(out.Order)}, nil
}

func (c *cOrder) Cancel(ctx context.Context, req *backendApi.OrderCancelReq) (res *backendApi.OrderCancelRes, err error) {
	out, err := service.Order().ManageCancel(ctx, model.BackendOrderCancelInput{
		Id:                 req.Id,
		OperatorEmployeeId: currentEmployeeId(ctx),
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.OrderCancelRes{Order: toBackendApiOrderItem(out.Order)}, nil
}

func toBackendApiOrderItems(in []model.OrderBase) []backendApi.OrderItem {
	list := make([]backendApi.OrderItem, 0, len(in))
	for _, item := range in {
		list = append(list, toBackendApiOrderItem(item))
	}
	return list
}

func toBackendApiOrderItem(in model.OrderBase) backendApi.OrderItem {
	return backendApi.OrderItem{
		Id:          in.Id,
		OrderNo:     in.OrderNo,
		EmployeeId:  in.EmployeeId,
		TotalPoints: in.TotalPoints,
		Status:      in.Status,
		Remark:      in.Remark,
		CreatedAt:   in.CreatedAt,
	}
}

func toBackendApiOrderDetail(in model.OrderDetail) backendApi.OrderDetail {
	items := make([]backendApi.OrderGoodsItem, 0, len(in.Items))
	for _, item := range in.Items {
		items = append(items, backendApi.OrderGoodsItem{
			Id:            item.Id,
			OrderId:       item.OrderId,
			EmployeeId:    item.EmployeeId,
			GoodsId:       item.GoodsId,
			GoodsName:     item.GoodsName,
			GoodsImageUrl: item.GoodsImageUrl,
			PointsPrice:   item.PointsPrice,
			Count:         item.Count,
			TotalPoints:   item.TotalPoints,
			CreatedAt:     item.CreatedAt,
		})
	}
	return backendApi.OrderDetail{
		OrderItem: toBackendApiOrderItem(in.OrderBase),
		Items:     items,
	}
}
