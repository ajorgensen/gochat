package db

const (
	User      Role = "User"
	System    Role = "System"
	Assistant Role = "Assistant"
)

type Role string

type Message struct {
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	Message        string `json:"message"`
	Role           Role   `json:"role"`
}

const createMessage = `
INSERT INTO messages 
	(conversation_id, message_id, message, role) 
VALUES 
	(?, ?, ?, ?)
`

func (db *DB) CreateMessage(conversationID string, role Role, message string) error {
	messageID := UUID()

	stmt, err := db.Prepare(createMessage)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(conversationID, messageID, message, role)
	if err != nil {
		return err
	}

	return nil
}

const getMessages = `
SELECT 
	conversation_id, 
	message_id, 
	message,
	role
FROM messages 
WHERE 
	conversation_id = ? order by created_at asc
`

func (db *DB) GetMessages(conversationID string) ([]*Message, error) {
	stmt, err := db.Prepare(getMessages)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*Message, 0)
	for rows.Next() {
		var message Message

		err := rows.Scan(&message.ConversationID, &message.MessageID, &message.Message, &message.Role)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
