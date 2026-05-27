package backend

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type sOrder struct{}

var Order = sOrder{}

func New() *sOrder {
	return &sOrder{}
}

// Add 创建订单
func (c *sOrder) Add(ctx context.Context, req *backend.OrderAddReq) (res *backend.OrderAddRes, err error) {
	// 构建订单信息
	var orderAddInput model.OrderAddInput
	if err = gconv.Scan(req, &orderAddInput); err != nil {
		return nil, err
	}

	// 处理商品信息
	orderAddInput.OrderAddGoodsInfos = make([]*model.OrderAddGoodsInfo, 0)
	for _, item := range req.GoodsList {
		goodsInfo := &model.OrderAddGoodsInfo{
			GoodsId:        item.GoodsId,
			GoodsOptionsId: item.GoodsOptionsId,
			Count:          item.Count,
			Price:          item.Price,
			CouponPrice:    item.CouponPrice,
			ActualPrice:    item.ActualPrice,
			Remark:         item.Remark,
		}
		orderAddInput.OrderAddGoodsInfos = append(orderAddInput.OrderAddGoodsInfos, goodsInfo)
	}

	// 调用 service 层
	result, err := service.Order().Add(ctx, orderAddInput)
	if err != nil {
		return nil, err
	}

	return &backend.OrderAddRes{
		Id: result.Id,
	}, nil
}

// List 获取所有订单列表
func (s *sOrder) List(ctx context.Context, req *backend.OrderListReq) (res *backend.OrderListRes, err error) {
	// 转换请求参数
	var orderListInput model.OrderListInput
	if err = gconv.Scan(req, &orderListInput); err != nil {
		return nil, err
	}

	// 调用service层获取订单列表
	result, err := service.Order().List(ctx, orderListInput)
	if err != nil {
		return nil, err
	}

	// 创建响应
	res = &backend.OrderListRes{
		CommonPaginationRes: backend.CommonPaginationRes{
			Total: result.Total,
			Page:  result.Page,
			Size:  result.Size,
		},
	}

	return res, nil
}

// Detail 获取订单详情
func (s *sOrder) Detail(ctx context.Context, req *backend.OrderDetailReq) (res *backend.OrderDetailRes, err error) {
	// 转换请求参数
	orderDetailInput := model.OrderDetailInput{
		Id: req.Id,
	}

	// 调用service层获取订单详情
	result, err := service.Order().Detail(ctx, orderDetailInput)
	if err != nil {
		return nil, err
	}

	// 创建响应
	res = &backend.OrderDetailRes{}

	// 转换数据
	if err = gconv.Scan(result, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateOrderStatus 更新订单状态
func (s *sOrder) UpdateOrderStatus(ctx context.Context, req *backend.OrderUpdateStatusReq) (res *backend.OrderUpdateStatusRes, err error) {
	// 需要在service层实现UpdateOrderStatus方法
	// TODO: 实现更新订单状态的逻辑

	// 直接返回成功，不调用未定义的service方法
	return &backend.OrderUpdateStatusRes{
		Id: req.Id,
	}, nil
}

// Delete 删除订单
func (s *sOrder) Delete(ctx context.Context, req *backend.OrderDeleteReq) (res *backend.OrderDeleteRes, err error) {
	// 需要在service层实现Delete方法
	// TODO: 实现删除订单的逻辑

	// 直接返回成功，不调用未定义的service方法
	return &backend.OrderDeleteRes{
		Id: req.Id,
	}, nil
}

// Refund 订单退款
func (s *sOrder) Refund(ctx context.Context, req *backend.OrderRefundReq) (res *backend.OrderRefundRes, err error) {
	// 需要在service层实现Refund方法
	// TODO: 实现订单退款的逻辑

	// 直接返回成功，不调用未定义的service方法
	return &backend.OrderRefundRes{
		Id: req.Id,
	}, nil
}

func (a *sOrder) UpdateStatus(ctx context.Context, req *frontend.UpdateOrderStatusReq) (res *frontend.UpdateOrderStatusRes, err error) {
	// 调用秒杀服务
	result, err := service.Seckill().UpdateOrderStatus(ctx, model.SeckillUpdateOrderStatusInput{
		OrderNo: req.OrderNo,
		Status:  req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &frontend.UpdateOrderStatusRes{
		Success: result.Success,
	}, nil
}
