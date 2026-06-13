package base

import "github.com/gogf/gf/v2/frame/g"

type HealthReq struct {
	g.Meta `path:"/health" method:"get" tags:"Base" summary:"Health check"`
}

type HealthRes struct {
	Status string `json:"status"`
}
