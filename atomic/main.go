package main

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	jobChoice  string
	maxWorkers = 2
)

type obj struct {
	//
}

var (
	finish  chan struct{}
	done    chan struct{}
	jobChan chan struct{}
	start   chan struct{}
)

func init() {
	var err error
	usage := func() {
		println("syntax: " + os.Args[0] + " <job1|job2> <workers>")
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		usage()
	}
	if os.Args[1] != "job1" && os.Args[1] != "job2" {
		usage()
	}
	jobChoice = os.Args[1]
	maxWorkers, err = strconv.Atoi(os.Args[2])
	if err != nil {
		usage()
	}
	finish = make(chan struct{})
	done = make(chan struct{})
	jobChan = make(chan struct{})
	start = make(chan struct{})

	seen = make(map[int]bool)
}

func checkHeat() {
	if !atomic.CompareAndSwapUint32(&heatLocker, stateUnlocked, stateLocked) {
		return
	}
	defer atomic.StoreUint32(&heatLocker, stateUnlocked)
	// if working is zero, and finished is greater (hopefully not) or more likely equal to jobs, we are done
	if get(workingCt) == 0 && get(finishedCt) >= get(jobsCt) {
		// this will pop out at the bottom of main, allowing execution to complete and the program to exit
		finish <- obj{}
	}
}

func bigHomie() {
	<-start
	for {
		select {
		// job finished
		case <-done:
			// decrement working count
			dec(workingCt)
			// increase finished count
			inc(finishedCt)
		case <-jobChan:
			go work()
		default:
			// check if we're done
			checkHeat()
		}
	}
}

func work() {
	for {
		// we're at maximum capacity, don't start any new jobs
		if get(workingCt) >= maxWorkers {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// room for more workers available, continue execution
		break
	}
	// increase working count
	inc(workingCt)
	switch jobChoice {
	case "job1":
		go job1(get(jobsCt) - get(finishedCt))
	case "job2":
		go job2()
	default:
		panic("uhh? that job doesn't exist")
	}
}

var seen map[int]bool
var seenLock = stateUnlocked

func job1(input int) {
	var ok bool
	for !atomic.CompareAndSwapUint32(&seenLock, stateUnlocked, stateLocked) {
		randSleep()
	}
	if _, ok = seen[input]; !ok {
		seen[input] = true
		print(input)
		print(" ")
	}
	atomic.StoreUint32(&seenLock, stateUnlocked)

	// sleep just for visual effect
	time.Sleep(1 * time.Second)
	done <- obj{}
}

func job2() {
	fmt.Println(time.Now())
	done <- obj{}
}

func main() {
	go bigHomie()
	for n := 0; n < 100; n++ {
		inc(jobsCt)
		go func() {
			jobChan <- obj{}
		}()
	}
	println("\nrunning " + os.Args[1] + " with " + os.Args[2] + " worker(s)\n")
	if os.Args[1] == "job1" {
		print("remaining jobs: ")
	}
	start <- obj{}
	<-finish
	print("\n\ndone")
}
