package module

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // postgres golang driver

	"AzureWS/config"
	Gv "AzureWS/globalvariable/variable"
	"AzureWS/schemas/models"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
)

// Community Post Area
func GetAllCommunityPostFromDB(indexPost int, UserId string) ([]response.PostData, error) {
	db := config.CreateConnection()

	defer db.Close()

	var postDatas []response.PostData

	sqlStatement := `SELECT id,nickname,post_message,nation,image_url,created_date,is_edited FROM community_post LIMIT 10 OFFSET $1`

	var OffsetPost = indexPost * 10

	rows, err := db.Query(sqlStatement, OffsetPost)

	if err != nil {
		log.Fatalf("Cannot exec the query : %v", err)
	}

	for rows.Next() {
		var postDat response.PostData

		err = rows.Scan(&postDat.Id, &postDat.Nickname, &postDat.PostMessage, &postDat.Nation, &postDat.ImageUrl, &postDat.CreatedDate, &postDat.IsEdited)

		if err != nil {
			log.Fatalf("Cannot get all the post data. %v", err)
		}

		sqlStatement := `SELECT user_id FROM community_post WHERE user_id=$1 AND id = $2`

		row := db.QueryRow(sqlStatement, UserId, postDat.Id)

		var userID string

		if err := row.Scan(&userID); err != nil {
			if err == sql.ErrNoRows {
				// User ID not found
				postDat.OwnPost = "false"
			} else {
				return nil, err
			}
		} else {
			// User ID found
			postDat.OwnPost = "true"
		}

		var countComment int
		var countLike int
		err = db.QueryRow("SELECT COUNT(*) FROM community_post_comment WHERE post_id = $1", &postDat.Id).Scan(&countComment)

		if err != nil {
			// handle error
			return nil, err
		}
		err = db.QueryRow("SELECT COUNT(*) FROM community_post_like WHERE post_id = $1", &postDat.Id).Scan(&countLike)

		if err != nil {
			// handle error
			return nil, err
		}

		postDat.LikeCount = countLike
		postDat.CommentCount = countComment

		GetComments, errGetCmnt := GetAllPreviewCommentCommunityPostFromDB(postDat.Id)

		if errGetCmnt != nil {
			return nil, errGetCmnt
		}

		postDat.Comment = GetComments

		postDatas = append(postDatas, postDat)
	}

	defer rows.Close()

	return postDatas, err
}

func InsertCommunityPostToDB(userId string, postDataModels models.PostDataModels) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO community_post (user_id, nickname, post_message, nation, image_url, created_date, is_edited) VALUES ($1, $2, $3, $4, $5, $6, 'false')`

	_, err := db.Exec(sqlStatement, userId, postDataModels.Nickname, postDataModels.PostMessage, postDataModels.Nation, postDataModels.ImageUrl, Gv.FormatedTimeiso8601)

	if err != nil {

		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	fmt.Printf("Insert data single record into Dashboards data\n")

	return true, nil
}

func UpdateCommunityPostFromDB(userId string, updatePostData request.UpdatePostData) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE community_post SET post_message =$1, created_date =$2, is_edited = 'true' WHERE user_id = $3 AND id = $4`

	_, err := db.Exec(sqlStatement, updatePostData.PostMessage, Gv.FormatedTimeiso8601, userId, updatePostData.PostId)

	if err != nil {

		return false, fmt.Errorf("%s %v", "Failed :", err)
	}

	fmt.Printf("Insert data single record into Dashboards data\n")

	return true, nil
}

func DeleteCommunityPostFromDB(userId string, deletePostData request.DeletePostData) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM community_post WHERE user_id = $1 AND id = $2`

	_, err := db.Exec(sqlStatement, userId, deletePostData.PostId)

	sqlStatementDeleteData := `DELETE FROM community_post_comment WHERE post_id = $1`

	_, errDelData := db.Exec(sqlStatementDeleteData, deletePostData.PostId)

	if errDelData != nil {

		return false, fmt.Errorf("%s %v", "Failed :", errDelData)
	}

	if err != nil {

		return false, fmt.Errorf("%s %v", "Failed :", err)
	}

	return true, nil
}

//End Community Post Area

// Comment Community Post Area
func GetAllPreviewCommentCommunityPostFromDB(PostId int) ([]models.CommentDataModels, error) {
	db := config.CreateConnection()

	defer db.Close()

	var commentsData []models.CommentDataModels

	sqlStatement := `SELECT id,post_id,nickname,comment_body,time_comment,is_edited FROM community_post_comment WHERE post_id = $1 LIMIT 3 OFFSET 0`

	rows, err := db.Query(sqlStatement, PostId)

	if err != nil {
		log.Fatalf("Cannot exec the query : %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var data models.CommentDataModels

		err = rows.Scan(&data.Id, &data.PostId, &data.Nickname, &data.CommentBody, &data.TimeComment, &data.IsEdited)

		if err != nil {
			log.Fatalf("Cannot get all the comment data. %v", err)
		}

		commentsData = append(commentsData, data)
	}

	return commentsData, err
}

func GetSpecificCommentCommunityPostFromDB(PostId int, indexComment int, UserId string) ([]models.CommentDataModels, error) {
	db := config.CreateConnection()

	defer db.Close()

	var commentsData []models.CommentDataModels

	sqlStatement := `SELECT id,post_id,nickname,comment_body,time_comment,is_edited FROM community_post_comment WHERE post_id = $1 LIMIT 10 OFFSET $2`

	var OffsetComment = indexComment * 10

	rows, err := db.Query(sqlStatement, PostId, OffsetComment)

	if err != nil {
		log.Fatalf("Cannot exec the query : %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var data models.CommentDataModels

		err = rows.Scan(&data.Id, &data.PostId, &data.Nickname, &data.CommentBody, &data.TimeComment, &data.IsEdited)

		if err != nil {
			log.Fatalf("Cannot get all the comment data. %v", err)
		}

		sqlStatement := `SELECT user_id FROM community_post_comment WHERE user_id=$1 AND post_id = $2`

		row := db.QueryRow(sqlStatement, UserId, data.PostId)

		var userID string

		if err := row.Scan(&userID); err != nil {
			if err == sql.ErrNoRows {
				// User ID not found
				data.OwnComment = "false"
			} else {
				return nil, err
			}
		} else {
			// User ID found
			data.OwnComment = "true"
		}

		commentsData = append(commentsData, data)
	}

	return commentsData, err
}

func InsertCommentCommunityPostToDB(userId string, commentData request.CommentPost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO community_post_comment (user_id, post_id, nickname, comment_body, time_comment, is_edited) VALUES ($1, $2, $3, $4, $5, 'false')`

	_, err := db.Exec(sqlStatement, userId, commentData.PostId, commentData.Nickname, commentData.CommentBody, Gv.FormatedTimeiso8601)

	if err != nil {

		return false, fmt.Errorf("%s %v", "Failed :", err)
	}

	return true, nil
}

func UpdateCommentCommunityPostFromDB(userId string, updateCommentCommunityPost request.UpdateCommentPost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE community_post_comment SET comment_body =$1, time_comment =$2, is_edited = 'true' WHERE user_id = $3 AND post_id = $4 AND id = $5`

	_, err := db.Exec(sqlStatement, updateCommentCommunityPost.CommentBody, Gv.FormatedTimeiso8601, userId, updateCommentCommunityPost.PostId, updateCommentCommunityPost.CommentId)

	if err != nil {

		return false, fmt.Errorf("%s %v", "Failed :", err)
	}

	return true, nil
}

func DeleteCommentCommunityPostFromDB(userId string, deleteCommentData request.DeleteCommentPost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM community_post_comment WHERE user_id = $1 AND post_id = $2 AND id = $3`

	_, err := db.Exec(sqlStatement, userId, deleteCommentData.PostId, deleteCommentData.CommentId)

	if err != nil {

		return false, fmt.Errorf("%s %v", "Failed :", err)
	}

	return true, nil
}

//End Comment Community Post Area

// Like Community Post Area
func InsertLikeCommunityPostToDB(userId string, likeData request.LikePost) (bool, error) {
	db := config.CreateConnection()
	defer db.Close()

	// Check if like already exists for user and post
	var count int
	sqlStatement := `SELECT COUNT(*) FROM community_post_like WHERE user_id = $1 AND post_id = $2`
	err := db.QueryRow(sqlStatement, userId, likeData.PostId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if like already exists: %v", err)
	}

	if count > 0 {
		// Like already exists, delete it
		sqlStatement := `DELETE FROM community_post_like WHERE user_id = $1 AND post_id = $2`
		_, err := db.Exec(sqlStatement, userId, likeData.PostId)
		if err != nil {
			return false, fmt.Errorf("failed to delete existing like data: %v", err)
		}
	} else {
		// Like doesn't exist, insert new like data
		sqlStatement := `INSERT INTO community_post_like (user_id, post_id, is_like) VALUES ($1, $2, $3)`
		_, err := db.Exec(sqlStatement, userId, likeData.PostId, likeData.IsLike)
		if err != nil {
			return false, fmt.Errorf("failed to insert new like data: %v", err)
		}
	}

	return true, nil
}

//End Like Community Post Area
