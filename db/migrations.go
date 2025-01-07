package db

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS conversations (
		conversation_id TEXT NOT NULL,
		title TEXT NOT NULL,
		PRIMARY KEY (conversation_id)
	)`,

	`CREATE TABLE IF NOT EXISTS messages (
		conversation_id TEXT NOT NULL,
		message_id TEXT NOT NULL,
		message TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (conversation_id, message_id)
	)`,
}
