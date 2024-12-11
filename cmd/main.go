package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-queue/internal/app"
	"task-queue/internal/delivery"
	"task-queue/internal/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("Ошибка при получении переменных окружения %s", err.Error())
	}
	services := service.NewService()
	go services.StartProcessing()
	handlers := delivery.NewHandler(services)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(os.Getenv("HTTP_PORT"), handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Ошибка при работе http-сервера: %s", err.Error())
		}
	}()

	logrus.Println("Сервис очереди задач начал работу")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Сервис очереди задач завершил работу")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Произошла ошибка при завершении работы сервиса: %s", err.Error())
	}
}
