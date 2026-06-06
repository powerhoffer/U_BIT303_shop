package category

import (
	"context"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
)

type sCategory struct{}

func init() {
	service.RegisterCategory(New())
}

func New() *sCategory {
	return &sCategory{}
}

func (s *sCategory) List(ctx context.Context) (out model.CategoryListOutput, err error) {
	columns := dao.GoodsCategory.Columns()
	var categories []entity.GoodsCategory
	err = dao.GoodsCategory.Ctx(ctx).
		Where(columns.Status, consts.GoodsCategoryStatusEnabled).
		WhereNull(columns.DeletedAt).
		OrderAsc(columns.Sort).
		OrderAsc(columns.Id).
		Scan(&categories)
	if err != nil {
		return out, err
	}
	out.List = make([]model.CategoryItem, 0, len(categories))
	for _, category := range categories {
		out.List = append(out.List, model.CategoryItem{
			Id:     category.Id,
			Name:   category.Name,
			Sort:   category.Sort,
			Status: category.Status,
		})
	}
	out.Total = len(out.List)
	return out, nil
}
