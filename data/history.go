package data

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Exercises struct {
	id             int
	exercise_name  string
	duration_time  int
	timestamp_date int64
}

func database() (*gorm.DB, error) {
	dsn := "host=localhost user=cervante password=cervantepswd dbname=exercises port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
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

func createId(db *gorm.DB) int {

	id := 1

	var exercise Exercises

	db.Raw("SELECT id FROM history WHERE id = ?", id).Scan(&exercise.id)

	if exercise.id == 0 {

		return id
	}

	db.Raw("SELECT id FROM history ORDER BY id DESC LIMIT 1").Scan(&id)

	return id + 1
}
