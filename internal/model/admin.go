package model

import "time"

type AdminLoginInput struct {
	Username string
	Password string
	Remember bool
}

type AdminLoginOutput struct {
	Token    string
	ExpireAt time.Time
	Admin    AdminBase
}

type AdminInfoOutput struct {
	Admin AdminBase
}

type AdminUpdatePasswordInput struct {
	AdminId     uint
	OldPassword string
	NewPassword string
}

type AdminManageCreateInput struct {
	Username string
	Password string
	RealName string
	Phone    string
	Email    string
	IsSuper  int
	RoleIds  []uint
}

type AdminManageCreateOutput struct {
	Admin AdminBase
}

type AdminManageListInput struct {
	Page     int
	Size     int
	Username string
	RealName string
	Status   int
}

type AdminManageListOutput struct {
	List  []AdminBase
	Total int
	Page  int
	Size  int
}

type AdminManageDetailOutput struct {
	Admin AdminBase
}

type AdminManageUpdateInput struct {
	Id       uint
	RealName string
	Phone    string
	Email    string
	IsSuper  int
	RoleIds  []uint
}

type AdminManageUpdateOutput struct {
	Admin AdminBase
}

type AdminManageStatusInput struct {
	Id     uint
	Status int
}

type AdminManageResetPasswordInput struct {
	Id       uint
	Password string
}

type AdminManageRolesInput struct {
	Id      uint
	RoleIds []uint
}

type AdminBase struct {
	Id       uint
	Username string
	RealName string
	Phone    string
	Email    string
	Status   int
	IsSuper  int
	RoleIds  []uint
}
