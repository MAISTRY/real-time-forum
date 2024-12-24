package middleware

import (
	"RTF/DB"
	"database/sql"
	"encoding/json"
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

	// todo: add sender to the massage
	lastMessageQuery = `
		SELECT 
			m.sender_id,
			m.receiver_id,
			m.message,
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
			m.timestamp DESC
		LIMIT 1;
	`

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	syncronize sync.Mutex
	intUserID  int
	clients    = make(map[int]*websocket.Conn)
)

type Msg struct {
	Message    string `json:"message"`
	ReceiverId int    `json:"receiver_id"`
}

type UserStatus struct {
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	LastMessage string    `json:"last_message"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}

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

	intUserID, err = strconv.Atoi(userID)
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

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			return
		}
		fmt.Printf("Received from %s: %s\n", userID, data)

		err = MessageHandler(userID, msgType, data, db)
		if err != nil {
			log.Printf("Error handling message: %v\n", err)
			return
		}
	}

	syncronize.Lock()
	delete(clients, intUserID)
	syncronize.Unlock()
	fmt.Println("Client disconnected", userID)

}

func UsersHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	// Get all users
	var allusers []UserStatus
	userRow, err := db.Query(UsersQuery)
	if err != nil {
		http.Error(w, "Error querying UsersQuery", http.StatusInternalServerError)
		return
	}
	defer userRow.Close()

	for userRow.Next() {
		var userStatus UserStatus
		err = userRow.Scan(&userStatus.UserID, &userStatus.Username)
		if err != nil {
			log.Printf("Error scanning users %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
		if userStatus.UserID != intUserID {
			allusers = append(allusers, userStatus)
		}
	}

	// Get user last message for all users
	for i, userStatus := range allusers {
		rows, err := db.Query(lastMessageQuery, intUserID, userStatus.UserID, userStatus.UserID, intUserID)
		if err != nil {
			log.Printf("Error getting users %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusOK)
			return
		}
		defer rows.Close()

		if rows.Next() {
			var senderID, receiverID int
			var message string
			var timestamp time.Time
			err = rows.Scan(&senderID, &receiverID, &message, &timestamp)
			if err != nil {
				log.Printf("Error scanning users %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusOK)
				return
			}
			userStatus.LastMessage = message
			userStatus.Timestamp = timestamp
		} else {
			userStatus.LastMessage = fmt.Sprintf("Say hi to %s👋", userStatus.Username)
			userStatus.Timestamp = time.Time{} // Zero value for time
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allusers)

}

func MessageHandler(userID string, msgType int, message []byte, db *sql.DB) error {

	// MessageID, err := DB.InsertMessage(senderId, receiverId, string(message))
	return nil
}
