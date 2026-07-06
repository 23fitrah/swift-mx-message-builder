package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"swift-mx-message-builder/config"
	"swift-mx-message-builder/handlers"
	"swift-mx-message-builder/utils"
	"swift-mx-message-builder/worker"
)

func main() {
	cfg := config.Load()

	pool := worker.NewPool(cfg.WorkerCount, cfg.OutputDir, cfg.QueueSize)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "swift-mx-message-builder"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.Use(utils.AuthMiddlewareGin())
		{
			pacs008 := v1.Group("/pacs008")
			{
				pacs008.POST("/generate", handlers.Pacs008Handler(pool))
				pacs008.GET("/inquiry/:messageId", handlers.InquiryHandler(pool))
			}

			pacs009 := v1.Group("/pacs009")
			{
				pacs009.POST("/generate", handlers.Pacs009Handler(pool))
				pacs009.GET("/inquiry/:messageId", handlers.InquiryHandler(pool))
			}

			pacs002 := v1.Group("/pacs002")
			{
				pacs002.POST("/generate", handlers.Pacs002Handler(pool))
				pacs002.GET("/inquiry/:messageId", handlers.InquiryHandler(pool))
			}

			pacs004 := v1.Group("/pacs004")
			{
				pacs004.POST("/generate", handlers.Pacs004Handler(pool))
				pacs004.GET("/inquiry/:messageId", handlers.InquiryHandler(pool))
			}
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("swift-mx-message-builder listening on :%s (output dir: %s, workers: %d)\n",
			cfg.Port, cfg.OutputDir, cfg.WorkerCount)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	pool.Shutdown()
	log.Println("server exited cleanly")
}
