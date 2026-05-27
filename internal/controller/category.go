package controller

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

// 承上启下
// Category 内容管理
var Category = cCategory{}

type cCategory struct{}

func (a *cCategory) Create(ctx context.Context, req *backend.CategoryReq) (res *backend.CategoryRes, err error) {
	out, err := service.Category().Create(ctx, model.CategoryCreateInput{
		CategoryCreateUpdateBase: model.CategoryCreateUpdateBase{
			ParentId: req.ParentId,
			PicUrl:   req.PicUrl,
			Name:     req.Name,
			Sort:     req.Sort,
			Level:    req.Level,
		},
	})
	if err != nil {
		return nil, err
	}
	return &backend.CategoryRes{CategoryId: out.CategoryId}, nil
}

func (a *cCategory) Delete(ctx context.Context, req *backend.CategoryDeleteReq) (res *backend.CategoryDeleteRes, err error) {
	// 直接从请求 body 中解析 category_id，避免被 JWT 的 id 覆盖
	r := g.RequestFromCtx(ctx)
	var params struct {
		CategoryId uint `json:"category_id"`
	}
	_ = r.Parse(&params)
	categoryId := params.CategoryId
	g.Log().Infof(ctx, "Category Delete: params.CategoryId=%d, req.CategoryId=%d", params.CategoryId, req.CategoryId)
	if categoryId == 0 {
		categoryId = req.CategoryId
	}
	err = service.Category().Delete(ctx, categoryId)
	return
}

func (a *cCategory) Update(ctx context.Context, req *backend.CategoryUpdateReq) (res *backend.CategoryUpdateRes, err error) {
	// 直接从请求 body 中解析 category_id，避免被 JWT 的 id 覆盖
	r := g.RequestFromCtx(ctx)
	var params struct {
		CategoryId uint `json:"category_id"`
	}
	_ = r.Parse(&params)
	categoryId := params.CategoryId
	if categoryId == 0 {
		categoryId = req.CategoryId
	}
	err = service.Category().Update(ctx, model.CategoryUpdateInput{
		Id: categoryId,
		CategoryCreateUpdateBase: model.CategoryCreateUpdateBase{
			ParentId: req.ParentId,
			PicUrl:   req.PicUrl,
			Name:     req.Name,
			Sort:     req.Sort,
			Level:    req.Level,
		},
	})
	return &backend.CategoryUpdateRes{Id: categoryId}, nil
}

func (a *cCategory) List(ctx context.Context, req *backend.CategoryGetListCommonReq) (res *backend.CategoryGetListCommonRes, err error) {
	getListRes, err := service.Category().GetList(ctx, model.CategoryGetListInput{
		Page:     req.Page,
		Size:     req.Size,
		Sort:     req.Sort,
		ParentId: req.ParentId,
	})
	if err != nil {
		return nil, err
	}

	return &backend.CategoryGetListCommonRes{List: getListRes.List,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total}, nil
}

func (a *cCategory) ListAll(ctx context.Context, req *backend.CategoryGetListAllCommonReq) (res *backend.CategoryGetListAllCommonRes, err error) {
	getListRes, err := service.Category().GetListAll(ctx, model.CategoryGetListInput{})
	if err != nil {
		return nil, err
	}

	return &backend.CategoryGetListAllCommonRes{List: getListRes.List,
		Total: getListRes.Total}, nil
}
