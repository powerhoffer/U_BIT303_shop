package base

import "github.com/gogf/gf/v2/frame/g"

type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"基础" summary:"健康检查"`
}

type HealthRes struct {
	Status string `json:"status"`
}
