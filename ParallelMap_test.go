// Copyright 2014 Wei Shen (shenwei356@gmail.com). All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.
package pmap

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

func TestParallelMap(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// number of goroutines that will operate on ParallelMap
	N := runtime.NumCPU() * 30

	// constructor
	m := NewParallelMap()

	// In this exmaple, the Update function will be used.
	// To call this function, the UpdateValueFunc must be specified.
	m.SetUpdateValueFunc(func(oldValue interface{}, newValue interface{}) interface{} {
		return oldValue.(int) + newValue.(int)
	})

	// number of elements in map
	var n int = 1 << 9

	// create N goroutines which call Update function concurrently
	var wg sync.WaitGroup
	for i := 1; i <= N; i++ {
		wg.Add(1)

		go func() {
			defer func() {
				wg.Done()
			}()

			for j := 0; j < n; j++ {
				m.Update(j, 1)
			}
		}()
	}
	// wait for all operations to complement
	wg.Wait()

	// Stop the map backend
	m.Stop()

	// check length of map
	if len(m.Map) != n {
		t.Error("length error")
	}

	// check values
	for _, v := range m.Map {
		if v.(int) != int(N) {
			t.Error(fmt.Sprintf("value error: %d != %d", v.(int), int(N)))
		}
	}
}
