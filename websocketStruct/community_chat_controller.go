package websocketstruct

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Function to upgrade the connection into websocket
func CommunityChat(w http.ResponseWriter, r *http.Request) {
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Restore All previous chat that has been made to user after connecting
	err = RestoreHistoryChatFromDB(websocket)
    if err != nil {
        log.Println(err)
    }

	log.Println("Community chat Websocket Connected!")
	listen(websocket)
}

// Listening incoming users self message to store into DB
func listen(conn *websocket.Conn) {
	for {
		// read a message
		messageType, messageContent, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		timeReceive := time.Now()
		if err != nil {
			log.Println(err)
			return
		}
	
		// parse the message content as a ChatMessage
		var chatMsg ChatMessageModel
		err = json.Unmarshal(messageContent, &chatMsg)
		if err != nil {
			log.Printf("Error parsing message as JSON: %v\n", err)
			continue
		}
	
		// print out the message
		log.Printf("[%s] %s\n", chatMsg.Nickname, chatMsg.Message)
	
		// save the message to the database
		_, errSaveChat := SaveChatToDB(chatMsg)

		if errSaveChat != nil {
			log.Printf("Error : %v\n", errSaveChat)
			continue
		}
	
		// response message
		messageResponse := fmt.Sprintf("Your message is: %s. Time received : %v", messageContent, timeReceive)
	
		if err := conn.WriteMessage(messageType, []byte(messageResponse)); err != nil {
			log.Println(err)
			return
		}
	
	}
}









