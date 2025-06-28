package main

import (
	"fmt"
	"time"
)

type Bucket struct {
	totalCnt   int
	successCnt int
}

type RollingWindow struct {
	size    int
	buckets []Bucket
	offset  int

	interval time.Duration
	lastTime time.Time
}

func NewRollingWindow(size int, interval time.Duration) *RollingWindow {
	return &RollingWindow{
		size:     size,
		buckets:  make([]Bucket, size),
		offset:   -1,
		interval: interval,
		lastTime: time.Now(),
	}
}

// 距离上次最后一次更新window过去了n个span，所以如果这n个span中如果有数据，一定是过期的数据
func (rw *RollingWindow) span() int {
	return min(int(time.Since(rw.lastTime)/rw.interval), rw.size)
}

func (rw *RollingWindow) Add(totalCnt, successCnt int) {
	nSpan := rw.span()
	fmt.Println(nSpan)

	for i := 0; i < nSpan; i++ {
		rw.buckets[rw.offset+1+i%rw.size].reset()
	}

	rw.buckets[rw.offset+nSpan].totalCnt = totalCnt
	rw.buckets[rw.offset+nSpan].successCnt = successCnt

	rw.offset += nSpan
	rw.lastTime = time.Now()
}

func (rw *RollingWindow) Reduce() (int, int) {
	var totalCnt, successCnt int

	for _, bucket := range rw.buckets {
		totalCnt += bucket.totalCnt
		successCnt += bucket.successCnt
	}

	return totalCnt, successCnt
}

func (b *Bucket) reset() {
	b.successCnt = 0
	b.totalCnt = 0
}
