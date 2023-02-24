package websocketstruct

import (
	"AzureWS/config"
	Gv "AzureWS/globalvariable/variable"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Save users self chat to DB and become history
func SaveChatToDB(chatModel ChatMessageModel) (bool, error) {
	db := config.CreateConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO community_chat (user_id, nickname, message, timestamp, nation) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement, chatModel.UserId, chatModel.Nickname, chatModel.Message, Gv.FormatedTimeiso8601, chatModel.Nation)

	if err != nil {
		return false, fmt.Errorf("%s %v", "INSERT CHAT TO DB - Cannot execute query :", err)
	}

	return true, nil
}

// Restore saved chat from database as chat history
func RestoreHistoryChatFromDB(conn *websocket.Conn) error {
	// Connect to the database
	db := config.CreateConnection()

	defer db.Close()

	// Query the chat history
	rows, err := db.Query("SELECT nickname, message, timestamp, nation FROM community_chat ORDER BY timestamp")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Send each chat message to the WebSocket connection
	for rows.Next() {
		var nickname, message, timestamp, nation string

		err := rows.Scan(&nickname, &message, &timestamp, &nation)
		if err != nil {
			return err
		}
		chatMessage := ResponseChatMessage{nickname, message, timestamp, nation}
		err = conn.WriteJSON(chatMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

// Deleting old chat message if it more than 30 days, and it runs once everyday
func DeleteOldMessage() {
	ticker := time.NewTicker(24 * time.Hour) // ticker that runs once a day
	defer ticker.Stop()
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			db := config.CreateConnection()

			_, err := db.Exec("DELETE FROM community_chat WHERE timestamp < $1", Gv.FormatedTime30DayAgoIso8601)

			if err != nil {
				log.Println("Error deleting old chat messages: ", err)
			} else {
				log.Println("Successfully deleted old chat messages")
			}
		case <-quit:
			// stop the infinite loop when the quit channel receives a signal
			return

		}
	}
}
