package request

type RequestInsertFeedback struct {
	Token     string `json:"token"`
	Comment   string `json:"comment"`
	Timestamp string `json:"timestamp"`
}

type RequestEditFeedback struct {
	Id        int    `json:"id"`
	Token     string `json:"token"`
	Comment   string `json:"comment"`
	Timestamp string `json:"timestamp"`
}

type RequestGetAllFeedbackData struct {
	Token string `json:"token"`
	Index int    `json:"index"`
}

type RequestDeleteSingleFeedbackData struct {
	Token string `json:"token"`
	Id    int    `json:"id"`
}