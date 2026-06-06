package frontend

import (
	"context"

	frontendApi "bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

var Category = cCategory{}

type cCategory struct{}

func (c *cCategory) List(ctx context.Context, req *frontendApi.CategoryListReq) (res *frontendApi.CategoryListRes, err error) {
	out, err := service.Category().List(ctx)
	if err != nil {
		return nil, err
	}
	return &frontendApi.CategoryListRes{
		List:  toFrontendApiCategories(out.List),
		Total: out.Total,
	}, nil
}

func toFrontendApiCategories(in []model.CategoryItem) []frontendApi.CategoryItem {
	list := make([]frontendApi.CategoryItem, 0, len(in))
	for _, item := range in {
		list = append(list, frontendApi.CategoryItem{
			Id:     item.Id,
			Name:   item.Name,
			Sort:   item.Sort,
			Status: item.Status,
		})
	}
	return list
}
