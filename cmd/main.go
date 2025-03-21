package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	handler "github.com/tusmasoma/go-clean-arch/interfaces/handler"
	middleware "github.com/tusmasoma/go-clean-arch/interfaces/middleware"
	"github.com/tusmasoma/go-clean-arch/pkg/jwt"
	"github.com/tusmasoma/go-clean-arch/pkg/log"
	"github.com/tusmasoma/go-clean-arch/repository/mysql"
	"github.com/tusmasoma/go-clean-arch/usecase"

	"github.com/tusmasoma/go-clean-arch/pkg/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}
	var addr string
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()
	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()
	// --- DI ---
	serverConfig, err := config.NewServerConfig(mainCtx)
	if err != nil {
		log.Critical("Failed to initialize server config", log.Ferror(err))
		return
	}
	db, err := mysql.NewMySQLDB(mainCtx)
	if err != nil {
		log.Critical("Failed to initialize DB", log.Ferror(err))
		return
	}
	taskRepo := mysql.NewTaskRepository(db)
	userRepo := mysql.NewUserRepository(db)
	jwtGen := jwt.NewGenerator()
	taskUseCase := usecase.NewTaskUseCase(taskRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, jwtGen)
	taskHandler := handler.NewTaskHandler(taskUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	authMiddleware := middleware.NewAuthMiddleware(jwtGen)
	// --- Router ---
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"https://*", "http://*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:     []string{"Link", "Authorization"},
		AllowCredentials:   true,
		MaxAge:             serverConfig.PreflightCacheDurationSec,
		OptionsPassthrough: false,
	}))
	r.Use(middleware.Logging)
	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", userHandler.CreateUser)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/get", userHandler.GetUser)
				r.Put("/update", userHandler.UpdateUser)
			})
		})
		r.Route("/task", func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Get("/get", taskHandler.GetTask)
			r.Get("/list", taskHandler.ListTasks)
			r.Post("/create", taskHandler.CreateTask)
			r.Put("/update", taskHandler.UpdateTask)
			r.Delete("/delete", taskHandler.DeleteTask)
		})
	})
	// --- Server Run ---
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}
	log.Info("Server running...")
	// --- Graceful shutdown ---
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed", log.Ferror(err))
			return
		}
	}()
	<-signalCtx.Done()
	log.Info("Server stopping...")
	tctx, cancelShutdown := context.WithTimeout(context.Background(), serverConfig.GracefulShutdownTimeout)
	defer cancelShutdown()
	if err = srv.Shutdown(tctx); err != nil {
		log.Error("Failed to shutdown http server", log.Ferror(err))
	}
	log.Info("Server exited")
}
