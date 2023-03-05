package logging

import (
	"io/ioutil"
	"log"
	"net/http"

	"AzureWS/config"
	Gv "AzureWS/globalvariable/variable"
	"AzureWS/validation"
)

func LogsBadRequest(endPointName, method, requestBody, errorMessage string, responseCode int) {
	db := config.CreateConnection()

	defer db.Close()

	// Created Date for user

	sqlStatement := `INSERT INTO logging_bad_request (endpoint_name, method, request_body, timestamp, error_message) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement, endPointName, method, requestBody, Gv.FormattedTimeNowYYYYMMDDHHMM, errorMessage)

	if err != nil {
		log.Fatalf("\nINSERT LOGGING  - Cannot execute command : %v\n", err)
	}
}

func LogsNonOK(endPointName, method, requestBody, errorMessage, userId string, responseCode int) {
	db := config.CreateConnection()

	defer db.Close()

	// Created Date for user

	sqlStatement := `INSERT INTO logging_bad_request (user_id ,endpoint_name, method, response_code, timestamp, error_message) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sqlStatement, userId, endPointName, method, requestBody, Gv.FormattedTimeNowYYYYMMDDHHMM, errorMessage)

	if err != nil {
		log.Fatalf("\nINSERT LOGGING  - Cannot execute command : %v\n", err)
	}
}

func Logs200OK(endPointName, method, userId string, responseCode int) {
	db := config.CreateConnection()

	defer db.Close()

	// Created Date for user

	sqlStatement := `INSERT INTO logging_bad_request (user_id ,endpoint_name, method, response_code, timestamp, error_message) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sqlStatement, userId, endPointName, method, responseCode, Gv.FormattedTimeNowYYYYMMDDHHMM)

	if err != nil {
		log.Fatalf("\nINSERT LOGGING  - Cannot execute command : %v\n", err)
	}
}

func InsertLog(r *http.Request, endpointName, errorMessage, token string, methodOption, logsOption, responseCode int) {
	/*
		Method option value
		1 = GET
		2 = POST
		3 = PUT
		4 = DELETE

		Logs option value
		1 = LogsBadRequest only query params
		2 = LogsBadRequest
		3 = LogsNonOK
		4 = Logs200OK
	*/

	var method string

	switch methodOption {
	case 1:
		method = "GET"
	case 2:
		method = "POST"
	case 3:
		method = "PUT"
	case 4:
		method = "DELETE"
	}

	switch logsOption {
	case 1:
		queryParamsString := r.URL.RawQuery
		LogsBadRequest(endpointName, method, queryParamsString, errorMessage, responseCode)
	case 2:
		bodyString, _ := ioutil.ReadAll(r.Body)
		LogsBadRequest(endpointName, method, string(bodyString), errorMessage, responseCode)
	case 3:
		bodyString, _ := ioutil.ReadAll(r.Body)

		userId, _ := validation.ValidateTokenGetUuid(token)

		LogsNonOK(endpointName, method, string(bodyString), errorMessage, userId, responseCode)
	case 4:
		userId, _ := validation.ValidateTokenGetUuid(token)

		Logs200OK(endpointName, method, userId, responseCode)
	}
}
