package model

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
