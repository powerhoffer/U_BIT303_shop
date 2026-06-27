package model

type PermissionCreateInput struct {
	Name      string
	GroupName string
	Method    string
	Path      string
}

type PermissionCreateOutput struct {
	Permission PermissionBase
}

type PermissionListInput struct {
	Page      int
	Size      int
	Name      string
	GroupName string
	Method    string
	Status    int
}

type PermissionListOutput struct {
	List  []PermissionBase
	Total int
	Page  int
	Size  int
}

type PermissionDetailOutput struct {
	Permission PermissionBase
}

type PermissionUpdateInput struct {
	Id        uint
	Name      string
	GroupName string
	Method    string
	Path      string
}

type PermissionUpdateOutput struct {
	Permission PermissionBase
}

type PermissionStatusInput struct {
	Id     uint
	Status int
}

type PermissionBase struct {
	Id        uint
	Name      string
	GroupName string
	Method    string
	Path      string
	Status    int
}
