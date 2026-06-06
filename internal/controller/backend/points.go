package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Points = cPoints{}

type cPoints struct{}

func (c *cPoints) Balance(ctx context.Context, req *backendApi.PointsBalanceReq) (res *backendApi.PointsBalanceRes, err error) {
	out, err := service.Points().Balance(ctx, currentEmployeeId(ctx))
	if err != nil {
		return nil, err
	}
	return &backendApi.PointsBalanceRes{Balance: out.Balance}, nil
}

func (c *cPoints) Records(ctx context.Context, req *backendApi.PointsRecordsReq) (res *backendApi.PointsRecordsRes, err error) {
	out, err := service.Points().Records(ctx, model.PointsRecordsInput{
		EmployeeId: currentEmployeeId(ctx),
		Page:       req.Page,
		Size:       req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.PointsRecordsRes{
		List:  toApiRecords(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cPoints) ManageAdd(ctx context.Context, req *backendApi.PointsManageAddReq) (res *backendApi.PointsManageAddRes, err error) {
	out, err := service.Points().ManageAdd(ctx, model.PointsChangeInput{
		EmployeeId:         req.EmployeeId,
		OperatorEmployeeId: currentEmployeeId(ctx),
		Points:             req.Points,
		Remark:             req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.PointsManageAddRes{Balance: out.Balance}, nil
}

func (c *cPoints) ManageDeduct(ctx context.Context, req *backendApi.PointsManageDeductReq) (res *backendApi.PointsManageDeductRes, err error) {
	out, err := service.Points().ManageDeduct(ctx, model.PointsChangeInput{
		EmployeeId:         req.EmployeeId,
		OperatorEmployeeId: currentEmployeeId(ctx),
		Points:             req.Points,
		Remark:             req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &backendApi.PointsManageDeductRes{Balance: out.Balance}, nil
}

func (c *cPoints) ManageRecords(ctx context.Context, req *backendApi.PointsManageRecordsReq) (res *backendApi.PointsManageRecordsRes, err error) {
	var input model.PointsRecordsInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Points().ManageRecords(ctx, input)
	if err != nil {
		return nil, err
	}
	return &backendApi.PointsManageRecordsRes{
		List:  toApiRecords(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func toApiRecords(in []model.PointsRecordItem) []backendApi.PointsRecordItem {
	list := make([]backendApi.PointsRecordItem, 0, len(in))
	for _, item := range in {
		list = append(list, backendApi.PointsRecordItem{
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
