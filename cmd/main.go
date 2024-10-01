package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/iscritic/hot-coffee/internal/delivery"
	"github.com/iscritic/hot-coffee/internal/repository"
	"github.com/iscritic/hot-coffee/internal/service"
)

func main() {
	// Разбор аргументов командной строки
	help := flag.Bool("help", false, "Show this help message")
	port := flag.String("port", "8080", "Port number")
	dir := flag.String("dir", "./data", "Path to the data directory")
	flag.Parse()

	if *help {
		fmt.Println(`Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`)
		os.Exit(0)
	}

	// Инициализация репозиториев
	orderRepo := repository.NewOrderRepository(*dir + "/orders.json")
	// Здесь мы бы инициализировали menuRepo и inventoryRepo

	// Инициализация сервисов
	orderService := service.NewOrderService(orderRepo)
	// Здесь мы бы инициализировали menuService и inventoryService

	// Инициализация обработчиков
	mux := http.NewServeMux()
	handler.NewOrderHandler(mux, orderService)
	// Здесь мы бы инициализировали другие обработчики

	// Запуск сервера
	slog.Info("Server starting", "port", *port)
	err := http.ListenAndServe(":"+*port, mux)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
