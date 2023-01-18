package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakeDatabase_Put(t *testing.T) {
	var db = NewFakeDatabase()
	var id, putError = db.Put("")
	require.Nil(t, putError)
	require.Equal(t, uint(0), id)
}

func TestFakeDatabase_Get(t *testing.T) {
	var db = NewFakeDatabase()
	var id, putError = db.Put("Hello, World!")
	var r, err = db.Get(id)
	require.Nil(t, err)
	require.Nil(t, putError)
	require.Equal(t, "Hello, World!", r)
	r, err = db.Get(92)
	require.NotNil(t, err)
	require.Equal(t, "", r)
}
