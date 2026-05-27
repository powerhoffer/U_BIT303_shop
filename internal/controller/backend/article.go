package backend

import (
	"context"
	"bit303_shop/api/backend"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

// Article 内容管理
var Article = cArticle{}

type cArticle struct{}

func (a *cArticle) Create(ctx context.Context, req *backend.ArticleReq) (res *backend.ArticleRes, err error) {
	data := model.ArticleCreateInput{}
	err = gconv.Scan(req, &data)
	if err != nil {
		return nil, err
	}
	data.UserId = gconv.Int(ctx.Value(consts.CtxAdminId))
	out, err := service.Article().Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return &backend.ArticleRes{Id: out.Id}, nil
}

func (a *cArticle) Delete(ctx context.Context, req *backend.ArticleDeleteReq) (res *backend.ArticleDeleteRes, err error) {
	err = service.Article().Delete(ctx, model.ArticleDeleteInput{Id: req.ArticleId})
	return
}

func (a *cArticle) Update(ctx context.Context, req *backend.ArticleUpdateReq) (res *backend.ArticleUpdateRes, err error) {
	data := model.ArticleUpdateInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	// 使用 ArticleId 而不是 Id，避免被 JWT 的 id 覆盖
	data.Id = req.ArticleId
	//获取当前登录用户
	data.UserId = gconv.Int(ctx.Value(consts.CtxAdminId))
	// 后台管理员操作，设置 IsAdmin 为管理员标识，跳过权限检查
	data.IsAdmin = consts.ArticleIsAdmin
	err = service.Article().Update(ctx, data)
	if err != nil {
		return nil, err
	}
	return &backend.ArticleUpdateRes{Id: req.ArticleId}, nil
}

func (a *cArticle) List(ctx context.Context, req *backend.ArticleGetListCommonReq) (res *backend.ArticleGetListCommonRes, err error) {
	getListRes, err := service.Article().GetList(ctx, model.ArticleGetListInput{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &backend.ArticleGetListCommonRes{List: getListRes.List,
		Page:  getListRes.Page,
		Size:  getListRes.Size,
		Total: getListRes.Total}, nil
}
