package model

import "time"

type EmployeeRegisterInput struct {
	Username string
	Password string
	RealName string
	Phone    string
	Email    string
}

type EmployeeRegisterOutput struct {
	Employee EmployeeBase
}

type EmployeeLoginInput struct {
	Username string
	Password string
	Remember bool
}

type EmployeeLoginOutput struct {
	Token    string
	ExpireAt time.Time
	Employee EmployeeBase
}

type EmployeeInfoOutput struct {
	Employee EmployeeBase
}

type EmployeeUpdatePasswordInput struct {
	EmployeeId  uint
	OldPassword string
	NewPassword string
}

type EmployeeManageCreateInput struct {
	Username string
	Password string
	RealName string
	Phone    string
	Email    string
}

type EmployeeManageCreateOutput struct {
	Employee EmployeeBase
}

type EmployeeManageListInput struct {
	Page     int
	Size     int
	Username string
	RealName string
	Status   int
}

type EmployeeManageListOutput struct {
	List  []EmployeeBase
	Total int
	Page  int
	Size  int
}

type EmployeeManageDetailOutput struct {
	Employee EmployeeBase
}

type EmployeeManageUpdateInput struct {
	Id       uint
	RealName string
	Phone    string
	Email    string
}

type EmployeeManageUpdateOutput struct {
	Employee EmployeeBase
}

type EmployeeManageStatusInput struct {
	Id     uint
	Status int
}

type EmployeeManageResetPasswordInput struct {
	Id       uint
	Password string
}

type EmployeeManageDeleteInput struct {
	Id uint
}

type EmployeeBase struct {
	Id       uint
	Username string
	RealName string
	Phone    string
	Email    string
	Status   int
}
