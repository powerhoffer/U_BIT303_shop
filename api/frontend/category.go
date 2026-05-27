package frontend

import "github.com/gogf/gf/v2/frame/g"

// CategoryGetListReq 前台分类列表请求
type CategoryGetListReq struct {
	g.Meta   `path:"/category/list" method:"get" tags:"前台商品分类" summary:"前台商品分类列表接口"`
	Sort     int  `json:"sort" in:"query" dc:"排序类型"`
	ParentId uint `json:"parent_id" in:"query" dc:"父级ID"`
	Level    uint `json:"level" in:"query" dc:"分类级别"`
	CommonPaginationReq
}

type CategoryGetListRes struct {
	List  interface{} `json:"list" description:"列表"`
	Page  int         `json:"page" description:"分页码"`
	Size  int         `json:"size" description:"分页数量"`
	Total int         `json:"total" description:"数据总数"`
}

// CategoryGetListAllReq 前台分类全部列表请求
type CategoryGetListAllReq struct {
	g.Meta `path:"/category/list/all" method:"get" tags:"前台商品分类" summary:"前台商品分类全部列表"`
}

type CategoryGetListAllRes struct {
	List  interface{} `json:"list" description:"列表"`
	Total int         `json:"total" description:"数据总数"`
}

// CategoryGetHierarchicalReq 前台分类层级列表请求
type CategoryGetHierarchicalReq struct {
	g.Meta `path:"/category/hierarchical" method:"get" tags:"前台商品分类" summary:"前台商品分类层级列表"`
}

type CategoryGetHierarchicalRes struct {
	List interface{} `json:"list" description:"层级列表"`
}
