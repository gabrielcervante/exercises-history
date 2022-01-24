package data

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Exercise struct to help make queries in postgresql db using gorm
type Exercises struct {
	Id             int    `json:"id"`
	Exercise_name  string `json:"exerciseName"`
	Duration_time  int    `json:"durationTime"`
	Timestamp_date int64  `json:"timestamp"`
}

//database func is used to open the database
func database() (*gorm.DB, error) {
	dsn := "host=localhost user=cervant password=cervantepswd dbname=exercises port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

//GetOneExercise func is used to get the data of just 1 exercise
func GetOneExercise(id int, channelStatusCode chan int, chanExercise chan Exercises, exerciseCounting chan int) {

	//Opens db instance
	db, err := database()

	//Catch db error
	if err != nil {
		//channel status code is used to says if it's have success or failed in the request, 500 means internal server error and 200 means OK
		channelStatusCode <- 500
		return
	}

	//Get the details of a given exercise by id
	var exercise Exercises
	if err := db.Raw("SELECT id, exercise_name, duration_time,timestamp_date FROM history WHERE id = ?", id).Scan(&exercise).Error; err != nil {
		channelStatusCode <- 500
		return
	}

	//Gets the total number of exercises registered
	var exerciseCount int
	if err := db.Raw("SELECT count(*) FROM history").Scan(&exerciseCount).Error; err != nil {
		channelStatusCode <- 500
		return
	}

	//First channel returns the exercise, the second returns the exercise count and the last returns success to the request.
	chanExercise <- exercise
	exerciseCounting <- exerciseCount
	channelStatusCode <- 200

}

//GetExercises is a function to get all of the exercises
func GetExercises(channelStatusCode chan int, chanExercises chan []Exercises) {

	db, err := database()

	if err != nil {
		channelStatusCode <- 500
		return
	}

	//Getting all of the exercises using a slice o the Exercises struct
	var exercises []Exercises
	if err := db.Raw("SELECT id, exercise_name, duration_time,timestamp_date FROM history").Scan(&exercises).Error; err != nil {
		channelStatusCode <- 500
		return
	}

	//returning the slice of exercises and the status code
	chanExercises <- exercises
	channelStatusCode <- 200

}

//AddExercise adds a new exercise
func AddExercise(exerciseName string, durationTime int, channel chan int) {

	db, err := database()

	if err != nil {
		channel <- 500
		return
	}

	//Creating a new id without duplicates
	newExerciseId := createId(db, channel)

	//if the id is equals 0 then it will stop the request
	if newExerciseId == 0 {
		return
	}

	//Gets the current timestamp
	currentTimeStamp := time.Now().Unix()

	//Creates the INSERT to add the new exercise to the database
	var exercises Exercises

	if err := db.Raw("INSERT INTO history (id, exercise_name, duration_time, timestamp_date) VALUES (?,?,?,?)", newExerciseId, exerciseName, durationTime, currentTimeStamp).Scan(&exercises).Error; err != nil {
		channel <- 500
		return
	}

	channel <- 200

}

//UpdateExercise updates the name or duration of an exercise
func UpdateExercise(id int, exerciseName string, durationTime int, channel chan int, exerciseCounting chan int) {

	db, err := database()

	if err != nil {
		channel <- 500
		return
	}

	//New struct to get just the needed fields
	type Exerc struct {
		ExercName    string
		DurationTime int
		Id           int
	}

	/*The user can want to modify just the name, just the duration time or both, go doesn't have default value for parameters or optional parameter and how i
	don't wanned use spread parameter to have a bigger control over the function the consumer can or pass empty double quotes "" to says to the request not
	update the name or the zero value 0 to say to the request don't change the duration time, one of the two parameters has to be filles or both.
	*/
	var exercise Exerc

	switch {
	case exerciseName != "":
		if err := db.Raw("UPDATE history SET exercise_name = ? WHERE id = ?", exerciseName, id).Scan(&exercise.ExercName).Error; err != nil {
			channel <- 500
			return
		}
	case durationTime != 0:
		if err := db.Raw("UPDATE history SET duration_time = ? WHERE id = ?", durationTime, id).Scan(&exercise.DurationTime).Error; err != nil {
			channel <- 500
			return
		}
	}

	if exerciseName != "" && durationTime != 0 {
		if err := db.Raw("UPDATE history SET exercise_name = ?, duration_time = ? WHERE id = ?", exerciseName, durationTime, id).Scan(&exercise).Error; err != nil {
			channel <- 500
			return
		}
	}
	
	var exerciseCount int
        if err := db.Raw("SELECT count(*) FROM history").Scan(&exerciseCount).Error; err != nil {
                channel <- 500                                             return
        }                                                                                                                     exerciseCounting <- exerciseCount
	channel <- 200

}

//Delete exercise is used to delete an exercise
func DeleteExercise(id int, channel chan int, exerciseCounting chan int) {

	db, err := database()

	if err != nil {
		channel <- 500
		return
	}

	//Query to delete the exercise
	if err := db.Raw("DELETE from history WHERE id = ?", id).Scan(&id).Error; err != nil {
		channel <- 500
		return
	}

	//Changing all of the ids to decrease 1 number and don't have blank ids
	if err := db.Raw("UPDATE history SET id = id - 1 WHERE Id > ?", id-1).Scan(&id).Error; err != nil {
		channel <- 500
		return
	}

	//Counting the numbers of ids to not allow the user request to delete an id that doesn't exist
	var exerciseCount int
	if err := db.Raw("SELECT count(*) FROM history").Scan(&exerciseCount).Error; err != nil {
		channel <- 500
		return
	}

	exerciseCounting <- exerciseCount
	channel <- 200

}

//Create id is a function to create distinct ids automatically to the new exercises
func createId(db *gorm.DB, channel chan int) int {

	id := 1

	var exercise Exercises
	//The query is looking ofr an id that has the value 1 and if it doesn't exist so it will create the first exercise
	if err := db.Raw("SELECT id FROM history WHERE id = ?", id).Scan(&exercise.Id).Error; err != nil {
		channel <- 500
		return 0
	}

	//Checking if the id 1 has, if the id 1 doesn't exist means that the database have not records
	if exercise.Id == 0 {

		return id
	}

	//If the id 1 exist then get the lenght of the rows and add 1 number to get a new id
	if err := db.Raw("SELECT id FROM history ORDER BY id DESC LIMIT 1").Scan(&id).Error; err != nil {
		channel <- 500
		return 0
	}

	return id + 1
}
