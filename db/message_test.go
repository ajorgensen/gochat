package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMessage(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	converationID := "abc123"
	message := "Hello, world!"

	err = db.CreateMessage(converationID, message)
	assert.NoError(t, err)
}

func TestGetMessages(t *testing.T) {
	db, err := Connect(":memory:")
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}

	converationID := "abc123"
	message := "Hello, world!"

	err = db.CreateMessage(converationID, message)
	assert.NoError(t, err)

	conversation, err := db.GetMessages(converationID)
	assert.NoError(t, err)

	require.Equal(t, 1, len(conversation))
	assert.Equal(t, message, conversation[0].Message)
}
