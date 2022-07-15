package main

import (
	"context"
	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// test case1:
// curl -i http://localhost:8080/c

func main() {
	engine := gin.Default()

	engine.GET("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, "this is page A")
	})
	engine.GET("/b", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/a")
	})
	engine.Use(timeout.Timeout(
		timeout.WithTimeout(6 * time.Second),
	))
	engine.GET("/c", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/a")
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown:", err)
	}
}
