package backend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type AdminLoginReq struct {
	g.Meta   `path:"/admin/login" method:"post" tags:"Backend Admin" summary:"Admin login"`
	Username string `json:"username" v:"required#Username is required"`
	Password string `json:"password" v:"required#Password is required"`
	Remember bool   `json:"remember"`
}

type AdminLoginRes struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
	Admin    AdminBase `json:"admin"`
}

type AdminLogoutReq struct {
	g.Meta `path:"/admin/logout" method:"post" tags:"Backend Admin" summary:"Admin logout"`
}

type AdminLogoutRes struct {
	Message string `json:"message"`
}

type AdminInfoReq struct {
	g.Meta `path:"/admin/info" method:"get" tags:"Backend Admin" summary:"Current admin info"`
}

type AdminInfoRes struct {
	Admin AdminBase `json:"admin"`
}

type AdminUpdatePasswordReq struct {
	g.Meta      `path:"/admin/password" method:"post" tags:"Backend Admin" summary:"Change admin password"`
	OldPassword string `json:"old_password" v:"required#Current password is required"`
	NewPassword string `json:"new_password" v:"required|length:6,64#New password is required|New password length must be 6 to 64 characters"`
}

type AdminUpdatePasswordRes struct {
	Message string `json:"message"`
}

type AdminManageCreateReq struct {
	g.Meta   `path:"/admin/manage/create" method:"post" tags:"Backend Admin" summary:"Create admin"`
	Username string `json:"username" v:"required|length:3,64#Username is required|Username length must be 3 to 64 characters"`
	Password string `json:"password" v:"required|length:6,64#Password is required|Password length must be 6 to 64 characters"`
	RealName string `json:"real_name" v:"max-length:64#Name must be at most 64 characters"`
	Phone    string `json:"phone" v:"max-length:20#Phone must be at most 20 characters"`
	Email    string `json:"email" v:"email|max-length:128#Email format is invalid|Email must be at most 128 characters"`
	IsSuper  int    `json:"is_super" v:"in:0,1#Super flag must be 0 or 1"`
	RoleIds  []uint `json:"role_ids"`
}

type AdminManageCreateRes struct {
	Admin AdminBase `json:"admin"`
}

type AdminManageListReq struct {
	g.Meta   `path:"/admin/manage/list" method:"get" tags:"Backend Admin" summary:"Admin list"`
	Page     int    `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size     int    `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
	Username string `json:"username" in:"query"`
	RealName string `json:"real_name" in:"query"`
	Status   int    `json:"status" in:"query" d:"-1"`
}

type AdminManageListRes struct {
	List  []AdminBase `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

type AdminManageDetailReq struct {
	g.Meta `path:"/admin/manage/detail" method:"get" tags:"Backend Admin" summary:"Admin detail"`
	Id     uint `json:"id" in:"query" v:"required|min:1#Admin ID is required|Admin ID is invalid"`
}

type AdminManageDetailRes struct {
	Admin AdminBase `json:"admin"`
}

type AdminManageUpdateReq struct {
	g.Meta   `path:"/admin/manage/update" method:"post" tags:"Backend Admin" summary:"Update admin"`
	Id       uint   `json:"id" v:"required|min:1#Admin ID is required|Admin ID is invalid"`
	RealName string `json:"real_name" v:"max-length:64#Name must be at most 64 characters"`
	Phone    string `json:"phone" v:"max-length:20#Phone must be at most 20 characters"`
	Email    string `json:"email" v:"email|max-length:128#Email format is invalid|Email must be at most 128 characters"`
	IsSuper  int    `json:"is_super" v:"in:0,1#Super flag must be 0 or 1"`
	RoleIds  []uint `json:"role_ids"`
}

type AdminManageUpdateRes struct {
	Admin AdminBase `json:"admin"`
}

type AdminManageStatusReq struct {
	g.Meta `path:"/admin/manage/status" method:"post" tags:"Backend Admin" summary:"Enable or disable admin"`
	Id     uint `json:"id" v:"required|min:1#Admin ID is required|Admin ID is invalid"`
	Status int  `json:"status" v:"required|in:0,1#Status is required|Status must be 0 or 1"`
}

type AdminManageStatusRes struct {
	Message string `json:"message"`
}

type AdminManageResetPasswordReq struct {
	g.Meta   `path:"/admin/manage/reset-password" method:"post" tags:"Backend Admin" summary:"Reset admin password"`
	Id       uint   `json:"id" v:"required|min:1#Admin ID is required|Admin ID is invalid"`
	Password string `json:"password" v:"required|length:6,64#Password is required|Password length must be 6 to 64 characters"`
}

type AdminManageResetPasswordRes struct {
	Message string `json:"message"`
}

type AdminManageRolesReq struct {
	g.Meta  `path:"/admin/manage/roles" method:"post" tags:"Backend Admin" summary:"Assign admin roles"`
	Id      uint   `json:"id" v:"required|min:1#Admin ID is required|Admin ID is invalid"`
	RoleIds []uint `json:"role_ids"`
}

type AdminManageRolesRes struct {
	Message string `json:"message"`
}

type AdminBase struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	IsSuper  int    `json:"is_super"`
	RoleIds  []uint `json:"role_ids"`
}
