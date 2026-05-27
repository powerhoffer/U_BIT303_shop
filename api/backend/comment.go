package backend

import "github.com/gogf/gf/v2/frame/g"

// 后台评论列表请求
type CommentListReq struct {
	g.Meta `path:"/comment/list" method:"get" tags:"评论管理" summary:"评论列表"`
	Type   uint8 `json:"type" dc:"评论类型：1商品 2文章"`
	CommonPaginationReq
}

type CommentListRes struct {
	Page  int         `json:"page" description:"分页码"`
	Size  int         `json:"size" description:"分页数量"`
	Total int         `json:"total" description:"数据总数"`
	List  interface{} `json:"list" description:"列表"`
}

// 后台删除评论请求
type CommentDeleteReq struct {
	g.Meta    `path:"/comment/delete" method:"delete" tags:"评论管理" summary:"删除评论"`
	CommentId uint `json:"comment_id" v:"min:1#请选择需要删除的评论" dc:"评论id"`
}

type CommentDeleteRes struct{}
