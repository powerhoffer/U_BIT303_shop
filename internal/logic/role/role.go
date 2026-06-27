package role

import (
	"context"
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

type sRole struct{}

func init() {
	service.RegisterRole(New())
}

func New() *sRole {
	return &sRole{}
}

func (s *sRole) Create(ctx context.Context, in model.RoleCreateInput) (out model.RoleCreateOutput, err error) {
	columns := dao.AdminRole.Columns()
	count, err := dao.AdminRole.Ctx(ctx).Where(columns.Name, in.Name).WhereNull(columns.DeletedAt).Count()
	if err != nil {
		return out, err
	}
	if count > 0 {
		return out, errors.New("Role already exists")
	}
	id, err := dao.AdminRole.Ctx(ctx).Data(do.AdminRole{
		Name:        in.Name,
		Description: in.Description,
		Status:      consts.AdminRoleEnabled,
	}).InsertAndGetId()
	if err != nil {
		return out, err
	}
	out.Role = model.RoleBase{Id: uint(id), Name: in.Name, Description: in.Description, Status: consts.AdminRoleEnabled, PermissionIds: []uint{}}
	return out, nil
}

func (s *sRole) List(ctx context.Context, in model.RoleListInput) (out model.RoleListOutput, err error) {
	normalizePage(&in.Page, &in.Size, 50)
	columns := dao.AdminRole.Columns()
	m := dao.AdminRole.Ctx(ctx).WhereNull(columns.DeletedAt)
	if in.Name != "" {
		m = m.Where(columns.Name+" LIKE ?", "%"+in.Name+"%")
	}
	if in.Status == consts.AdminRoleDisabled || in.Status == consts.AdminRoleEnabled {
		m = m.Where(columns.Status, in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.RoleListOutput{List: make([]model.RoleBase, 0), Total: total, Page: in.Page, Size: in.Size}
	if total == 0 {
		return out, nil
	}
	var roles []entity.AdminRole
	if err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&roles); err != nil {
		return out, err
	}
	for _, item := range roles {
		out.List = append(out.List, s.toRoleBase(ctx, item))
	}
	return out, nil
}

func (s *sRole) Detail(ctx context.Context, id uint) (out model.RoleDetailOutput, err error) {
	role, err := s.getRoleById(ctx, id)
	if err != nil {
		return out, err
	}
	if role.Id == 0 {
		return out, errors.New("Role does not exist")
	}
	out.Role = s.toRoleBase(ctx, role)
	return out, nil
}

func (s *sRole) Update(ctx context.Context, in model.RoleUpdateInput) (out model.RoleUpdateOutput, err error) {
	role, err := s.getRoleById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	if role.Id == 0 {
		return out, errors.New("Role does not exist")
	}
	_, err = dao.AdminRole.Ctx(ctx).
		Where(dao.AdminRole.Columns().Id, in.Id).
		Data(g.Map{dao.AdminRole.Columns().Name: in.Name, dao.AdminRole.Columns().Description: in.Description}).
		Update()
	if err != nil {
		return out, err
	}
	role, err = s.getRoleById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	out.Role = s.toRoleBase(ctx, role)
	return out, nil
}

func (s *sRole) UpdateStatus(ctx context.Context, in model.RoleStatusInput) error {
	if in.Status != consts.AdminRoleDisabled && in.Status != consts.AdminRoleEnabled {
		return errors.New("Status must be 0 or 1")
	}
	role, err := s.getRoleById(ctx, in.Id)
	if err != nil {
		return err
	}
	if role.Id == 0 {
		return errors.New("Role does not exist")
	}
	_, err = dao.AdminRole.Ctx(ctx).
		Where(dao.AdminRole.Columns().Id, in.Id).
		Data(g.Map{dao.AdminRole.Columns().Status: in.Status}).
		Update()
	return err
}

func (s *sRole) AssignPermissions(ctx context.Context, in model.RolePermissionsInput) error {
	role, err := s.getRoleById(ctx, in.Id)
	if err != nil {
		return err
	}
	if role.Id == 0 {
		return errors.New("Role does not exist")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.AdminRolePermission.Ctx(ctx).TX(tx).Where(dao.AdminRolePermission.Columns().RoleId, in.Id).Delete()
		if err != nil {
			return err
		}
		for _, permissionId := range uniqueUintSlice(in.PermissionIds) {
			_, err = dao.AdminRolePermission.Ctx(ctx).TX(tx).Data(do.AdminRolePermission{RoleId: in.Id, PermissionId: permissionId}).Insert()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *sRole) getRoleById(ctx context.Context, id uint) (role entity.AdminRole, err error) {
	columns := dao.AdminRole.Columns()
	err = dao.AdminRole.Ctx(ctx).Where(columns.Id, id).WhereNull(columns.DeletedAt).Scan(&role)
	return
}

func (s *sRole) toRoleBase(ctx context.Context, role entity.AdminRole) model.RoleBase {
	return model.RoleBase{Id: role.Id, Name: role.Name, Description: role.Description, Status: role.Status, PermissionIds: rolePermissionIds(ctx, role.Id)}
}

func rolePermissionIds(ctx context.Context, roleId uint) []uint {
	var relations []entity.AdminRolePermission
	err := dao.AdminRolePermission.Ctx(ctx).Where(dao.AdminRolePermission.Columns().RoleId, roleId).Scan(&relations)
	if err != nil {
		return []uint{}
	}
	ids := make([]uint, 0, len(relations))
	for _, relation := range relations {
		ids = append(ids, relation.PermissionId)
	}
	return ids
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

func uniqueUintSlice(in []uint) []uint {
	seen := make(map[uint]struct{})
	out := make([]uint, 0, len(in))
	for _, id := range in {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
