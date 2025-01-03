package middleware

import (
	"RTF/DB"
	"RTF/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// *need to be in the websocket
// netifications
// chat
// status

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	syncronize sync.Mutex
	clients    = make(map[int]*websocket.Conn)
)

type (
	WebSocketMessage struct {
		Type       string    `json:"type"`
		Receiver   string    `json:"receiver,omitempty"`
		Message    string    `json:"message,omitempty"`
		FirstUser  int       `json:"firstUser,omitempty"`
		SecondUser int       `json:"secondUser,omitempty"`
		Timestamp  time.Time `json:"timestamp,omitempty"`
	}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	userID, err := DB.GetUserIDByCookie(r, db)
	if err != nil {
		log.Printf("The Error getting user ID %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting user ID to int %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	defer conn.Close()

	syncronize.Lock()
	clients[intUserID] = conn
	syncronize.Unlock()
	fmt.Printf("Client ID:%s connected\n", userID)
	handlers.BroadcastUserList(db, &syncronize, clients)

	for {
		var WebMsg WebSocketMessage
		if err := conn.ReadJSON(&WebMsg); err != nil {
			log.Printf("Read Error: %v\n", err)
			break
		}

		switch WebMsg.Type {
		case "GetMessages":
			messages, err := handlers.ShowAllMessages(db, intUserID, WebMsg.SecondUser)
			if err != nil {
				log.Printf("Error getting users %v\n", err)
				continue
			}

			if err := conn.WriteJSON(map[string]interface{}{
				"type":       "getMessages",
				"Sender":     intUserID,
				"Receiver":   WebMsg.Receiver,
				"ReceiverID": WebMsg.SecondUser,
				"messages":   messages,
			}); err != nil {
				log.Printf("Error Sending message %v\n", err)
			}

		case "SendMessage":
			found := false
			var SecondUserConn *websocket.Conn
			for user, conn := range clients {
				if user == WebMsg.SecondUser {
					found = true
					SecondUserConn = conn
				}
			}
			if found {
				message, err := handlers.MessageHandler(db, intUserID, WebMsg.SecondUser, WebMsg.Message)
				if err != nil {
					log.Printf("Error getting Message from User%v\n", err)
					continue
				}

				if err := conn.WriteJSON(map[string]interface{}{
					"type":       "SendMessage",
					"Sender":     intUserID,
					"Receiver":   WebMsg.Receiver,
					"ReceiverID": WebMsg.SecondUser,
					"messages":   message,
				}); err != nil {
					log.Printf("Error Sending message %v\n", err)
				}

				if err := SecondUserConn.WriteJSON(map[string]interface{}{
					"type":       "SendMessage",
					"Sender":     WebMsg.SecondUser,
					"Receiver":   WebMsg.Receiver,
					"ReceiverID": intUserID,
					"messages":   message,
				}); err != nil {
					log.Printf("Error Sending message %v\n", err)
				}
			} else {
				if err := conn.WriteJSON(map[string]interface{}{
					"type":       "Offline",
					"Sender":     intUserID,
					"Receiver":   WebMsg.Receiver,
					"ReceiverID": WebMsg.SecondUser,
				}); err != nil {
					log.Printf("Error Sending message %v\n", err)
				}
			}
		case "loadUsers":
			handlers.BroadcastUserList(db, &syncronize, clients)
		}
	}

	syncronize.Lock()
	delete(clients, intUserID)
	syncronize.Unlock()
	fmt.Println("Client disconnected", userID)
	handlers.BroadcastUserList(db, &syncronize, clients)

}
