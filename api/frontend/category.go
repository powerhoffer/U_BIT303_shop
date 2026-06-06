package frontend

import "github.com/gogf/gf/v2/frame/g"

type CategoryListReq struct {
	g.Meta `path:"/category/list" method:"get" tags:"前台商品分类" summary:"前台商品分类列表"`
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
