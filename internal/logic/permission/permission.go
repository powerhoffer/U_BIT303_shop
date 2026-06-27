package permission

import (
	"context"
	"errors"
	"strings"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sPermission struct{}

func init() {
	service.RegisterPermission(New())
}

func New() *sPermission {
	return &sPermission{}
}

func (s *sPermission) Create(ctx context.Context, in model.PermissionCreateInput) (out model.PermissionCreateOutput, err error) {
	in.Method = strings.ToUpper(in.Method)
	columns := dao.AdminPermission.Columns()
	count, err := dao.AdminPermission.Ctx(ctx).
		Where(columns.Method, in.Method).
		Where(columns.Path, in.Path).
		WhereNull(columns.DeletedAt).
		Count()
	if err != nil {
		return out, err
	}
	if count > 0 {
		return out, errors.New("Permission already exists")
	}
	id, err := dao.AdminPermission.Ctx(ctx).Data(do.AdminPermission{
		Name:      in.Name,
		GroupName: in.GroupName,
		Method:    in.Method,
		Path:      in.Path,
		Status:    consts.AdminPermissionEnabled,
	}).InsertAndGetId()
	if err != nil {
		return out, err
	}
	out.Permission = model.PermissionBase{Id: uint(id), Name: in.Name, GroupName: in.GroupName, Method: in.Method, Path: in.Path, Status: consts.AdminPermissionEnabled}
	return out, nil
}

func (s *sPermission) List(ctx context.Context, in model.PermissionListInput) (out model.PermissionListOutput, err error) {
	normalizePage(&in.Page, &in.Size, 100)
	columns := dao.AdminPermission.Columns()
	m := dao.AdminPermission.Ctx(ctx).WhereNull(columns.DeletedAt)
	if in.Name != "" {
		m = m.Where(columns.Name+" LIKE ?", "%"+in.Name+"%")
	}
	if in.GroupName != "" {
		m = m.Where(columns.GroupName, in.GroupName)
	}
	if in.Method != "" {
		m = m.Where(columns.Method, strings.ToUpper(in.Method))
	}
	if in.Status == consts.AdminPermissionDisabled || in.Status == consts.AdminPermissionEnabled {
		m = m.Where(columns.Status, in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.PermissionListOutput{List: make([]model.PermissionBase, 0), Total: total, Page: in.Page, Size: in.Size}
	if total == 0 {
		return out, nil
	}
	var permissions []entity.AdminPermission
	if err = m.Page(in.Page, in.Size).OrderAsc(columns.GroupName).OrderAsc(columns.Id).Scan(&permissions); err != nil {
		return out, err
	}
	for _, item := range permissions {
		out.List = append(out.List, toPermissionBase(item))
	}
	return out, nil
}

func (s *sPermission) Detail(ctx context.Context, id uint) (out model.PermissionDetailOutput, err error) {
	permission, err := s.getPermissionById(ctx, id)
	if err != nil {
		return out, err
	}
	if permission.Id == 0 {
		return out, errors.New("Permission does not exist")
	}
	out.Permission = toPermissionBase(permission)
	return out, nil
}

func (s *sPermission) Update(ctx context.Context, in model.PermissionUpdateInput) (out model.PermissionUpdateOutput, err error) {
	permission, err := s.getPermissionById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	if permission.Id == 0 {
		return out, errors.New("Permission does not exist")
	}
	in.Method = strings.ToUpper(in.Method)
	_, err = dao.AdminPermission.Ctx(ctx).
		Where(dao.AdminPermission.Columns().Id, in.Id).
		Data(g.Map{
			dao.AdminPermission.Columns().Name:      in.Name,
			dao.AdminPermission.Columns().GroupName: in.GroupName,
			dao.AdminPermission.Columns().Method:    in.Method,
			dao.AdminPermission.Columns().Path:      in.Path,
		}).Update()
	if err != nil {
		return out, err
	}
	permission, err = s.getPermissionById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	out.Permission = toPermissionBase(permission)
	return out, nil
}

func (s *sPermission) UpdateStatus(ctx context.Context, in model.PermissionStatusInput) error {
	if in.Status != consts.AdminPermissionDisabled && in.Status != consts.AdminPermissionEnabled {
		return errors.New("Status must be 0 or 1")
	}
	permission, err := s.getPermissionById(ctx, in.Id)
	if err != nil {
		return err
	}
	if permission.Id == 0 {
		return errors.New("Permission does not exist")
	}
	_, err = dao.AdminPermission.Ctx(ctx).
		Where(dao.AdminPermission.Columns().Id, in.Id).
		Data(g.Map{dao.AdminPermission.Columns().Status: in.Status}).
		Update()
	return err
}

func (s *sPermission) getPermissionById(ctx context.Context, id uint) (permission entity.AdminPermission, err error) {
	columns := dao.AdminPermission.Columns()
	err = dao.AdminPermission.Ctx(ctx).Where(columns.Id, id).WhereNull(columns.DeletedAt).Scan(&permission)
	return
}

func toPermissionBase(permission entity.AdminPermission) model.PermissionBase {
	return model.PermissionBase{Id: permission.Id, Name: permission.Name, GroupName: permission.GroupName, Method: permission.Method, Path: permission.Path, Status: permission.Status}
}

func normalizePage(page *int, size *int, max int) {
	if *page < 1 {
		*page = 1
	}
	if *size < 1 {
		*size = 10
	}
	if *size > max {
		*size = max
	}
}
