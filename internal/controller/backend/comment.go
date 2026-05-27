package backend

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

// 保留 model 包引用用于 List 方法

// Comment 后台评论管理
var Comment = cComment{}

type cComment struct{}

// List 评论列表
func (c *cComment) List(ctx context.Context, req *backend.CommentListReq) (res *backend.CommentListRes, err error) {
	getListRes, err := service.Comment().GetList(ctx, model.CommentListInput{
		Page: req.Page,
		Size: req.Size,
		Type: req.Type,
	})
	if err != nil {
		return nil, err
	}
	return &backend.CommentListRes{
		List:  getListRes.List,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total,
	}, nil
}

// Delete 删除评论
func (c *cComment) Delete(ctx context.Context, req *backend.CommentDeleteReq) (res *backend.CommentDeleteRes, err error) {
	// 直接从请求中获取 comment_id，避免被 JWT 的 id 覆盖
	r := g.RequestFromCtx(ctx)
	commentId := r.Get("comment_id").Uint()
	err = service.Comment().AdminDeleteComment(ctx, commentId)
	return &backend.CommentDeleteRes{}, err
}
