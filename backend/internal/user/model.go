package user

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type CreateUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginUserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
	IsAdmin     bool   `json:"is_admin"`
}
