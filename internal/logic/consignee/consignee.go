package consignee

import (
	"context"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type sConsignee struct{}

func init() {
	service.RegisterConsignee(New())
}

func New() *sConsignee {
	return &sConsignee{}
}

// GetList 查询内容列表
func (s *sConsignee) GetList(ctx context.Context, in model.ConsigneeGetListInput) (out *model.ConsigneeGetListOutput, err error) {
	// 获取当前用户ID
	userId := gconv.Uint(ctx.Value(consts.CtxUserId))

	var (
		m = dao.ConsigneeInfo.Ctx(ctx).Where(dao.ConsigneeInfo.Columns().UserId, userId)
	)
	out = &model.ConsigneeGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}

	// 分配查询
	listModel := m.Page(in.Page, in.Size)

	// 执行查询
	var list []*entity.ConsigneeInfo
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

// Add 添加收货地址
func (s *sConsignee) Add(ctx context.Context, in model.AddConsigneeInput) (out *model.AddConsigneeOutput, err error) {
	// 如果设置为默认地址，先将该用户的其他地址设置为非默认
	if in.IsDefault == 1 {
		_, err = dao.ConsigneeInfo.Ctx(ctx).
			Where(dao.ConsigneeInfo.Columns().UserId, in.UserId).
			Data(dao.ConsigneeInfo.Columns().IsDefault, 0).
			Update()
		if err != nil {
			return out, err
		}
	}

	// 插入新地址
	id, err := dao.ConsigneeInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}

	return &model.AddConsigneeOutput{Id: uint(id)}, nil
}

// Update 更新收货地址
func (s *sConsignee) Update(ctx context.Context, in model.UpdateConsigneeInput) (out *model.UpdateConsigneeOutput, err error) {
	// 如果设置为默认地址，先将该用户的其他地址设置为非默认
	if in.IsDefault == 1 {
		_, err = dao.ConsigneeInfo.Ctx(ctx).
			Where(dao.ConsigneeInfo.Columns().UserId, in.UserId).
			Where(dao.ConsigneeInfo.Columns().Id+"!=", in.Id).
			Data(dao.ConsigneeInfo.Columns().IsDefault, 0).
			Update()
		if err != nil {
			return out, err
		}
	}

	// 更新地址信息
	_, err = dao.ConsigneeInfo.Ctx(ctx).
		Where(dao.ConsigneeInfo.Columns().Id, in.Id).
		Where(dao.ConsigneeInfo.Columns().UserId, in.UserId).
		Data(in).
		Update()
	if err != nil {
		return out, err
	}

	return &model.UpdateConsigneeOutput{Id: in.Id}, nil
}

// Delete 删除收货地址
func (s *sConsignee) Delete(ctx context.Context, in model.DeleteConsigneeInput) (out *model.DeleteConsigneeOutput, err error) {
	// 删除地址
	_, err = dao.ConsigneeInfo.Ctx(ctx).
		Where(dao.ConsigneeInfo.Columns().Id, in.Id).
		Delete()
	if err != nil {
		return out, err
	}

	return &model.DeleteConsigneeOutput{Id: in.Id}, nil
}

// AdminGetList 后台管理列表（包含用户信息）
func (s *sConsignee) AdminGetList(ctx context.Context, page, size int) (out *model.ConsigneeAdminListOutput, err error) {
	out = &model.ConsigneeAdminListOutput{
		Page: page,
		Size: size,
	}

	m := dao.ConsigneeInfo.Ctx(ctx)

	// 查询总数
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}

	// 联表查询用户信息
	type ConsigneeWithUser struct {
		entity.ConsigneeInfo
		UserName string `json:"user_name"`
	}

	var list []ConsigneeWithUser
	err = m.LeftJoin("user_info", "consignee_info.user_id = user_info.id").
		Fields("consignee_info.*, user_info.name as user_name").
		Page(page, size).
		OrderDesc(dao.ConsigneeInfo.Columns().Id).
		Scan(&list)
	if err != nil {
		return out, err
	}

	// 转换为输出格式
	out.List = make([]model.ConsigneeAdminListItem, 0, len(list))
	for _, item := range list {
		// 拼接完整地址
		address := item.Province + item.City + item.Town
		if item.Street != "" && item.Street != "null" {
			address += item.Street
		}
		address += item.Detail

		out.List = append(out.List, model.ConsigneeAdminListItem{
			Id:        uint(item.Id),
			UserName:  item.UserName,
			IsDefault: item.IsDefault,
			Name:      item.Name,
			Phone:     item.Phone,
			Address:   address,
		})
	}

	return out, nil
}

// AdminDelete 后台管理删除（不检查用户ID）
func (s *sConsignee) AdminDelete(ctx context.Context, id uint) error {
	_, err := dao.ConsigneeInfo.Ctx(ctx).
		Where(dao.ConsigneeInfo.Columns().Id, id).
		Delete()
	return err
}
