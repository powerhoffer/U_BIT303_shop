package controller

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"bit303_shop/api/backend"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
)

// 承上启下
// Goods 内容管理
var Goods = cGoods{}

type cGoods struct{}

func (a *cGoods) Create(ctx context.Context, req *backend.GoodsReq) (res *backend.GoodsRes, err error) {
	data := model.GoodsCreateInput{}
	err = gconv.Scan(req, &data)
	if err != nil {
		return nil, err
	}
	out, err := service.Goods().Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return &backend.GoodsRes{Id: out.Id}, nil
}

func (a *cGoods) Delete(ctx context.Context, req *backend.GoodsDeleteReq) (res *backend.GoodsDeleteRes, err error) {
	err = service.Goods().Delete(ctx, req.Id)
	return
}

func (a *cGoods) Update(ctx context.Context, req *backend.GoodsUpdateReq) (res *backend.GoodsUpdateRes, err error) {
	data := model.GoodsUpdateInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	err = service.Goods().Update(ctx, data)
	return &backend.GoodsUpdateRes{Id: req.Id}, nil
}

func (a *cGoods) List(ctx context.Context, req *backend.GoodsGetListCommonReq) (res *backend.GoodsGetListCommonRes, err error) {
	getListRes, err := service.Goods().GetList(ctx, model.GoodsGetListInput{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &backend.GoodsGetListCommonRes{List: getListRes.List,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total}, nil
}

func (*cGoods) Detail(ctx context.Context, req *frontend.GoodsDetailReq) (res *frontend.GoodsDetailRes, err error) {
	detail, err := service.Goods().Detail(ctx, model.GoodsDetailInput{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	res = &frontend.GoodsDetailRes{}
	err = gconv.Struct(detail, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *cGoods) GetLevelList(ctx context.Context, req *frontend.GoodsGetLevelListReq) (res *frontend.GoodsGetLevelListRes, err error) {
	getListRes, err := service.Goods().GetLevelList(ctx, model.GoodsGetLevelListInput{
		Page:    req.Page,
		Size:    req.Size,
		LevelId: req.LevelId,
	})
	if err != nil {
		return nil, err
	}
	
	// 转换类型
	list := make([]frontend.GoodsInfoBase, 0, len(getListRes.List))
	for _, item := range getListRes.List {
		goodsBase := frontend.GoodsInfoBase{
			Id:               item.Id,
			PicUrl:           item.PicUrl,
			Name:             item.Name,
			Price:            item.Price,
			Level1CategoryId: item.Level1CategoryId,
			Level2CategoryId: item.Level2CategoryId,
			Level3CategoryId: item.Level3CategoryId,
			Brand:            item.Brand,
			Stock:            item.Stock,
			Sale:             item.Sale,
			Tags:             item.Tags,
			DetailInfo:       item.DetailInfo,
			CreatedAt:        item.CreatedAt,
		}
		list = append(list, goodsBase)
	}
	
	return &frontend.GoodsGetLevelListRes{
		List:  list,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total,
	}, nil
}
