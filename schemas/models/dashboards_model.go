package models

type DashboardsDataModel struct {
	ProfileCount *int    `json:"profileCount,omitempty"`
	ProfileRank  *string `json:"profileRank,omitempty"`
	ThermalCount *int    `json:"thermalCount,omitempty"`
	ThermalRank  *string `json:"thermalRank,omitempty"`
	DozeCount    *int    `json:"dozeCount,omitempty"`
	DozeRank     *string `json:"dozeRank,omitempty"`
}