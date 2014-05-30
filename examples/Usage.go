package main

import (
	"fmt"

	"github.com/shenwei356/pmap"
)

func main() {

	// constructor
	m := pmap.NewParallelMap()

	// Set
	m.Set("age", 26)

	// Get
	if age, ok := m.Get("age"); ok {
		fmt.Printf("age: %d\n", age)
	}

	// Update
	//
	// To call this function, the UpdateValueFunc must be specified.
	m.SetUpdateValueFunc(func(oldValue interface{}, newValue interface{}) interface{} {
		return oldValue.(int) + newValue.(int)
	})

	m.Update("age", 1)
	if age, ok := m.Get("age"); ok {
		fmt.Printf("age: %d\n", age)
	}

	// Stop the map backend
	m.Stop()
}
