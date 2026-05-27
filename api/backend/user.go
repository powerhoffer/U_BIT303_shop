package backend

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// 用户列表请求
type UserGetListReq struct {
	g.Meta `path:"/user/list" method:"get" tags:"用户管理" summary:"用户列表接口"`
	CommonPaginationReq
}

type UserGetListRes struct {
	List  []UserListItem `json:"list" description:"列表"`
	Page  int            `json:"page" description:"分页码"`
	Size  int            `json:"size" description:"分页数量"`
	Total int            `json:"total" description:"数据总数"`
}

type UserListItem struct {
	Id        int         `json:"id"`
	Name      string      `json:"name"`
	Avatar    string      `json:"avatar"`
	Sex       int         `json:"sex"`
	Status    int         `json:"status"` // 1正常 2冻结
	Sign      string      `json:"sign"`
	CreatedAt *gtime.Time `json:"created_at"`
	UpdatedAt *gtime.Time `json:"updated_at"`
}

// 冻结/解冻用户请求
type UserUpdateStatusReq struct {
	g.Meta `path:"/user/status" method:"post" tags:"用户管理" summary:"冻结/解冻用户"`
	Id     uint `json:"id" v:"required#用户ID不能为空" dc:"用户ID"`
	Status int  `json:"status" v:"required|in:0,1#状态不能为空|状态值只能是0或1" dc:"状态 1正常 0冻结"`
}

type UserUpdateStatusRes struct{}

// 删除用户请求
type UserDeleteReq struct {
	g.Meta `path:"/user/delete" method:"delete" tags:"用户管理" summary:"删除用户"`
	Id     uint `json:"id" v:"required#用户ID不能为空" dc:"用户ID"`
}

type UserDeleteRes struct{}

// 提升用户为管理员请求
type UserPromoteToAdminReq struct {
	g.Meta   `path:"/user/promote" method:"post" tags:"用户管理" summary:"提升用户为管理员"`
	UserId   uint   `json:"user_id" v:"required#用户ID不能为空" dc:"用户ID"`
	RoleIds  string `json:"role_ids" dc:"角色IDs，多个用逗号分隔"`
	Password string `json:"password" v:"required#密码不能为空" dc:"管理员密码"`
}

type UserPromoteToAdminRes struct {
	AdminId int `json:"admin_id" dc:"新创建的管理员ID"`
}
