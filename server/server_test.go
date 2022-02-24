package server

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/iKayrat/rest-counter/db"
	"github.com/stretchr/testify/require"
)

func mockRedis(t *testing.T) *miniredis.Miniredis {
	s, err := miniredis.Run()
	require.NoError(t, err)

	return s
}

func NewTestServer(t *testing.T) *Server {
	redis := mockRedis(t)
	testStore, err := db.NewStore(redis.Addr())
	require.NoError(t, err)

	server := NewServer(testStore)

	return server
}
