package category

import "github.com/gogf/gf/v2/frame/g"

type ListReq struct {
	g.Meta `path:"/list" method:"get" tags:"商品分类" summary:"商品分类列表"`
}

type ListRes struct {
	List  []CategoryItem `json:"list"`
	Total int            `json:"total"`
}

type CategoryItem struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Sort   uint   `json:"sort"`
	Status int    `json:"status"`
}
