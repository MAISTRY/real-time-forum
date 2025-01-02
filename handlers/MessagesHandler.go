package handlers

import (
	"RTF/DB"
	"RTF/utils"
	"database/sql"
	"errors"
	"log"
)

const (
	OneMessagesQuery = `
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
			m.id = ?
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
)

func MessageHandler(db *sql.DB, senderId int, receiverId int, msg string) (Msg, error) {

	var Message Msg

	if msg == "" {
		log.Println("Message is empty")
		return Message, errors.New("MESSAGES IS EMPTY")
	}

	message := utils.InputSanitizer(msg)

	if senderId == 0 || receiverId == 0 {
		log.Println("Invalid userID provided")
		return Message, errors.New("ID IS EMPTY")
	}

	MessageID, err := DB.InsertMessage(senderId, receiverId, message)
	if err != nil {
		log.Printf("Error inserting message %v\n", err)
		return Message, err
	}

	if err = db.QueryRow(OneMessagesQuery, MessageID).Scan(&Message.FirstUser, &Message.SecondUser, &Message.Message, &Message.Sender, &Message.Receiver, &Message.Timestamp); err != nil {
		return Message, err
	}

	return Message, nil
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
