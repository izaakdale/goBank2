package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/izaakdale/goBank2/db/sqlc"
	"github.com/izaakdale/goBank2/util"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTokenDuration: time.Minute,
		TokenSymmetricKey:   util.RandomString(32),
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
