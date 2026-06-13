package backend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type PointsBalanceReq struct {
	g.Meta `path:"/points/balance" method:"get" tags:"Backend Credits" summary:"Current employee credit balance"`
}

type PointsBalanceRes struct {
	Balance uint `json:"balance"`
}

type PointsRecordsReq struct {
	g.Meta `path:"/points/records" method:"get" tags:"Backend Credits" summary:"Current employee credit records"`
	Page   int `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size   int `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
}

type PointsRecordsRes struct {
	List  []PointsRecordItem `json:"list"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

type PointsManageAddReq struct {
	g.Meta     `path:"/points/manage/add" method:"post" tags:"Backend Credits" summary:"Add employee credits"`
	EmployeeId uint   `json:"employee_id" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
	Points     uint   `json:"points" v:"required|min:1#Credits are required|Credits must be greater than 0"`
	Remark     string `json:"remark" v:"max-length:255#Remark must be at most 255 characters"`
}

type PointsManageAddRes struct {
	Balance uint `json:"balance"`
}

type PointsManageDeductReq struct {
	g.Meta     `path:"/points/manage/deduct" method:"post" tags:"Backend Credits" summary:"Deduct employee credits"`
	EmployeeId uint   `json:"employee_id" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
	Points     uint   `json:"points" v:"required|min:1#Credits are required|Credits must be greater than 0"`
	Remark     string `json:"remark" v:"max-length:255#Remark must be at most 255 characters"`
}

type PointsManageDeductRes struct {
	Balance uint `json:"balance"`
}

type PointsManageRecordsReq struct {
	g.Meta     `path:"/points/manage/records" method:"get" tags:"Backend Credits" summary:"Employee credit records"`
	EmployeeId uint `json:"employee_id" in:"query" v:"required|min:1#Employee ID is required|Employee ID is invalid"`
	Page       int  `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size       int  `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
}

type PointsManageRecordsRes struct {
	List  []PointsRecordItem `json:"list"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

type PointsRecordItem struct {
	Id                 uint      `json:"id"`
	EmployeeId         uint      `json:"employee_id"`
	ChangeType         int       `json:"change_type"`
	Points             uint      `json:"points"`
	BeforeBalance      uint      `json:"before_balance"`
	AfterBalance       uint      `json:"after_balance"`
	OperatorEmployeeId uint      `json:"operator_employee_id"`
	Remark             string    `json:"remark"`
	CreatedAt          time.Time `json:"created_at"`
}
