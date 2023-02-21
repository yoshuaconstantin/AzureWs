package request

type RequestLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestChangePasswordData struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type RequestTokenData struct {
	Token string `json:"token"`
}