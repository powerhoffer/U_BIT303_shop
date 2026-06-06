package points

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type BalanceReq struct {
	g.Meta `path:"/balance" method:"get" tags:"积分管理" summary:"当前员工积分余额"`
}

type BalanceRes struct {
	Balance uint `json:"balance"`
}

type RecordsReq struct {
	g.Meta `path:"/records" method:"get" tags:"积分管理" summary:"当前员工积分流水"`
	Page   int `json:"page" in:"query" d:"1" v:"min:1#分页号码错误"`
	Size   int `json:"size" in:"query" d:"10" v:"max:50#分页数量最大50条"`
}

type RecordsRes struct {
	List  []RecordItem `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

type ManageAddReq struct {
	g.Meta     `path:"/manage/add" method:"post" tags:"积分管理" summary:"管理员增加员工积分"`
	EmployeeId uint   `json:"employee_id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Points     uint   `json:"points" v:"required|min:1#积分不能为空|积分必须大于0"`
	Remark     string `json:"remark" v:"max-length:255#备注最多255位"`
}

type ManageAddRes struct {
	Balance uint `json:"balance"`
}

type ManageDeductReq struct {
	g.Meta     `path:"/manage/deduct" method:"post" tags:"积分管理" summary:"管理员扣除员工积分"`
	EmployeeId uint   `json:"employee_id" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Points     uint   `json:"points" v:"required|min:1#积分不能为空|积分必须大于0"`
	Remark     string `json:"remark" v:"max-length:255#备注最多255位"`
}

type ManageDeductRes struct {
	Balance uint `json:"balance"`
}

type ManageRecordsReq struct {
	g.Meta     `path:"/manage/records" method:"get" tags:"积分管理" summary:"管理员查看员工积分流水"`
	EmployeeId uint `json:"employee_id" in:"query" v:"required|min:1#员工ID不能为空|员工ID错误"`
	Page       int  `json:"page" in:"query" d:"1" v:"min:1#分页号码错误"`
	Size       int  `json:"size" in:"query" d:"10" v:"max:50#分页数量最大50条"`
}

type ManageRecordsRes struct {
	List  []RecordItem `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

type RecordItem struct {
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
