package backend

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ConsigneeListReq struct {
	g.Meta `path:"/consignee/list" tags:"收货地址列表" method:"get" summary:"收货地址列表"`
	CommonPaginationReq
}

type ConsigneeListRes struct {
	List  interface{} `json:"list" description:"列表"`
	Page  int         `json:"page" description:"分页码"`
	Size  int         `json:"size" description:"分页数量"`
	Total int         `json:"total" description:"数据总数"`
}

// 添加收货地址
type ConsigneeAddReq struct {
	g.Meta    `path:"/consignee/add" tags:"收货地址管理" method:"post" summary:"添加收货地址"`
	UserId    uint   `json:"user_id" v:"required#用户ID必填" dc:"用户ID"`
	IsDefault int    `json:"is_default" dc:"默认地址1 非默认0"`
	Name      string `json:"name" v:"required#收货人名字必填" dc:"收货人名字"`
	Phone     string `json:"phone" v:"required#收货人手机号必填" dc:"收货人手机号"`
	Province  string `json:"province" v:"required#省必填" dc:"省"`
	City      string `json:"city" v:"required#城市必填" dc:"城市"`
	Town      string `json:"town" v:"required#县区必填" dc:"县区"`
	Street    string `json:"street" dc:"街道乡镇"`
	Detail    string `json:"detail" v:"required#地址详情必填" dc:"地址详情"`
}

type ConsigneeAddRes struct {
	Id uint `json:"id"`
}

// 更新收货地址
type ConsigneeUpdateReq struct {
	g.Meta      `path:"/consignee/update" tags:"收货地址管理" method:"post" summary:"更新收货地址"`
	ConsigneeId uint   `json:"consignee_id" v:"required#收货地址ID必填" dc:"收货地址ID"`
	UserId      uint   `json:"user_id" dc:"用户ID"`
	IsDefault   int    `json:"is_default" dc:"默认地址1 非默认0"`
	Name        string `json:"name" dc:"收货人名字"`
	Phone       string `json:"phone" dc:"收货人手机号"`
	Province    string `json:"province" dc:"省"`
	City        string `json:"city" dc:"城市"`
	Town        string `json:"town" dc:"县区"`
	Street      string `json:"street" dc:"街道乡镇"`
	Detail      string `json:"detail" dc:"地址详情"`
}

type ConsigneeUpdateRes struct {
	Id uint `json:"id"`
}

// 删除收货地址
type ConsigneeDeleteReq struct {
	g.Meta      `path:"/consignee/delete" tags:"收货地址管理" method:"delete" summary:"删除收货地址"`
	ConsigneeId uint `json:"consignee_id" v:"required#收货地址ID必填" dc:"收货地址ID"`
}

type ConsigneeDeleteRes struct{}

// type ConsigneeDetailReq struct {
// 	g.Meta `path:"/consignee/detail" tags:"订单详情" method:"get" summary:"订单详情"`
// 	Id     uint `json:"id"`
// }

// type ConsigneeInfoBase struct {
// 	Id        int         `json:"id"         dc:""`
// 	UserId    int         `json:"userId"     dc:"用户id"`
// 	IsDefault int         `json:"is_default" dc:"默认地址1  非默认0"`
// 	Name      string      `json:"name"       dc:"收货人名字"`
// 	Phone     string      `json:"phone"      dc:"收货人手机号"`
// 	Province  string      `json:"province"   dc:"省"`
// 	City      string      `json:"city"   	 dc:"城市"`
// 	Town      string      `json:"town"       dc:"县区"`
// 	Street    int         `json:"street"     dc:"街道乡镇"`
// 	Detail    string      `json:"detail"     dc:"地址详情"`
// 	CreatedAt *gtime.Time `json:"created_at"        dc:""`
// 	UpdatedAt *gtime.Time `json:"updated_at"        dc:""`
// 	DeletedAt *gtime.Time `json:"deleted_at"        dc:""`
// }
