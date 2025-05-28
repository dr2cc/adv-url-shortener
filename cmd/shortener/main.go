package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"adv-url-shortener/internal/config"
	"adv-url-shortener/internal/http-server/handlers"
	mwLogger "adv-url-shortener/internal/http-server/middleware/logger"
	"adv-url-shortener/internal/lib/logger/sl"
	"adv-url-shortener/internal/storage/sqlite"

	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
)

// Объявить переменные окружения из iter5 так:
// $env:SERVER_ADDRESS = "localhost:8089"
// $env:BASE_URL  = "http://localhost:9999"

// Если использую local.yaml , то перед запуском нужно установить переменную окружения CONFIG_PATH
//
// $env:CONFIG_PATH = "C:\__git\adv-url-shortener\config\local.yaml"
// $env:CONFIG_PATH = "C:\Mega\__git\adv-url-shortener\config\local.yaml"  (на ноуте)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// обрабатываем аргументы командной строки
	config.ParseFlags()

	cfg := config.MustLoad()

	//
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	//
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	//переопределяем внутренний логгер
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// // adv #start#
	// router.Post("/", save.New(log, storage))
	// // adv #end#

	// // Примитивное хранилище - map
	// storageInstance := maps.NewStorage()

	// sqlite.New или "подключает" файл db , а если его нет то создает
	storageInstance, err := sqlite.New("./storage.db")
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	router.Post("/", handlers.PostHandler(log, storageInstance))
	router.Get("/{id}", handlers.GetHandler(log, storageInstance))

	// // примитивный запуск сервера
	//return http.ListenAndServe(config.FlagRunAddr, router)

	// adv #server#
	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
		//ReadTimeout:  cfg.HTTPServer.Timeout,
		//WriteTimeout: cfg.HTTPServer.Timeout,
		//IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
