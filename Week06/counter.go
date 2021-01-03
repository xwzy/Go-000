package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Counter struct {
	// bucketTimePeriodLength the time duration for a single bucket
	bucketTimePeriodLength time.Duration

	// windowSize the number of bucket in sliding window
	windowSize int32

	// buckets bucket list
	buckets []int64

	// position current pointer position in buckets list
	position int32

	// lastTimestamp last bucket create time
	lastTimestamp time.Time

	//// mutex mutex for read list
	//mutex sync.Mutex
}

// DefaultCounter default construct
func DefaultCounter() *Counter {
	return NewCounter(time.Second, 10)
}

// NewCounter counter factory
func NewCounter(bucketTimePeriodLength time.Duration, windowSize int32) *Counter {
	c := &Counter{
		bucketTimePeriodLength: bucketTimePeriodLength,
		windowSize:             windowSize,
		buckets:                make([]int64, windowSize),
		position:               0,
		lastTimestamp:          time.Now(),
	}
	go c.bucketMaintainer()

	return c
}

// bucketMaintainer background task to maintain position value
func (c *Counter) bucketMaintainer() {
	for {
		if time.Now().Sub(c.lastTimestamp) > c.bucketTimePeriodLength {
			atomic.StoreInt32(&c.position, (c.position+1)%c.windowSize)
			//c.mutex.Lock()
			atomic.StoreInt64(&c.buckets[c.position], 0)
			//c.mutex.Unlock()
			c.lastTimestamp = time.Now()
		} else {
			time.Sleep(c.bucketTimePeriodLength / 10)
		}
	}
}

// Count count one
func (c *Counter) CountOne() {
	atomic.AddInt64(&c.buckets[c.position], 1)
}

// CountFor count for some value
func (c *Counter) CountFor(value int64) {
	atomic.AddInt64(&c.buckets[c.position], value)
}

// GetTotal get current counter value
func (c *Counter) GetTotal() int64 {
	var res int64 = 0
	//fmt.Print(c.lastTimestamp)
	//c.mutex.Lock()
	for i := 0; i < (int)(c.windowSize); i++ {
		fmt.Print(" ")
		item := atomic.LoadInt64(&c.buckets[i])
		fmt.Print(item)
		res += item
	}
	//c.mutex.Unlock()
	fmt.Println()
	return res
}
