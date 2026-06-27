package consts

const (
	ProjectName  = "bit303_shop"
	ProjectUsage = "bit303_shop"
	ProjectBrief = "BIT303 employee welfare mall backend scaffold"
)

const (
	EmployeeStatusDisabled  = 0
	EmployeeStatusNormal    = 1
	CtxEmployeeId           = "employee_id"
	CtxEmployeeUsername     = "employee_username"
	AdminStatusDisabled     = 0
	AdminStatusNormal       = 1
	AdminRoleDisabled       = 0
	AdminRoleEnabled        = 1
	AdminPermissionDisabled = 0
	AdminPermissionEnabled  = 1
	CtxAdminId              = "admin_id"
	CtxAdminUsername        = "admin_username"
	CtxAdminIsSuper         = "admin_is_super"
)

const (
	PointsAccountStatusDisabled = 0
	PointsAccountStatusNormal   = 1
	PointsChangeTypeAdd         = 1
	PointsChangeTypeDeduct      = 2
)

const (
	GoodsCategoryStatusDisabled = 0
	GoodsCategoryStatusEnabled  = 1
)

const (
	GoodsStatusOffShelf = 0
	GoodsStatusOnShelf  = 1
)

const (
	OrderStatusPending   = 1
	OrderStatusCompleted = 2
	OrderStatusCancelled = 3
)
