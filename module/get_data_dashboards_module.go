package module

import (
	"AzureWS/config"
	"AzureWS/schemas/models"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq" // postgres golang driver
)

func InitDashboardsDataSetToDB(userId string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO dashboards_data (user_id, profilecount, profilerank, thermalcount, thermalrank, dozecount, dozerank ) VALUES ($1, 0,'Nubie',0,'Nubie',0,'Nubie')`

	_, err := db.Exec(sqlStatement, userId)

	fmt.Printf("\nUser id %v\n", userId)

	if err != nil {
		log.Fatalf("tidak bisa mengeksekusi query init dashboards. %v\n", err)
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	fmt.Printf("Insert data single record into Dashboards data\n")

	return true, nil
}

func GetDashboardsDataFromDB(userId string) ([]models.DashboardsDataModel, error) {

	db := config.CreateConnection()

	defer db.Close()

	var dashboardsData []models.DashboardsDataModel

	sqlStatement := `SELECT profilecount, profilerank, thermalcount, thermalrank, dozecount, dozerank FROM dashboards_data WHERE user_id = $1`

	rows, err := db.Query(sqlStatement, userId)

	if err != nil {
		log.Fatalf("\nGET DASHBOARDS DATA ERROR - tidak bisa mengeksekusi query. %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		var dashboardData models.DashboardsDataModel

		err = rows.Scan(&dashboardData.ProfileCount, &dashboardData.ProfileRank, &dashboardData.ThermalCount, &dashboardData.ThermalRank, &dashboardData.DozeCount, &dashboardData.DozeRank)

		if err != nil {
			log.Fatalf("\nGET DASHBOARDS DATA ERROR - tidak bisa mengambil data semua dashboards. %v\n", err)
		}

		dashboardsData = append(dashboardsData, dashboardData)

	}

	// return empty buku atau jika error
	return dashboardsData, err
}

func UpdateDashboardsDataFromDB(userId string, mode string) (bool, error) {

	db := config.CreateConnection()

	defer db.Close()

	// map mode to column name
	modeMapCount := map[string]string{
		"profiles": "ProfileCount",
		"thermal":  "ThermalCount",
		"doze":     "DozeCount",
	}

	modeMapRank := map[string]string{
		"profiles": "ProfileRank",
		"thermal":  "ThermalRank",
		"doze":     "DozeRank",
	}

	column, ok := modeMapCount[mode]
	if !ok {
		return false, fmt.Errorf("%s %s", "\nUPDATE DASHBOARDS DATA ERROR - Invalid mode\n", userId)
	}

	columnRank, ok := modeMapRank[mode]
	if !ok {
		return false, fmt.Errorf("%s %s", "\nUPDATE DASHBOARDS DATA ERROR - Invalid mode\n", mode)
	}

	sqlStatement := `SELECT ` + column + ` FROM dashboards_data WHERE user_id = $1`

	var result sql.NullString
	err := db.QueryRow(sqlStatement, userId).Scan(&result)

	if err == sql.ErrNoRows {
		return false, err
	}

	if err != nil {
		log.Fatalf("\nUPDATE DASHBOARDS DATA ERROR - Error executing the SQL statement: %v\n", err)
		return false, err
	}

	if result.Valid {

		resultInt, err := strconv.Atoi(result.String)

		if err != nil {
			return false, err
		}

		var updatedRank string

		resultInt = resultInt + 1

		if resultInt > 20 {
			updatedRank = "Enjoyer"
		} else if resultInt > 60 {
			updatedRank = "Madness"
		} else if resultInt > 200 {
			updatedRank = "Crazy"
		} else if resultInt > 500 {
			updatedRank = "Legends"
		} else if resultInt > 2000 {
			updatedRank = "PRO"
		} else {
			updatedRank = "Nubie"
		}

		sqlStatement := `UPDATE dashboards_data SET ` + column + ` = $1, ` + columnRank + ` = $2  WHERE user_id = $3`

		//exec result + 1 each mode tapped
		res, errUpdate := db.Exec(sqlStatement, resultInt, updatedRank, userId)

		if errUpdate != nil {
			log.Fatalf("\nUPDATE DASHBOARDS DATA ERROR - Error executing the SQL statement: %v\n", err)
			return false, errUpdate
		}

		rowsAffected, err := res.RowsAffected()
		if errUpdate != nil {
			return false, err
		}

		// if update count succes, then next
		if rowsAffected == 1 {
			return true, nil
		} else {
			return false, fmt.Errorf("%s %d", "\nUPDATE DASHBOARDS DATA ERROR - Expected to affect 1 row, but affected \n", rowsAffected)
		}

	} else {
		return false, err
	}
}
