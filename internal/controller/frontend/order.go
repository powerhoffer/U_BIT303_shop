package frontend

import (
	"context"
	"bit303_shop/api/frontend"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

type sOrder struct{}

var Order = sOrder{}

// Add 创建订单
func (s *sOrder) Add(ctx context.Context, req *frontend.OrderAddReq) (res *frontend.OrderAddRes, err error) {
	// 获取当前用户ID
	userId := gconv.Int(ctx.Value(consts.CtxUserId))
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 转换请求参数
	var orderAddInput model.OrderAddInput
	if err = gconv.Scan(req, &orderAddInput); err != nil {
		return nil, err
	}
	orderAddInput.UserId = uint(userId)

	// 调用服务层
	result, err := service.Order().Add(ctx, orderAddInput)
	if err != nil {
		return nil, err
	}

	return &frontend.OrderAddRes{
		Id: result.Id,
	}, nil
}

// List 获取当前用户的订单列表
func (s *sOrder) List(ctx context.Context, req *frontend.OrderListReq) (res *frontend.OrderListRes, err error) {
	// 获取当前用户ID
	userId := gconv.Int(ctx.Value(consts.CtxUserId))
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 转换请求参数
	var orderListInput model.OrderListInput
	if err = gconv.Scan(req, &orderListInput); err != nil {
		return nil, err
	}
	orderListInput.UserId = userId

	// 调用service层获取订单列表
	result, err := service.Order().List(ctx, orderListInput)
	if err != nil {
		return nil, err
	}

	// 创建响应
	res = &frontend.OrderListRes{
		CommonPaginationRes: frontend.CommonPaginationRes{
			Page:  result.Page,
			Size:  result.Size,
			Total: result.Total,
		},
		List: make([]frontend.OrderInfo, 0),
	}

	// 转换订单列表
	if err = gconv.Scan(result.List, &res.List); err != nil {
		return nil, err
	}

	return res, nil
}

// Detail 获取当前用户的订单详情
func (s *sOrder) Detail(ctx context.Context, req *frontend.OrderDetailReq) (res *frontend.OrderDetailRes, err error) {
	// 获取当前用户ID
	userId := gconv.Int(ctx.Value(consts.CtxUserId))
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 转换请求参数
	orderDetailInput := model.OrderDetailInput{
		Id: req.Id,
	}

	// 调用service层获取订单详情
	result, err := service.Order().Detail(ctx, orderDetailInput)
	if err != nil {
		return nil, err
	}

	// 创建响应并转换数据
	res = &frontend.OrderDetailRes{}

	// 填充订单信息
	var orderInfo frontend.OrderInfo
	if err = gconv.Scan(result, &orderInfo); err != nil {
		return nil, err
	}
	res.OrderInfo = orderInfo

	// 填充商品信息
	if err = gconv.Scan(result.GoodsInfo, &res.GoodsInfo); err != nil {
		return nil, err
	}

	// 验证订单所属权
	if res.OrderInfo.UserId != uint(userId) {
		return nil, gerror.New("无权查看该订单")
	}

	return res, nil
}

// Cancel 取消订单
func (s *sOrder) Cancel(ctx context.Context, req *frontend.OrderCancelReq) (res *frontend.OrderCancelRes, err error) {
	// 获取当前用户ID
	userId := gconv.Int(ctx.Value(consts.CtxUserId))
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 验证订单所属权和状态
	orderDetailInput := model.OrderDetailInput{
		Id: req.Id,
	}

	orderDetail, err := service.Order().Detail(ctx, orderDetailInput)
	if err != nil {
		return nil, err
	}

	var orderInfo frontend.OrderInfo
	if err = gconv.Scan(orderDetail, &orderInfo); err != nil {
		return nil, err
	}

	if orderInfo.UserId != uint(userId) {
		return nil, gerror.New("无权操作该订单")
	}

	// 将字符串状态检查改为数字状态检查
	statusValue := gconv.Int(orderInfo.Status)
	if statusValue != 1 { // 1表示待支付状态
		return nil, gerror.New("订单状态不允许取消")
	}

	// 调用service层取消订单
	result, err := service.Order().Cancel(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Pay 订单支付
func (s *sOrder) Pay(ctx context.Context, req *frontend.OrderPayReq) (res *frontend.OrderPayRes, err error) {
	// 获取当前用户ID
	userId := gconv.Int(ctx.Value(consts.CtxUserId))
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 验证订单所属权和状态
	orderDetailInput := model.OrderDetailInput{
		Id: req.Id,
	}

	orderDetail, err := service.Order().Detail(ctx, orderDetailInput)
	if err != nil {
		return nil, err
	}

	var orderInfo frontend.OrderInfo
	if err = gconv.Scan(orderDetail, &orderInfo); err != nil {
		return nil, err
	}

	if orderInfo.UserId != uint(userId) {
		return nil, gerror.New("无权操作该订单")
	}

	// 将字符串状态检查改为数字状态检查
	statusValue := gconv.Int(orderInfo.Status)
	if statusValue != 1 { // 1表示待支付状态
		return nil, gerror.New("订单状态不允许支付")
	}

	// 调用service层支付订单
	result, err := service.Order().Pay(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Confirm 确认收货
func (s *sOrder) Confirm(ctx context.Context, req *frontend.OrderConfirmReq) (res *frontend.OrderConfirmRes, err error) {
	// 获取当前用户ID
	userId := gconv.Int(ctx.Value(consts.CtxUserId))
	if userId == 0 {
		return nil, gerror.New("用户未登录")
	}

	// 验证订单所属权和状态
	orderDetailInput := model.OrderDetailInput{
		Id: req.Id,
	}

	orderDetail, err := service.Order().Detail(ctx, orderDetailInput)
	if err != nil {
		return nil, err
	}

	var orderInfo frontend.OrderInfo
	if err = gconv.Scan(orderDetail, &orderInfo); err != nil {
		return nil, err
	}

	if orderInfo.UserId != uint(userId) {
		return nil, gerror.New("无权操作该订单")
	}

	// 将字符串状态检查改为数字状态检查
	statusValue := gconv.Int(orderInfo.Status)
	if statusValue != 3 { // 3表示已发货状态
		return nil, gerror.New("订单状态不允许确认收货")
	}

	// 调用service层确认收货
	result, err := service.Order().Confirm(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
