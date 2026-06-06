package backend

import (
	"context"

	backendApi "bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

var Category = cCategory{}

type cCategory struct{}

func (c *cCategory) List(ctx context.Context, req *backendApi.CategoryListReq) (res *backendApi.CategoryListRes, err error) {
	out, err := service.Category().List(ctx)
	if err != nil {
		return nil, err
	}
	return &backendApi.CategoryListRes{
		List:  toBackendApiCategories(out.List),
		Total: out.Total,
	}, nil
}

func toBackendApiCategories(in []model.CategoryItem) []backendApi.CategoryItem {
	list := make([]backendApi.CategoryItem, 0, len(in))
	for _, item := range in {
		list = append(list, backendApi.CategoryItem{
			Id:     item.Id,
			Name:   item.Name,
			Sort:   item.Sort,
			Status: item.Status,
		})
	}
	return list
}
