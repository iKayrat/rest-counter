package db

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)

	store, err := NewStore(s.Addr())
	require.NoError(t, err)

	ping, err := store.Client.Ping().Result()
	require.NoError(t, err)

	require.Equal(t, ping, "PONG")
}
