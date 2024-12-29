package middleware

import (
	"RTF/DB"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sort"
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
	UsersQuery = `
        SELECT 
			UserID, 
			Username
        FROM
            User
		ORDER BY
			Username ASC;			
	`

	lastMessageQuery = `
		SELECT 
			m.sender_id,
			m.receiver_id,
			m.message,
			sender.Username AS sender,
			m.timestamp
		FROM 
			Messages m
		JOIN 
			User sender 
			ON m.sender_id = sender.UserID
		WHERE 
			(m.sender_id = ? AND m.receiver_id = ?)
			OR
			(m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY 
			m.timestamp DESC
		LIMIT 1;
	`

	AllMessagesQuery = `
		SELECT 
			m.sender_id,
			m.receiver_id,
			m.message,
			sender.Username AS sender,
			receiver.Username AS receiver,
			m.timestamp
		FROM 
			Messages m
		JOIN 
			User sender 
			ON m.sender_id = sender.UserID
		JOIN 
			User receiver 
			ON m.receiver_id = receiver.UserID
		WHERE 
			(m.sender_id = ? AND m.receiver_id = ?)
			OR
			(m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY 
			m.timestamp ASC
	`

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	syncronize sync.Mutex
	clients    = make(map[int]*websocket.Conn)
)

type (
	Msg struct {
		Message    string    `json:"message"`
		FirstUser  int       `json:"FirstUser"`
		SecondUser int       `json:"SecondUser"`
		Sender     string    `json:"Sender"`
		Receiver   string    `json:"Receiver"`
		Timestamp  time.Time `json:"timestamp"`
	}

	UserStatus struct {
		UserID      int       `json:"userID"`
		Username    string    `json:"username"`
		LastMessage string    `json:"lastMessage"`
		Sender      string    `json:"sender"`
		Status      string    `json:"status"`
		Timestamp   time.Time `json:"timestamp"`
	}

	WebSocketMessage struct {
		Type       string    `json:"type"`
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
	BroadcastUserList(db)

	for {
		var WebMsg WebSocketMessage
		if err := conn.ReadJSON(&WebMsg); err != nil {
			log.Printf("Read Error: %v\n", err)
			break
		}

		switch WebMsg.Type {
		case "GetMessages":
			messages, err := ShowAllMessages(db, intUserID, WebMsg.SecondUser)
			log.Print(intUserID)
			log.Print(WebMsg.SecondUser)

			if err != nil {
				log.Printf("Error getting users %v\n", err)
				continue
			}

			if err := conn.WriteJSON(map[string]interface{}{
				"type":     "getMessages",
				"Sender":   intUserID,
				"Receiver": WebMsg.SecondUser,
				"messages": messages,
			}); err != nil {
				log.Printf("Error Sending message %v\n", err)
			}
		}

	}

	syncronize.Lock()
	delete(clients, intUserID)
	syncronize.Unlock()
	fmt.Println("Client disconnected", userID)
	BroadcastUserList(db)

}

func BroadcastUserList(db *sql.DB) {
	syncronize.Lock()
	defer syncronize.Unlock()

	for user, conn := range clients {
		users, err := getUsers(db, user)
		if err != nil {
			log.Printf("Error getting users %v\n", err)
			continue
		}

		if err := conn.WriteJSON(map[string]interface{}{
			"type":  "loadUsersResponse",
			"users": users,
		}); err != nil {
			log.Printf("Error broadcasting user list %v\n", err)
			conn.Close()
			delete(clients, user)
		}
	}
}
func getUsers(db *sql.DB, userID int) ([]UserStatus, error) {

	if userID == 0 {
		log.Println("No userID provided")
		return nil, nil
	}

	// Get all users
	var allusers []UserStatus
	userRow, err := db.Query(UsersQuery)
	if err != nil {
		return nil, err
	}
	defer userRow.Close()

	for userRow.Next() {
		var userStatus UserStatus
		err = userRow.Scan(&userStatus.UserID, &userStatus.Username)
		if err != nil {
			log.Printf("Error scanning users %v\n", err)
			return nil, err
		}
		if userStatus.UserID != userID {
			allusers = append(allusers, userStatus)
		}
	}

	// Get user last message for all users
	for i, userStatus := range allusers {
		rows, err := db.Query(lastMessageQuery, userID, userStatus.UserID, userStatus.UserID, userID)
		if err != nil {
			log.Printf("Error getting users %v\n", err)
			return nil, err
		}
		defer rows.Close()

		if rows.Next() {
			var senderID, receiverID int
			var message, sender string
			var timestamp time.Time
			err = rows.Scan(&senderID, &receiverID, &message, &sender, &timestamp)
			if err != nil {
				log.Printf("Error scanning users %v\n", err)
				return nil, err
			}
			userStatus.LastMessage = message
			userStatus.Timestamp = timestamp
			userStatus.Sender = sender
		} else {
			userStatus.LastMessage = "Say hi 👋"
			userStatus.Timestamp = time.Time{}
		}

		userStatus.Status = "offline"
		for client := range clients {
			if userStatus.UserID == client {
				userStatus.Status = "online"
				break
			}
		}

		allusers[i] = userStatus
	}

	sort.Slice(allusers, func(i, j int) bool {
		if allusers[i].Timestamp.IsZero() && allusers[j].Timestamp.IsZero() {
			return allusers[i].Username < allusers[j].Username
		}
		if allusers[i].Timestamp.IsZero() {
			return false
		}
		if allusers[j].Timestamp.IsZero() {
			return true
		}
		return allusers[i].Timestamp.After(allusers[j].Timestamp)
	})

	return allusers, nil
}

func ShowAllMessages(db *sql.DB, SenderID int, ReceiverID int) ([]Msg, error) {

	if SenderID == 0 || ReceiverID == 0 {
		log.Println("Invalid userID provided")
		return nil, nil
	}

	var allMessages []Msg
	rows, err := db.Query(AllMessagesQuery, SenderID, ReceiverID, ReceiverID, SenderID)
	if err != nil {
		log.Printf("Error getting users %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var msg Msg
		err = rows.Scan(&msg.FirstUser, &msg.SecondUser, &msg.Message, &msg.Sender, &msg.Receiver, &msg.Timestamp)
		if err != nil {
			log.Printf("Error scanning users %v\n", err)
			return nil, err
		}
		allMessages = append(allMessages, msg)
	}

	return allMessages, nil

}

func MessageHandler(userID string, msgType string, message []byte, db *sql.DB) error {

	// MessageID, err := DB.InsertMessage(senderId, receiverId, string(message))
	// var msg Msg
	// if err := json.Unmarshal(message, &msg); err != nil {
	// 	return err
	// }

	// switch msgType {
	// case "message":
	// 	syncronize.Lock()
	// 	if conn, ok := clients[msg.Sender]; ok {
	// 		conn.WriteJSON(map[string]interface{}{
	// 			"type":       "messageUpdate",
	// 			"firstUser":  msg.FirstUser,
	// 			"secondUser": msg.SecondUser,
	// 		})
	// 	}

	// }
	return nil
}
