package handlers

import (
	"database/sql"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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
)

func BroadcastUserList(db *sql.DB, syncronize *sync.Mutex, clients map[int]*websocket.Conn) {
	syncronize.Lock()
	defer syncronize.Unlock()

	for user, conn := range clients {
		users, err := GetUsers(db, user, clients)
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
func GetUsers(db *sql.DB, userID int, clients map[int]*websocket.Conn) ([]UserStatus, error) {

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
