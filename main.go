package main

import (
	"github.com/gabrielcervante/exercises-history/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	exerciseHandler := handlers.NewExercise()

	router.POST("/", exerciseHandler.AddExercise)
	router.GET("/", exerciseHandler.GetExercises)

	router.Run()
}
