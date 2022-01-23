package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gabrielcervante/exercises-history/data"
	"github.com/gin-gonic/gin"
)

type Exercise struct {
}

func NewExercise() *Exercise {
	return &Exercise{}
}

func (e *Exercise) GetExercises(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": data.GetExercises(),
	})

}

func (e *Exercise) AddExercise(c *gin.Context) {

	type NewExercises struct {
		ExerciseName string `json:"exerciseName"`
		DurationTime int    `json:"durationTime"`
	}

	var newExercises NewExercises

	c.BindJSON(&newExercises)

	if newExercises.ExerciseName == "" || newExercises.DurationTime == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "No exercise name or time provided",
		})
		return
	}

	if !isExerciseNameValid(newExercises.ExerciseName) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Sorry, exercise name is not using allowed letters",
		})

		return
	}

	data.AddExercise(newExercises.ExerciseName, newExercises.DurationTime)
	c.IndentedJSON(http.StatusOK, gin.H{
		"Success": newExercises,
	})
}

func (e *Exercise) UpdateExercise(c *gin.Context) {

	type ExerciseUpdate struct {
		Id           int    `json:"id"`
		ExerciseName string `json:"exerciseName"`
		DurationTime int    `json:"durationTime"`
	}

	var exerciseUpdate ExerciseUpdate

	c.BindJSON(&exerciseUpdate)

	if exerciseUpdate.Id == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "No id provided",
		})

		return
	}

	if exerciseUpdate.ExerciseName == "" && exerciseUpdate.DurationTime == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "You need to provide at least one of the exerciseName or durationTime fields",
		})

		return
	}

	if exerciseUpdate.ExerciseName != "" {
		if !isExerciseNameValid(exerciseUpdate.ExerciseName) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": "Sorry, exercise name is not using allowed letters",
			})

		}

		return
	}

	data.UpdateExercise(exerciseUpdate.Id, exerciseUpdate.ExerciseName, exerciseUpdate.DurationTime)

	c.IndentedJSON(http.StatusOK, gin.H{"success": "Your exercise has been updated"})

}

func (e *Exercise) DeleteExercise(c *gin.Context) {

	paramId := c.Query("id")

	if !isIdValid(paramId) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Sorry, no valid id value provided",
		})

		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		return
	}

	data.DeleteExercise(id)

	c.IndentedJSON(http.StatusOK, gin.H{"success": "Exercise successfull deleted"})
}

func isExerciseNameValid(exerciseName string) bool {

	if len(exerciseName) < 1 || len(exerciseName) > 30 {
		return false
	}

	checkExerciseName := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return checkExerciseName.MatchString(exerciseName)

}

func isIdValid(id string) bool {

	if len(id) < 1 {
		return false
	}
	checkId := regexp.MustCompile(`^[1-9]+$`)
	return checkId.MatchString(id)
}
