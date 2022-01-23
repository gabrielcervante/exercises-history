package data

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Exercises struct {
	Id             int    `json:"id"`
	Exercise_name  string `json:"exerciseName"`
	Duration_time  int    `json:"durationTime"`
	Timestamp_date int64  `json:"timestamp"`
}

func database() (*gorm.DB, error) {
	dsn := "host=localhost user=cervante password=cervantepswd dbname=exercises port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func GetExercises() []Exercises {

	db, err := database()

	if err != nil {

	}

	var exercises []Exercises
	db.Raw("SELECT id, exercise_name, duration_time,timestamp_date FROM history").Scan(&exercises)

	return exercises
}

func AddExercise(exerciseName string, durationTime int) {

	db, err := database()

	if err != nil {
		return
	}

	newExerciseId := createId(db)

	currentTimeStamp := time.Now().Unix()

	var exercises Exercises

	db.Raw("INSERT INTO history (id, exercise_name, duration_time, timestamp_date) VALUES (?,?,?,?)", newExerciseId, exerciseName, durationTime, currentTimeStamp).Scan(&exercises)

}

func UpdateExercise(id int, exerciseName string, durationTime int) {

	db, err := database()

	if err != nil {
		return
	}

	type Exerc struct {
		ExercName string
		DurationTime int
		Id        int
	}

	var exercise Exerc

	switch {
	case exerciseName != "":
		db.Raw("UPDATE history SET exercise_name = ? WHERE id = ?", exerciseName, id).Scan(&exercise.ExercName)
	case durationTime != 0:
		db.Raw("UPDATE history SET duration_time = ? WHERE id = ?", durationTime, id).Scan(&exercise.DurationTime)
	}

	if exerciseName != "" && durationTime != 0 {
		db.Raw("UPDATE history SET exercise_name = ?, duration_time = ? WHERE id = ?", exerciseName, durationTime, id).Scan(&exercise)
	}

}

func DeleteExercise(id int) {

	db, err := database()

	if err != nil {
		return
	}

	db.Raw("DELETE from history WHERE id = ?", id).Scan(&id)

	db.Raw("UPDATE history SET id = id - 1 WHERE Id > ?", id-1).Scan(&id)

}

func createId(db *gorm.DB) int {

	id := 1

	var exercise Exercises

	db.Raw("SELECT id FROM history WHERE id = ?", id).Scan(&exercise.Id)

	if exercise.Id == 0 {

		return id
	}

	db.Raw("SELECT id FROM history ORDER BY id DESC LIMIT 1").Scan(&id)

	return id + 1
}
