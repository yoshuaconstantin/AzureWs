package models

type UserModel struct {
	ID       *int64  `json:"id,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	UserId   *string `json:"user_id,omitempty"`
}

type TokenWithJwtModel struct {
	Token	string `json:"token"`
	JWT		string `json:"jwt"`
}