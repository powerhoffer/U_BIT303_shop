package backend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type EmployeeRegisterReq struct {
	g.Meta   `path:"/employee/register" method:"post" tags:"Backend Employees" summary:"Register employee"`
	Username string `json:"username" v:"required|length:3,64#Username is required|Username length must be 3 to 64 characters"`
	Password string `json:"password" v:"required|length:6,64#Password is required|Password length must be 6 to 64 characters"`
	RealName string `json:"real_name" v:"max-length:64#Name must be at most 64 characters"`
	Phone    string `json:"phone" v:"max-length:20#Phone must be at most 20 characters"`
	Email    string `json:"email" v:"email|max-length:128#Email format is invalid|Email must be at most 128 characters"`
}

type EmployeeRegisterRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeLoginReq struct {
	g.Meta   `path:"/employee/login" method:"post" tags:"Backend Employees" summary:"Employee login"`
	Username string `json:"username" v:"required#Username is required"`
	Password string `json:"password" v:"required#Password is required"`
	Remember bool   `json:"remember"`
}

type EmployeeLoginRes struct {
	Token    string       `json:"token"`
	ExpireAt time.Time    `json:"expire_at"`
	Employee EmployeeBase `json:"employee"`
}

type EmployeeLogoutReq struct {
	g.Meta `path:"/employee/logout" method:"post" tags:"Backend Employees" summary:"Employee logout"`
}

type EmployeeLogoutRes struct {
	Message string `json:"message"`
}

type EmployeeInfoReq struct {
	g.Meta `path:"/employee/info" method:"get" tags:"Backend Employees" summary:"Current employee info"`
}

type EmployeeInfoRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeUpdatePasswordReq struct {
	g.Meta      `path:"/employee/password" method:"post" tags:"Backend Employees" summary:"Change password"`
	OldPassword string `json:"old_password" v:"required#Current password is required"`
	NewPassword string `json:"new_password" v:"required|length:6,64#New password is required|New password length must be 6 to 64 characters"`
}

type EmployeeUpdatePasswordRes struct {
	Message string `json:"message"`
}

type EmployeeManageCreateReq struct {
	g.Meta   `path:"/employee/manage/create" method:"post" tags:"Backend Employees" summary:"Create employee"`
	Username string `json:"username" v:"required|length:3,64#Username is required|Username length must be 3 to 64 characters"`
	Password string `json:"password" v:"required|length:6,64#Password is required|Password length must be 6 to 64 characters"`
	RealName string `json:"real_name" v:"max-length:64#Name must be at most 64 characters"`
	Phone    string `json:"phone" v:"max-length:20#Phone must be at most 20 characters"`
	Email    string `json:"email" v:"email|max-length:128#Email format is invalid|Email must be at most 128 characters"`
}

type EmployeeManageCreateRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeManageListReq struct {
	g.Meta   `path:"/employee/manage/list" method:"get" tags:"Backend Employees" summary:"Employee list"`
	Page     int    `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size     int    `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
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
	g.Meta `path:"/employee/manage/detail" method:"get" tags:"Backend Employees" summary:"Employee detail"`
	Id     uint `json:"id" in:"query" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
}

type EmployeeManageDetailRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeManageUpdateReq struct {
	g.Meta   `path:"/employee/manage/update" method:"post" tags:"Backend Employees" summary:"Update employee"`
	Id       uint   `json:"id" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
	RealName string `json:"real_name" v:"max-length:64#Name must be at most 64 characters"`
	Phone    string `json:"phone" v:"max-length:20#Phone must be at most 20 characters"`
	Email    string `json:"email" v:"email|max-length:128#Email format is invalid|Email must be at most 128 characters"`
}

type EmployeeManageUpdateRes struct {
	Employee EmployeeBase `json:"employee"`
}

type EmployeeManageStatusReq struct {
	g.Meta `path:"/employee/manage/status" method:"post" tags:"Backend Employees" summary:"Enable or disable employee"`
	Id     uint   `json:"id" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
	Status string `json:"status" v:"required|in:0,1#Status is required|Status must be 0 or 1"`
}

type EmployeeManageStatusRes struct {
	Message string `json:"message"`
}

type EmployeeManageResetPasswordReq struct {
	g.Meta   `path:"/employee/manage/reset-password" method:"post" tags:"Backend Employees" summary:"Reset employee password"`
	Id       uint   `json:"id" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
	Password string `json:"password" v:"required|length:6,64#Password is required|Password length must be 6 to 64 characters"`
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
