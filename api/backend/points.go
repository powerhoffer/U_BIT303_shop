package backend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type PointsBalanceReq struct {
	g.Meta `path:"/points/balance" method:"get" tags:"后台积分" summary:"当前员工积分余额"`
}

type PointsBalanceRes struct {
	Balance uint `json:"balance"`
}

type PointsRecordsReq struct {
	g.Meta `path:"/points/records" method:"get" tags:"后台积分" summary:"当前员工积分流水"`
	Page   int `json:"page" in:"query" d:"1" v:"min:1#分页页码错误"`
	Size   int `json:"size" in:"query" d:"10" v:"max:50#分页数量最多50条"`
}

type PointsRecordsRes struct {
	List  []PointsRecordItem `json:"list"`
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

type PointsManageAddReq struct {
	g.Meta     `path:"/points/manage/add" method:"post" tags:"后台积分" summary:"管理员增加员工积分"`
	EmployeeId uint   `json:"employee_id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Points     uint   `json:"points" v:"required|min:1#积分不能为空|积分必须大于0"`
	Remark     string `json:"remark" v:"max-length:255#备注最多255位"`
}

type PointsManageAddRes struct {
	Balance uint `json:"balance"`
}

type PointsManageDeductReq struct {
	g.Meta     `path:"/points/manage/deduct" method:"post" tags:"后台积分" summary:"管理员扣除员工积分"`
	EmployeeId uint   `json:"employee_id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Points     uint   `json:"points" v:"required|min:1#积分不能为空|积分必须大于0"`
	Remark     string `json:"remark" v:"max-length:255#备注最多255位"`
}

type PointsManageDeductRes struct {
	Balance uint `json:"balance"`
}

type PointsManageRecordsReq struct {
	g.Meta     `path:"/points/manage/records" method:"get" tags:"后台积分" summary:"管理员查看员工积分流水"`
	EmployeeId uint `json:"employee_id" in:"query" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Page       int  `json:"page" in:"query" d:"1" v:"min:1#分页页码错误"`
	Size       int  `json:"size" in:"query" d:"10" v:"max:50#分页数量最多50条"`
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
