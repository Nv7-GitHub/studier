package main

import (
	"fmt"
	"os"
)

type ModelState int

const (
	ModelStateFileInput ModelState = iota
)

// TODO: Better error handling
func (m *Model) HandleErr(msg string) {
	// Red
	fmt.Println(ErrStyle.Render(msg))
	os.Exit(1)
}

type Model struct {
	State ModelState
}
