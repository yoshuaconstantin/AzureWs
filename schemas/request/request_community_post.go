package request

type RequestInsertCommunityPost struct {
	Token 			string `json:"token"`
	PostMessage		string `json:"post_message"`
}

type RequestUpdateCommunityPost struct {
	Token string `json:"token"`
	Data  []UpdatePostData `json:"data"`
}

type UpdatePostData struct {
	PostId 		int `json:"post_id"`
	PostMessage		string `json:"post_message"`
}

type RequestDeleteCommunityPost struct {
	Token string `json:"token"`
	Data  []DeletePostData `json:"data"`
}

type DeletePostData struct {
	PostId 		int `json:"post_id"`
}

type RequestInsertLikeCommunityPost struct {
	Token string `json:"token"`
	Data  []LikePost `json:"data"`
}

type LikePost struct {
	PostId 		int `json:"post_id"`
	IsLike		string `json:"is_like"`
}

type RequestGetCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []CommentData `json:"data"`
}

type CommentData struct {
	PostId	int `json:"post_id"`
	Index 	int `json:"index"`
}

type RequestInsertCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []CommentPost `json:"data"`
}

type CommentPost struct {
	PostId			int `json:"post_id"`
	Nickname		string `json:"nickname"`
	CommentBody		string `json:"comment_body"`
}

type RequestUpdateCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []UpdateCommentPost `json:"data"`
}

type UpdateCommentPost struct {
	PostId			int `json:"post_id"`
	CommentId		int `json:"comment_id"`
	CommentBody		string `json:"comment_body"`
}

type RequestDeleteCommentCommunityPost struct {
	Token string `json:"token"`
	Data  []DeleteCommentPost `json:"data"`
}

type DeleteCommentPost struct {
	PostId			int `json:"post_id"`
	CommentId		int `json:"comment_id"`
}