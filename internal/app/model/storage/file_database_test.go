package storage

import (
	"context"
	"strconv"
	"testing"

	"github.com/kerelape/urlshortener/internal/app"
)

func BenchmarkFileDatabase(b *testing.B) {
	ids := make([]uint, 0)
	b.Run("Put", func(b *testing.B) {
		b.StopTimer()
		database, err := OpenFileDatabase(b.TempDir()+"test.db", true, 0o644, 1024)
		if err != nil {
			b.Fatal("Failed to open file database", err)
		}
		token := app.NewToken()

		for i := 0; i < b.N; i++ {
			b.StartTimer()
			id, err := database.Put(context.Background(), token, strconv.Itoa(i))
			b.StopTimer()
			if err != nil {
				b.Fatal(err)
			}
			ids = append(ids, id)
		}
	})
}
