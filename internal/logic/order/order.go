package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sOrder struct{}

func init() {
	service.RegisterOrder(New())
}

func New() *sOrder {
	return &sOrder{}
}

func (s *sOrder) Create(ctx context.Context, in model.OrderCreateInput) (out model.OrderCreateOutput, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cartColumns := dao.CartInfo.Columns()
		var carts []entity.CartInfo
		if err := dao.CartInfo.Ctx(ctx).
			Where(cartColumns.EmployeeId, in.EmployeeId).
			WhereNull(cartColumns.DeletedAt).
			LockUpdate().
			OrderAsc(cartColumns.Id).
			Scan(&carts); err != nil {
			return err
		}
		if len(carts) == 0 {
			return errors.New("Cart is empty")
		}

		goodsList := make([]entity.GoodsInfo, 0, len(carts))
		totalPoints := uint(0)
		goodsColumns := dao.GoodsInfo.Columns()
		for _, cart := range carts {
			var goods entity.GoodsInfo
			if err := dao.GoodsInfo.Ctx(ctx).
				Where(goodsColumns.Id, cart.GoodsId).
				WhereNull(goodsColumns.DeletedAt).
				LockUpdate().
				Scan(&goods); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			if goods.Id == 0 || goods.Status != consts.GoodsStatusOnShelf {
				return errors.New("Goods does not exist or is off shelf")
			}
			if goods.Stock < cart.Count {
				return errors.New("Insufficient goods stock")
			}
			goodsList = append(goodsList, goods)
			totalPoints += goods.PointsPrice * cart.Count
		}
		if totalPoints == 0 {
			return errors.New("Order points must be greater than 0")
		}

		account, err := s.getPointsAccountForUpdate(ctx, in.EmployeeId)
		if err != nil {
			return err
		}
		if account.Id == 0 || account.Status != consts.PointsAccountStatusNormal {
			return errors.New("Credit account does not exist or is disabled")
		}
		if account.Balance < totalPoints {
			return errors.New("Insufficient credit balance")
		}

		orderNo := generateOrderNo()
		orderId, err := dao.OrderInfo.Ctx(ctx).Data(do.OrderInfo{
			OrderNo:     orderNo,
			EmployeeId:  in.EmployeeId,
			TotalPoints: totalPoints,
			Status:      consts.OrderStatusPending,
			Remark:      in.Remark,
		}).InsertAndGetId()
		if err != nil {
			return err
		}

		goodsById := make(map[uint]entity.GoodsInfo, len(goodsList))
		for _, goods := range goodsList {
			goodsById[goods.Id] = goods
		}
		for _, cart := range carts {
			goods := goodsById[cart.GoodsId]
			if _, err = dao.OrderItem.Ctx(ctx).Data(do.OrderItem{
				OrderId:       uint(orderId),
				EmployeeId:    in.EmployeeId,
				GoodsId:       goods.Id,
				GoodsName:     goods.Name,
				GoodsImageUrl: goods.ImageUrl,
				PointsPrice:   goods.PointsPrice,
				Count:         cart.Count,
				TotalPoints:   goods.PointsPrice * cart.Count,
			}).Insert(); err != nil {
				return err
			}
			if _, err = dao.GoodsInfo.Ctx(ctx).
				Where(goodsColumns.Id, goods.Id).
				Data(g.Map{goodsColumns.Stock: goods.Stock - cart.Count}).
				Update(); err != nil {
				return err
			}
		}

		beforeBalance := account.Balance
		afterBalance := beforeBalance - totalPoints
		if _, err = dao.EmployeePointsAccount.Ctx(ctx).
			Where(dao.EmployeePointsAccount.Columns().Id, account.Id).
			Data(g.Map{dao.EmployeePointsAccount.Columns().Balance: afterBalance}).
			Update(); err != nil {
			return err
		}
		if _, err = dao.EmployeePointsRecord.Ctx(ctx).Data(do.EmployeePointsRecord{
			EmployeeId:         in.EmployeeId,
			ChangeType:         consts.PointsChangeTypeDeduct,
			Points:             totalPoints,
			BeforeBalance:      beforeBalance,
			AfterBalance:       afterBalance,
			OperatorEmployeeId: in.EmployeeId,
			Remark:             "Order redemption: " + orderNo,
		}).Insert(); err != nil {
			return err
		}

		cartIds := make([]uint, 0, len(carts))
		for _, cart := range carts {
			cartIds = append(cartIds, cart.Id)
		}
		if _, err = dao.CartInfo.Ctx(ctx).
			WhereIn(cartColumns.Id, cartIds).
			Where(cartColumns.EmployeeId, in.EmployeeId).
			Data(g.Map{cartColumns.DeletedAt: gtime.Now()}).
			Update(); err != nil {
			return err
		}

		out.Order = model.OrderBase{
			Id:          uint(orderId),
			OrderNo:     orderNo,
			EmployeeId:  in.EmployeeId,
			TotalPoints: totalPoints,
			Status:      consts.OrderStatusPending,
			Remark:      in.Remark,
			CreatedAt:   time.Now(),
		}
		return nil
	})
	return out, err
}

func (s *sOrder) List(ctx context.Context, in model.OrderListInput) (out model.OrderListOutput, err error) {
	if in.Page < 1 {
		in.Page = 1
	}
	if in.Size < 1 {
		in.Size = 10
	}
	if in.Size > 50 {
		in.Size = 50
	}
	columns := dao.OrderInfo.Columns()
	m := dao.OrderInfo.Ctx(ctx).
		Where(columns.EmployeeId, in.EmployeeId).
		WhereNull(columns.DeletedAt)
	total, err := m.Count()
	if err != nil {
		return out, err
	}
	out = model.OrderListOutput{List: make([]model.OrderBase, 0), Total: total, Page: in.Page, Size: in.Size}
	if total == 0 {
		return out, nil
	}
	var orders []entity.OrderInfo
	if err = m.Page(in.Page, in.Size).OrderDesc(columns.Id).Scan(&orders); err != nil {
		return out, err
	}
	for _, order := range orders {
		out.List = append(out.List, toOrderBase(order))
	}
	return out, nil
}

func (s *sOrder) Detail(ctx context.Context, in model.OrderDetailInput) (out model.OrderDetailOutput, err error) {
	order, err := s.getEmployeeOrder(ctx, in.EmployeeId, in.Id)
	if err != nil {
		return out, err
	}
	items, err := s.getOrderItems(ctx, in.EmployeeId, in.Id)
	if err != nil {
		return out, err
	}
	out.Order = model.OrderDetail{OrderBase: toOrderBase(order), Items: items}
	return out, nil
}

func (s *sOrder) Cancel(ctx context.Context, in model.OrderCancelInput) (out model.OrderCancelOutput, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		order, err := s.getEmployeeOrderForUpdate(ctx, in.EmployeeId, in.Id)
		if err != nil {
			return err
		}
		if order.Status != consts.OrderStatusPending {
			return errors.New("Only pending orders can be cancelled")
		}
		items, err := s.getOrderItems(ctx, in.EmployeeId, in.Id)
		if err != nil {
			return err
		}
		goodsColumns := dao.GoodsInfo.Columns()
		for _, item := range items {
			var goods entity.GoodsInfo
			if err = dao.GoodsInfo.Ctx(ctx).
				Where(goodsColumns.Id, item.GoodsId).
				WhereNull(goodsColumns.DeletedAt).
				LockUpdate().
				Scan(&goods); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			if goods.Id > 0 {
				if _, err = dao.GoodsInfo.Ctx(ctx).
					Where(goodsColumns.Id, goods.Id).
					Data(g.Map{goodsColumns.Stock: goods.Stock + item.Count}).
					Update(); err != nil {
					return err
				}
			}
		}
		account, err := s.getPointsAccountForUpdate(ctx, in.EmployeeId)
		if err != nil {
			return err
		}
		if account.Id == 0 {
			return errors.New("Credit account does not exist")
		}
		beforeBalance := account.Balance
		afterBalance := beforeBalance + order.TotalPoints
		if _, err = dao.EmployeePointsAccount.Ctx(ctx).
			Where(dao.EmployeePointsAccount.Columns().Id, account.Id).
			Data(g.Map{dao.EmployeePointsAccount.Columns().Balance: afterBalance}).
			Update(); err != nil {
			return err
		}
		if _, err = dao.EmployeePointsRecord.Ctx(ctx).Data(do.EmployeePointsRecord{
			EmployeeId:         in.EmployeeId,
			ChangeType:         consts.PointsChangeTypeAdd,
			Points:             order.TotalPoints,
			BeforeBalance:      beforeBalance,
			AfterBalance:       afterBalance,
			OperatorEmployeeId: in.EmployeeId,
			Remark:             "Cancel order refund: " + order.OrderNo,
		}).Insert(); err != nil {
			return err
		}
		if _, err = dao.OrderInfo.Ctx(ctx).
			Where(dao.OrderInfo.Columns().Id, order.Id).
			Data(g.Map{dao.OrderInfo.Columns().Status: consts.OrderStatusCancelled}).
			Update(); err != nil {
			return err
		}
		order.Status = consts.OrderStatusCancelled
		out.Order = toOrderBase(order)
		return nil
	})
	return out, err
}

func (s *sOrder) getEmployeeOrder(ctx context.Context, employeeId, id uint) (order entity.OrderInfo, err error) {
	columns := dao.OrderInfo.Columns()
	err = dao.OrderInfo.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.EmployeeId, employeeId).
		WhereNull(columns.DeletedAt).
		Scan(&order)
	if errors.Is(err, sql.ErrNoRows) {
		return order, errors.New("Order does not exist")
	}
	if err != nil {
		return order, err
	}
	if order.Id == 0 {
		return order, errors.New("Order does not exist")
	}
	return order, nil
}

func (s *sOrder) getEmployeeOrderForUpdate(ctx context.Context, employeeId, id uint) (order entity.OrderInfo, err error) {
	columns := dao.OrderInfo.Columns()
	err = dao.OrderInfo.Ctx(ctx).
		Where(columns.Id, id).
		Where(columns.EmployeeId, employeeId).
		WhereNull(columns.DeletedAt).
		LockUpdate().
		Scan(&order)
	if errors.Is(err, sql.ErrNoRows) {
		return order, errors.New("Order does not exist")
	}
	if err != nil {
		return order, err
	}
	if order.Id == 0 {
		return order, errors.New("Order does not exist")
	}
	return order, nil
}

func (s *sOrder) getOrderItems(ctx context.Context, employeeId, orderId uint) ([]model.OrderGoodsItem, error) {
	columns := dao.OrderItem.Columns()
	var items []entity.OrderItem
	if err := dao.OrderItem.Ctx(ctx).
		Where(columns.OrderId, orderId).
		Where(columns.EmployeeId, employeeId).
		OrderAsc(columns.Id).
		Scan(&items); err != nil {
		return nil, err
	}
	out := make([]model.OrderGoodsItem, 0, len(items))
	for _, item := range items {
		out = append(out, toOrderGoodsItem(item))
	}
	return out, nil
}

func (s *sOrder) getPointsAccountForUpdate(ctx context.Context, employeeId uint) (account entity.EmployeePointsAccount, err error) {
	columns := dao.EmployeePointsAccount.Columns()
	err = dao.EmployeePointsAccount.Ctx(ctx).
		Where(columns.EmployeeId, employeeId).
		WhereNull(columns.DeletedAt).
		LockUpdate().
		Scan(&account)
	if errors.Is(err, sql.ErrNoRows) {
		return account, nil
	}
	return account, err
}

func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("OR%s%09d", now.Format("20060102150405"), now.Nanosecond())
}

func toOrderBase(order entity.OrderInfo) model.OrderBase {
	item := model.OrderBase{
		Id:          order.Id,
		OrderNo:     order.OrderNo,
		EmployeeId:  order.EmployeeId,
		TotalPoints: order.TotalPoints,
		Status:      order.Status,
		Remark:      order.Remark,
	}
	if order.CreatedAt != nil {
		item.CreatedAt = order.CreatedAt.Time
	}
	return item
}

func toOrderGoodsItem(item entity.OrderItem) model.OrderGoodsItem {
	out := model.OrderGoodsItem{
		Id:            item.Id,
		OrderId:       item.OrderId,
		EmployeeId:    item.EmployeeId,
		GoodsId:       item.GoodsId,
		GoodsName:     item.GoodsName,
		GoodsImageUrl: item.GoodsImageUrl,
		PointsPrice:   item.PointsPrice,
		Count:         item.Count,
		TotalPoints:   item.TotalPoints,
	}
	if item.CreatedAt != nil {
		out.CreatedAt = item.CreatedAt.Time
	}
	return out
}
