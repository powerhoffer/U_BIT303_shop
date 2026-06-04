package employee

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type RegisterReq struct {
	g.Meta   `path:"/register" method:"post" tags:"员工账号" summary:"员工注册"`
	Username string `json:"username" v:"required|length:3,64#账号不能为空|账号长度应为3到64位"`
	Password string `json:"password" v:"required|length:6,64#密码不能为空|密码长度应为6到64位"`
	RealName string `json:"real_name" v:"max-length:64#姓名最多64位"`
	Phone    string `json:"phone" v:"max-length:20#手机号最多20位"`
	Email    string `json:"email" v:"email|max-length:128#邮箱格式不正确|邮箱最多128位"`
}

type RegisterRes struct {
	Employee EmployeeBase `json:"employee"`
}

type LoginReq struct {
	g.Meta   `path:"/login" method:"post" tags:"员工账号" summary:"员工登录"`
	Username string `json:"username" v:"required#账号不能为空"`
	Password string `json:"password" v:"required#密码不能为空"`
	Remember bool   `json:"remember"`
}

type LoginRes struct {
	Token    string       `json:"token"`
	ExpireAt time.Time    `json:"expire_at"`
	Employee EmployeeBase `json:"employee"`
}

type LogoutReq struct {
	g.Meta `path:"/logout" method:"post" tags:"员工账号" summary:"员工登出"`
}

type LogoutRes struct {
	Message string `json:"message"`
}

type InfoReq struct {
	g.Meta `path:"/info" method:"get" tags:"员工账号" summary:"当前员工信息"`
}

type InfoRes struct {
	Employee EmployeeBase `json:"employee"`
}

type UpdatePasswordReq struct {
	g.Meta      `path:"/password" method:"post" tags:"员工账号" summary:"修改密码"`
	OldPassword string `json:"old_password" v:"required#旧密码不能为空"`
	NewPassword string `json:"new_password" v:"required|length:6,64#新密码不能为空|新密码长度应为6到64位"`
}

type UpdatePasswordRes struct {
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
