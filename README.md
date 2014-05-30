pmap (ParallelMap)
==================

A lock-free parallel map in go.

ParallelMap uses a backend goroutine for the sequential excution of 
functions including Get, Set, Update and custom function, 
which was inspired by section 14.17 in book *The Way to Go*.

Install
-------
This package is "go-gettable", just:

    go get github.com/shenwei356/pmap

Usage
-------
    
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
 

Documentation
-------------

[See documentation on gowalker for more detail](http://gowalker.org/github.com/shenwei356/pmap).

[MIT License](https://github.com/shenwei356/pmap/blob/master/LICENSE)