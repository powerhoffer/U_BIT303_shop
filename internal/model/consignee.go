package model

import "github.com/gogf/gf/v2/os/gtime"

// ConsigneeGetListInput 获取内容列表
type ConsigneeGetListInput struct {
	Page int // 分页号码
	Size int // 分页数量，最大50
	Sort int // 排序类型(0:最新, 默认。1:活跃, 2:热度)
}

// ConsigneeGetListOutput 查询列表结果
type ConsigneeGetListOutput struct {
	List  []ConsigneeGetListOutputItem `json:"list" description:"列表"`
	Page  int                          `json:"page" description:"分页码"`
	Size  int                          `json:"size" description:"分页数量"`
	Total int                          `json:"total" description:"数据总数"`
}

type ConsigneeGetListOutputItem struct {
	Id        uint        `json:"id"` // 自增ID
	UserId    uint        `json:"userId"     dc:"用户id"`
	IsDefault int         `json:"is_default" dc:"默认地址1  非默认0"`
	Name      string      `json:"name"       dc:"收货人名字"`
	Phone     string      `json:"phone"      dc:"收货人手机号"`
	Province  string      `json:"province"   dc:"省"`
	City      string      `json:"city"   dc:"城市"`
	Town      string      `json:"town"   dc:"县区"`
	Street    string      `json:"street"     dc:"街道乡镇"`
	Detail    string      `json:"detail"     dc:"地址详情"`
	CreatedAt *gtime.Time `json:"created_at"` // 创建时间
	UpdatedAt *gtime.Time `json:"updated_at"` // 修改时间
	DeletedAt *gtime.Time `json:"deleted_at" `
}

// 后台管理列表输出（包含用户信息）
type ConsigneeAdminListOutput struct {
	List  []ConsigneeAdminListItem `json:"list" description:"列表"`
	Page  int                      `json:"page" description:"分页码"`
	Size  int                      `json:"size" description:"分页数量"`
	Total int                      `json:"total" description:"数据总数"`
}

type ConsigneeAdminListItem struct {
	Id        uint   `json:"id"`
	UserName  string `json:"user_name"` // 用户名
	IsDefault int    `json:"is_default"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"` // 完整地址
}

// AddConsigneeInput 添加收货地址输入
type AddConsigneeInput struct {
	UserId    uint   `json:"userId"`
	IsDefault int    `json:"is_default"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Town      string `json:"town"`
	Street    string `json:"street"`
	Detail    string `json:"detail"`
}

// AddConsigneeOutput 添加收货地址输出
type AddConsigneeOutput struct {
	Id uint `json:"id"`
}

// UpdateConsigneeInput 更新收货地址输入
type UpdateConsigneeInput struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"userId"`
	IsDefault int    `json:"is_default"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Town      string `json:"town"`
	Street    string `json:"street"`
	Detail    string `json:"detail"`
}

// UpdateConsigneeOutput 更新收货地址输出
type UpdateConsigneeOutput struct {
	Id uint `json:"id"`
}

// DeleteConsigneeInput 删除收货地址输入
type DeleteConsigneeInput struct {
	Id uint `json:"id"`
}

// DeleteConsigneeOutput 删除收货地址输出
type DeleteConsigneeOutput struct {
	Id uint `json:"id"`
}
