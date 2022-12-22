package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakeDatabase_Put(t *testing.T) {
	require.Equal(t, uint(0), new(FakeDatabase).Put(""))
}

func TestFakeDatabase_Get(t *testing.T) {
	var db = NewFakeDatabase()
	var r, err = db.Get(db.Put("Hello, World!"))
	require.Nil(t, err)
	require.Equal(t, "Hello, World!", r)
	r, err = db.Get(92)
	require.NotNil(t, err)
	require.Equal(t, "", r)
}
