package goods

import (
	"context"
	"database/sql"
	"errors"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sGoods struct{}

func init() {
	service.RegisterGoods(New())
}

func New() *sGoods {
	return &sGoods{}
}

func (s *sGoods) Create(ctx context.Context, in model.GoodsCreateInput) (out model.GoodsCreateOutput, err error) {
	if err = s.checkCategory(ctx, in.CategoryId); err != nil {
		return out, err
	}
	lastInsertId, err := dao.GoodsInfo.Ctx(ctx).Data(do.GoodsInfo{
		CategoryId:  in.CategoryId,
		Name:        in.Name,
		ImageUrl:    in.ImageUrl,
		PointsPrice: in.PointsPrice,
		Stock:       in.Stock,
		Description: in.Description,
		Status:      consts.GoodsStatusOnShelf,
	}).InsertAndGetId()
	if err != nil {
		return out, err
	}
	out.Id = uint(lastInsertId)
	return out, nil
}

func (s *sGoods) List(ctx context.Context, in model.GoodsListInput) (out model.GoodsListOutput, err error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 {
		in.Size = 10
	}
	if in.Size > 50 {
		in.Size = 50
	}

	columns := dao.GoodsInfo.Columns()
	m := dao.GoodsInfo.Ctx(ctx).WhereNull(columns.DeletedAt)
	if in.CategoryId > 0 {
		m = m.Where(columns.CategoryId, in.CategoryId)
	}
	if in.Name != "" {
		m = m.Where(columns.Name+" LIKE ?", "%"+in.Name+"%")
	}
	if in.Status == consts.GoodsStatusOffShelf || in.Status == consts.GoodsStatusOnShelf {
		m = m.Where(columns.Status, in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.GoodsListOutput{
		List:  make([]model.GoodsItem, 0),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}
	if total == 0 {
		return out, nil
	}

	var goodsList []entity.GoodsInfo
	err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&goodsList)
	if err != nil {
		return out, err
	}
	out.List, err = s.toGoodsItems(ctx, goodsList)
	return out, err
}

func (s *sGoods) Detail(ctx context.Context, id uint) (out model.GoodsDetailOutput, err error) {
	goods, err := s.getGoodsById(ctx, id)
	if err != nil {
		return out, err
	}
	if goods.Id == 0 {
		return out, errors.New("Goods does not exist")
	}
	items, err := s.toGoodsItems(ctx, []entity.GoodsInfo{goods})
	if err != nil {
		return out, err
	}
	out.Goods = items[0]
	return out, nil
}

func (s *sGoods) Update(ctx context.Context, in model.GoodsUpdateInput) (out model.GoodsUpdateOutput, err error) {
	goods, err := s.getGoodsById(ctx, in.Id)
	if err != nil {
		return out, err
	}
	if goods.Id == 0 {
		return out, errors.New("Goods does not exist")
	}
	if err = s.checkCategory(ctx, in.CategoryId); err != nil {
		return out, err
	}
	_, err = dao.GoodsInfo.Ctx(ctx).
		Where(dao.GoodsInfo.Columns().Id, in.Id).
		Data(g.Map{
			dao.GoodsInfo.Columns().CategoryId:  in.CategoryId,
			dao.GoodsInfo.Columns().Name:        in.Name,
			dao.GoodsInfo.Columns().ImageUrl:    in.ImageUrl,
			dao.GoodsInfo.Columns().PointsPrice: in.PointsPrice,
			dao.GoodsInfo.Columns().Stock:       in.Stock,
			dao.GoodsInfo.Columns().Description: in.Description,
		}).
		Update()
	if err != nil {
		return out, err
	}
	detail, err := s.Detail(ctx, in.Id)
	if err != nil {
		return out, err
	}
	out.Goods = detail.Goods
	return out, nil
}

func (s *sGoods) UpdateStatus(ctx context.Context, in model.GoodsStatusInput) error {
	if in.Status != consts.GoodsStatusOffShelf && in.Status != consts.GoodsStatusOnShelf {
		return errors.New("Status must be 0 or 1")
	}
	goods, err := s.getGoodsById(ctx, in.Id)
	if err != nil {
		return err
	}
	if goods.Id == 0 {
		return errors.New("Goods does not exist")
	}
	_, err = dao.GoodsInfo.Ctx(ctx).
		Where(dao.GoodsInfo.Columns().Id, in.Id).
		Data(g.Map{dao.GoodsInfo.Columns().Status: in.Status}).
		Update()
	return err
}

func (s *sGoods) FrontendList(ctx context.Context, in model.FrontendGoodsListInput) (out model.FrontendGoodsListOutput, err error) {
	goodsOut, err := s.List(ctx, model.GoodsListInput{
		Page:       in.Page,
		Size:       in.Size,
		CategoryId: in.CategoryId,
		Name:       in.Name,
		Status:     consts.GoodsStatusOnShelf,
	})
	if err != nil {
		return out, err
	}
	out = model.FrontendGoodsListOutput{
		List:  make([]model.FrontendGoodsListItem, 0, len(goodsOut.List)),
		Total: goodsOut.Total,
		Page:  goodsOut.Page,
		Size:  goodsOut.Size,
	}
	for _, item := range goodsOut.List {
		out.List = append(out.List, model.FrontendGoodsListItem{
			Id:           item.Id,
			CategoryId:   item.CategoryId,
			CategoryName: item.CategoryName,
			Name:         item.Name,
			ImageUrl:     item.ImageUrl,
			PointsPrice:  item.PointsPrice,
			Stock:        item.Stock,
		})
	}
	return out, nil
}

func (s *sGoods) FrontendDetail(ctx context.Context, id uint) (out model.FrontendGoodsDetailOutput, err error) {
	goods, err := s.getGoodsById(ctx, id)
	if err != nil {
		return out, err
	}
	if goods.Id == 0 || goods.Status != consts.GoodsStatusOnShelf {
		return out, errors.New("Goods does not exist or is off shelf")
	}
	items, err := s.toGoodsItems(ctx, []entity.GoodsInfo{goods})
	if err != nil {
		return out, err
	}
	item := items[0]
	out.Goods = model.FrontendGoodsDetailItem{
		Id:           item.Id,
		CategoryId:   item.CategoryId,
		CategoryName: item.CategoryName,
		Name:         item.Name,
		ImageUrl:     item.ImageUrl,
		PointsPrice:  item.PointsPrice,
		Stock:        item.Stock,
		Description:  item.Description,
	}
	return out, nil
}

func (s *sGoods) checkCategory(ctx context.Context, categoryId uint) error {
	columns := dao.GoodsCategory.Columns()
	count, err := dao.GoodsCategory.Ctx(ctx).
		Where(columns.Id, categoryId).
		Where(columns.Status, consts.GoodsCategoryStatusEnabled).
		WhereNull(columns.DeletedAt).
		Count()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Category does not exist or is disabled")
	}
	return nil
}

func (s *sGoods) getGoodsById(ctx context.Context, id uint) (goods entity.GoodsInfo, err error) {
	columns := dao.GoodsInfo.Columns()
	err = dao.GoodsInfo.Ctx(ctx).
		Where(columns.Id, id).
		WhereNull(columns.DeletedAt).
		Scan(&goods)
	if errors.Is(err, sql.ErrNoRows) {
		return goods, nil
	}
	return
}

func (s *sGoods) toGoodsItems(ctx context.Context, goodsList []entity.GoodsInfo) ([]model.GoodsItem, error) {
	categoryIds := make([]uint, 0, len(goodsList))
	for _, goods := range goodsList {
		categoryIds = append(categoryIds, goods.CategoryId)
	}

	categoryNames := make(map[uint]string)
	if len(categoryIds) > 0 {
		var categories []entity.GoodsCategory
		err := dao.GoodsCategory.Ctx(ctx).
			WhereIn(dao.GoodsCategory.Columns().Id, categoryIds).
			Scan(&categories)
		if err != nil {
			return nil, err
		}
		for _, category := range categories {
			categoryNames[category.Id] = category.Name
		}
	}

	items := make([]model.GoodsItem, 0, len(goodsList))
	for _, goods := range goodsList {
		items = append(items, model.GoodsItem{
			Id:           goods.Id,
			CategoryId:   goods.CategoryId,
			CategoryName: categoryNames[goods.CategoryId],
			Name:         goods.Name,
			ImageUrl:     goods.ImageUrl,
			PointsPrice:  goods.PointsPrice,
			Stock:        goods.Stock,
			Description:  goods.Description,
			Status:       goods.Status,
		})
	}
	return items, nil
}
