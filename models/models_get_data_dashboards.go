package models

import (
	"log"
	"AzureWS/config"

	_ "github.com/lib/pq" // postgres golang driver
)


type DashboardsData struct {
	ID				*int64		`json:"id,omitempty"`
	ProfileCount	*int		`json:"profileCount,omitempty"`
	ProfileRank		*string		`json:"profileRank,omitempty"`
	ThermalCount	*int		`json:"thermalCount,omitempty"`
	ThermalRank		*string		`json:"thermalRank,omitempty"`
	DozeCount		*int		`json:"dozeCount,omitempty"`
	DozeRank		*string		`json:"dozeRank,omitempty"`
	Token			*string		`json:"token,omitempty"`
}


func GetDashboardsData(token string) ([]DashboardsData, error) {
	
	db := config.CreateConnection()

	defer db.Close()

	var dashboardsData []DashboardsData

	sqlStatement := `SELECT * FROM dashboards_data WHERE token = $1`

	
	rows, err := db.Query(sqlStatement, token)

	if err != nil {
		log.Fatalf("tidak bisa mengeksekusi query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var dashboardData DashboardsData

		err = rows.Scan(&dashboardData.ID, &dashboardData.ProfileCount, &dashboardData.ThermalCount, &dashboardData.ThermalRank, &dashboardData.DozeCount, &dashboardData.DozeRank)

		if err != nil {
			log.Fatalf("tidak bisa mengambil data semua dashboards. %v", err)
		}

		dashboardsData = append(dashboardsData, dashboardData)

	}

	// return empty buku atau jika error
	return dashboardsData, err
}

