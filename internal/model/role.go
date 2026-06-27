package model

type RoleCreateInput struct {
	Name        string
	Description string
}

type RoleCreateOutput struct {
	Role RoleBase
}

type RoleListInput struct {
	Page   int
	Size   int
	Name   string
	Status int
}

type RoleListOutput struct {
	List  []RoleBase
	Total int
	Page  int
	Size  int
}

type RoleDetailOutput struct {
	Role RoleBase
}

type RoleUpdateInput struct {
	Id          uint
	Name        string
	Description string
}

type RoleUpdateOutput struct {
	Role RoleBase
}

type RoleStatusInput struct {
	Id     uint
	Status int
}

type RolePermissionsInput struct {
	Id            uint
	PermissionIds []uint
}

type RoleBase struct {
	Id            uint
	Name          string
	Description   string
	Status        int
	PermissionIds []uint
}
