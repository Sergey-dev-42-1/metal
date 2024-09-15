package main

import (
	"context"
	"fmt"
	"log"
	"metal/internal/pkg/domain/repositories"
	"metal/internal/pkg/domain/repositories/interfaces"
	"metal/internal/pkg/gzip"
	"metal/internal/pkg/logger"
	"metal/internal/server/presentation/controller"
	"metal/internal/server/presentation/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	parseFlags()

	// Logger
	zlog, err := logger.New("debug")
	if err != nil {
		fmt.Println("Error while trying to initiate logger")
	}
	// Repository
	var repo interfaces.MetricsStorage
	if connectionURL == "" {
		repo = repositories.NewMemStorage(fileStoragePath, zlog)
		if restore {
			if err := repo.Restore(); err != nil {
				zlog.Errorln("Couldn't restore saved data", err)
			}
		}
		if storeInterval > 0 {
			zlog.Infof("Initiate saving loop with %v seconds timeout", storeInterval)
			go func() {
				for {
					time.Sleep(time.Duration(storeInterval) * time.Second)
					repo.Save()
				}
			}()
		}
	} else {

		repo = repositories.NewSQLStorage(connectionURL, zlog)
	}
	// Router and middlewares
	r := router.Router()
	r.Use(logger.Logger())
	r.Use(gzip.GzipHandler())
	mc := controller.New(r, zlog, repo)
	r = mc.AddRoutes()
	//Не будет работать если запускать сервер не из корневой папки
	r.LoadHTMLGlob("internal/server/presentation/templates/*.html")

	srv := &http.Server{
		Addr:    startAddress.String(),
		Handler: r.Handler(),
	}
	//Server initialization
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Fatalf("Couldn't start the server: %s\n", err)
		}
	}()

	//Graceful shutdown
	// Из документации Gin https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	repo.Save()
	// catching ctx.Done(). timeout of 3 seconds.

	<-ctx.Done()
	log.Println("timeout of 3 seconds.")
}
