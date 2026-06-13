package frontend

import "github.com/gogf/gf/v2/frame/g"

type CategoryListReq struct {
	g.Meta `path:"/category/list" method:"get" tags:"Frontend Categories" summary:"Frontend category list"`
}

type CategoryListRes struct {
	List  []CategoryItem `json:"list"`
	Total int            `json:"total"`
}

type CategoryItem struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Sort   uint   `json:"sort"`
	Status int    `json:"status"`
}
