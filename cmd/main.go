package main

import (
	"GoExamGatewayAPI/internal/config"
	"GoExamGatewayAPI/internal/logger"
	"GoExamGatewayAPI/internal/server"
	"GoExamGatewayAPI/internal/stopsignal"
	"log/slog"
)

func main() {

	// Инициализируем конфиг файл и логгер.
	logger.SetupLogger()
	cfg := config.MustLoad()
	slog.Debug("config file and logger initialized")

	// Инициализируем сервер, объявляем обработчики API и запускаем сервер.
	srv := server.New(cfg)
	srv.API()
	srv.Middleware()
	srv.Start()
	slog.Info("Server started")

	// Блокируем выполнение основной горутины и ожидаем сигнала прерывания.
	stopsignal.Stop()

	// После сигнала прерывания останавливаем сервер.
	srv.Shutdown()

	slog.Info("Server stopped")
}
