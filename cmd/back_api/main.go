package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OoThan/usermanagement/cmd/back_api/handler"
	"github.com/OoThan/usermanagement/internal/ds"
	"github.com/OoThan/usermanagement/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.String("port", "7001", "default port is 7001")
	flag.Parse()

	addr := net.JoinHostPort("", *port)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	ds, err := ds.NewDataSource()
	if err != nil {
		logger.Sugar.Panic(err.Error())
	}

	h := handler.NewHandler(
		&handler.HConfig{
			R:  router,
			DS: ds,
		},
	)
	h.Register()

	server := http.Server{
		Addr:           addr,
		Handler:        h.R,
		ReadTimeout:    time.Duration(time.Minute * 3),
		WriteTimeout:   time.Duration(time.Minute * 3),
		MaxHeaderBytes: 10 << 20, // 10MB
	}

	go func() {
		logger.Sugar.Infof("server started listening on port: %v", *port)
		if err := server.ListenAndServe(); err != nil {
			logger.Sugar.Panicf("server failed to initialized on port: %v", *port)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Sugar.Panicf("Failed to shutdown server: %v", err.Error())
	}
}
