package main

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	
	"github.com/shenwei356/pmap"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// number of goroutines that will operate on ParallelMap
	N := runtime.NumCPU() * 30

	// constructor
	m := pmap.NewParallelMap()

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

	// do something else
	length := len(m.Map)
	fmt.Printf("%d elements in map\n", length)

	keys := make([]int, length)
	i := 0
	for k, _ := range m.Map {
		keys[i] = k.(int)
		i++
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("%d => %d\n", k, m.Map[k])
	}
}
