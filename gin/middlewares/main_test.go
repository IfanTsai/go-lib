package middlewares_test

import (
	"go-lib/randutils"
	"go-lib/user/token"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	tokenMaker token.Maker
	router     *gin.Engine
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestServer(t *testing.T) *Server {
	t.Helper()

	server, err := newServer()
	require.NoError(t, err)

	return server
}

// NewServer creates a new HTTP server.
func newServer() (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(randutils.RandomString(32))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create token")
	}
	server := &Server{
		router:     gin.New(),
		tokenMaker: tokenMaker,
	}

	return server, nil
}
