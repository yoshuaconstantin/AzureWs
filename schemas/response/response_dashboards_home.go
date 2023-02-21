package response

import (
	"AzureWS/schemas/models"
)

type ResponseDashboardsData struct {
	Status  int                          `json:"status"`
	Message string                       `json:"message"`
	Data    []models.DashboardsDataModel `json:"data"`
}
