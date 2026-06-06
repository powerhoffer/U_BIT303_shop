package model

import "time"

type PointsBalanceOutput struct {
	Balance uint
}

type PointsRecordsInput struct {
	EmployeeId uint
	Page       int
	Size       int
}

type PointsRecordsOutput struct {
	List  []PointsRecordItem
	Total int
	Page  int
	Size  int
}

type PointsChangeInput struct {
	EmployeeId         uint
	OperatorEmployeeId uint
	Points             uint
	Remark             string
}

type PointsChangeOutput struct {
	Balance uint
}

type PointsRecordItem struct {
	Id                 uint
	EmployeeId         uint
	ChangeType         int
	Points             uint
	BeforeBalance      uint
	AfterBalance       uint
	OperatorEmployeeId uint
	Remark             string
	CreatedAt          time.Time
}
