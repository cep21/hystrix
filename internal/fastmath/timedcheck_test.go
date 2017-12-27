package fastmath

import (
	"sync"
	"testing"
	"time"
)

func TestTimedCheck_Empty(t *testing.T) {
	x := TimedCheck{}
	now := time.Now()
	if !x.Check(now) {
		t.Error("First check should pass on empty object")
	}
}

func TestTimedCheck_Check(t *testing.T) {
	x := TimedCheck{}
	x.SetSleepDuration(time.Second)
	now := time.Now()
	x.SleepStart(now)
	if x.Check(now) {
		t.Error("Should not check at first")
	}
	if x.Check(now.Add(time.Millisecond * 999)) {
		t.Error("Should not check close to end")
	}
	if !x.Check(now.Add(time.Second)) {
		t.Error("Should check at barrier")
	}
	if x.Check(now.Add(time.Second)) {
		t.Error("Should only check once")
	}
	if x.Check(now.Add(time.Second + time.Millisecond)) {
		t.Error("Should only double check")
	}
	if !x.Check(now.Add(time.Second * 2)) {
		t.Error("Should check again at 2 sec")
	}
}

func TestTimedCheckRaces(t *testing.T) {
	x := TimedCheck{}
	x.SetSleepDuration(time.Nanosecond * 100)
	endTime := time.Now().Add(time.Millisecond * 50)
	wg := sync.WaitGroup{}
	doTillTime(endTime, &wg, func() {
		x.Check(time.Now())
	})
	doTillTime(endTime, &wg, func() {
		x.SetEventCountToAllow(2)
	})
	doTillTime(endTime, &wg, func() {
		x.SetSleepDuration(time.Millisecond * 100)
	})
	wg.Wait()
}