package user

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
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

type sUser struct{}

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}

// 注册
func (s *sUser) Register(ctx context.Context, in model.RegisterInput) (out model.RegisterOutput, err error) {
	// 校验用户名是否已经存在
	count, err := dao.UserInfo.Ctx(ctx).Where(dao.UserInfo.Columns().Name, in.Name).Count()
	if err != nil {
		return out, err
	}
	if count > 0 {
		return out, errors.New("用户名已存在")
	}

	//处理加密盐和密码的逻辑
	UserSalt := grand.S(10)
	in.Password = utility.EncryptPassword(in.Password, UserSalt)
	in.UserSalt = UserSalt
	in.Status = 1 // 默认状态为正常
	//插入数据返回id
	lastInsertID, err := dao.UserInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.RegisterOutput{Id: uint(lastInsertID)}, err
}

// 修改密码
func (*sUser) UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error) {
	//	验证密保问题
	userInfo := do.UserInfo{}
	userId := gconv.Uint(service.UserAuth().GetIdentity(ctx))
	err = dao.UserInfo.Ctx(ctx).WherePri(userId).Scan(&userInfo)
	if err != nil {
		return model.UpdatePasswordOutput{}, err
	}
	if gconv.String(userInfo.SecretAnswer) != in.SecretAnswer {
		return out, errors.New(consts.ErrSecretAnswerMsg)
	}
	userSalt := grand.S(10)
	in.UserSalt = userSalt
	in.Password = utility.EncryptPassword(in.Password, userSalt)
	_, err = dao.UserInfo.Ctx(ctx).WherePri(userId).Update(in)
	if err != nil {
		return model.UpdatePasswordOutput{}, err
	}
	return model.UpdatePasswordOutput{Id: userId}, nil
}

// GetList 获取用户列表（后台管理）
func (s *sUser) GetList(ctx context.Context, in model.UserGetListInput) (out *model.UserGetListOutput, err error) {
	m := dao.UserInfo.Ctx(ctx)
	out = &model.UserGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}

	// 分页查询
	listModel := m.Page(in.Page, in.Size).OrderDesc(dao.UserInfo.Columns().Id)

	// 查询总数
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		out.List = make([]model.UserInfoItem, 0)
		return out, err
	}

	// 查询列表
	var list []entity.UserInfo
	if err = listModel.Scan(&list); err != nil {
		return out, err
	}

	// 转换为输出格式
	out.List = make([]model.UserInfoItem, 0, len(list))
	for _, item := range list {
		out.List = append(out.List, model.UserInfoItem{
			Id:        item.Id,
			Name:      item.Name,
			Avatar:    item.Avatar,
			Sex:       item.Sex,
			Status:    item.Status,
			Sign:      item.Sign,
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		})
	}

	return out, nil
}

// UpdateStatus 更新用户状态（冻结/解冻）
func (s *sUser) UpdateStatus(ctx context.Context, id uint, status int) error {
	_, err := dao.UserInfo.Ctx(ctx).
		Where(dao.UserInfo.Columns().Id, id).
		Data(g.Map{dao.UserInfo.Columns().Status: status}).
		Update()
	return err
}

// Delete 删除用户
func (s *sUser) Delete(ctx context.Context, id uint) error {
	_, err := dao.UserInfo.Ctx(ctx).
		Where(dao.UserInfo.Columns().Id, id).
		Delete()
	return err
}

// GetById 根据ID获取用户信息
func (s *sUser) GetById(ctx context.Context, id uint) (*model.UserInfoItem, error) {
	var user entity.UserInfo
	err := dao.UserInfo.Ctx(ctx).Where(dao.UserInfo.Columns().Id, id).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, nil
	}
	return &model.UserInfoItem{
		Id:        user.Id,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Sex:       user.Sex,
		Status:    user.Status,
		Sign:      user.Sign,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
