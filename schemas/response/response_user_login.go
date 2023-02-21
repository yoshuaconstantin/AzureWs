package response

import (
	"AzureWS/schemas/models"
)

type ResponseUserLoginWithJWT struct {
	Message  string `json:"message,omitempty"`
	Status   int    `json:"status,omitempty"`
	JwtToken string `json:"jtwToken,omitempty"`
}

type ResponseUserLogin struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    []models.UserModel `json:"data"`
}

type ResponseUserTokenAndJwt struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Token    string `json:"token"`
	JwtToken string `json:"jwtToken"`
}
