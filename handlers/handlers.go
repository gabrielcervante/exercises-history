package main //handlers

import (
	"fmt"
)

type Exercise struct {
}

func NewExercise() *Exercise {
	return &Exercise{}
}

func main() {
	fmt.Println("vim-go")
}
