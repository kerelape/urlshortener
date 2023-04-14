package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakeDatabase_Put(t *testing.T) {
	db := NewFakeDatabase()
	id, putError := db.Put(context.Background(), "")
	require.Nil(t, putError)
	require.Equal(t, uint(0), id)
}

func TestFakeDatabase_Get(t *testing.T) {
	db := NewFakeDatabase()
	id, putError := db.Put(context.Background(), "Hello, World!")
	r, err := db.Get(context.Background(), id)
	require.Nil(t, err)
	require.Nil(t, putError)
	require.Equal(t, "Hello, World!", r)
	r, err = db.Get(context.Background(), 92)
	require.NotNil(t, err)
	require.Equal(t, "", r)
}
