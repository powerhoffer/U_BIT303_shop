package points

import (
	"context"
	"database/sql"
	"errors"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sPoints struct{}

func init() {
	service.RegisterPoints(New())
}

func New() *sPoints {
	return &sPoints{}
}

func (s *sPoints) Balance(ctx context.Context, employeeId uint) (out model.PointsBalanceOutput, err error) {
	account, err := s.getAccount(ctx, employeeId)
	if err != nil {
		return out, err
	}
	if account.Id == 0 {
		return out, nil
	}
	out.Balance = account.Balance
	return out, nil
}

func (s *sPoints) Records(ctx context.Context, in model.PointsRecordsInput) (out model.PointsRecordsOutput, err error) {
	return s.recordsByEmployee(ctx, in)
}

func (s *sPoints) ManageRecords(ctx context.Context, in model.PointsRecordsInput) (out model.PointsRecordsOutput, err error) {
	if err = s.checkEmployee(ctx, in.EmployeeId); err != nil {
		return out, err
	}
	return s.recordsByEmployee(ctx, in)
}

func (s *sPoints) ManageAdd(ctx context.Context, in model.PointsChangeInput) (out model.PointsChangeOutput, err error) {
	return s.changePoints(ctx, in, consts.PointsChangeTypeAdd)
}

func (s *sPoints) ManageDeduct(ctx context.Context, in model.PointsChangeInput) (out model.PointsChangeOutput, err error) {
	return s.changePoints(ctx, in, consts.PointsChangeTypeDeduct)
}

func (s *sPoints) changePoints(ctx context.Context, in model.PointsChangeInput, changeType int) (out model.PointsChangeOutput, err error) {
	if in.Points == 0 {
		return out, errors.New("Credits must be greater than 0")
	}
	if err = s.checkEmployee(ctx, in.EmployeeId); err != nil {
		return out, err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		account, err := s.getAccountForUpdate(ctx, in.EmployeeId)
		if err != nil {
			return err
		}
		if account.Id == 0 {
			if changeType == consts.PointsChangeTypeDeduct {
				return errors.New("Insufficient credit balance")
			}
			if _, err = dao.EmployeePointsAccount.Ctx(ctx).Data(do.EmployeePointsAccount{
				EmployeeId: in.EmployeeId,
				Balance:    0,
				Status:     consts.PointsAccountStatusNormal,
			}).Insert(); err != nil {
				return err
			}
			account, err = s.getAccountForUpdate(ctx, in.EmployeeId)
			if err != nil {
				return err
			}
		}
		if account.Status != consts.PointsAccountStatusNormal {
			return errors.New("Credit account is disabled")
		}

		beforeBalance := account.Balance
		afterBalance := beforeBalance
		switch changeType {
		case consts.PointsChangeTypeAdd:
			afterBalance = beforeBalance + in.Points
		case consts.PointsChangeTypeDeduct:
			if beforeBalance < in.Points {
				return errors.New("Insufficient credit balance")
			}
			afterBalance = beforeBalance - in.Points
		default:
			return errors.New("Credit change type is invalid")
		}

		if _, err = dao.EmployeePointsAccount.Ctx(ctx).
			Where(dao.EmployeePointsAccount.Columns().Id, account.Id).
			Data(g.Map{dao.EmployeePointsAccount.Columns().Balance: afterBalance}).
			Update(); err != nil {
			return err
		}
		if _, err = dao.EmployeePointsRecord.Ctx(ctx).Data(do.EmployeePointsRecord{
			EmployeeId:         in.EmployeeId,
			ChangeType:         changeType,
			Points:             in.Points,
			BeforeBalance:      beforeBalance,
			AfterBalance:       afterBalance,
			OperatorEmployeeId: in.OperatorEmployeeId,
			Remark:             in.Remark,
		}).Insert(); err != nil {
			return err
		}
		out.Balance = afterBalance
		return nil
	})
	return out, err
}

func (s *sPoints) recordsByEmployee(ctx context.Context, in model.PointsRecordsInput) (out model.PointsRecordsOutput, err error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 {
		in.Size = 10
	}
	if in.Size > 50 {
		in.Size = 50
	}

	columns := dao.EmployeePointsRecord.Columns()
	m := dao.EmployeePointsRecord.Ctx(ctx).Where(columns.EmployeeId, in.EmployeeId)
	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.PointsRecordsOutput{
		List:  make([]model.PointsRecordItem, 0),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}
	if total == 0 {
		return out, nil
	}

	var records []entity.EmployeePointsRecord
	if err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&records); err != nil {
		return out, err
	}
	for _, record := range records {
		out.List = append(out.List, toRecordItem(record))
	}
	return out, nil
}

func (s *sPoints) checkEmployee(ctx context.Context, employeeId uint) error {
	columns := dao.EmployeeInfo.Columns()
	count, err := dao.EmployeeInfo.Ctx(ctx).
		Where(columns.Id, employeeId).
		Where(columns.Status, consts.EmployeeStatusNormal).
		WhereNull(columns.DeletedAt).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Employee account does not exist or is disabled")
	}
	return nil
}

func (s *sPoints) getAccount(ctx context.Context, employeeId uint) (account entity.EmployeePointsAccount, err error) {
	columns := dao.EmployeePointsAccount.Columns()
	err = dao.EmployeePointsAccount.Ctx(ctx).
		Where(columns.EmployeeId, employeeId).
		WhereNull(columns.DeletedAt).
		Scan(&account)
	if errors.Is(err, sql.ErrNoRows) {
		return account, nil
	}
	return
}

func (s *sPoints) getAccountForUpdate(ctx context.Context, employeeId uint) (account entity.EmployeePointsAccount, err error) {
	columns := dao.EmployeePointsAccount.Columns()
	err = dao.EmployeePointsAccount.Ctx(ctx).
		Where(columns.EmployeeId, employeeId).
		WhereNull(columns.DeletedAt).
		LockUpdate().
		Scan(&account)
	if errors.Is(err, sql.ErrNoRows) {
		return account, nil
	}
	return
}

func toRecordItem(record entity.EmployeePointsRecord) model.PointsRecordItem {
	item := model.PointsRecordItem{
		Id:                 record.Id,
		EmployeeId:         record.EmployeeId,
		ChangeType:         record.ChangeType,
		Points:             record.Points,
		BeforeBalance:      record.BeforeBalance,
		AfterBalance:       record.AfterBalance,
		OperatorEmployeeId: record.OperatorEmployeeId,
		Remark:             record.Remark,
	}
	if record.CreatedAt != nil {
		item.CreatedAt = record.CreatedAt.Time
	}
	return item
}
