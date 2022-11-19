package reentrant

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Static constants.
const (
	MAX_RECURSION int32 = 0
)

// A Mutex is a reentrant mutual exclusion lock.
//
// The zero value for a Mutex is an unlocked mutex. And a Mutex must not be copied
// after first use.
type Mutex struct {
	sync.Mutex

	goid      int64
	recursion int32
}

// Lock locks this Mutex.
//
// If the lock is already in use by the calling goroutine, it only checks the
// recursion. Otherwise, the calling goroutine blocks until the mutex is available.
func (me *Mutex) Lock() {
	self := GetCurrentGoroutineID()
	if atomic.LoadInt64(&me.goid) == self {
		me.recursion++
		if MAX_RECURSION > 0 && me.recursion > MAX_RECURSION {
			panic(fmt.Sprintf("max recursion reached: %d", me.recursion))
		}
		return
	}

	me.Mutex.Lock()
	atomic.StoreInt64(&me.goid, self)
	me.recursion = 1
}

// Unlock unlocks this Mutex.
//
// It is a run-time error if this is not locked on entry to Unlock. A locked Mutex
// is not associated with a particular goroutine. It is allowed for one goroutine
// to lock a Mutex and then arrange for another goroutine to unlock it.
func (me *Mutex) Unlock() {
	me.recursion--
	if me.recursion == 0 {
		atomic.StoreInt64(&me.goid, 0)
		me.Mutex.Unlock()
	}
}
