package request

type RequestInsertCommunityPost struct {
	Token string `json:"token"`
	Data  []PostData `json:"data"`
}

type PostData struct {
	Nickname 		string `json:"nickname"`
	PostMessage		string `json:"PostMessage"`
	Nation	 		string `json:"nation,omitempty"`
	ImageUrl 		string `json:"image_url,omitempty"`
}

type RequestUpdateCommunityPost struct {
	Token string `json:"token"`
	Data  []PostData `json:"data"`
}

type UpdatePostData struct {
	PostId 		string `json:"post_id"`
	PostMessage		string `json:"post_message"`
	TimePost		string `json:"time_post"`
}

type RequestDeleteCommunityPost struct {
	Token string `json:"token"`
	Data  []DeletePostData `json:"data"`
}

type DeletePostData struct {
	PostId 		string `json:"post_id"`
}

type RequestInsertLikeCommunityPost struct {
	Token string `json:"token"`
	Data  []LikePost `json:"data"`
}

type LikePost struct {
	PostId 		string `json:"post_id"`
	Nickname		string `json:"nickname"`
	IsLike	 		string `json:"is_like"`
}

type RequestInsertCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []CommentPost `json:"data"`
}

type CommentPost struct {
	PostId			string `json:"post_id"`
	Nickname		string `json:"nickname"`
	CommentBody		string `json:"comment_body"`
	TimeComment		string `json:"time_comment"`
}

type RequestUpdateCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []UpdateCommentPost `json:"data"`
}

type UpdateCommentPost struct {
	PostId			string `json:"post_id"`
	CommentId		string `json:"comment_id"`
	CommentBody		string `json:"comment_body"`
	TimeComment		string `json:"time_comment"`
}

type RequestDeleteCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []UpdateCommentPost `json:"data"`
}

type DeleteCommentPost struct {
	PostId			string `json:"post_id"`
	CommentId		string `json:"comment_id"`
}