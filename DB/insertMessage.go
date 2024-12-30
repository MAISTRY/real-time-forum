package DB

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	InsertMessageQuery = `
        INSERT INTO "Messages" (sender_id, receiver_id, message)
        VALUES (?,?,?)
    `
)

func InsertMessage(senderId int, receiverId int, message string) (int64, error) {

	db, err := sql.Open("sqlite3", "meow.db")
	if err != nil {
		log.Printf("Error getting user %v\n", err)
		return -1, fmt.Errorf("error Opening database: %v", err)
	}
	defer db.Close()

	result, err := db.Exec(InsertMessageQuery, senderId, receiverId, message)
	if err != nil {
		return -1, fmt.Errorf("error insert in the database: %v", err)
	}

	MessageID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return MessageID, nil
}
