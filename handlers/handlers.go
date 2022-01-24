package handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gabrielcervante/exercises-history/data"
	"github.com/gin-gonic/gin"
)

//Struct to create methods and handle it
type Exercise struct {
}

//NewExercise is used to instance the methods below
func NewExercise() *Exercise {
	return &Exercise{}
}

//GetOneExercise in handle is used to get one exercise
func (e *Exercise) GetOneExercise(c *gin.Context) {

	//Gettings the paramter id
	paramId := c.Param("id")

	//Checking if the id is valid
	if !isIdValid(paramId) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Sorry, no valid id value provided",
		})

		return
	}

	//Converting parameter to int
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry, something went wrong",
		})
		return
	}

	//Creating channels to fill the func
	exerciseCount := make(chan int, 1)

	hasDbError := make(chan int, 1)
	getExercise := make(chan data.Exercises, 1)

	//Goroutine on the func instanced in history
	go data.GetOneExercise(id, hasDbError, getExercise, exerciseCount)

	//Gettings channels values
	errorStatusCode := <-hasDbError

	maxExercises := <-exerciseCount

	//Checking if the requested id is valid
	if id > maxExercises || maxExercise == 0{
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "The id provided is not associated to an exercise"})
		return
	}

	//Checking the channel to see if the function returned an error
	if errorStatusCode == 200 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"success": <-getExercise,
		})

	} else if errorStatusCode == 500 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Sorry, something went wrong"})

	}
}

//GetExercises is the handle to get all of the exercises
func (e *Exercise) GetExercises(c *gin.Context) {

	hasDbError := make(chan int, 1)

	getExercises := make(chan []data.Exercises, 1)

	go data.GetExercises(hasDbError, getExercises)

	errorStatusCode := <-hasDbError

	if errorStatusCode == 200 {

		c.IndentedJSON(http.StatusOK, gin.H{
			"success": <-getExercises,
		})
		return
	} else if errorStatusCode == 500 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Sorry, something went wrong"})
		return
	}

}

//AddExercise fun is used to handle the AddExercise function
func (e *Exercise) AddExercise(c *gin.Context) {

	//Struct to get the correct json fields
	type NewExercises struct {
		ExerciseName string `json:"exerciseName"`
		DurationTime int    `json:"durationTime"`
	}

	//Binding the json to get the fields
	var newExercises NewExercises

	c.BindJSON(&newExercises)

	//Checks to see if the consumer is using the api properly
	if newExercises.ExerciseName == "" || newExercises.DurationTime == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "No exercise name or time provided",
		})
		return
	}

	//Checking if the exercise name is valid and hasn't errors
	if !isExerciseNameValid(newExercises.ExerciseName) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Sorry, exercise name is not using allowed letters",
		})

		return
	}

	hasDbError := make(chan int, 1)

	go data.AddExercise(newExercises.ExerciseName, newExercises.DurationTime, hasDbError)

	errorStatusCode := <-hasDbError

	if errorStatusCode == 200 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"success": newExercises,
		})
		return
	} else if errorStatusCode == 500 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry, something went wrong",
		})
		return
	}

}

//UpdateExercise is used to handle the UpdateExercise method
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

	hasDbError := make(chan int, 1)
	exerciseCount := make(chan int, 1)
	
	go data.UpdateExercise(exerciseUpdate.Id, exerciseUpdate.ExerciseName, exerciseUpdate.DurationTime, hasDbError, exerciseCount)

	maxExercises := <-exerciseCount
	errorStatusCode := <-hasDbError

	if id > maxExercises || maxExercise == 0 {
                c.IndentedJSON(http.StatusNotFound, gin.H{"error": "The id provided is not associated to an exercise"}
)                                                                          return
        }

	if errorStatusCode == 200 {
		c.IndentedJSON(http.StatusOK, gin.H{"success": exerciseUpdate})
		return
	} else if errorStatusCode == 500 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry, something went wrong",
		})
		return
	}
}

//DeleteExercise is used to delete some exercise with the given id
func (e *Exercise) DeleteExercise(c *gin.Context) {

	//QueryParameter to delete the exercise by using the id
	paramId := c.Query("id")

	if !isIdValid(paramId) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Sorry, no valid id value provided",
		})

		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry, something went wrong",
		})
		return
	}

	hasDbError := make(chan int, 1)
	exerciseCount := make(chan int, 1)

	go data.DeleteExercise(id, hasDbError, exerciseCount)

	maxExercises := <-exerciseCount
	errorStatusCode := <-hasDbError

	if id > maxExercises || maxExercise == 0{
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "The id provided is not associated to an exercise"})
		return
	}

	if errorStatusCode == 200 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"success": "Exercise successfull deleted",
			"id":      id,
		})
		return
	} else if errorStatusCode == 500 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry, something went wrong",
		})
		return
	}

}

//Regexp func to validate the exercise name
func isExerciseNameValid(exerciseName string) bool {

	if len(exerciseName) < 1 || len(exerciseName) > 30 {
		return false
	}

	checkExerciseName := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return checkExerciseName.MatchString(exerciseName)

}

//Regexp func to validade exercise id
func isIdValid(id string) bool {

	if len(id) < 1 {
		return false
	}
	checkId := regexp.MustCompile(`^(?:[1-9]|\d\d\d*)$`)
	return checkId.MatchString(id)
}
