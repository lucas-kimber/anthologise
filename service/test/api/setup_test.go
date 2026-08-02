package api_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Setenv("APP_ENV", "test")

	slog.SetLogLoggerLevel(slog.LevelDebug)

	exitCode := m.Run()
	os.Exit(exitCode)
}
