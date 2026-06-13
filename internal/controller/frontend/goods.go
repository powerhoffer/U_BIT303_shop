package frontend

import (
	"context"

	frontendApi "bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

var Goods = cGoods{}

type cGoods struct{}

func (c *cGoods) List(ctx context.Context, req *frontendApi.GoodsListReq) (res *frontendApi.GoodsListRes, err error) {
	var input model.FrontendGoodsListInput
	if err = gconv.Struct(req, &input); err != nil {
		return nil, err
	}
	out, err := service.Goods().FrontendList(ctx, input)
	if err != nil {
		return nil, err
	}
	return &frontendApi.GoodsListRes{
		List:  toFrontendApiGoodsListItems(out.List),
		Total: out.Total,
		Page:  out.Page,
		Size:  out.Size,
	}, nil
}

func (c *cGoods) Detail(ctx context.Context, req *frontendApi.GoodsDetailReq) (res *frontendApi.GoodsDetailRes, err error) {
	out, err := service.Goods().FrontendDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &frontendApi.GoodsDetailRes{Goods: toFrontendApiGoodsDetailItem(out.Goods)}, nil
}

func toFrontendApiGoodsListItems(in []model.FrontendGoodsListItem) []frontendApi.GoodsListItem {
	list := make([]frontendApi.GoodsListItem, 0, len(in))
	for _, item := range in {
		list = append(list, frontendApi.GoodsListItem{
			Id:           item.Id,
			CategoryId:   item.CategoryId,
			CategoryName: item.CategoryName,
			Name:         item.Name,
			ImageUrl:     item.ImageUrl,
			PointsPrice:  item.PointsPrice,
			Stock:        item.Stock,
		})
	}
	return list
}

func toFrontendApiGoodsDetailItem(in model.FrontendGoodsDetailItem) frontendApi.GoodsDetailItem {
	return frontendApi.GoodsDetailItem{
		Id:           in.Id,
		CategoryId:   in.CategoryId,
		CategoryName: in.CategoryName,
		Name:         in.Name,
		ImageUrl:     in.ImageUrl,
		PointsPrice:  in.PointsPrice,
		Stock:        in.Stock,
		Description:  in.Description,
	}
}
