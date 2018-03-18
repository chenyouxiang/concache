package main

import (
	"math/rand"
	"sync"
	"fmt"
	"time"
	"sync/atomic"
	"os"
)

func main() {
	fmt.Println("numThreads numTrial totalOps totalDur, opsPerSecond")
	for numThreads := 1; numThreads < 9; numThreads++ {
		for trialNumber := 1; trialNumber <= 3; trialNumber++ {
			if len(os.Args) == 2 {
				val, dur := trial(numThreads, 5, os.Args[1])
				fmt.Println(numThreads, trialNumber, val, dur, float64(val)/dur.Seconds())
			} else {
				fmt.Println("Not proper number of arguments.")
			}
		}
	}
}

func trial (numThreads int, threadDuration int, readWrite string) (uint64, time.Duration) {
	var data = make(map[int]int)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	var ops uint64

	rand.Seed(time.Now().UnixNano()) //generate seed

	wg.Add(numThreads) //reader, writer


	timeStart := time.Now()

	for i:=0; i < numThreads; i++ {
		go func(from int) {
			defer wg.Done()
			var numOperations uint64 = 0

			for time.Now().Before(timeStart.Add((time.Duration(threadDuration) * time.Second))) {
				//just some random key/values
				for i := 0; i < 10000; i++ {
					var constant = rand.Int()%2 //read or write
					if readWrite == "rw" {
						if constant % 2 == 0 {
							mutex.Lock()
							data[constant] = constant
							mutex.Unlock()
							numOperations += 1
						} else {
							mutex.Lock()
							_ = data[constant]
							mutex.Unlock()
							numOperations += 1
						}
					} else if readWrite == "w" {
						mutex.Lock()
						data[constant] = constant
						mutex.Unlock()
						numOperations += 1
					} else if readWrite == "r" {
						mutex.Lock()
						_ = data[constant]
						mutex.Unlock()
						numOperations += 1
					} else {
						fmt.Println("Not proper choice.")
						break
					}	
				}
			}
			// fmt.Println("Number of Operations from Writer #", from, ": ", numOperations)
			atomic.AddUint64(&ops, numOperations)
		} (i)
	}
	wg.Wait() //wait for the goroutines to finish
	totalDuration := time.Since(timeStart)
	opsFinal := atomic.LoadUint64(&ops)
	
	// fmt.Println(opsFinal)

    return opsFinal, totalDuration
}