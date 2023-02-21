package response

import (
	"AzureWS/schemas/models"
)

type ResponseGetAllFeedBackUser struct {
	Status  int                              `json:"status"`
	Message string                           `json:"message"`
	Data    []models.ReturnFeedBackUserModel `json:"data"`
}