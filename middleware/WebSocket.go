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
			sender.Username AS sender,
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
	clients    = make(map[int]*websocket.Conn)
)

type Msg struct {
	Message    string `json:"message"`
	ReceiverId int    `json:"receiver_id"`
}

type UserStatus struct {
	UserID      int       `json:"userID"`
	Username    string    `json:"username"`
	LastMessage string    `json:"lastMessage"`
	Sender      string    `json:"sender"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
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

	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %s disconnected: %v\n", userID, err)
			return
		}

		err = MessageHandler(userID, msgType, data, db)
		if err != nil {
			log.Printf("Error handling message: %v\n", err)
			return
		}
	}

	// syncronize.Lock()
	// delete(clients, intUserID)
	// syncronize.Unlock()
	// fmt.Println("Client disconnected", userID)

}

func UsersHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	defer db.Close()

	userID := r.URL.Query().Get("user")
	if userID == "" {
		log.Println("No userID provided")
		http.Error(w, "No userID provided", http.StatusBadRequest)
		return
	}

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error converting user ID to int %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusOK)
		return
	}
	if intUserID == 0 {
		log.Println("Invalid userID provided")
		http.Error(w, "Invalid userID provided", http.StatusBadRequest)
		return
	}

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
			var message, sender string
			var timestamp time.Time
			err = rows.Scan(&senderID, &receiverID, &message, &sender, &timestamp)
			if err != nil {
				log.Printf("Error scanning users %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusOK)
				return
			}
			userStatus.LastMessage = message
			userStatus.Timestamp = timestamp
			userStatus.Sender = sender
		} else {
			userStatus.LastMessage = "Say hi 👋"
			userStatus.Timestamp = time.Time{} // todo: the value need to be removed
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
