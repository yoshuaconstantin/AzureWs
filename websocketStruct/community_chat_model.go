package websocketstruct

type ChatMessageModel struct{
	UserId string		`json:"user_id"`
	Nickname string		`json:"nickname,omitempty"`
	Message string		`json:"message"`
	Nation string 		`json:"nation,omitempty"`
}

type ResponseChatMessage struct{
	Nickname string		`json:"nickname,omitempty"`
	Message string		`json:"message"`
	Timestamp string	`json:"timestamp"`
	Nation string 		`json:"nation,omitempty"`
}