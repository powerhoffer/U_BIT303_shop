package admin

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
	"bit303_shop/utility"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"
)

type sAdmin struct{}

func init() {
	service.RegisterAdmin(New())
}

func New() *sAdmin {
	return &sAdmin{}
}

func (s *sAdmin) Login(ctx context.Context, in model.AdminLoginInput) (out model.AdminLoginOutput, err error) {
	admin, err := s.getAdminByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return out, errors.New("Invalid username or password")
		}
		return out, err
	}
	if admin.Id == 0 {
		return out, errors.New("Invalid username or password")
	}
	if admin.Status != consts.AdminStatusNormal {
		return out, errors.New("Admin account is disabled")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(in.Password)); err != nil {
		return out, errors.New("Invalid username or password")
	}

	out.Token, out.ExpireAt, err = utility.GenerateAdminToken(admin.Id, admin.Username, admin.IsSuper, in.Remember)
	if err != nil {
		return out, err
	}
	_, err = dao.AdminInfo.Ctx(ctx).
		Where(dao.AdminInfo.Columns().Id, admin.Id).
		Data(g.Map{dao.AdminInfo.Columns().LastLoginAt: gtime.Now()}).
		Update()
	if err != nil {
		return out, err
	}
	out.Admin = s.toAdminBase(ctx, admin)
	return out, nil
}

func (s *sAdmin) Info(ctx context.Context, adminId uint) (out model.AdminInfoOutput, err error) {
	admin, err := s.getAdminById(ctx, adminId)
	if err != nil {
		return out, err
	}
	if admin.Id == 0 {
		return out, errors.New("Admin account does not exist")
	}
	out.Admin = s.toAdminBase(ctx, admin)
	return out, nil
}

func (s *sAdmin) UpdatePassword(ctx context.Context, in model.AdminUpdatePasswordInput) error {
	admin, err := s.getAdminById(ctx, in.AdminId)
	if err != nil {
		return err
	}
	if admin.Id == 0 {
		return errors.New("Admin account does not exist")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(in.OldPassword)); err != nil {
		return errors.New("Current password is incorrect")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.AdminInfo.Ctx(ctx).
		Where(dao.AdminInfo.Columns().Id, admin.Id).
		Data(g.Map{dao.AdminInfo.Columns().PasswordHash: string(passwordHash)}).
		Update()
	return err
}

func (s *sAdmin) ManageCreate(ctx context.Context, in model.AdminManageCreateInput) (out model.AdminManageCreateOutput, err error) {
	columns := dao.AdminInfo.Columns()
	count, err := dao.AdminInfo.Ctx(ctx).Where(columns.Username, in.Username).WhereNull(columns.DeletedAt).Count()
	if err != nil {
		return out, err
	}
	if count > 0 {
		return out, errors.New("Admin account already exists")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return out, err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := dao.AdminInfo.Ctx(ctx).TX(tx).Data(do.AdminInfo{
			Username:     in.Username,
			PasswordHash: string(passwordHash),
			RealName:     in.RealName,
			Phone:        in.Phone,
			Email:        in.Email,
			Status:       consts.AdminStatusNormal,
			IsSuper:      normalizeFlag(in.IsSuper),
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		if err = s.assignRoles(ctx, tx, uint(id), in.RoleIds); err != nil {
			return err
		}
		out.Admin = model.AdminBase{
			Id:       uint(id),
			Username: in.Username,
			RealName: in.RealName,
			Phone:    in.Phone,
			Email:    in.Email,
			Status:   consts.AdminStatusNormal,
			IsSuper:  normalizeFlag(in.IsSuper),
			RoleIds:  uniqueUintSlice(in.RoleIds),
		}
		return nil
	})
	return out, err
}

func (s *sAdmin) ManageList(ctx context.Context, in model.AdminManageListInput) (out model.AdminManageListOutput, err error) {
	normalizePage(&in.Page, &in.Size, 50)
	columns := dao.AdminInfo.Columns()
	m := dao.AdminInfo.Ctx(ctx).WhereNull(columns.DeletedAt)
	if in.Username != "" {
		m = m.Where(columns.Username+" LIKE ?", "%"+in.Username+"%")
	}
	if in.RealName != "" {
		m = m.Where(columns.RealName+" LIKE ?", "%"+in.RealName+"%")
	}
	if in.Status == consts.AdminStatusDisabled || in.Status == consts.AdminStatusNormal {
		m = m.Where(columns.Status, in.Status)
	}
	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.AdminManageListOutput{List: make([]model.AdminBase, 0), Total: total, Page: in.Page, Size: in.Size}
	if total == 0 {
		return out, nil
	}
	var admins []entity.AdminInfo
	if err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&admins); err != nil {
		return out, err
	}
	for _, item := range admins {
		out.List = append(out.List, s.toAdminBase(ctx, item))
	}
	return out, nil
}

func (s *sAdmin) ManageDetail(ctx context.Context, id uint) (out model.AdminManageDetailOutput, err error) {
	admin, err := s.getAdminById(ctx, id)
	if err != nil {
		return out, err
	}
	if admin.Id == 0 {
		return out, errors.New("Admin account does not exist")
	}
	out.Admin = s.toAdminBase(ctx, admin)
	return out, nil
}

func (s *sAdmin) ManageUpdate(ctx context.Context, in model.AdminManageUpdateInput) (out model.AdminManageUpdateOutput, err error) {
	admin, err := s.getAdminById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	if admin.Id == 0 {
		return out, errors.New("Admin account does not exist")
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.AdminInfo.Ctx(ctx).TX(tx).
			Where(dao.AdminInfo.Columns().Id, in.Id).
			Data(g.Map{
				dao.AdminInfo.Columns().RealName: in.RealName,
				dao.AdminInfo.Columns().Phone:    in.Phone,
				dao.AdminInfo.Columns().Email:    in.Email,
				dao.AdminInfo.Columns().IsSuper:  normalizeFlag(in.IsSuper),
			}).Update()
		if err != nil {
			return err
		}
		return s.assignRoles(ctx, tx, in.Id, in.RoleIds)
	})
	if err != nil {
		return out, err
	}
	admin, err = s.getAdminById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	out.Admin = s.toAdminBase(ctx, admin)
	return out, nil
}

func (s *sAdmin) ManageUpdateStatus(ctx context.Context, in model.AdminManageStatusInput) error {
	if in.Status != consts.AdminStatusDisabled && in.Status != consts.AdminStatusNormal {
		return errors.New("Status must be 0 or 1")
	}
	admin, err := s.getAdminById(ctx, in.Id)
	if err != nil {
		return err
	}
	if admin.Id == 0 {
		return errors.New("Admin account does not exist")
	}
	_, err = dao.AdminInfo.Ctx(ctx).
		Where(dao.AdminInfo.Columns().Id, in.Id).
		Data(g.Map{dao.AdminInfo.Columns().Status: in.Status}).
		Update()
	return err
}

func (s *sAdmin) ManageResetPassword(ctx context.Context, in model.AdminManageResetPasswordInput) error {
	admin, err := s.getAdminById(ctx, in.Id)
	if err != nil {
		return err
	}
	if admin.Id == 0 {
		return errors.New("Admin account does not exist")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.AdminInfo.Ctx(ctx).
		Where(dao.AdminInfo.Columns().Id, in.Id).
		Data(g.Map{dao.AdminInfo.Columns().PasswordHash: string(passwordHash)}).
		Update()
	return err
}

func (s *sAdmin) ManageAssignRoles(ctx context.Context, in model.AdminManageRolesInput) error {
	admin, err := s.getAdminById(ctx, in.Id)
	if err != nil {
		return err
	}
	if admin.Id == 0 {
		return errors.New("Admin account does not exist")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		return s.assignRoles(ctx, tx, in.Id, in.RoleIds)
	})
}

func (s *sAdmin) getAdminByUsername(ctx context.Context, username string) (admin entity.AdminInfo, err error) {
	columns := dao.AdminInfo.Columns()
	err = dao.AdminInfo.Ctx(ctx).Where(columns.Username, username).WhereNull(columns.DeletedAt).Scan(&admin)
	return
}

func (s *sAdmin) getAdminById(ctx context.Context, id uint) (admin entity.AdminInfo, err error) {
	columns := dao.AdminInfo.Columns()
	err = dao.AdminInfo.Ctx(ctx).Where(columns.Id, id).WhereNull(columns.DeletedAt).Scan(&admin)
	return
}

func (s *sAdmin) assignRoles(ctx context.Context, tx gdb.TX, adminId uint, roleIds []uint) error {
	roleIds = uniqueUintSlice(roleIds)
	_, err := dao.AdminRoleRelation.Ctx(ctx).TX(tx).Where(dao.AdminRoleRelation.Columns().AdminId, adminId).Delete()
	if err != nil {
		return err
	}
	for _, roleId := range roleIds {
		_, err = dao.AdminRoleRelation.Ctx(ctx).TX(tx).Data(do.AdminRoleRelation{AdminId: adminId, RoleId: roleId}).Insert()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sAdmin) toAdminBase(ctx context.Context, admin entity.AdminInfo) model.AdminBase {
	return model.AdminBase{
		Id:       admin.Id,
		Username: admin.Username,
		RealName: admin.RealName,
		Phone:    admin.Phone,
		Email:    admin.Email,
		Status:   admin.Status,
		IsSuper:  admin.IsSuper,
		RoleIds:  adminRoleIds(ctx, admin.Id),
	}
}

func adminRoleIds(ctx context.Context, adminId uint) []uint {
	var relations []entity.AdminRoleRelation
	err := dao.AdminRoleRelation.Ctx(ctx).
		Where(dao.AdminRoleRelation.Columns().AdminId, adminId).
		Scan(&relations)
	if err != nil {
		return []uint{}
	}
	ids := make([]uint, 0, len(relations))
	for _, relation := range relations {
		ids = append(ids, relation.RoleId)
	}
	return ids
}

func normalizeFlag(flag int) int {
	if flag == 1 {
		return 1
	}
	return 0
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
