package order

import (
	"context"
	"encoding/json"
	"fmt"
	"bit303_shop/api/backend"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
	"bit303_shop/utility"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sOrder struct{}

func init() {
	service.RegisterOrder(New())
}

func New() *sOrder {
	return &sOrder{}
}

// 下单
func (s *sOrder) Add(ctx context.Context, in model.OrderAddInput) (out *model.OrderAddOutput, err error) {
	in.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))
	in.Number = utility.GetOrderNum()
	out = &model.OrderAddOutput{}
	//官方建议的事务闭包处理
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		//生成主订单
		lastInsertId, err := dao.OrderInfo.Ctx(ctx).InsertAndGetId(in)
		if err != nil {
			return err
		}
		//生成商品订单
		for _, info := range in.OrderAddGoodsInfos {
			info.OrderId = gconv.Int(lastInsertId)
			_, err := dao.OrderGoodsInfo.Ctx(ctx).Insert(info)
			if err != nil {
				return err
			}
		}
		//更新商品销量和库存，todo 后期接入消息
		for _, info := range in.OrderAddGoodsInfos {
			//商品增加销量
			_, err := dao.GoodsInfo.Ctx(ctx).WherePri(info.GoodsId).Increment(dao.GoodsInfo.Columns().Sale, info.Count)
			if err != nil {
				return err
			}
			//商品减少库存
			_, err2 := dao.GoodsInfo.Ctx(ctx).WherePri(info.GoodsId).Decrement(dao.GoodsInfo.Columns().Stock, info.Count)
			if err2 != nil {
				return err
			}
			//商品规格减少库存
			_, err3 := dao.GoodsOptionsInfo.Ctx(ctx).WherePri(info.GoodsOptionsId).Decrement(dao.GoodsOptionsInfo.Columns().Stock, info.Count)
			if err3 != nil {
				return err
			}
		}
		out.Id = uint(lastInsertId)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return
}

func (s *sOrder) List(ctx context.Context, in model.OrderListInput) (out *model.OrderListOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	whereCondition := s.orderListCondition(in)
	m := dao.OrderInfo.Ctx(ctx).Where(whereCondition)
	//2. 实例化响应结构体
	out = &model.OrderListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.OrderListOutputItem, 0, in.Size)
	//6. 把查询到的结果赋值到响应结构体中
	if err := listModel.Scan(&out.List); err != nil {
		return out, err
	}
	return
}

// todo 优化这里的代码
func (s *sOrder) orderListCondition(in model.OrderListInput) *gmap.Map {
	m := gmap.New()

	if in.Number != "" {
		m.Set(dao.OrderInfo.Columns().Number+" like ", "%"+in.Number+"%")
	}

	if in.UserId != 0 {
		m.Set(dao.OrderInfo.Columns().UserId, in.UserId)
	}

	if in.PayType != 0 {
		m.Set(dao.OrderInfo.Columns().PayType, in.PayType)
	}

	if in.PayAtGte != "" {
		m.Set(dao.OrderInfo.Columns().PayAt+" >= ", gtime.New(in.PayAtGte).StartOfDay())
	}

	if in.PayAtLte != "" {
		m.Set(dao.OrderInfo.Columns().PayAt+" <= ", gtime.New(in.PayAtLte).EndOfDay())
	}

	if in.Status != 0 {
		m.Set(dao.OrderInfo.Columns().Status, in.Status)
	}

	if in.ConsigneePhone != "" {
		m.Set(dao.OrderInfo.Columns().ConsigneePhone+" like ", "%"+in.ConsigneePhone+"%")
	}

	if in.PriceGte != 0 {
		m.Set(dao.OrderInfo.Columns().Price+" >= ", in.PriceGte)
	}

	if in.PriceLte != 0 {
		m.Set(dao.OrderInfo.Columns().Price+" <= ", in.PriceLte)
	}

	if in.DateGte != "" {
		m.Set(dao.OrderInfo.Columns().CreatedAt+" >= ", gtime.New(in.DateGte).StartOfDay())
	}

	if in.DateLte != "" {
		m.Set(dao.OrderInfo.Columns().CreatedAt+" <= ", gtime.New(in.DateLte).EndOfDay())
	}

	return m
}

func (s *sOrder) Detail(ctx context.Context, in model.OrderDetailInput) (out *model.OrderDetailOutput, err error) {
	err = dao.OrderInfo.Ctx(ctx).WithAll().WherePri(in.Id).Scan(&out)

	return
}

// 为测试代码实现的方法

// Pay 支付订单
func (s *sOrder) Pay(ctx context.Context, in *frontend.OrderPayReq) (out *frontend.OrderPayRes, err error) {
	// 验证订单
	order, err := s.getOrderById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 检查订单状态
	statusValue := gconv.Int(order.Status)
	if statusValue != 1 { // 1表示待支付状态
		return nil, gerror.New("订单状态不允许支付")
	}

	// 更新订单状态
	_, err = dao.OrderInfo.Ctx(ctx).Data(g.Map{
		dao.OrderInfo.Columns().Status:  2, // 已支付待发货
		dao.OrderInfo.Columns().PayType: in.PayType,
		dao.OrderInfo.Columns().PayAt:   gtime.Now(),
	}).Where(dao.OrderInfo.Columns().Id, in.Id).Update()

	if err != nil {
		return nil, err
	}

	return &frontend.OrderPayRes{
		PayUrl: "https://example.com/pay", // 模拟支付URL，实际项目中应返回真实支付链接
	}, nil
}

// Cancel 取消订单
func (s *sOrder) Cancel(ctx context.Context, in *frontend.OrderCancelReq) (out *frontend.OrderCancelRes, err error) {
	// 验证订单
	order, err := s.getOrderById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 检查订单状态
	statusValue := gconv.Int(order.Status)
	if statusValue != 1 { // 1表示待支付状态
		return nil, gerror.New("订单状态不允许取消")
	}

	// 更新订单状态
	_, err = dao.OrderInfo.Ctx(ctx).Data(g.Map{
		dao.OrderInfo.Columns().Status: 6, // 已取消
	}).Where(dao.OrderInfo.Columns().Id, in.Id).Update()

	if err != nil {
		return nil, err
	}

	// 恢复商品库存
	err = s.restoreGoodsStock(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &frontend.OrderCancelRes{
		Id: in.Id,
	}, nil
}

// Confirm 确认收货
func (s *sOrder) Confirm(ctx context.Context, in *frontend.OrderConfirmReq) (out *frontend.OrderConfirmRes, err error) {
	// 验证订单
	order, err := s.getOrderById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 检查订单状态
	statusValue := gconv.Int(order.Status)
	if statusValue != 3 { // 3表示已发货状态
		return nil, gerror.New("订单状态不允许确认收货")
	}

	// 更新订单状态
	_, err = dao.OrderInfo.Ctx(ctx).Data(g.Map{
		dao.OrderInfo.Columns().Status: 4, // 已收货待评价
	}).Where(dao.OrderInfo.Columns().Id, in.Id).Update()

	if err != nil {
		return nil, err
	}

	return &frontend.OrderConfirmRes{
		Id: in.Id,
	}, nil
}

// getOrderById 获取订单信息
func (s *sOrder) getOrderById(ctx context.Context, id uint) (order *model.OrderDetailOutput, err error) {
	err = dao.OrderInfo.Ctx(ctx).WithAll().WherePri(id).Scan(&order)
	return
}

// UpdateOrderStatus 更新订单状态(后台)
func (s *sOrder) UpdateOrderStatus(ctx context.Context, in *backend.OrderUpdateStatusReq) (out *backend.OrderUpdateStatusRes, err error) {
	// 验证订单
	order, err := s.getOrderById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 验证状态值
	if in.Status < 1 || in.Status > 5 {
		return nil, gerror.New("无效的订单状态")
	}

	// 更新订单状态
	_, err = dao.OrderInfo.Ctx(ctx).Data(g.Map{
		dao.OrderInfo.Columns().Status: in.Status,
	}).Where(dao.OrderInfo.Columns().Id, in.Id).Update()

	if err != nil {
		return nil, err
	}

	return &backend.OrderUpdateStatusRes{
		Id: in.Id,
	}, nil
}

// Delete 删除订单(后台)
func (s *sOrder) Delete(ctx context.Context, in *backend.OrderDeleteReq) (out *backend.OrderDeleteRes, err error) {
	// 验证订单
	order, err := s.getOrderById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 使用事务处理
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除订单相关商品
		_, err := dao.OrderGoodsInfo.Ctx(ctx).Where(dao.OrderGoodsInfo.Columns().OrderId, in.Id).Delete()
		if err != nil {
			return err
		}

		// 删除主订单
		_, err = dao.OrderInfo.Ctx(ctx).Where(dao.OrderInfo.Columns().Id, in.Id).Delete()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &backend.OrderDeleteRes{
		Id: in.Id,
	}, nil
}

// Refund 订单退款(后台)
func (s *sOrder) Refund(ctx context.Context, in *backend.OrderRefundReq) (out *backend.OrderRefundRes, err error) {
	// 验证订单
	order, err := s.getOrderById(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}

	// 验证订单状态，只有已支付的订单才能退款
	if order.Status != 2 {
		return nil, gerror.New("只有已支付的订单才能退款")
	}

	// 使用事务处理
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新订单状态为已退款
		_, err := dao.OrderInfo.Ctx(ctx).Data(g.Map{
			dao.OrderInfo.Columns().Status: 6, // 已退款状态
			dao.OrderInfo.Columns().Remark: in.Reason,
		}).Where(dao.OrderInfo.Columns().Id, in.Id).Update()

		if err != nil {
			return err
		}

		// 获取订单商品
		var orderGoods []*entity.OrderGoodsInfo
		err = dao.OrderGoodsInfo.Ctx(ctx).Where(dao.OrderGoodsInfo.Columns().OrderId, in.Id).Scan(&orderGoods)
		if err != nil {
			return err
		}

		// 恢复商品库存
		for _, goods := range orderGoods {
			// 商品库存恢复
			_, err := dao.GoodsInfo.Ctx(ctx).WherePri(goods.GoodsId).Increment(dao.GoodsInfo.Columns().Stock, goods.Count)
			if err != nil {
				return err
			}

			// 商品规格库存恢复
			if goods.GoodsOptionsId > 0 {
				_, err := dao.GoodsOptionsInfo.Ctx(ctx).WherePri(goods.GoodsOptionsId).Increment(dao.GoodsOptionsInfo.Columns().Stock, goods.Count)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &backend.OrderRefundRes{
		Id: in.Id,
	}, nil
}

// CreateSeckillOrder 创建秒杀订单
func (s *sOrder) CreateSeckillOrder(ctx context.Context, in model.SeckillOrderInput) (out *model.OrderAddOutput, err error) {
	g.Log().Info(ctx, fmt.Sprintf("开始创建秒杀订单, 订单号=%s, 用户ID=%d, 商品ID=%d",
		in.OrderNo, in.UserId, in.GoodsId))

	// 转换为标准订单输入
	orderInput := model.OrderAddInput{
		UserId:      in.UserId,
		Number:      in.OrderNo,
		Price:       int(in.Price * float64(in.Count) * 100), // 转换为分
		ActualPrice: int(in.Price * float64(in.Count) * 100), // 转换为分
		Remark:      "秒杀订单",
		OrderAddGoodsInfos: []*model.OrderAddGoodsInfo{
			{
				GoodsId:        int(in.GoodsId),
				GoodsOptionsId: int(in.GoodsOptionsId),
				Count:          int(in.Count),
				Price:          int(in.Price * 100), // 转换为分
				ActualPrice:    int(in.Price * 100), // 转换为分
			},
		},
	}

	// 调用标准下单逻辑
	return s.Add(ctx, orderInput)
}

// ProcessOrderMessage 处理订单消息
func (s *sOrder) ProcessOrderMessage(ctx context.Context, message []byte) error {
	// 解析订单消息
	var orderMsg model.SeckillOrderMsg
	if err := json.Unmarshal(message, &orderMsg); err != nil {
		return fmt.Errorf("解析订单消息失败: %v", err)
	}

	g.Log().Info(ctx, fmt.Sprintf("处理秒杀订单消息: 订单号=%s, 用户ID=%d",
		orderMsg.OrderNo, orderMsg.UserId))

	// 转换为秒杀订单输入
	input := model.SeckillOrderInput{
		UserId:         orderMsg.UserId,
		GoodsId:        orderMsg.GoodsId,
		GoodsOptionsId: orderMsg.GoodsOptionsId,
		Count:          uint(orderMsg.Count),
		Price:          orderMsg.Price,
		OrderNo:        orderMsg.OrderNo,
	}

	// 创建订单
	output, err := s.CreateSeckillOrder(ctx, input)
	if err != nil {
		return fmt.Errorf("处理订单消息创建订单失败: %v", err)
	}

	g.Log().Info(ctx, fmt.Sprintf("订单消息处理成功: 订单ID=%d, 订单号=%s",
		output.Id, orderMsg.OrderNo))

	return nil
}

// restoreGoodsStock 恢复商品库存
func (s *sOrder) restoreGoodsStock(ctx context.Context, orderId uint) error {
	// 查询订单商品
	var orderGoods []*model.OrderAddGoodsInfo
	err := dao.OrderGoodsInfo.Ctx(ctx).Where(dao.OrderGoodsInfo.Columns().OrderId, orderId).Scan(&orderGoods)
	if err != nil {
		return err
	}

	// 使用事务处理
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, item := range orderGoods {
			// 恢复商品库存
			_, err := dao.GoodsInfo.Ctx(ctx).WherePri(item.GoodsId).Increment(dao.GoodsInfo.Columns().Stock, item.Count)
			if err != nil {
				return err
			}

			// 恢复商品规格库存
			_, err = dao.GoodsOptionsInfo.Ctx(ctx).WherePri(item.GoodsOptionsId).Increment(dao.GoodsOptionsInfo.Columns().Stock, item.Count)
			if err != nil {
				return err
			}

			// 减少销量
			_, err = dao.GoodsInfo.Ctx(ctx).WherePri(item.GoodsId).Decrement(dao.GoodsInfo.Columns().Sale, item.Count)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
