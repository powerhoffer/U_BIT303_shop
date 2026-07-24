package points

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"sort"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

const maxBatchPointsEmployees = 200

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

func (s *sPoints) ManageBatchAdd(ctx context.Context, in model.PointsBatchAddInput) (out model.PointsBatchAddOutput, err error) {
	if in.OperatorAdminId == 0 {
		return out, errors.New("Admin identity is required")
	}
	if in.Points == 0 {
		return out, errors.New("Credits must be greater than 0")
	}
	employeeIds, err := normalizeBatchEmployeeIds(in.EmployeeIds)
	if err != nil {
		return out, err
	}

	results := make([]model.PointsBatchAddResultItem, 0, len(employeeIds))
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if err := s.validateEmployeesForUpdate(ctx, employeeIds); err != nil {
			return err
		}
		for _, employeeId := range employeeIds {
			changeOut, err := s.applyPointsChange(ctx, model.PointsChangeInput{
				EmployeeId:      employeeId,
				OperatorAdminId: in.OperatorAdminId,
				Points:          in.Points,
				Remark:          in.Remark,
			}, consts.PointsChangeTypeAdd)
			if err != nil {
				return err
			}
			results = append(results, model.PointsBatchAddResultItem{
				EmployeeId: employeeId,
				Balance:    changeOut.Balance,
			})
		}
		return nil
	})
	if err != nil {
		return out, err
	}
	out = model.PointsBatchAddOutput{
		ProcessedCount: len(results),
		TotalPoints:    uint64(len(results)) * uint64(in.Points),
		List:           results,
	}
	return out, nil
}

func (s *sPoints) ManageDeduct(ctx context.Context, in model.PointsChangeInput) (out model.PointsChangeOutput, err error) {
	return s.changePoints(ctx, in, consts.PointsChangeTypeDeduct)
}

func (s *sPoints) changePoints(ctx context.Context, in model.PointsChangeInput, changeType int) (out model.PointsChangeOutput, err error) {
	if in.OperatorAdminId == 0 {
		return out, errors.New("Admin identity is required")
	}
	if in.Points == 0 {
		return out, errors.New("Credits must be greater than 0")
	}
	if err = s.checkEmployee(ctx, in.EmployeeId); err != nil {
		return out, err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		changeOut, changeErr := s.applyPointsChange(ctx, in, changeType)
		if changeErr != nil {
			return changeErr
		}
		out = changeOut
		return nil
	})
	return out, err
}

func (s *sPoints) applyPointsChange(ctx context.Context, in model.PointsChangeInput, changeType int) (out model.PointsChangeOutput, err error) {
	account, err := s.getAccountForUpdate(ctx, in.EmployeeId)
	if err != nil {
		return out, err
	}
	if account.Id == 0 {
		if changeType == consts.PointsChangeTypeDeduct {
			return out, errors.New("Insufficient credit balance")
		}
		if _, err = dao.EmployeePointsAccount.Ctx(ctx).Data(do.EmployeePointsAccount{
			EmployeeId: in.EmployeeId,
			Balance:    0,
			Status:     consts.PointsAccountStatusNormal,
		}).Insert(); err != nil {
			return out, err
		}
		account, err = s.getAccountForUpdate(ctx, in.EmployeeId)
		if err != nil {
			return out, err
		}
		if account.Id == 0 {
			return out, errors.New("Credit account could not be initialized")
		}
	}
	if account.Status != consts.PointsAccountStatusNormal {
		return out, errors.New("Credit account is disabled")
	}

	beforeBalance := account.Balance
	afterBalance := beforeBalance
	switch changeType {
	case consts.PointsChangeTypeAdd:
		if in.Points > uint(math.MaxUint32)-beforeBalance {
			return out, errors.New("Credit balance exceeds allowed maximum")
		}
		afterBalance = beforeBalance + in.Points
	case consts.PointsChangeTypeDeduct:
		if beforeBalance < in.Points {
			return out, errors.New("Insufficient credit balance")
		}
		afterBalance = beforeBalance - in.Points
	default:
		return out, errors.New("Credit change type is invalid")
	}

	if _, err = dao.EmployeePointsAccount.Ctx(ctx).
		Where(dao.EmployeePointsAccount.Columns().Id, account.Id).
		Data(g.Map{dao.EmployeePointsAccount.Columns().Balance: afterBalance}).
		Update(); err != nil {
		return out, err
	}
	if _, err = dao.EmployeePointsRecord.Ctx(ctx).Data(do.EmployeePointsRecord{
		EmployeeId:         in.EmployeeId,
		ChangeType:         changeType,
		Points:             in.Points,
		BeforeBalance:      beforeBalance,
		AfterBalance:       afterBalance,
		OperatorEmployeeId: in.OperatorEmployeeId,
		OperatorAdminId:    in.OperatorAdminId,
		Remark:             in.Remark,
	}).Insert(); err != nil {
		return out, err
	}
	out.Balance = afterBalance
	return out, nil
}

func normalizeBatchEmployeeIds(employeeIds []uint) ([]uint, error) {
	if len(employeeIds) == 0 {
		return nil, errors.New("At least one employee is required")
	}
	if len(employeeIds) > maxBatchPointsEmployees {
		return nil, errors.New("A batch can contain at most 200 employees")
	}
	seen := make(map[uint]struct{}, len(employeeIds))
	ids := make([]uint, 0, len(employeeIds))
	for _, employeeId := range employeeIds {
		if employeeId == 0 {
			return nil, errors.New("Employee ID is invalid")
		}
		if _, exists := seen[employeeId]; exists {
			return nil, errors.New("Employee IDs must not contain duplicates")
		}
		seen[employeeId] = struct{}{}
		ids = append(ids, employeeId)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}

func (s *sPoints) validateEmployeesForUpdate(ctx context.Context, employeeIds []uint) error {
	columns := dao.EmployeeInfo.Columns()
	var employees []entity.EmployeeInfo
	if err := dao.EmployeeInfo.Ctx(ctx).
		WhereIn(columns.Id, employeeIds).
		Where(columns.Status, consts.EmployeeStatusNormal).
		WhereNull(columns.DeletedAt).
		OrderAsc(columns.Id).
		LockUpdate().
		Scan(&employees); err != nil {
		return err
	}
	if len(employees) != len(employeeIds) {
		return errors.New("One or more employee accounts do not exist or are disabled")
	}
	return nil
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
		OperatorAdminId:    record.OperatorAdminId,
		Remark:             record.Remark,
	}
	if record.CreatedAt != nil {
		item.CreatedAt = record.CreatedAt.Time
	}
	return item
}
