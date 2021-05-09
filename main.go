package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/amrHassanAbdallah/tahweelaway/api"
	"github.com/amrHassanAbdallah/tahweelaway/persistence"
	"github.com/amrHassanAbdallah/tahweelaway/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"time"
)

type Logger = *zap.SugaredLogger

var log Logger

func init() {
	core := zap.NewProductionConfig()
	core.EncoderConfig.TimeKey = "timestamp"
	core.EncoderConfig.MessageKey = "message"
	core.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	customLog, _ := core.Build()
	log = customLog.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func main() {
	var (
		postgresURL = flag.String("postgresql-connection", "postgresql://root:secret@postgres:5432/tahweelaway?sslmode=disable", "DSN for postgres")
		listen      = flag.String("listen", ":8111", "Listen specified address.")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	conn, err := sql.Open("postgres", *postgresURL)
	if err != nil {
		log.Fatalw("cannot connect to db:", err)
	}

	persistenceLayer, err := persistence.NewStore(ctx, conn)
	if err != nil {
		log.Fatalw("failed to initate the persistence layer", err)
	}
	thService := service.NewService(persistenceLayer)
	server := api.NewServer(thService)
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Timeout(time.Duration(20) * time.Second))
		api.HandlerFromMux(server, r)
	})
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Handler: r,
		Addr:    *listen,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ES is a http.Handler, so you can pass it directly to your mux

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("server failed to start", "error", err)
		}
	}()
	log.Infof("server started on port %v", *listen)

	<-done
	log.Info("server terminating...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorw("server terminating failed", "err", err)
	}

	// cleanup: closing db...etc
	cancel()
	log.Info("exited")

}
