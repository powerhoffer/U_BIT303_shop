package role

import (
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/net/context"
)

type sPermission struct{}

func init() {
	service.RegisterPermission(New())
}

func New() *sPermission {
	return &sPermission{}
}

func (s *sPermission) Create(ctx context.Context, in model.PermissionCreateInput) (out model.PermissionCreateOutput, err error) {
	//插入数据返回id
	lastInsertID, err := dao.PermissionInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.PermissionCreateOutput{PermissionId: uint(lastInsertID)}, err
}

// Delete 删除
func (s *sPermission) Delete(ctx context.Context, id uint) error {
	// 删除内容
	_, err := dao.PermissionInfo.Ctx(ctx).Where(g.Map{
		dao.PermissionInfo.Columns().Id: id,
	}).Unscoped().Delete()
	return err
}

// Update 修改
func (s *sPermission) Update(ctx context.Context, in model.PermissionUpdateInput) error {
	_, err := dao.PermissionInfo.
		Ctx(ctx).
		Data(in).
		FieldsEx(dao.PermissionInfo.Columns().Id).
		Where(dao.PermissionInfo.Columns().Id, in.Id).
		Update()
	return err
}

// GetList 查询内容列表
func (s *sPermission) GetList(ctx context.Context, in model.PermissionGetListInput) (out *model.PermissionGetListOutput, err error) {
	var (
		m = dao.PermissionInfo.Ctx(ctx)
	)
	out = &model.PermissionGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}

	// 分配查询
	listModel := m.Page(in.Page, in.Size)

	// 执行查询
	var list []*entity.PermissionInfo
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}
	// 没有数据
	if len(list) == 0 {
		return out, nil
	}
	out.Total, err = m.Count()
	if err != nil {
		return out, err
	}
	//不指定item的键名用：Scan
	if err := listModel.Scan(&out.List); err != nil {
		return out, err
	}
	return
}

// GetPathsByRoleIds 根据角色ID列表获取所有权限路径
func (s *sPermission) GetPathsByRoleIds(ctx context.Context, roleIds []int) ([]string, error) {
	if len(roleIds) == 0 {
		return []string{}, nil
	}

	// 1. 根据 role_ids 查询 role_permission_info 表获取 permission_ids
	var rolePermissions []entity.RolePermissionInfo
	err := dao.RolePermissionInfo.Ctx(ctx).
		WhereIn(dao.RolePermissionInfo.Columns().RoleId, roleIds).
		Scan(&rolePermissions)
	if err != nil {
		return nil, err
	}

	if len(rolePermissions) == 0 {
		g.Log().Debug(ctx, "GetPathsByRoleIds - no role_permission found for roleIds:", roleIds)
		return []string{}, nil
	}

	// 提取 permission_ids 并去重
	permissionIdMap := make(map[int]bool)
	for _, rp := range rolePermissions {
		permissionIdMap[rp.PermissionId] = true
	}
	permissionIds := make([]int, 0, len(permissionIdMap))
	for id := range permissionIdMap {
		permissionIds = append(permissionIds, id)
	}

	g.Log().Debug(ctx, "GetPathsByRoleIds - permissionIds:", permissionIds)

	// 2. 根据 permission_ids 查询 permission_info 表获取 path
	var permissions []entity.PermissionInfo
	err = dao.PermissionInfo.Ctx(ctx).
		WhereIn(dao.PermissionInfo.Columns().Id, permissionIds).
		Scan(&permissions)
	if err != nil {
		return nil, err
	}

	// 3. 提取所有 path
	paths := make([]string, 0, len(permissions))
	for _, p := range permissions {
		if p.Path != "" {
			paths = append(paths, p.Path)
		}
	}

	g.Log().Debug(ctx, "GetPathsByRoleIds - paths:", paths)

	return paths, nil
}
