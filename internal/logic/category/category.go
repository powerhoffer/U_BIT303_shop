package category

import (
	"context"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sCategory struct{}

func init() {
	service.RegisterCategory(New())
}

func New() *sCategory {
	return &sCategory{}
}

func (s *sCategory) Create(ctx context.Context, in model.CategoryCreateInput) (out model.CategoryCreateOutput, err error) {
	lastInsertID, err := dao.CategoryInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.CategoryCreateOutput{CategoryId: int(lastInsertID)}, err
}

// Delete 删除
func (s *sCategory) Delete(ctx context.Context, id uint) (err error) {
	_, err = dao.CategoryInfo.Ctx(ctx).Where(g.Map{
		dao.CategoryInfo.Columns().Id: id,
	}).Delete()
	if err != nil {
		return err
	}
	return
}

// Update 修改
func (s *sCategory) Update(ctx context.Context, in model.CategoryUpdateInput) error {
	_, err := dao.CategoryInfo.
		Ctx(ctx).
		Data(in).
		FieldsEx(dao.CategoryInfo.Columns().Id).
		Where(dao.CategoryInfo.Columns().Id, in.Id).
		Update()
	return err
}

// GetList 查询分类列表
func (s *sCategory) GetList(ctx context.Context, in model.CategoryGetListInput) (out *model.CategoryGetListOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	m := dao.CategoryInfo.Ctx(ctx)

	// 如果指定了ParentId，则添加过滤条件
	if in.ParentId > 0 {
		m = m.Where(dao.CategoryInfo.Columns().ParentId, in.ParentId)
	}

	//2. 实例化响应结构体
	out = &model.CategoryGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.CategoryGetListOutputItem, 0, in.Size)
	//6. 把查询到的结果赋值到响应结构体中
	if err := listModel.Scan(&out.List); err != nil {
		return out, err
	}
	return
}

// GetList 查询分类全部信息-不翻页
func (s *sCategory) GetListAll(ctx context.Context, in model.CategoryGetListInput) (out *model.CategoryGetListOutput, err error) {
	var (
		m = dao.CategoryInfo.Ctx(ctx)
	)
	out = &model.CategoryGetListOutput{}

	listModel := m
	// 排序方式
	listModel = listModel.OrderDesc(dao.CategoryInfo.Columns().Sort)

	// 执行查询
	var list []*entity.CategoryInfo
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
	if err := listModel.Scan(&out.List); err != nil {
		return out, err
	}
	return
}

// GetHierarchicalList 获取分类层级列表
func (s *sCategory) GetHierarchicalList(ctx context.Context) (out interface{}, err error) {
	// 先获取所有分类
	var (
		m = dao.CategoryInfo.Ctx(ctx)
	)
	listModel := m.OrderDesc(dao.CategoryInfo.Columns().Sort)
	var list []*entity.CategoryInfo
	if err := listModel.Scan(&list); err != nil {
		return nil, err
	}
	// 构建层级结构
	// 这里可以根据业务需求实现分类的层级关系构建
	// 例如：将分类列表转换为树形结构
	return list, nil
}
