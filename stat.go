package main

import (
	"fmt"
	"sync"
	"time"
)

type Stat struct {
	sync.RWMutex
	LastTime    int64
	LastItemCnt int64
	LastPkgCnt  int64

	ThisTime    int64
	ThisItemCnt int64
	ThisPkgCnt  int64
}

func NewStat() *Stat {
	return &Stat{
		LastTime:    time.Now().UnixNano(),
		LastPkgCnt:  0,
		LastItemCnt: 0,

		ThisTime:    time.Now().UnixNano(),
		ThisPkgCnt:  0,
		ThisItemCnt: 0,
	}
}

func (s *Stat) Incr(itemCnt int64) {
	s.Lock()
	s.ThisTime = time.Now().UnixNano()
	s.ThisItemCnt += itemCnt
	s.ThisPkgCnt += 1
	s.Unlock()
}

func (s *Stat) Stats() {
	s.Lock()
	now := time.Now()
	duringTime := now.UnixNano() - s.LastTime
	duringPkgCnt := s.ThisPkgCnt - s.LastPkgCnt
	duringItemCnt := s.ThisItemCnt - s.LastItemCnt
	itemRate := float64(0)
	if duringTime != 0 {
		itemRate = float64(duringItemCnt) * 1e9 / float64(duringTime)
	}

	s.LastTime = now.UnixNano()
	s.LastPkgCnt = s.ThisPkgCnt
	s.LastItemCnt = s.ThisItemCnt

	fmt.Printf("NowTime: %d, LastSend: %d, PkgCnt: %d, ItemCnt: %d, ItemRate: %0.2f/s\n",
		now.Unix(),
		s.ThisTime/1e9,
		duringPkgCnt,
		duringItemCnt,
		itemRate,
	)
	s.Unlock()
}
