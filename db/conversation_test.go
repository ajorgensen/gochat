package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindConversation(t *testing.T) {
	dbc, err := Connect(":memory:")
	require.NoError(t, err)

	id, err := dbc.CreateConversation("test")
	require.NoError(t, err)

	c, err := dbc.FindConversation(id)
	require.NoError(t, err)

	assert.Equal(t, "test", c.Title)
}

func TestSelectConversations(t *testing.T) {
	dbc, err := Connect(":memory:")
	require.NoError(t, err)

	_, err = dbc.CreateConversation("test")
	require.NoError(t, err)

	_, err = dbc.CreateConversation("other conversation")
	require.NoError(t, err)

	conversations, err := dbc.SelectConversations()
	require.NoError(t, err)

	assert.Len(t, conversations, 2)
}
