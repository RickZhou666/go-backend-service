package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.EnvGinMode)
	os.Exit(m.Run())
}
