package backend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type EmployeeRegisterReq struct {
	g.Meta   `path:"/employee/register" method:"post" tags:"后台员工" summary:"员工注册"`
	Username string `json:"username" v:"required|length:3,64#账号不能为空|账号长度应为3到64位"`
	Password string `json:"password" v:"required|length:6,64#密码不能为空|密码长度应为6到64位"`
	RealName string `json:"real_name" v:"max-length:64#姓名最多64位"`
	Phone    string `json:"phone" v:"max-length:20#手机号最多20位"`
	Email    string `json:"email" v:"email|max-length:128#邮箱格式不正确|邮箱最多128位"`
}

type EmployeeRegisterRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeLoginReq struct {
	g.Meta   `path:"/employee/login" method:"post" tags:"后台员工" summary:"员工登录"`
	Username string `json:"username" v:"required#账号不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
	Remember bool   `json:"remember"`
}

type EmployeeLoginRes struct {
	Token    string       `json:"token"`
	ExpireAt time.Time    `json:"expire_at"`
	Employee EmployeeBase `json:"employee"`
}

type EmployeeLogoutReq struct {
	g.Meta `path:"/employee/logout" method:"post" tags:"后台员工" summary:"员工登出"`
}

type EmployeeLogoutRes struct {
	Message string `json:"message"`
}

type EmployeeInfoReq struct {
	g.Meta `path:"/employee/info" method:"get" tags:"后台员工" summary:"当前员工信息"`
}

type EmployeeInfoRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeUpdatePasswordReq struct {
	g.Meta      `path:"/employee/password" method:"post" tags:"后台员工" summary:"修改密码"`
	OldPassword string `json:"old_password" v:"required#旧密码不能为空"`
	NewPassword string `json:"new_password" v:"required|length:6,64#新密码不能为空|新密码长度应为6到64位"`
}

type EmployeeUpdatePasswordRes struct {
	Message string `json:"message"`
}

type EmployeeManageCreateReq struct {
	g.Meta   `path:"/employee/manage/create" method:"post" tags:"后台员工" summary:"新增员工"`
	Username string `json:"username" v:"required|length:3,64#账号不能为空|账号长度应为3到64位"`
	Password string `json:"password" v:"required|length:6,64#密码不能为空|密码长度应为6到64位"`
	RealName string `json:"real_name" v:"max-length:64#姓名最多64位"`
	Phone    string `json:"phone" v:"max-length:20#手机号最多20位"`
	Email    string `json:"email" v:"email|max-length:128#邮箱格式不正确|邮箱最多128位"`
}

type EmployeeManageCreateRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeManageListReq struct {
	g.Meta   `path:"/employee/manage/list" method:"get" tags:"后台员工" summary:"员工列表"`
	Page     int    `json:"page" in:"query" d:"1" v:"min:1#分页页码错误"`
	Size     int    `json:"size" in:"query" d:"10" v:"max:50#分页数量最多50条"`
	Username string `json:"username" in:"query"`
	RealName string `json:"real_name" in:"query"`
	Status   int    `json:"status" in:"query" d:"-1"`
}

type EmployeeManageListRes struct {
	List  []EmployeeBase `json:"list"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
}

type EmployeeManageDetailReq struct {
	g.Meta `path:"/employee/manage/detail" method:"get" tags:"后台员工" summary:"员工详情"`
	Id     uint `json:"id" in:"query" v:"required|min:1#员工ID不能为空|员工ID错误"`
}

type EmployeeManageDetailRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeManageUpdateReq struct {
	g.Meta   `path:"/employee/manage/update" method:"post" tags:"后台员工" summary:"编辑员工信息"`
	Id       uint   `json:"id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	RealName string `json:"real_name" v:"max-length:64#姓名最多64位"`
	Phone    string `json:"phone" v:"max-length:20#手机号最多20位"`
	Email    string `json:"email" v:"email|max-length:128#邮箱格式不正确|邮箱最多128位"`
}

type EmployeeManageUpdateRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeManageStatusReq struct {
	g.Meta `path:"/employee/manage/status" method:"post" tags:"后台员工" summary:"启用或禁用员工"`
	Id     uint   `json:"id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Status string `json:"status" v:"required|in:0,1#状态不能为空|状态只能是0或1"`
}

type EmployeeManageStatusRes struct {
	Message string `json:"message"`
}

type EmployeeManageResetPasswordReq struct {
	g.Meta   `path:"/employee/manage/reset-password" method:"post" tags:"后台员工" summary:"重置员工密码"`
	Id       uint   `json:"id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Password string `json:"password" v:"required|length:6,64#密码不能为空|密码长度应为6到64位"`
}

type EmployeeManageResetPasswordRes struct {
	Message string `json:"message"`
}

type EmployeeBase struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}
