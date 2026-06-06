package controller

import (
	"context"

	"bit303_shop/api/points"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

var Points = cPoints{}

type cPoints struct{}

func (c *cPoints) Balance(ctx context.Context, req *points.BalanceReq) (res *points.BalanceRes, err error) {
	out, err := service.Points().Balance(ctx, pointsCurrentEmployeeId(ctx))
	if err != nil {
		return nil, err
	}
	return &points.BalanceRes{Balance: out.Balance}, nil
}

func (c *cPoints) Records(ctx context.Context, req *points.RecordsReq) (res *points.RecordsRes, err error) {
	out, err := service.Points().Records(ctx, model.PointsRecordsInput{
		EmployeeId: pointsCurrentEmployeeId(ctx),
		Page:       req.Page,
		Size:       req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &points.RecordsRes{
		List:  toApiRecords(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cPoints) ManageAdd(ctx context.Context, req *points.ManageAddReq) (res *points.ManageAddRes, err error) {
	out, err := service.Points().ManageAdd(ctx, model.PointsChangeInput{
		EmployeeId:         req.EmployeeId,
		OperatorEmployeeId: pointsCurrentEmployeeId(ctx),
		Points:             req.Points,
		Remark:             req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &points.ManageAddRes{Balance: out.Balance}, nil
}

func (c *cPoints) ManageDeduct(ctx context.Context, req *points.ManageDeductReq) (res *points.ManageDeductRes, err error) {
	out, err := service.Points().ManageDeduct(ctx, model.PointsChangeInput{
		EmployeeId:         req.EmployeeId,
		OperatorEmployeeId: pointsCurrentEmployeeId(ctx),
		Points:             req.Points,
		Remark:             req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &points.ManageDeductRes{Balance: out.Balance}, nil
}

func (c *cPoints) ManageRecords(ctx context.Context, req *points.ManageRecordsReq) (res *points.ManageRecordsRes, err error) {
	var input model.PointsRecordsInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Points().ManageRecords(ctx, input)
	if err != nil {
		return nil, err
	}
	return &points.ManageRecordsRes{
		List:  toApiRecords(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func pointsCurrentEmployeeId(ctx context.Context) uint {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar(consts.CtxEmployeeId).Uint()
}

func toApiRecords(in []model.PointsRecordItem) []points.RecordItem {
	list := make([]points.RecordItem, 0, len(in))
	for _, item := range in {
		list = append(list, points.RecordItem{
			Id:                 item.Id,
			EmployeeId:         item.EmployeeId,
			ChangeType:         item.ChangeType,
			Points:             item.Points,
			BeforeBalance:      item.BeforeBalance,
			AfterBalance:       item.AfterBalance,
			OperatorEmployeeId: item.OperatorEmployeeId,
			Remark:             item.Remark,
			CreatedAt:          item.CreatedAt,
		})
	}
	return list
}
