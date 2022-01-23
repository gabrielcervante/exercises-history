package main

import (
	"github.com/gabrielcervante/exercises-history/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	exerciseHandler := handlers.NewExercise()

	router.GET("/", exerciseHandler.GetExercises)
	router.POST("/", exerciseHandler.AddExercise)
	router.PUT("/", exerciseHandler.UpdateExercise)
	router.DELETE("/", exerciseHandler.DeleteExercise)

	router.Run()
}
