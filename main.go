package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrielcervante/exercises-history/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	//Setting gin to release mode
	gin.SetMode(gin.ReleaseMode)

	//The code below is used to log the gin output to a "log" file
	gin.DisableConsoleColor()

	//Opening a file with permissions to read, create and append
	f, err := os.OpenFile("gin.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	//Catching possible file error
	if err != nil {
		fmt.Println(err)
		return

	}

	//Saying to gin write output to the file
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//Closing the file when all of the things already is closed
	defer f.Close()

	//Creating gin routes
	router := gin.Default()

	//Getting my handlers from the handlers package
	exerciseHandler := handlers.NewExercise()

	router.GET("/exercises", exerciseHandler.GetExercises)
	router.GET("/exercises/:id", exerciseHandler.GetOneExercise)
	router.POST("/add", exerciseHandler.AddExercise)
	router.PUT("/update", exerciseHandler.UpdateExercise)
	router.DELETE("/delete", exerciseHandler.DeleteExercise)

	//Configuration to the server
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//Open the server and catch error
	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}

	}()

	//Creates a channel to a signal
	quit := make(chan os.Signal)

	//Notify the server to quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	//Start the server quit with 5 times of wait
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Quits the server
	if err := server.Shutdown(ctx); err != nil {

	}

}
