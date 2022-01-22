package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gabrielcervante/exercises-history/data"
	"github.com/gin-gonic/gin"
)

type NewExercises struct {
	ExerciseName string `json:"exerciseName"`
	ExerciseTime int    `json:"exerciseTime"`
}

type ExerciseId struct {
	Id int `json:"id"`
}

type Exercise struct {
}

func NewExercise() *Exercise {
	return &Exercise{}
}

func (e *Exercise) GetExercises(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{
		"exercisesHistory": data.GetExercises(),
	})

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

func (e *Exercise) DeleteExercise(c *gin.Context) {

	paramId := c.Query("id")

	if paramId == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"response": "Error, no value in id input"})
		return
	}

	if paramId == "0" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"response": "Sorry, value 0 is not a valid input"})
		return
	}

	id, err := strconv.Atoi(paramId)
	fmt.Println(id)
	if err != nil {
		return
	}

	data.DeleteExercise(id)

}
