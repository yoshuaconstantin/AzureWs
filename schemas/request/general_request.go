package request

type RequestToken struct {
	Token string `json:"token"`
}

type RequestTokenWithIndex struct {
	Token string	`json:"token"`
	Index int		`json:"index"`
}