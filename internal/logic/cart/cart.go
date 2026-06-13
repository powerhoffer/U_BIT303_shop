package cart

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
	"github.com/gogf/gf/v2/os/gtime"
)

type sCart struct{}

func init() {
	service.RegisterCart(New())
}

func New() *sCart {
	return &sCart{}
}

func (s *sCart) Add(ctx context.Context, in model.CartAddInput) (out model.CartAddOutput, err error) {
	if in.Count <= 0 {
		return out, errors.New("Goods count must be greater than 0")
	}
	if _, err = s.getAvailableGoods(ctx, in.GoodsId); err != nil {
		return out, err
	}

	cartColumns := dao.CartInfo.Columns()
	var cart entity.CartInfo
	err = dao.CartInfo.Ctx(ctx).
		Where(cartColumns.EmployeeId, in.EmployeeId).
		Where(cartColumns.GoodsId, in.GoodsId).
		WhereNull(cartColumns.DeletedAt).
		Scan(&cart)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return out, err
	}
	if cart.Id == 0 {
		id, err := dao.CartInfo.Ctx(ctx).Data(do.CartInfo{
			EmployeeId: in.EmployeeId,
			GoodsId:    in.GoodsId,
			Count:      in.Count,
		}).InsertAndGetId()
		if err != nil {
			return out, err
		}
		out.Id = uint(id)
		return out, nil
	}

	_, err = dao.CartInfo.Ctx(ctx).
		Where(cartColumns.Id, cart.Id).
		Data(g.Map{cartColumns.Count: cart.Count + uint(in.Count)}).
		Update()
	if err != nil {
		return out, err
	}
	out.Id = cart.Id
	return out, nil
}

func (s *sCart) List(ctx context.Context, in model.CartListInput) (out model.CartListOutput, err error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 {
		in.Size = 10
	}
	if in.Size > 50 {
		in.Size = 50
	}

	cartColumns := dao.CartInfo.Columns()
	m := dao.CartInfo.Ctx(ctx).
		Where(cartColumns.EmployeeId, in.EmployeeId).
		WhereNull(cartColumns.DeletedAt)

	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.CartListOutput{
		List:  make([]model.CartItem, 0),
		Total: total,
		Page:  in.Page,
		Size:  in.Size,
	}
	if total == 0 {
		return out, nil
	}

	var carts []entity.CartInfo
	err = m.Page(in.Page, in.Size).OrderDesc(cartColumns.Id).Scan(&carts)
	if err != nil {
		return out, err
	}
	out.List, err = s.toCartItems(ctx, carts)
	return out, err
}

func (s *sCart) Update(ctx context.Context, in model.CartUpdateInput) (out model.CartUpdateOutput, err error) {
	if in.Count < 0 {
		return out, errors.New("Goods count cannot be negative")
	}
	cart, err := s.getEmployeeCart(ctx, in.EmployeeId, in.Id)
	if err != nil {
		return out, err
	}
	if in.Count == 0 {
		err = s.softDelete(ctx, in.EmployeeId, in.Id)
		return model.CartUpdateOutput{Id: cart.Id}, err
	}
	_, err = s.getAvailableGoods(ctx, cart.GoodsId)
	if err != nil {
		return out, err
	}
	_, err = dao.CartInfo.Ctx(ctx).
		Where(dao.CartInfo.Columns().Id, cart.Id).
		Data(g.Map{dao.CartInfo.Columns().Count: in.Count}).
		Update()
	if err != nil {
		return out, err
	}
	out.Id = cart.Id
	return out, nil
}

func (s *sCart) Remove(ctx context.Context, in model.CartRemoveInput) (out model.CartRemoveOutput, err error) {
	err = s.softDelete(ctx, in.EmployeeId, in.Id)
	if err != nil {
		return out, err
	}
	out.Id = in.Id
	return out, nil
}

func (s *sCart) getAvailableGoods(ctx context.Context, goodsId uint) (goods entity.GoodsInfo, err error) {
	columns := dao.GoodsInfo.Columns()
	err = dao.GoodsInfo.Ctx(ctx).
		Where(columns.Id, goodsId).
		Where(columns.Status, consts.GoodsStatusOnShelf).
		WhereNull(columns.DeletedAt).
		Scan(&goods)
	if errors.Is(err, sql.ErrNoRows) {
		return goods, errors.New("Goods does not exist or is off shelf")
	}
	if err != nil {
		return goods, err
	}
	if goods.Id == 0 {
		return goods, errors.New("Goods does not exist or is off shelf")
	}
	return goods, nil
}

func (s *sCart) getEmployeeCart(ctx context.Context, employeeId, id uint) (cart entity.CartInfo, err error) {
	columns := dao.CartInfo.Columns()
	err = dao.CartInfo.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.EmployeeId, employeeId).
		WhereNull(columns.DeletedAt).
		Scan(&cart)
	if errors.Is(err, sql.ErrNoRows) {
		return cart, errors.New("Cart item does not exist")
	}
	if err != nil {
		return cart, err
	}
	if cart.Id == 0 {
		return cart, errors.New("Cart item does not exist")
	}
	return cart, nil
}

func (s *sCart) softDelete(ctx context.Context, employeeId, id uint) error {
	if _, err := s.getEmployeeCart(ctx, employeeId, id); err != nil {
		return err
	}
	_, err := dao.CartInfo.Ctx(ctx).
		Where(dao.CartInfo.Columns().Id, id).
		Where(dao.CartInfo.Columns().EmployeeId, employeeId).
		Data(g.Map{dao.CartInfo.Columns().DeletedAt: gtime.Now()}).
		Update()
	return err
}

func (s *sCart) toCartItems(ctx context.Context, carts []entity.CartInfo) ([]model.CartItem, error) {
	goodsIds := make([]uint, 0, len(carts))
	for _, cart := range carts {
		goodsIds = append(goodsIds, cart.GoodsId)
	}

	goodsMap := make(map[uint]entity.GoodsInfo)
	categoryIds := make([]uint, 0)
	if len(goodsIds) > 0 {
		var goodsList []entity.GoodsInfo
		err := dao.GoodsInfo.Ctx(ctx).
			WhereIn(dao.GoodsInfo.Columns().Id, goodsIds).
			WhereNull(dao.GoodsInfo.Columns().DeletedAt).
			Scan(&goodsList)
		if err != nil {
			return nil, err
		}
		for _, goods := range goodsList {
			goodsMap[goods.Id] = goods
			categoryIds = append(categoryIds, goods.CategoryId)
		}
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

	items := make([]model.CartItem, 0, len(carts))
	for _, cart := range carts {
		goods := goodsMap[cart.GoodsId]
		items = append(items, model.CartItem{
			Id:           cart.Id,
			GoodsId:      cart.GoodsId,
			CategoryId:   goods.CategoryId,
			CategoryName: categoryNames[goods.CategoryId],
			GoodsName:    goods.Name,
			ImageUrl:     goods.ImageUrl,
			PointsPrice:  goods.PointsPrice,
			Stock:        goods.Stock,
			Count:        cart.Count,
			TotalPoints:  goods.PointsPrice * cart.Count,
		})
	}
	return items, nil
}
