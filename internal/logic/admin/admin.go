package admin

import (
	"context"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
	"bit303_shop/utility"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/ghtml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
)

type sAdmin struct{}

func init() {
	service.RegisterAdmin(New())
}

func New() *sAdmin {
	return &sAdmin{}
}

func (s *sAdmin) Create(ctx context.Context, in *model.AdminCreateInput) (out model.AdminCreateOutput, err error) {
	// 不允许HTML代码
	if err = ghtml.SpecialCharsMapOrStruct(in); err != nil {
		return out, err
	}
	//处理加密盐和密码的逻辑
	UserSalt := grand.S(10)
	in.Password = utility.EncryptPassword(in.Password, UserSalt)
	in.UserSalt = UserSalt
	//插入数据返回id
	lastInsertID, err := dao.AdminInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.AdminCreateOutput{AdminId: int(lastInsertID)}, err
}

func (s *sAdmin) GetUserByUserNamePassword(ctx context.Context, in model.UserLoginInput) map[string]interface{} {
	//验证账号密码是否正确
	adminInfo := entity.AdminInfo{}
	err := dao.AdminInfo.Ctx(ctx).Where(dao.AdminInfo.Columns().Name, in.Name).Scan(&adminInfo)
	if err != nil {
		return nil
	}
	if utility.EncryptPassword(in.Password, adminInfo.UserSalt) != adminInfo.Password {
		return nil
	} else {
		return g.Map{
			"id":       adminInfo.Id,
			"username": adminInfo.Name,
		}
	}
}

// Delete 删除
func (s *sAdmin) Delete(ctx context.Context, id uint) error {
	return dao.AdminInfo.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除内容
		_, err := dao.AdminInfo.Ctx(ctx).Where(g.Map{
			dao.AdminInfo.Columns().Id: id,
		}).Delete()
		return err
	})
}

// Update 修改
func (s *sAdmin) Update(ctx context.Context, in model.AdminUpdateInput) error {
	//判断是否修改了密码
	if in.Password != "" {
		UserSalt := grand.S(10)
		in.Password = utility.EncryptPassword(in.Password, UserSalt)
		in.UserSalt = UserSalt
	}
	//更新操作
	_, err := dao.AdminInfo.
		Ctx(ctx).
		Data(in).
		OmitEmpty(). //注意：使用OmitEmpty()实现部分更新
		Where(dao.AdminInfo.Columns().Id, in.Id).
		Update()
	return err
}

//1.获得*gdb.Model对象，方面后续调用
//2. 实例化响应结构体
//3. 分页查询
//4. 再查询count，判断有无数据
//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
//6. 把查询到的结果赋值到响应结构体中

// GetList 查询内容列表
func (s *sAdmin) GetList(ctx context.Context, in model.AdminGetListInput) (out *model.AdminGetListOutput, err error) {
	//1.获得*gdb.Model对象，方便后续调用
	m := dao.AdminInfo.Ctx(ctx)
	//2. 实例化响应结构体
	out = &model.AdminGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		//解决空数据返回[] 而不是返回nil null的问题
		out.List = make([]model.AdminGetListOutputItem, 0)
		return out, err
	}

	// 查询管理员列表
	var list []entity.AdminInfo
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}

	// 构建输出列表，解析 role_ids 为数组
	out.List = make([]model.AdminGetListOutputItem, 0, len(list))
	for _, admin := range list {
		// 解析 role_ids 字符串为数组
		var roleIdArray []int
		if admin.RoleIds != "" {
			roleIdStrs := strings.Split(admin.RoleIds, ",")
			for _, idStr := range roleIdStrs {
				idStr = strings.TrimSpace(idStr)
				if idStr != "" {
					if id, err := strconv.Atoi(idStr); err == nil {
						roleIdArray = append(roleIdArray, id)
					}
				}
			}
		}

		out.List = append(out.List, model.AdminGetListOutputItem{
			Id:          uint(admin.Id),
			Name:        admin.Name,
			RoleIds:     admin.RoleIds,
			RoleIdArray: roleIdArray,
			IsAdmin:     admin.IsAdmin,
			CreatedAt:   admin.CreatedAt,
			UpdatedAt:   admin.UpdatedAt,
		})
	}
	return
}

func (s *sAdmin) GetAdminByNamePassword(ctx context.Context, in model.UserLoginInput) map[string]interface{} {
	//验证账号密码是否正确
	adminInfo := entity.AdminInfo{}
	err := dao.AdminInfo.Ctx(ctx).Where("name", in.Name).Scan(&adminInfo)
	if err != nil {
		return nil
	}
	if utility.EncryptPassword(in.Password, adminInfo.UserSalt) != adminInfo.Password {
		return nil
	} else {
		return g.Map{
			"id":       adminInfo.Id,
			"username": adminInfo.Name,
		}
	}
}

// GetAdminByNamePasswordWithRoles 登录验证并返回包含角色信息的数据（用于JWT存储）
func (s *sAdmin) GetAdminByNamePasswordWithRoles(ctx context.Context, in model.UserLoginInput) map[string]interface{} {
	adminInfo := entity.AdminInfo{}
	err := dao.AdminInfo.Ctx(ctx).Where("name", in.Name).Scan(&adminInfo)
	if err != nil {
		return nil
	}
	if utility.EncryptPassword(in.Password, adminInfo.UserSalt) != adminInfo.Password {
		return nil
	}
	return g.Map{
		"id":       adminInfo.Id,
		"username": adminInfo.Name,
		"is_admin": adminInfo.IsAdmin,
		"role_ids": adminInfo.RoleIds,
	}
}

// GetById 根据ID获取管理员信息
func (s *sAdmin) GetById(ctx context.Context, id int) (*model.AdminInfo, error) {
	var adminInfo entity.AdminInfo
	err := dao.AdminInfo.Ctx(ctx).Where(dao.AdminInfo.Columns().Id, id).Scan(&adminInfo)
	if err != nil {
		return nil, err
	}
	if adminInfo.Id == 0 {
		return nil, nil
	}
	return &model.AdminInfo{
		Id:      adminInfo.Id,
		Name:    adminInfo.Name,
		RoleIds: adminInfo.RoleIds,
		IsAdmin: adminInfo.IsAdmin,
	}, nil
}
