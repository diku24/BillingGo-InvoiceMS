package main

import (
	"InvoiceMS/config"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(serverPort ...string) *gin.Engine {
	configuredPort := ""
	if len(serverPort) > 0 {
		configuredPort = serverPort[0]
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/health/live", LiveHealthHandler())
	router.GET("/ping", ServerPingHandler(configuredPort))
	return router
}

func LiveHealthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func ServerPingHandler(serverPort string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		message := "Server is running"
		if serverPort != "" {
			message = fmt.Sprintf("Server is running on Port %s", serverPort)
		}
		ctx.JSON(http.StatusOK, message)
	}
}

func main() {

	cfg, err := config.LoadDefault()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	router := NewRouter(cfg.ServerPort)

	serve := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router.Handler(),
	}

	listener, err := net.Listen(cfg.NetworkProtocol, cfg.ServerPort)
	if err != nil {
		log.Fatalf("listen on %s %s: %v", cfg.NetworkProtocol, cfg.ServerPort, err)
	}

	go func() {
		// service connections
		if err := serve.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	//wait for interrupt signal to gracefully shutdown the server with
	//a timeout of 5 seconds
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't  need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serve.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Server ...")
	}

	//catching ctx.Done(). timeout of 5 seconds
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")
	}
	log.Println("Server Existing ...")
}
