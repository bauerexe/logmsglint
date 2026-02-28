package a

import (
	"fmt"
	"log/slog"

	"go.uber.org/zap"
)

func slogCases(logger *slog.Logger) {
	slog.Info("Starting server")       // want "logmsg-lowercase"
	slog.Error("ошибка подключения")   // want "logmsg-english"
	slog.Warn("server started!!!")     // want "logmsg-nospecial"
	slog.Debug("bearer token exposed") // want "logmsg-sensitive"

	logger.Info("Another start") // want "logmsg-lowercase"

	msg := "Starting server"
	slog.Info(msg)
	slog.Info(fmt.Sprintf("%s", msg))
}

func zapCases() {
	logger := zap.NewNop()
	logger.Info("Starting worker")   // want "logmsg-lowercase"
	logger.Error("api_key detected") // want "logmsg-sensitive"
	logger.Warn("operation...")      // want "logmsg-nospecial"

	sugar := logger.Sugar()
	sugar.Debug("Тест") // want "logmsg-lowercase" "logmsg-english"
}
