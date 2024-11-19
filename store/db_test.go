package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConn(t *testing.T) {
	t.Run("Should connect to database", func(t *testing.T) {
		db, err := NewConn()
		assert.Nil(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = db.PingContext(ctx)
		assert.Nil(t, err)
		cancel()
	})
}
