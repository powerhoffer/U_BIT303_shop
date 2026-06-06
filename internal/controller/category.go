package controller

import (
	"context"

	"bit303_shop/api/category"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

var Category = cCategory{}

type cCategory struct{}

func (c *cCategory) List(ctx context.Context, req *category.ListReq) (res *category.ListRes, err error) {
	out, err := service.Category().List(ctx)
	if err != nil {
		return nil, err
	}
	return &category.ListRes{
		List:  toApiCategories(out.List),
		Total: out.Total,
	}, nil
}

func toApiCategories(in []model.CategoryItem) []category.CategoryItem {
	list := make([]category.CategoryItem, 0, len(in))
	for _, item := range in {
		list = append(list, category.CategoryItem{
			Id:     item.Id,
			Name:   item.Name,
			Sort:   item.Sort,
			Status: item.Status,
		})
	}
	return list
}
