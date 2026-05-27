package praise

import (
	"context"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sPraise struct{}

func init() {
	service.RegisterPraise(New())
}

func New() *sPraise {
	return &sPraise{}
}

func (*sPraise) AddPraise(ctx context.Context, in model.AddPraiseInput) (res *model.AddPraiseOutput, err error) {
	in.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))
	id, err := dao.PraiseInfo.Ctx(ctx).InsertAndGetId(in)
	if err != nil {
		return &model.AddPraiseOutput{}, err
	}
	return &model.AddPraiseOutput{Id: gconv.Uint(id)}, nil
}

// 兼容处理：优先根据点赞id删除，点赞id为0；再根据对象id和type删除
func (*sPraise) DeletePraise(ctx context.Context, in model.DeletePraiseInput) (res *model.DeletePraiseOutput, err error) {
	//优先根据点赞id删除
	if in.Id != 0 {
		_, err = dao.PraiseInfo.Ctx(ctx).WherePri(in.Id).Delete()
		if err != nil {
			return nil, err
		}
		return &model.DeletePraiseOutput{Id: gconv.Uint(in.Id)}, nil
	} else {
		//	点赞id为0；再根据对象id和type删除
		in.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))
		// 先查询点赞记录ID
		var praiseInfo entity.PraiseInfo
		err = dao.PraiseInfo.Ctx(ctx).Where(dao.PraiseInfo.Columns().UserId, in.UserId).
			Where(dao.PraiseInfo.Columns().ObjectId, in.ObjectId).
			Where(dao.PraiseInfo.Columns().Type, in.Type).
			Scan(&praiseInfo)
		if err != nil {
			return &model.DeletePraiseOutput{}, err
		}
		
		// 删除点赞记录
		_, err = dao.PraiseInfo.Ctx(ctx).OmitEmpty(). //注意：需要过滤空值
			Where(in).Delete()
		if err != nil {
			return &model.DeletePraiseOutput{}, err
		}
		return &model.DeletePraiseOutput{Id: gconv.Uint(praiseInfo.Id)}, nil
	}
}

// 列表
// GetList 查询内容列表
func (*sPraise) GetList(ctx context.Context, in model.PraiseListInput) (out *model.PraiseListOutput, err error) {
	//1.获得*gdb.Model对象，方便后续调用
	userId := gconv.Uint(ctx.Value(consts.CtxUserId))
	m := dao.PraiseInfo.Ctx(ctx).Where(dao.PraiseInfo.Columns().Type, in.Type).
		Where(dao.PraiseInfo.Columns().UserId, userId)
	//2. 实例化响应结构体
	out = &model.PraiseListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 先查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		out.List = make([]model.PraiseListOutputItem, 0, 0)
		return out, err
	}
	
	//4. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.PraiseListOutputItem, 0, in.Size)
	
	//5. 分页查询
	listModel := m.Page(in.Page, in.Size).OrderDesc(dao.PraiseInfo.Columns().CreatedAt)
	
	//6. 根据类型查询数据
	if in.Type == consts.CollectionTypeGoods {
		// 查询点赞的商品
		var praises []entity.PraiseInfo
		err = listModel.Scan(&praises)
		if err != nil {
			return out, err
		}
		
		// 查询商品详情
		for _, praise := range praises {
			var goods entity.GoodsInfo
			err = dao.GoodsInfo.Ctx(ctx).WherePri(praise.ObjectId).Scan(&goods)
			if err != nil {
				continue // 如果商品不存在，跳过
			}
			
			goodsItem := model.GoodsItem{
				Id:     gconv.Uint(goods.Id),
				Name:   goods.Name,
				PicUrl: goods.PicUrl,
				Price:  goods.Price,
			}
			
			out.List = append(out.List, model.PraiseListOutputItem{
				Id:        praise.Id,
				UserId:    praise.UserId,
				ObjectId:  praise.ObjectId,
				Type:      praise.Type,
				Goods:     goodsItem,
				CreatedAt: praise.CreatedAt,
				UpdatedAt: praise.UpdatedAt,
			})
		}
	} else if in.Type == consts.CollectionTypeArticle {
		// 查询点赞的文章
		var praises []entity.PraiseInfo
		err = listModel.Scan(&praises)
		if err != nil {
			return out, err
		}
		
		// 查询文章详情
		for _, praise := range praises {
			var article entity.ArticleInfo
			err = dao.ArticleInfo.Ctx(ctx).WherePri(praise.ObjectId).Scan(&article)
			if err != nil {
				continue // 如果文章不存在，跳过
			}
			
			articleItem := model.ArticleItem{
				Id:     gconv.Uint(article.Id),
				Title:  article.Title,
				Desc:   article.Desc,
				PicUrl: article.PicUrl,
			}
			
			out.List = append(out.List, model.PraiseListOutputItem{
				Id:        praise.Id,
				UserId:    praise.UserId,
				ObjectId:  praise.ObjectId,
				Type:      praise.Type,
				Article:   articleItem,
				CreatedAt: praise.CreatedAt,
				UpdatedAt: praise.UpdatedAt,
			})
		}
	}
	
	return
}

// 抽取获得收藏数量的方法 for 商品详情&文章详情
func (*sPraise) PraiseCount(ctx context.Context, objectId uint, collectionType uint8) (count int, err error) {
	condition := g.Map{
		dao.PraiseInfo.Columns().ObjectId: objectId,
		dao.PraiseInfo.Columns().Type:     collectionType,
	}
	count, err = dao.PraiseInfo.Ctx(ctx).Where(condition).Count()
	if err != nil {
		return 0, err
	}
	return
}

// 抽取方法 判断当前用户是否点赞 for 商品详情&文章详情
func (*sPraise) CheckIsPraise(ctx context.Context, in model.CheckIsCollectInput) (bool, error) {
	// 如果用户ID为0，表示用户未登录，直接返回false
	if in.UserId == 0 {
		return false, nil
	}
	
	condition := g.Map{
		dao.PraiseInfo.Columns().UserId:   in.UserId,
		dao.PraiseInfo.Columns().ObjectId: in.ObjectId,
		dao.PraiseInfo.Columns().Type:     in.Type,
	}
	count, err := dao.PraiseInfo.Ctx(ctx).Where(condition).Count()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
