package db

type Conversation struct {
	ConversationID string `json:"conversation_id"`
	Title          string `json:"title"`
}

const createConversation = `
INSERT INTO conversations
	(conversation_id, title)
VALUES
	(?, ?)
`

func (db *DB) CreateConversation(title string) (string, error) {
	stmt, err := db.Prepare(createConversation)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	id := UUID()
	if _, err = stmt.Exec(id, title); err != nil {
		return "", err
	}

	return id, nil
}

const findConversation = `
SELECT
	conversation_id,
	title
FROM conversations
WHERE
	conversation_id = ?
LIMIT 1
`

func (db *DB) FindConversation(id string) (*Conversation, error) {
	stmt, err := db.Prepare(findConversation)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	var conversation Conversation
	err = rows.Scan(
		&conversation.ConversationID,
		&conversation.Title,
	)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

const selectConversations = `
SELECT
	conversation_id,
	title
FROM conversations
`

func (db *DB) SelectConversations() ([]*Conversation, error) {
	stmt, err := db.Prepare(selectConversations)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	conversations := make([]*Conversation, 0)
	for rows.Next() {
		var conversation Conversation
		err = rows.Scan(
			&conversation.ConversationID,
			&conversation.Title,
		)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, &conversation)
	}

	return conversations, nil
}
