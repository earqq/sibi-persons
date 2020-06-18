package graph

import (
	"crypto/rand"
	"fmt"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func GetRandomNumber() string {
	var array = make([]byte, 4)
	if _, err := rand.Read(array); err != nil {
		return ""
	}
	Random := fmt.Sprintf("%X", array)
	return Random
}
