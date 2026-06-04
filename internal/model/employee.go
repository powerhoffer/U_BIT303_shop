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

type EmployeeBase struct {
	Id       uint
	Username string
	RealName string
	Phone    string
	Email    string
	Status   int
}
