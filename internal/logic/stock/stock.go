package stock

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sStock struct{}

func init() {
	service.RegisterStock(New())
}

func New() *sStock {
	return &sStock{}
}

func (s *sStock) Adjust(ctx context.Context, in model.StockAdjustInput) (out model.StockAdjustOutput, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		columns := dao.GoodsInfo.Columns()
		var goods entity.GoodsInfo
		if err := dao.GoodsInfo.Ctx(ctx).
			Where(columns.Id, in.GoodsId).
			WhereNull(columns.DeletedAt).
			LockUpdate().
			Scan(&goods); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if goods.Id == 0 {
			return errors.New("Goods does not exist")
		}

		beforeStock := goods.Stock
		afterStock := beforeStock
		changeType := 0
		changeQuantity := 0
		switch in.Action {
		case consts.StockActionIncrease:
			if in.Quantity > uint(math.MaxUint32)-beforeStock {
				return errors.New("Stock exceeds limit")
			}
			afterStock += in.Quantity
			changeType = consts.StockChangeTypeAdminIncrease
			changeQuantity = int(in.Quantity)
		case consts.StockActionDecrease:
			if beforeStock < in.Quantity {
				return errors.New("Insufficient goods stock")
			}
			afterStock -= in.Quantity
			changeType = consts.StockChangeTypeAdminDecrease
			changeQuantity = -int(in.Quantity)
		default:
			return errors.New("Action must be increase or decrease")
		}

		if _, err := dao.GoodsInfo.Ctx(ctx).
			Where(columns.Id, goods.Id).
			Data(g.Map{columns.Stock: afterStock}).
			Update(); err != nil {
			return err
		}
		if err := s.RecordChange(ctx, model.StockRecordInput{
			GoodsId:        goods.Id,
			GoodsName:      goods.Name,
			ChangeType:     changeType,
			ChangeQuantity: changeQuantity,
			BeforeStock:    beforeStock,
			AfterStock:     afterStock,
			BizType:        consts.StockBizTypeStockAdjust,
			BizId:          goods.Id,
			OperatorType:   consts.StockOperatorTypeAdmin,
			OperatorId:     in.OperatorAdminId,
			Remark:         in.Remark,
		}); err != nil {
			return err
		}
		out = model.StockAdjustOutput{GoodsId: goods.Id, Stock: afterStock}
		return nil
	})
	return out, err
}

func (s *sStock) RecordChange(ctx context.Context, in model.StockRecordInput) error {
	_, err := dao.GoodsStockRecord.Ctx(ctx).Data(do.GoodsStockRecord{
		GoodsId:        in.GoodsId,
		GoodsName:      in.GoodsName,
		ChangeType:     in.ChangeType,
		ChangeQuantity: in.ChangeQuantity,
		BeforeStock:    in.BeforeStock,
		AfterStock:     in.AfterStock,
		BizType:        in.BizType,
		BizId:          in.BizId,
		OperatorType:   in.OperatorType,
		OperatorId:     in.OperatorId,
		Remark:         in.Remark,
	}).Insert()
	return err
}

func (s *sStock) Records(ctx context.Context, in model.StockRecordsInput) (out model.StockRecordsOutput, err error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 {
		in.Size = 10
	}
	if in.Size > 50 {
		in.Size = 50
	}

	startTime, endTime, err := parseStockDateRange(in.StartTime, in.EndTime)
	if err != nil {
		return out, err
	}

	columns := dao.GoodsStockRecord.Columns()
	m := dao.GoodsStockRecord.Ctx(ctx)
	if in.GoodsId > 0 {
		m = m.Where(columns.GoodsId, in.GoodsId)
	}
	if in.ChangeType >= consts.StockChangeTypeInitial && in.ChangeType <= consts.StockChangeTypeOrderCancelRestore {
		m = m.Where(columns.ChangeType, in.ChangeType)
	}
	if startTime != nil {
		m = m.Where(columns.CreatedAt+" >= ?", startTime)
	}
	if endTime != nil {
		m = m.Where(columns.CreatedAt+" < ?", endTime)
	}

	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.StockRecordsOutput{List: make([]model.StockRecordItem, 0), Total: total, Page: in.Page, Size: in.Size}
	if total == 0 {
		return out, nil
	}

	var records []entity.GoodsStockRecord
	if err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&records); err != nil {
		return out, err
	}
	for _, record := range records {
		item := model.StockRecordItem{
			Id:             record.Id,
			GoodsId:        record.GoodsId,
			GoodsName:      record.GoodsName,
			ChangeType:     record.ChangeType,
			ChangeQuantity: record.ChangeQuantity,
			BeforeStock:    record.BeforeStock,
			AfterStock:     record.AfterStock,
			BizType:        record.BizType,
			BizId:          record.BizId,
			OperatorType:   record.OperatorType,
			OperatorId:     record.OperatorId,
			Remark:         record.Remark,
		}
		if record.CreatedAt != nil {
			item.CreatedAt = record.CreatedAt.Time
		}
		out.List = append(out.List, item)
	}
	return out, nil
}

func parseStockDateRange(startValue, endValue string) (startTime, endTime *time.Time, err error) {
	const dateLayout = "2006-01-02"
	if startValue != "" {
		start, parseErr := time.ParseInLocation(dateLayout, startValue, time.Local)
		if parseErr != nil {
			return nil, nil, errors.New("Start time must use YYYY-MM-DD format")
		}
		startTime = &start
	}
	if endValue != "" {
		end, parseErr := time.ParseInLocation(dateLayout, endValue, time.Local)
		if parseErr != nil {
			return nil, nil, errors.New("End time must use YYYY-MM-DD format")
		}
		end = end.AddDate(0, 0, 1)
		endTime = &end
	}
	if startTime != nil && endTime != nil && !startTime.Before(*endTime) {
		return nil, nil, errors.New("Start time must not be later than end time")
	}
	return startTime, endTime, nil
}
