package response

import "AzureWS/schemas/models"

type ResponseGetCommunityPost struct {
	Message string	`json:"message,omitempty"`
	Status  int		`json:"status,omitempty"`
	Data []PostData `json:"data"`
}

type PostData struct{
	Id			int `json:"id"`
	OwnPost			string `json:"own_post"`
	Nickname 		string `json:"nickname"`
	PostMessage		string `json:"post_message"`
	Nation	 		string `json:"nation,omitempty"`
	ImageUrl 		string `json:"image_url,omitempty"`
	CreatedDate		string `json:"created_date"`
	IsEdited		string `json:"is_edited"`
	LikeCount		int		`json:"like_count"`
	CommentCount	int		`json:"comment_count"`
	Comment 		[]models.CommentDataModels `json:"comments,omitempty"`
}

type ResponseGetCommunityPostComment struct {
	Message string	`json:"message,omitempty"`
	Status  int		`json:"status,omitempty"`
	Data []models.CommentDataModels `json:"data"`
}

