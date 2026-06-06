package employee

import (
	"context"
	"errors"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
	"bit303_shop/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"
)

type sEmployee struct{}

func init() {
	service.RegisterEmployee(New())
}

func New() *sEmployee {
	return &sEmployee{}
}

func (s *sEmployee) Register(ctx context.Context, in model.EmployeeRegisterInput) (out model.EmployeeRegisterOutput, err error) {
	columns := dao.EmployeeInfo.Columns()
	count, err := dao.EmployeeInfo.Ctx(ctx).
		Where(columns.Username, in.Username).
		WhereNull(columns.DeletedAt).
		Count()
	if err != nil {
		return out, err
	}
	if count > 0 {
		return out, errors.New("员工账号已存在")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return out, err
	}

	lastInsertId, err := dao.EmployeeInfo.Ctx(ctx).Data(do.EmployeeInfo{
		Username:     in.Username,
		PasswordHash: string(passwordHash),
		RealName:     in.RealName,
		Phone:        in.Phone,
		Email:        in.Email,
		Status:       consts.EmployeeStatusNormal,
	}).InsertAndGetId()
	if err != nil {
		return out, err
	}

	out.Employee = model.EmployeeBase{
		Id:       uint(lastInsertId),
		Username: in.Username,
		RealName: in.RealName,
		Phone:    in.Phone,
		Email:    in.Email,
		Status:   consts.EmployeeStatusNormal,
	}
	return out, nil
}

func (s *sEmployee) Login(ctx context.Context, in model.EmployeeLoginInput) (out model.EmployeeLoginOutput, err error) {
	employee, err := s.getNormalEmployeeByUsername(ctx, in.Username)
	if err != nil {
		return out, err
	}
	if employee.Id == 0 {
		return out, errors.New("账号或者密码不正确")
	}
	if employee.Status != consts.EmployeeStatusNormal {
		return out, errors.New("员工账号已禁用")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(in.Password)); err != nil {
		return out, errors.New("账号或者密码不正确")
	}

	out.Token, out.ExpireAt, err = utility.GenerateEmployeeToken(employee.Id, employee.Username, in.Remember)
	if err != nil {
		return out, err
	}
	_, err = dao.EmployeeInfo.Ctx(ctx).
		Where(dao.EmployeeInfo.Columns().Id, employee.Id).
		Data(g.Map{dao.EmployeeInfo.Columns().LastLoginAt: gtime.Now()}).
		Update()
	if err != nil {
		return out, err
	}

	out.Employee = toEmployeeBase(employee)
	return out, nil
}

func (s *sEmployee) Info(ctx context.Context, employeeId uint) (out model.EmployeeInfoOutput, err error) {
	employee, err := s.getEmployeeById(ctx, employeeId)
	if err != nil {
		return out, err
	}
	if employee.Id == 0 {
		return out, errors.New("员工账号不存在")
	}
	out.Employee = toEmployeeBase(employee)
	return out, nil
}

func (s *sEmployee) UpdatePassword(ctx context.Context, in model.EmployeeUpdatePasswordInput) error {
	employee, err := s.getEmployeeById(ctx, in.EmployeeId)
	if err != nil {
		return err
	}
	if employee.Id == 0 {
		return errors.New("员工账号不存在")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(in.OldPassword)); err != nil {
		return errors.New("旧密码不正确")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.EmployeeInfo.Ctx(ctx).
		Where(dao.EmployeeInfo.Columns().Id, employee.Id).
		Data(g.Map{dao.EmployeeInfo.Columns().PasswordHash: string(passwordHash)}).
		Update()
	return err
}

func (s *sEmployee) ManageCreate(ctx context.Context, in model.EmployeeManageCreateInput) (out model.EmployeeManageCreateOutput, err error) {
	registerOut, err := s.Register(ctx, model.EmployeeRegisterInput{
		Username: in.Username,
		Password: in.Password,
		RealName: in.RealName,
		Phone:    in.Phone,
		Email:    in.Email,
	})
	if err != nil {
		return out, err
	}
	out.Employee = registerOut.Employee
	return out, nil
}

func (s *sEmployee) ManageList(ctx context.Context, in model.EmployeeManageListInput) (out model.EmployeeManageListOutput, err error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 {
		in.Size = 10
	}
	if in.Size > 50 {
		in.Size = 50
	}

	columns := dao.EmployeeInfo.Columns()
	m := dao.EmployeeInfo.Ctx(ctx).WhereNull(columns.DeletedAt)
	if in.Username != "" {
		m = m.Where(columns.Username+" LIKE ?", "%"+in.Username+"%")
	}
	if in.RealName != "" {
		m = m.Where(columns.RealName+" LIKE ?", "%"+in.RealName+"%")
	}
	if in.Status == consts.EmployeeStatusDisabled || in.Status == consts.EmployeeStatusNormal {
		m = m.Where(columns.Status, in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.EmployeeManageListOutput{
		List:  make([]model.EmployeeBase, 0),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}
	if total == 0 {
		return out, nil
	}

	var employees []entity.EmployeeInfo
	err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&employees)
	if err != nil {
		return out, err
	}
	for _, employee := range employees {
		out.List = append(out.List, toEmployeeBase(employee))
	}
	return out, nil
}

func (s *sEmployee) ManageDetail(ctx context.Context, id uint) (out model.EmployeeManageDetailOutput, err error) {
	employee, err := s.getEmployeeById(ctx, id)
	if err != nil {
		return out, err
	}
	if employee.Id == 0 {
		return out, errors.New("员工账号不存在")
	}
	out.Employee = toEmployeeBase(employee)
	return out, nil
}

func (s *sEmployee) ManageUpdate(ctx context.Context, in model.EmployeeManageUpdateInput) (out model.EmployeeManageUpdateOutput, err error) {
	employee, err := s.getEmployeeById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	if employee.Id == 0 {
		return out, errors.New("员工账号不存在")
	}
	_, err = dao.EmployeeInfo.Ctx(ctx).
		Where(dao.EmployeeInfo.Columns().Id, in.Id).
		Data(g.Map{
			dao.EmployeeInfo.Columns().RealName: in.RealName,
			dao.EmployeeInfo.Columns().Phone:    in.Phone,
			dao.EmployeeInfo.Columns().Email:    in.Email,
		}).
		Update()
	if err != nil {
		return out, err
	}
	employee, err = s.getEmployeeById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	out.Employee = toEmployeeBase(employee)
	return out, nil
}

func (s *sEmployee) ManageUpdateStatus(ctx context.Context, in model.EmployeeManageStatusInput) error {
	if in.Status != consts.EmployeeStatusDisabled && in.Status != consts.EmployeeStatusNormal {
		return errors.New("状态只能是0或1")
	}
	employee, err := s.getEmployeeById(ctx, in.Id)
	if err != nil {
		return err
	}
	if employee.Id == 0 {
		return errors.New("员工账号不存在")
	}
	_, err = dao.EmployeeInfo.Ctx(ctx).
		Where(dao.EmployeeInfo.Columns().Id, in.Id).
		Data(g.Map{dao.EmployeeInfo.Columns().Status: in.Status}).
		Update()
	return err
}

func (s *sEmployee) ManageResetPassword(ctx context.Context, in model.EmployeeManageResetPasswordInput) error {
	employee, err := s.getEmployeeById(ctx, in.Id)
	if err != nil {
		return err
	}
	if employee.Id == 0 {
		return errors.New("员工账号不存在")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.EmployeeInfo.Ctx(ctx).
		Where(dao.EmployeeInfo.Columns().Id, in.Id).
		Data(g.Map{dao.EmployeeInfo.Columns().PasswordHash: string(passwordHash)}).
		Update()
	return err
}

func (s *sEmployee) getNormalEmployeeByUsername(ctx context.Context, username string) (employee entity.EmployeeInfo, err error) {
	columns := dao.EmployeeInfo.Columns()
	err = dao.EmployeeInfo.Ctx(ctx).
		Where(columns.Username, username).
		WhereNull(columns.DeletedAt).
		Scan(&employee)
	return
}

func (s *sEmployee) getEmployeeById(ctx context.Context, employeeId uint) (employee entity.EmployeeInfo, err error) {
	columns := dao.EmployeeInfo.Columns()
	err = dao.EmployeeInfo.Ctx(ctx).
		Where(columns.Id, employeeId).
		WhereNull(columns.DeletedAt).
		Scan(&employee)
	return
}

func toEmployeeBase(employee entity.EmployeeInfo) model.EmployeeBase {
	return model.EmployeeBase{
		Id:       employee.Id,
		Username: employee.Username,
		RealName: employee.RealName,
		Phone:    employee.Phone,
		Email:    employee.Email,
		Status:   employee.Status,
	}
}
