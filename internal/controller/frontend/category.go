package frontend

import (
	"context"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

// Category 前台商品分类控制器
var Category = cCategory{}

type cCategory struct{}

// List 前台分类列表
func (c *cCategory) List(ctx context.Context, req *frontend.CategoryGetListReq) (res *frontend.CategoryGetListRes, err error) {
	getListRes, err := service.Category().GetList(ctx, model.CategoryGetListInput{
		Page:     req.Page,
		Size:     req.Size,
		Sort:     req.Sort,
		ParentId: req.ParentId,
	})
	if err != nil {
		return nil, err
	}

	return &frontend.CategoryGetListRes{
		List:  getListRes.List,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total,
	}, nil
}

// ListAll 前台分类全部列表
func (c *cCategory) ListAll(ctx context.Context, req *frontend.CategoryGetListAllReq) (res *frontend.CategoryGetListAllRes, err error) {
	getListRes, err := service.Category().GetListAll(ctx, model.CategoryGetListInput{})
	if err != nil {
		return nil, err
	}

	return &frontend.CategoryGetListAllRes{
		List:  getListRes.List,
		Total: getListRes.Total,
	}, nil
}

// Hierarchical 前台分类层级列表
func (c *cCategory) Hierarchical(ctx context.Context, req *frontend.CategoryGetHierarchicalReq) (res *frontend.CategoryGetHierarchicalRes, err error) {
	getListRes, err := service.Category().GetListAll(ctx, model.CategoryGetListInput{})
	if err != nil {
		return nil, err
	}

	return &frontend.CategoryGetHierarchicalRes{
		List: getListRes.List,
	}, nil
}
