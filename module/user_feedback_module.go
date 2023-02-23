package module

import (
	"AzureWS/config"
	"AzureWS/schemas/models"
	Gv "AzureWS/globalvariable"

)

// Insert Feedback to database using non-model data
func InsertFeedbackUserToDB(UserId, Nickname, Comment string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO user_feedback (user_id, nickname, comment, timestamp, is_edited) VALUES ($1, $2, $3, $4, 'false')`

	_, err := db.Exec(sqlStatement, UserId, Nickname, Comment, Gv.FormatedTime)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Fetch All feedback users with pagination offset and 20 for the limit
func GetFeedBackUserDataFromDB(offset int) ([]models.ReturnFeedBackUserModel, error) {
	db := config.CreateConnection()
	defer db.Close()

	query := "SELECT id, nickname, comment, timestamp, is_edited FROM user_feedback LIMIT 20 OFFSET $1"
	rows, err := db.Query(query, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := []models.ReturnFeedBackUserModel{}
	for rows.Next() {
		var feedback models.ReturnFeedBackUserModel
		if err := rows.Scan(&feedback.Id, &feedback.Nickname, &feedback.Comment, &feedback.Timestamp, &feedback.IsEdited); err != nil {
			return nil, err
		}

		feedbacks = append(feedbacks, feedback)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedbacks, nil
}

// Edit Feedback if match with userId
func EditFeedBackUserFromDB(id int, comment, userId string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE user_feedback SET comment = $1, timestamp = $2 is_edit = 'true' WHERE id = $3 AND user_id = $4`

	_, err := db.Exec(sqlStatement, id, comment, Gv.FormatedTime, userId)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete users based UserId feedback data
func DeleteFeedBackUserFromDB(id int, userId string) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM user_feedback WHERE id = $1 AND user_id = $2`

	_, err := db.Exec(sqlStatement, id, userId)

	if err != nil {
		return false, err
	}

	return true, nil
}
