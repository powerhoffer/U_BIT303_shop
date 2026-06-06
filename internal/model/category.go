package model

type CategoryListOutput struct {
	List  []CategoryItem
	Total int
}

type CategoryItem struct {
	Id     uint
	Name   string
	Sort   uint
	Status int
}
