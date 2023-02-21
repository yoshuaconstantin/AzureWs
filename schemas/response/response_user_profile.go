package response

import (
	"AzureWS/schemas/models"
)


type ResponseGetDataUserProfile struct {
	Message string                           `json:"message,omitempty"`
	Status  int                              `json:"status,omitempty"`
	Data    []models.GetUserProfileDataModel `json:"data,omitempty"`
}
