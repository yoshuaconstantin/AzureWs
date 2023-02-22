package module

import (
	"AzureWS/config"
	"AzureWS/schemas/request"
	"AzureWS/schemas/response"
	"log"

	"fmt"

	_ "github.com/lib/pq" // postgres golang driver
)

//Community Post Area
func GetAllCommunityPostFromDB(indexPost, indexComment int) ([]response.PostData, error){
	db := config.CreateConnection()

	defer db.Close()

	var postDatas []response.PostData
	
	sqlStatement := `SELECT * FROM community_post LIMIT 20 OFFSET $1`

	var OffsetPost = indexPost * 10

	rows, err := db.Query(sqlStatement, OffsetPost)

	if err != nil {
		log.Fatalf("Cannot exec the query : %v", err)
	}
	
	defer rows.Close()

	for rows.Next() {
		var postDat response.PostData

		err = rows.Scan( &postDat.PostId, &postDat.Nickname, &postDat.PostMessage, &postDat.Nation, &postDat.ImageUrl, &postDat.LikeCount, &postDat.CommentCount)

		if err != nil {
			log.Fatalf("Cannot get all the post data. %v", err)
		}

		GetComments, errGetCmnt := GetAllCommentCommunityPostFromDB(postDat.PostId)

		if errGetCmnt != nil {
			return nil, errGetCmnt
		}

		postDat.Comment = GetComments

		postDatas = append(postDatas, postDat)
	}
	
	return postDatas, err
}

func InsertCommunityPostToDB(userId string, postData request.PostData) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO community_post (user_id, Nickname, PostMessage, Nation, ImageUrl) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement, userId, postData.Nickname, postData.PostMessage, postData.Nation, postData.ImageUrl)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	fmt.Printf("Insert data single record into Dashboards data\n")

	return true, nil
}

func UpdateCommunityPostFromDB(userId string, updatePostData request.UpdatePostData) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE community_post SET post_message =$1, time_comment =$2, is_edited = 'true' WHERE user_id = $3 AND post_id = $4`

	_, err := db.Exec(sqlStatement, userId, updatePostData.PostMessage, updatePostData.TimePost, userId, updatePostData.PostId)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	fmt.Printf("Insert data single record into Dashboards data\n")

	return true, nil
}

func DeleteCommunityPostFromDB(userId string, deletePostData request.DeletePostData) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM community_post WHERE user_id = $1 AND post_id = $2`

	_, err := db.Exec(sqlStatement, userId, deletePostData.PostId)

	sqlStatementDeleteData := `DELETE FROM community_post_comment WHERE post_id = $1`

	_, errDelData := db.Exec(sqlStatementDeleteData, deletePostData.PostId)

	if errDelData != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	return true, nil
}
//End Community Post Area

//Comment Community Post Area
func GetAllCommentCommunityPostFromDB(PostId string) ([]response.Comment, error){
	db := config.CreateConnection()

	defer db.Close()

	var commentsData []response.Comment

	sqlStatement := `SELECT * FROM community_post_comment WHERE post_id = $1 LIMIT 20 OFFSET 0`


	rows, err := db.Query(sqlStatement, PostId)

	if err != nil {
		log.Fatalf("Cannot exec the query : %v", err)
	}
	
	defer rows.Close()

	for rows.Next() {
		var data response.Comment

		err = rows.Scan( &data.CommentId, &data.Nickname, &data.Message, &data.TimeComment)

		if err != nil {
			log.Fatalf("Cannot get all the comment data. %v", err)
		}

		commentsData = append(commentsData, data)
	}

	return commentsData, err
}

func GetSpecificCommentCommunityPostFromDB(PostId string, indexComment int) ([]response.Comment, error){
	db := config.CreateConnection()

	defer db.Close()

	var commentsData []response.Comment

	sqlStatement := `SELECT * FROM community_post_comment WHERE post_id = $1 LIMIT 20 OFFSET $2`

	var OffsetComment = indexComment * 10

	rows, err := db.Query(sqlStatement, PostId, OffsetComment)

	if err != nil {
		log.Fatalf("Cannot exec the query : %v", err)
	}
	
	defer rows.Close()

	for rows.Next() {
		var data response.Comment

		err = rows.Scan( &data.CommentId, &data.Nickname, &data.Message, &data.TimeComment)

		if err != nil {
			log.Fatalf("Cannot get all the comment data. %v", err)
		}

		commentsData = append(commentsData, data)
	}

	return commentsData, err
}

func InserCommentCommunityPostToDB(userId string, commentData request.CommentPost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO community_post_comment (user_id, post_id, nickname, comment_body, time_comment) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement, userId, commentData.PostId, commentData.Nickname, commentData.CommentBody, commentData.TimeComment)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	return true, nil
}

func UpdateCommentCommunityPostFromDB(userId string, updateCommentCommunityPost request.UpdateCommentPost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `UPDATE community_post_comment SET comment_body =$1, time_comment =$2, is_edited = 'true' WHERE user_id = $3 AND post_id = $4 AND comment_id = $5`

	_, err := db.Exec(sqlStatement, updateCommentCommunityPost.CommentBody, updateCommentCommunityPost.TimeComment, userId, updateCommentCommunityPost.PostId, updateCommentCommunityPost.CommentId)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	return true, nil
}

func DeleteCommentCommunityPostFromDB(userId string, deleteCommentData request.DeleteCommentPost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM community_post_comment WHERE user_id = $1 AND post_id = $2 AND comment_id = $3`

	_, err := db.Exec(sqlStatement, userId, deleteCommentData.PostId, deleteCommentData.CommentId, userId)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	return true, nil
}
//End Comment Community Post Area

//Like Community Post Area
func InserLikeCommunityPostToDB(userId string, likeData request.LikePost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO community_post_like (user_id, post_id, is_like) VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, userId, likeData.PostId, likeData.IsLike)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	return true, nil
}

func DeleteLikeCommunityPostToDB(userId string, likeData request.LikePost) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM community_post_like WHERE user_id = $1 AND post_id = $2 AND comment_id = $3`

	_, err := db.Exec(sqlStatement, userId, likeData.PostId, likeData.IsLike)


	if err != nil {
		
		return false, fmt.Errorf("%s", "Failed, try again later")
	}

	return true, nil
}
//End Like Community Post Area