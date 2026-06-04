// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"bit303_shop/internal/model"
	"context"
)

type (
	IEmployee interface {
		Register(ctx context.Context, in model.EmployeeRegisterInput) (out model.EmployeeRegisterOutput, err error)
		Login(ctx context.Context, in model.EmployeeLoginInput) (out model.EmployeeLoginOutput, err error)
		Info(ctx context.Context, employeeId uint) (out model.EmployeeInfoOutput, err error)
		UpdatePassword(ctx context.Context, in model.EmployeeUpdatePasswordInput) error
	}
)

var (
	localEmployee IEmployee
)

func Employee() IEmployee {
	if localEmployee == nil {
		panic("implement not found for interface IEmployee, forgot register?")
	}
	return localEmployee
}

func RegisterEmployee(i IEmployee) {
	localEmployee = i
}
