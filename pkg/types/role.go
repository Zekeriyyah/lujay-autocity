package types


type Role string

const (
	RoleAdmin   Role = "admin"
	RoleSeller  Role = "seller"
	RoleBuyer   Role = "buyer"
)