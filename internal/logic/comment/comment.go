package comment

import (
	"context"
	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model"
	"bit303_shop/internal/model/do"
	"bit303_shop/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sComment struct{}

func init() {
	service.RegisterComment(New())
}

func New() *sComment {
	return &sComment{}
}

func (*sComment) AddComment(ctx context.Context, in model.AddCommentInput) (res *model.AddCommentOutput, err error) {
	userId := ctx.Value(consts.CtxUserId)
	in.UserId = gconv.Uint(userId)
	g.Log().Infof(ctx, "添加评论，用户ID: %v, 参数: %+v", userId, in)
	id, err := dao.CommentInfo.Ctx(ctx).InsertAndGetId(in)
	if err != nil {
		g.Log().Errorf(ctx, "添加评论失败: %v", err)
		return &model.AddCommentOutput{}, err
	}
	g.Log().Infof(ctx, "添加评论成功，ID: %v", id)
	return &model.AddCommentOutput{Id: gconv.Uint(id)}, nil
}

// 兼容处理：优先根据收藏id删除，收藏id为0；再根据对象id和type删除
func (*sComment) DeleteComment(ctx context.Context, in model.DeleteCommentInput) (res *model.DeleteCommentOutput, err error) {
	condition := g.Map{
		dao.CommentInfo.Columns().Id:     in.Id,
		dao.CommentInfo.Columns().UserId: ctx.Value(consts.CtxUserId),
	}
	_, err = dao.CommentInfo.Ctx(ctx).Where(condition).Delete()
	if err != nil {
		return nil, err
	}
	return &model.DeleteCommentOutput{Id: gconv.Uint(in.Id)}, nil
}

// AdminDeleteComment 后台管理员删除评论（不检查用户ID）
func (*sComment) AdminDeleteComment(ctx context.Context, id uint) error {
	_, err := dao.CommentInfo.Ctx(ctx).Where(dao.CommentInfo.Columns().Id, id).Delete()
	return err
}

// GetList 查询内容列表 TODO:评论列表的查询逻辑处理按uid查.应该还要加按parent_id查,
func (*sComment) GetList(ctx context.Context, in model.CommentListInput) (out *model.CommentListOutput, err error) {
	//1.获得*gdb.Model对象，方面后续调用
	m := dao.CommentInfo.Ctx(ctx)
	// 只有指定了 Type 和 ObjectId 时才过滤（前台按文章/商品查询）
	if in.Type > 0 {
		m = m.Where(dao.CommentInfo.Columns().Type, in.Type)
	}
	if in.ObjectId > 0 {
		m = m.Where(dao.CommentInfo.Columns().ObjectId, in.ObjectId)
	}
	//2. 实例化响应结构体
	out = &model.CommentListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	//3. 分页查询
	listModel := m.Page(in.Page, in.Size)
	//4. 再查询count，判断有无数据
	out.Total, err = m.Count()
	if err != nil || out.Total == 0 {
		return out, err
	}
	//5. 延迟初始化list切片 确定有数据，再按期望大小初始化切片容量
	out.List = make([]model.CommentBase, 0, in.Size)
	//6. 使用Join明确加载用户关联数据
	g.Log().Infof(ctx, "查询评论列表，参数: %+v", in)

	// 先获取所有评论数据，包含用户信息
	type CommentWithUser struct {
		do.CommentInfo
		UserName   string `json:"user_name"`
		UserAvatar string `json:"user_avatar"`
	}

	var commentsWithUser []CommentWithUser
	if err := listModel.LeftJoin("user_info", "comment_info.user_id = user_info.id").
		Fields("comment_info.*, user_info.name as user_name, user_info.avatar as user_avatar").
		Scan(&commentsWithUser); err != nil {
		g.Log().Errorf(ctx, "查询评论列表失败: %v", err)
		return out, err
	}

	// 创建评论映射，便于查找父评论
	commentMap := make(map[int]*CommentWithUser)
	for i := range commentsWithUser {
		commentId := gconv.Int(commentsWithUser[i].Id)
		commentMap[commentId] = &commentsWithUser[i]
	}

	// 转换为最终的CommentBase结构并处理回复关系
	for _, comment := range commentsWithUser {
		commentBase := model.CommentBase{
			CommentInfo: comment.CommentInfo,
			User: model.UserInfoBase{
				Id:     uint(gconv.Int(comment.UserId)),
				Name:   comment.UserName,
				Avatar: comment.UserAvatar,
			},
		}

		// 处理回复关系
		parentId := gconv.Int(comment.ParentId)
		if parentId != 0 {
			if parentComment, exists := commentMap[parentId]; exists {
				commentBase.ReplyTo = parentComment.UserName
				g.Log().Infof(ctx, "设置回复关系: 评论ID=%v 回复用户=%s", comment.Id, commentBase.ReplyTo)
			} else {
				g.Log().Warningf(ctx, "未找到父评论: 评论ID=%v, ParentId=%v", comment.Id, comment.ParentId)
			}
		}

		out.List = append(out.List, commentBase)
		g.Log().Infof(ctx, "评论: ID=%v, UserName=%s, ParentId=%v, ReplyTo=%s",
			comment.Id, comment.UserName, comment.ParentId, commentBase.ReplyTo)
	}

	// 记录调试信息 - 检查parent_id的值
	for _, comment := range out.List {
		g.Log().Infof(ctx, "评论数据检查: ID=%d, ParentId=%d (类型: %T), UserName=%s, Content=%s",
			comment.Id, comment.ParentId, comment.ParentId, comment.User.Name, comment.Content)
	}

	// 统计parent_id为0的评论数量
	var directCommentCount int
	for _, comment := range out.List {
		if comment.ParentId == 0 {
			directCommentCount++
		}
	}
	g.Log().Infof(ctx, "直接评论数量(ParentId=0): %d, 总评论数: %d", directCommentCount, len(out.List))

	g.Log().Infof(ctx, "查询到评论数量: %d", len(out.List))
	return
}

// 抽取获得收藏数量的方法 for 商品详情&文章详情
func CommentCount(ctx context.Context, objectId uint, collectionType uint8) (count int, err error) {
	condition := g.Map{
		dao.CommentInfo.Columns().ObjectId: objectId,
		dao.CommentInfo.Columns().Type:     collectionType,
	}
	count, err = dao.CommentInfo.Ctx(ctx).Where(condition).Count()
	if err != nil {
		return 0, err
	}
	return
}

// 抽取方法 判断当前用户是否收藏 for 商品详情&文章详情
func CheckIsComment(ctx context.Context, in model.CheckIsCollectInput) (bool, error) {
	condition := g.Map{
		dao.CommentInfo.Columns().UserId:   ctx.Value(consts.CtxUserId),
		dao.CommentInfo.Columns().ObjectId: in.ObjectId,
		dao.CommentInfo.Columns().Type:     in.Type,
	}
	count, err := dao.CommentInfo.Ctx(ctx).Where(condition).Count()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}
