package collection

import (
	"context"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sCollection struct{}

func init() {
	service.RegisterCollection(New())
}

func New() *sCollection {
	return &sCollection{}
}

func (*sCollection) AddCollection(ctx context.Context, in model.AddCollectionInput) (res *model.AddCollectionOutput, err error) {
	in.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))

	// 先检查是否已经收藏过
	count, err := dao.CollectionInfo.Ctx(ctx).Where(dao.CollectionInfo.Columns().UserId, in.UserId).
		Where(dao.CollectionInfo.Columns().ObjectId, in.ObjectId).
		Where(dao.CollectionInfo.Columns().Type, in.Type).
		Count()
	if err != nil {
		return &model.AddCollectionOutput{}, err
	}
	if count > 0 {
		// 已经收藏过，直接返回，不重复添加
		return &model.AddCollectionOutput{}, nil
	}

	// 没有收藏过，添加收藏
	id, err := dao.CollectionInfo.Ctx(ctx).InsertAndGetId(in)
	if err != nil {
		return &model.AddCollectionOutput{}, err
	}
	return &model.AddCollectionOutput{Id: gconv.Uint(id)}, nil
}

// 兼容处理：优先根据收藏id删除，收藏id为0；再根据对象id和type删除
func (*sCollection) DeleteCollection(ctx context.Context, in model.DeleteCollectionInput) (res *model.DeleteCollectionOutput, err error) {
	//优先根据收藏id删除
	if in.Id != 0 {
		_, err = dao.CollectionInfo.Ctx(ctx).WherePri(in.Id).Delete()
		if err != nil {
			return nil, err
		}
		return &model.DeleteCollectionOutput{Id: gconv.Uint(in.Id)}, nil
	} else {
		//	收藏id为0；再根据对象id和type删除
		in.UserId = gconv.Uint(ctx.Value(consts.CtxUserId))
		// 先查询收藏记录ID
		var collectionInfo entity.CollectionInfo
		err = dao.CollectionInfo.Ctx(ctx).Where(dao.CollectionInfo.Columns().UserId, in.UserId).
			Where(dao.CollectionInfo.Columns().ObjectId, in.ObjectId).
			Where(dao.CollectionInfo.Columns().Type, in.Type).
			Scan(&collectionInfo)
		if err != nil {
			// 如果记录不存在，直接返回成功
			if err.Error() == gcode.CodeNotFound.String() {
				return &model.DeleteCollectionOutput{Id: 0}, nil
			}
			return &model.DeleteCollectionOutput{}, err
		}

		// 删除收藏记录
		_, err = dao.CollectionInfo.Ctx(ctx).Where(dao.CollectionInfo.Columns().UserId, in.UserId).
			Where(dao.CollectionInfo.Columns().ObjectId, in.ObjectId).
			Where(dao.CollectionInfo.Columns().Type, in.Type).
			Delete()
		if err != nil {
			return &model.DeleteCollectionOutput{}, err
		}
		return &model.DeleteCollectionOutput{Id: gconv.Uint(collectionInfo.Id)}, nil
	}
}

// 列表
// GetList 查询内容列表
func (*sCollection) GetList(ctx context.Context, in model.CollectionListInput) (out *model.CollectionListOutput, err error) {
	//1.获得*gdb.Model对象，方便后续调用
	userId := gconv.Uint(ctx.Value(consts.CtxUserId))
	m := dao.CollectionInfo.Ctx(ctx).Where(dao.CollectionInfo.Columns().Type, in.Type).
		Where(dao.CollectionInfo.Columns().UserId, userId)
	//2. 实例化响应结构体
	out = &model.CollectionListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 先查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		out.List = make([]model.CollectionListOutputItem, 0, 0)
		return out, err
	}

	//4. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.CollectionListOutputItem, 0, in.Size)

	//5. 分页查询
	listModel := m.Page(in.Page, in.Size).OrderDesc(dao.CollectionInfo.Columns().CreatedAt)

	//6. 根据类型查询数据
	if in.Type == consts.CollectionTypeGoods {
		// 查询收藏的商品
		var collections []entity.CollectionInfo
		err = listModel.Scan(&collections)
		if err != nil {
			return out, err
		}

		// 查询商品详情
		for _, collection := range collections {
			var goods entity.GoodsInfo
			err = dao.GoodsInfo.Ctx(ctx).WherePri(collection.ObjectId).Scan(&goods)
			if err != nil {
				continue // 如果商品不存在，跳过
			}

			goodsItem := model.GoodsItem{
				Id:     gconv.Uint(goods.Id),
				Name:   goods.Name,
				PicUrl: goods.PicUrl,
				Price:  goods.Price,
			}

			out.List = append(out.List, model.CollectionListOutputItem{
				Id:        collection.Id,
				UserId:    collection.UserId,
				ObjectId:  collection.ObjectId,
				Type:      collection.Type,
				Goods:     goodsItem,
				CreatedAt: collection.CreatedAt,
				UpdatedAt: collection.UpdatedAt,
			})
		}
	} else if in.Type == consts.CollectionTypeArticle {
		// 查询收藏的文章
		var collections []entity.CollectionInfo
		err = listModel.Scan(&collections)
		if err != nil {
			return out, err
		}

		// 查询文章详情
		for _, collection := range collections {
			var article entity.ArticleInfo
			err = dao.ArticleInfo.Ctx(ctx).WherePri(collection.ObjectId).Scan(&article)
			if err != nil {
				continue // 如果文章不存在，跳过
			}

			articleItem := model.ArticleItem{
				Id:     gconv.Uint(article.Id),
				Title:  article.Title,
				Desc:   article.Desc,
				PicUrl: article.PicUrl,
			}

			out.List = append(out.List, model.CollectionListOutputItem{
				Id:        collection.Id,
				UserId:    collection.UserId,
				ObjectId:  collection.ObjectId,
				Type:      collection.Type,
				Article:   articleItem,
				CreatedAt: collection.CreatedAt,
				UpdatedAt: collection.UpdatedAt,
			})
		}
	}

	return
}

// 抽取获得收藏数量的方法 for 商品详情&文章详情
func (s *sCollection) CollectionCount(ctx context.Context, objectId uint, collectionType uint8) (count int, err error) {
	condition := g.Map{
		dao.CollectionInfo.Columns().ObjectId: objectId,
		dao.CollectionInfo.Columns().Type:     collectionType,
	}
	count, err = dao.CollectionInfo.Ctx(ctx).Where(condition).Count()
	if err != nil {
		return 0, err
	}
	return
}

// 抽取方法 判断当前用户是否收藏 for 商品详情&文章详情
func (s *sCollection) CheckIsCollect(ctx context.Context, in model.CheckIsCollectInput) (bool, error) {
	// 如果用户ID为0，表示用户未登录，直接返回false
	if in.UserId == 0 {
		return false, nil
	}

	condition := g.Map{
		dao.CollectionInfo.Columns().UserId:   in.UserId,
		dao.CollectionInfo.Columns().ObjectId: in.ObjectId,
		dao.CollectionInfo.Columns().Type:     in.Type,
	}
	count, err := dao.CollectionInfo.Ctx(ctx).Where(condition).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
