package handlers

import (
	"net/http"

	"github.com/gabrielcervante/exercises-history/data"
	"github.com/gin-gonic/gin"
)

type NewExercises struct {
	ExerciseName string `json:"exerciseName"`
	ExerciseTime int    `json:"exerciseTime"`
}

type Exercise struct {
}

func NewExercise() *Exercise {
	return &Exercise{}
}

func (e *Exercise) AddExercise(c *gin.Context) {

	var newExercises NewExercises

	c.BindJSON(&newExercises)

	if newExercises.ExerciseName == "" || newExercises.ExerciseTime == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "No exercise name or time provided",
		})
		return
	}

	data.AddExercise(newExercises.ExerciseName, newExercises.ExerciseTime)
	c.IndentedJSON(http.StatusOK, gin.H{
		"Success": newExercises,
	})
}
