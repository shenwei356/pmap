// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.
//
// ParallelMap - A lock-free parallel map in go.
//
// ParallelMap uses a backend goroutine for the sequential excution of
// functions including Get, Set, Update and custom function,
// which was inspired by section 14.17 in book *The Way to Go*.
//
// Example:
//
//    import (
//        "fmt"
//        "runtime"
//        "sort"
//        "sync"
//
//        "github.com/shenwei356/pmap"
//    )
//
//    func main() {
//        runtime.GOMAXPROCS(runtime.NumCPU())
//
//        // number of goroutines that will operate on ParallelMap
//        N := runtime.NumCPU() * 30
//
//        // constructor
//        m := pmap.NewParallelMap()
//
//        // In this exmaple, the Update function will be used.
//        // To call this function, the UpdateValueFunc must be specified.
//        m.SetUpdateValueFunc(func(oldValue interface{}, newValue interface{}) interface{} {
//            return oldValue.(int) + newValue.(int)
//        })
//
//        // number of elements in map
//        var n int = 1 << 9
//
//        // create N goroutines which call Update function concurrently
//        var wg sync.WaitGroup
//        for i := 1; i <= N; i++ {
//            wg.Add(1)
//
//            go func() {
//                defer func() {
//                    wg.Done()
//                }()
//
//                for j := 0; j < n; j++ {
//                    m.Update(j, 1)
//                }
//            }()
//        }
//        // wait for all operations to complement
//        wg.Wait()
//
//        // Stop the map backend
//        m.Stop()
//
//        // do something else
//        length := len(m.Map)
//        fmt.Printf("%d elements in map\n", length)
//    }
//
package pmap

import (
	"fmt"
	"os"
	"sync"
)

// ParallelMap
type ParallelMap struct {
	// map
	Map map[interface{}]interface{}

	// backend goroutine for sequential operations
	Op chan func() error
	// waitgroup for operations
	wg sync.WaitGroup

	// function to update value
	UpdateValueFunc func(interface{}, interface{}) interface{}
}

// Constructor of ParallelMap
func NewParallelMap() *ParallelMap {
	this := new(ParallelMap)
	this.Map = make(map[interface{}]interface{})
	this.Op = make(chan func() error)

	// by default, the Update function is equal to Set function.
	this.UpdateValueFunc = func(oldValue interface{}, newValue interface{}) interface{} {
		return newValue
	}

	go this.backend()
	return this
}

// Run operation channel as backend
func (this *ParallelMap) backend() {
	var f func() error
	var err error
	for {
		select {
		case f = <-this.Op:
			err = f()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
	}
}

// Stop the map backend
func (this *ParallelMap) Stop() {
	this.wg.Wait()
}

// Getting element of the map is executed sequentially
func (this *ParallelMap) Get(key interface{}) (interface{}, bool) {
	this.wg.Add(1)

	c1 := make(chan interface{})
	c2 := make(chan bool)
	this.Op <- func() error {
		value, ok := this.Map[key]

		c1 <- value
		c2 <- ok
		this.wg.Done()
		return nil
	}
	return <-c1, <-c2
}

// Setting operation is executed sequentially to ensure the
// operation is atomic.
func (this *ParallelMap) Set(key interface{}, value interface{}) {
	this.wg.Add(1)

	c := make(chan bool)
	this.Op <- func() error {
		this.Map[key] = value

		c <- true
		this.wg.Done()
		return nil
	}
	<-c
}

// To use Update function, a custom UpdateValueFunc must be set.
// By default, the Update function is equal to Set function.
//
// The default UpdateValueFunc is:
//
//    this.UpdateValueFunc = func(oldValue interface{}, newValue interface{}) interface{} {
//        return newValue
//    }
func (this *ParallelMap) SetUpdateValueFunc(f func(interface{}, interface{}) interface{}) {
	this.UpdateValueFunc = f
}

// Update function.
// To use Update function, a custom UpdateValueFunc must be set.
func (this *ParallelMap) Update(key interface{}, value interface{}) {
	this.wg.Add(1)

	c := make(chan bool)
	this.Op <- func() error {
		oldValue, ok := this.Map[key]
		if ok {
			this.Map[key] = this.UpdateValueFunc(oldValue, value)
		} else {
			this.Map[key] = value
		}

		c <- true
		this.wg.Done()
		return nil
	}
	<-c
}

// Execute a custom function.
//
// Example: An element increasing function
//
//    m.ExecuteFunc(func() error {
//        if v, ok := m.Map[i]; ok {
//            m.Map[i] = v.(int) + 1
//        } else {
//            m.Map[i] = int(1)
//        }
//        return nil
//    })
func (this *ParallelMap) ExecuteFunc(f func() error) {
	this.wg.Add(1)

	c := make(chan bool)
	this.Op <- func() error {
		err := f()

		c <- true
		this.wg.Done()
		return err
	}
	<-c
}
