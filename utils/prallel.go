package utils

import (
	"sync"
)

func PFor[T any](arr []T, fn func(i int, el T)) {
	var wg sync.WaitGroup
	wg.Add(len(arr))
	for i, _ := range arr {
		go func(i int) {
			fn(i, arr[i])
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func PForLim[T any](arr []T, lim int, fn func(int, T)) {
	var wg sync.WaitGroup
	wg.Add(len(arr))
	limits := make(chan struct{}, lim)
	for i, _ := range arr {
		limits <- struct{}{}
		go func(i int) {
			fn(i, arr[i])
			<-limits
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func ParallelFor(begin, end int, fn func(i int)) {
	var wg sync.WaitGroup
	wg.Add(end - begin)
	for i := begin; i < end; i++ {
		//go func(i int) {
		//	fn(i)
		//	wg.Done()
		//}(i)
		fn(i)
		wg.Done()
	}
	wg.Wait()
}

func ParallelLimFor(begin, end, lim int, fn func(i int)) {
	var wg sync.WaitGroup
	wg.Add(end - begin)
	limits := make(chan struct{}, lim)
	for i := begin; i < end; i++ {
		limits <- struct{}{}
		go func(i int) {
			fn(i)
			<-limits
			wg.Done()
		}(i)
	}
	wg.Wait()
}
