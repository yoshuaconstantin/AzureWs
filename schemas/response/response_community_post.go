package response

type ResponseGetCommunityPost struct {
	Message string	`json:"message,omitempty"`
	Status  int		`json:"status,omitempty"`
	Data []PostData `json:"data"`
}

type PostData struct{
	PostId			string `json:"post_id"`
	Nickname 		string `json:"nickname"`
	PostMessage		string `json:"PostMessage"`
	Nation	 		string `json:"nation,omitempty"`
	ImageUrl 		string `json:"image_url,omitempty"`
	LikeCount		int		`json:"like_count,omitempty"`
	CommentCount	int		`json:"comment_count,omitempty"`
	Comment 		[]Comment `json:"comments,omitempty"`
}

type Comment struct {
	CommentId string `json:"comment_id"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
	TimeComment string `json:"time_comment"`
}