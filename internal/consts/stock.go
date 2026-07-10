package consts

const (
	StockChangeTypeInitial            = 1
	StockChangeTypeAdminIncrease      = 2
	StockChangeTypeAdminDecrease      = 3
	StockChangeTypeOrderDeduct        = 4
	StockChangeTypeOrderCancelRestore = 5
)

const (
	StockOperatorTypeSystem   = 0
	StockOperatorTypeAdmin    = 1
	StockOperatorTypeEmployee = 2
)

const (
	StockActionIncrease = "increase"
	StockActionDecrease = "decrease"
)

const (
	StockBizTypeInitial             = "initial"
	StockBizTypeGoodsCreate         = "goods_create"
	StockBizTypeGoodsUpdate         = "goods_update"
	StockBizTypeStockAdjust         = "stock_adjust"
	StockBizTypeOrderCreate         = "order_create"
	StockBizTypeEmployeeOrderCancel = "employee_order_cancel"
	StockBizTypeAdminOrderCancel    = "admin_order_cancel"
)
