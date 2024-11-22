package model

type SignUpRequest struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignOutRequest struct {
	Username string `json:"username"`
}

type GetPermissionHandlerRequest struct {
	Method string
	Path   string
}

type PermissionHandler struct {
	ID     int64
	Method string
	Path   string
	Roles  []int64
}
