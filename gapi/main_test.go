package gapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/tonisco/simple-bank-go/db/sqlc"
	"github.com/tonisco/simple-bank-go/util"
	"github.com/tonisco/simple-bank-go/worker"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}
