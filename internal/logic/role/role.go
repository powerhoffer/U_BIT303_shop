package role

import (
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/context"
)

type sRole struct{}

func init() {
	service.RegisterRole(New())
}

func New() *sRole {
	return &sRole{}
}

func (s *sRole) Create(ctx context.Context, in model.RoleCreateInput) (out model.RoleCreateOutput, err error) {
	//插入数据返回id
	lastInsertID, err := dao.RoleInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.RoleCreateOutput{RoleId: uint(lastInsertID)}, err
}

// 角色添加权限
func (s *sRole) AddPermission(ctx context.Context, in model.RoleAddPermissionInput) (out model.RoleAddPermissionOutput, err error) {
	id, err := dao.RolePermissionInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return model.RoleAddPermissionOutput{}, err
	}
	return model.RoleAddPermissionOutput{Id: uint(id)}, err
}

// Delete 删除
func (s *sRole) Delete(ctx context.Context, id uint) error {
	// 删除内容
	_, err := dao.RoleInfo.Ctx(ctx).Where(g.Map{
		dao.RoleInfo.Columns().Id: id,
	}).Unscoped().Delete()
	return err
}

// 角色删除权限
func (s *sRole) DeletePermission(ctx context.Context, in model.RoleDeletePermissionInput) error {
	_, err := dao.RolePermissionInfo.Ctx(ctx).Where(g.Map{
		dao.RolePermissionInfo.Columns().RoleId:       in.RoleId,
		dao.RolePermissionInfo.Columns().PermissionId: in.PermissionId,
	}).Delete()
	if err != nil {
		return err
	}
	return err
}

// Update 修改
func (s *sRole) Update(ctx context.Context, in model.RoleUpdateInput) error {
	_, err := dao.RoleInfo.
		Ctx(ctx).
		Data(in).
		FieldsEx(dao.RoleInfo.Columns().Id).
		Where(dao.RoleInfo.Columns().Id, in.Id).
		Update()
	return err
}

// GetList 查询内容列表
func (s *sRole) GetList(ctx context.Context, in model.RoleGetListInput) (out *model.RoleGetListOutput, err error) {
	var (
		m = dao.RoleInfo.Ctx(ctx)
	)
	out = &model.RoleGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}

	// 分配查询
	listModel := m.Page(in.Page, in.Size)

	// 执行查询
	var list []*entity.RoleInfo
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}
	// 没有数据
	if len(list) == 0 {
		out.List = make([]model.RoleGetListOutputItem, 0)
		return out, nil
	}
	out.Total, err = m.Count()
	if err != nil {
		return out, err
	}

	// 构建输出列表，包含权限信息
	out.List = make([]model.RoleGetListOutputItem, 0, len(list))
	for _, role := range list {
		// 获取该角色的权限列表
		permissions, _ := s.GetPermissionsByRoleId(ctx, uint(role.Id))
		out.List = append(out.List, model.RoleGetListOutputItem{
			Id:          uint(role.Id),
			Name:        role.Name,
			Desc:        role.Desc,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			Permissions: permissions,
		})
	}
	return
}

// AddPermissions 批量添加权限
func (s *sRole) AddPermissions(ctx context.Context, roleId uint, permissionIds []uint) error {
	if len(permissionIds) == 0 {
		return nil
	}
	// 批量插入
	data := make([]g.Map, 0, len(permissionIds))
	for _, pid := range permissionIds {
		data = append(data, g.Map{
			dao.RolePermissionInfo.Columns().RoleId:       roleId,
			dao.RolePermissionInfo.Columns().PermissionId: pid,
		})
	}
	_, err := dao.RolePermissionInfo.Ctx(ctx).Data(data).Insert()
	return err
}

// DeletePermissions 批量删除权限
func (s *sRole) DeletePermissions(ctx context.Context, roleId uint, permissionIds []uint) error {
	if len(permissionIds) == 0 {
		return nil
	}
	_, err := dao.RolePermissionInfo.Ctx(ctx).
		Where(dao.RolePermissionInfo.Columns().RoleId, roleId).
		WhereIn(dao.RolePermissionInfo.Columns().PermissionId, permissionIds).
		Delete()
	return err
}

// GetPermissionsByRoleId 根据角色ID获取权限列表
func (s *sRole) GetPermissionsByRoleId(ctx context.Context, roleId uint) ([]model.RolePermissionItem, error) {
	// 查询角色权限关联
	var rolePermissions []entity.RolePermissionInfo
	err := dao.RolePermissionInfo.Ctx(ctx).
		Where(dao.RolePermissionInfo.Columns().RoleId, roleId).
		Scan(&rolePermissions)
	if err != nil {
		return nil, err
	}
	if len(rolePermissions) == 0 {
		return []model.RolePermissionItem{}, nil
	}

	// 提取权限ID
	permissionIds := make([]int, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		permissionIds = append(permissionIds, rp.PermissionId)
	}

	// 查询权限详情
	var permissions []entity.PermissionInfo
	err = dao.PermissionInfo.Ctx(ctx).
		WhereIn(dao.PermissionInfo.Columns().Id, permissionIds).
		Scan(&permissions)
	if err != nil {
		return nil, err
	}

	// 转换输出
	result := make([]model.RolePermissionItem, 0, len(permissions))
	for _, p := range permissions {
		result = append(result, model.RolePermissionItem{
			Id:   uint(p.Id),
			Name: p.Name,
			Path: p.Path,
		})
	}
	return result, nil
}
