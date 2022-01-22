package data

import (
	"testing"
)

func TestGetData(t *testing.T) {

	//fmt.Println(GetExercises())

	//DeleteExercise(3)

	err := AddExercise("execs", 50)

	if err != nil {

		return

	}

	//UpdateExercise(3, "gosrt", 777)
}
