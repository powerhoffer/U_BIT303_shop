// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"bit303_shop/internal/model"
)

type (
	ICollection interface {
		AddCollection(ctx context.Context, in model.AddCollectionInput) (res *model.AddCollectionOutput, err error)
		// 兼容处理：优先根据收藏id删除，收藏id为0；再根据对象id和type删除
		DeleteCollection(ctx context.Context, in model.DeleteCollectionInput) (res *model.DeleteCollectionOutput, err error)
		// 列表
		// GetList 查询内容列表
		GetList(ctx context.Context, in model.CollectionListInput) (out *model.CollectionListOutput, err error)
		// 获取收藏数量
		CollectionCount(ctx context.Context, objectId uint, collectionType uint8) (count int, err error)
		// 判断当前用户是否收藏
		CheckIsCollect(ctx context.Context, in model.CheckIsCollectInput) (bool, error)
	}
)

var (
	localCollection ICollection
)

func Collection() ICollection {
	if localCollection == nil {
		panic("implement not found for interface ICollection, forgot register?")
	}
	return localCollection
}

func RegisterCollection(i ICollection) {
	localCollection = i
}
