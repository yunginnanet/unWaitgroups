package main

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type counterID uint32

const (
	finishedCt counterID = iota // 0
	workingCt                   // 1
	jobsCt                      // 2
	lastNumber                  // 3
)

const (
	stateUnlocked uint32 = iota
	stateLocked
)

var (
	setLocker  = stateUnlocked
	incLocker  = stateUnlocked
	decLocker  = stateUnlocked
	heatLocker = stateUnlocked
)

type counter struct {
	value atomic.Value
}

var counters map[counterID]*counter

func randSleep() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

func init() {

	// instantiate our map
	counters = make(map[counterID]*counter)

	// instantiate value pointer to each of our counters, referenced by it's respective counterID via our map
	for ctr := uint32(0); ctr < 4; ctr++ {
		counters[counterID(ctr)] = &counter{value: atomic.Value{}}
		// store 0 in each counter otherwise we will return nil when trying to do math and panic
		counters[counterID(ctr)].value.Store(0)
	}
}

func get(ctr counterID) int {
	return counters[ctr].value.Load().(int)
}

func set(ctr counterID, val int) {
	for !atomic.CompareAndSwapUint32(&setLocker, stateUnlocked, stateLocked) {
		randSleep()
	}
	defer atomic.StoreUint32(&setLocker, stateUnlocked)
	counters[ctr].value.Store(val)
}

func inc(ctr counterID) {
	for !atomic.CompareAndSwapUint32(&incLocker, stateUnlocked, stateLocked) {
		randSleep()
	}
	defer atomic.StoreUint32(&incLocker, stateUnlocked)
	counters[ctr].value.Store(get(ctr) + 1)
}

func dec(ctr counterID) {
	for !atomic.CompareAndSwapUint32(&decLocker, stateUnlocked, stateLocked) {
		randSleep()
	}
	defer atomic.StoreUint32(&decLocker, stateUnlocked)
	counters[ctr].value.Store(get(ctr) - 1)
}
