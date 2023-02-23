package models

type PostDataModels struct {
	Nickname 		string `json:"nickname"`
	PostMessage		string `json:"post_message"`
	Nation	 		string `json:"nation,omitempty"`
	ImageUrl 		string `json:"image_url,omitempty"`
}

type CommentDataModels struct {
	Id 		string `json:"id"`
	PostId 		int `json:"post_id"`
	Nickname 		string `json:"nickname"`
	CommentBody		string `json:"comment_body"`
	TimeComment	 		string `json:"time_comment,omitempty"`
	IsEdited 		string `json:"is_edited,omitempty"`
	OwnComment		string `json:"own_comment"`
}