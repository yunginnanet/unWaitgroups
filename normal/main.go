package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	jobChoice  string
	maxWorkers = 2
	finished   int
	working    int
	jobs       int
	l          *sync.RWMutex
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
	l = &sync.RWMutex{}
}

func checkHeat() {
	l.RLock()
	defer l.RUnlock()
	if working == 0 && finished >= jobs {
		finish <- obj{}
	}
}

func bigHomie() {
	<-start
	for {
		select {
		case <-done:
			l.Lock()
			working--
			finished++
			l.Unlock()
		case <-jobChan:
			go work()
		default:
			checkHeat()
		}
	}
}

func work() {
	for {
		l.RLock()
		if working > maxWorkers {
			l.RUnlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}
		l.RUnlock()
		break
	}
	l.Lock()
	working++
	switch jobChoice {
	case "job1":
		go job1(jobs - finished)
	case "job2":
		go job2()
	}
	l.Unlock()
}

func job1(input int) {
	l.Lock()

	if input != lastcount {
		print(input)
		print(" ")
	}

	lastcount = input
	l.Unlock()

	// sleep just for visual effect
	time.Sleep(1 * time.Second)

	done <- obj{}
}

func job2() {
	fmt.Println(time.Now())
	// sleep just for visual effect
	done <- obj{}
}

func main() {
	go bigHomie()
	for n := 0; n < 100; n++ {
		jobs++
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
