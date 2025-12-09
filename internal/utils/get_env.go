package utils

import (
	"log/slog"
	"os"

	"github.com/hafiztri123/kki-be/internal/constants"
)




func GetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	} 

	slog.Error(constants.MsgEnvNotFound, "env", key)
	panic(constants.MsgEnvNotFound)
}