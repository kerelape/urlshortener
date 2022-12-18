package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFakeDatabase_Put(t *testing.T) {
	require.Equal(t, uint(1), new(FakeDatabase).Put(""))
}

func TestFakeDatabase_Get(t *testing.T) {
	var db = new(FakeDatabase)
	var r, err = db.Get(db.Put("Hello, World!"))
	require.Nil(t, err)
	require.Equal(t, "Hello, World!", r)
}
