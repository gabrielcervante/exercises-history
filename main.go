package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrielcervante/exercises-history/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	gin.DisableConsoleColor()

	f, err := os.OpenFile("gin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		fmt.Println(err)
		return

	}

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	defer f.Close()

	router := gin.Default()

	exerciseHandler := handlers.NewExercise()

	router.GET("/exercises", exerciseHandler.GetExercises)
	router.POST("/add", exerciseHandler.AddExercise)
	router.PUT("/update", exerciseHandler.UpdateExercise)
	router.DELETE("/delete", exerciseHandler.DeleteExercise)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {

		}

	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

	}

}
