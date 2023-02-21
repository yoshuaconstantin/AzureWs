package models

type FeedbackUserModel struct {
	UserId    string `json:"user_id"`
	Comment   string `json:"comment,omitempty"`
	Timestamp string `json:"timestamp"`
}

type ReturnFeedBackUserModel struct {
	Id        int    `json:"id"`
	Nickname  string `json:"nickname,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	IsEdited  string `json:"is_edited"`
}