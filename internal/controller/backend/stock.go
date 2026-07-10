package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Stock = cStock{}

type cStock struct{}

func (c *cStock) Adjust(ctx context.Context, req *backendApi.StockAdjustReq) (res *backendApi.StockAdjustRes, err error) {
	var input model.StockAdjustInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	input.OperatorAdminId = currentAdminId(ctx)
	out, err := service.Stock().Adjust(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.StockAdjustRes{GoodsId: out.GoodsId, Stock: out.Stock}, nil
}

func (c *cStock) Records(ctx context.Context, req *backendApi.StockRecordsReq) (res *backendApi.StockRecordsRes, err error) {
	var input model.StockRecordsInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Stock().Records(ctx, input)
	if err != nil {
		return nil, err
	}
	list := make([]backendApi.StockRecordItem, 0, len(out.List))
	for _, item := range out.List {
		list = append(list, backendApi.StockRecordItem{
			Id:             item.Id,
			GoodsId:        item.GoodsId,
			GoodsName:      item.GoodsName,
			ChangeType:     item.ChangeType,
			ChangeQuantity: item.ChangeQuantity,
			BeforeStock:    item.BeforeStock,
			AfterStock:     item.AfterStock,
			BizType:        item.BizType,
			BizId:          item.BizId,
			OperatorType:   item.OperatorType,
			OperatorId:     item.OperatorId,
			Remark:         item.Remark,
			CreatedAt:      item.CreatedAt,
		})
	}
	return &backendApi.StockRecordsRes{List: list, Total: out.Total, Page: out.Page, Size: out.Size}, nil
}
