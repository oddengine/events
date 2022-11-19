package events

import (
	"fmt"
	"runtime/debug"
	"unsafe"

	"github.com/oddcancer/events/reentrant"
	"github.com/oddcancer/log"
)

// Static constants.
const (
	MAX_RECURSION int32 = 8
)

// EventTarget is the base class for all classes that dispatch events.
// It uses array instead of list, which causes the add/remove method expensive.
// However, it is possible to clone the listener group fast while triggering an event.
// And, the frequency of triggering event is much higher than that of add/remove.
type EventTarget struct {
	mtx reentrant.Mutex

	logger    log.ILogger
	listeners map[string]map[uintptr]*EventListener
	recursion int32
}

// Init this class.
func (me *EventTarget) Init(logger log.ILogger) *EventTarget {
	me.logger = logger
	me.listeners = make(map[string]map[uintptr]*EventListener)
	return me
}

// AddEventListener registers an event listener object with an EventTarget object so that the listener receives notification of an event.
func (me *EventTarget) AddEventListener(event string, listener *EventListener) {
	if event == "" || listener == nil {
		me.logger.Debugf(1, "Event type or listener not present: type=%s, listener=%p", event, listener)
		return
	}

	me.mtx.Lock()
	defer me.mtx.Unlock()
	me.addEventListener(event, listener)
}

func (me *EventTarget) addEventListener(event string, listener *EventListener) {
	m := me.listeners[event]
	if m == nil {
		m = make(map[uintptr]*EventListener)
		me.listeners[event] = m
	}

	me.logger.Debugf(1, "Adding event listener: type=%s, listener=%p", event, listener)
	m[uintptr(unsafe.Pointer(listener))] = listener
}

// RemoveEventListener removes an event listener from the EventTarget object.
func (me *EventTarget) RemoveEventListener(event string, listener *EventListener) {
	if event == "" || listener == nil {
		me.logger.Debugf(1, "Event type or listener not present: type=%s, listener=%p", event, listener)
		return
	}

	me.mtx.Lock()
	defer me.mtx.Unlock()
	me.removeEventListener(event, listener)
}

func (me *EventTarget) removeEventListener(event string, listener *EventListener) {
	m := me.listeners[event]
	if m == nil {
		me.logger.Debugf(0, "No listener[s] found: type=%s", event)
		return
	}

	me.logger.Debugf(1, "Removing event listener: type=%s, listener=%p", event, listener)
	delete(m, uintptr(unsafe.Pointer(listener)))
}

// DispatchEvent dispatches an event into the event flow.
func (me *EventTarget) DispatchEvent(e IEvent) EventResult {
	defer func() {
		if err := recover(); err != nil {
			me.logger.Errorf("Failed to handle event: type=%s, %v", e.Type(), err)
			debug.PrintStack()
		}
	}()

	me.mtx.Lock()
	defer me.mtx.Unlock()

	e.SetCurrentTarget(me)
	me.logger.Debugf(0, "Dispatching event: %s", e.Type())

	// Check recursion.
	me.recursion++
	defer func() {
		me.recursion--
	}()

	if MAX_RECURSION > 0 && me.recursion > MAX_RECURSION {
		panic(fmt.Sprintf("max recursion reached: %d", me.recursion))
	}

	// Make a copy of the typed listener map.
	m := me.listeners[e.Type()]
	if m == nil {
		me.logger.Debugf(0, "No listener[s] found: type=%s", e.Type())
		return NotCanceled
	}
	c := make(map[uintptr]*EventListener, len(m))
	for i := range m {
		c[i] = m[i]
	}

	// Loop to invoke the handlers.
	for _, listener := range c {
		listener.Invoke(e)

		if listener.options.Once {
			me.logger.Debugf(1, "Removing event listener: type=%s, listener=%p", e.Type(), listener)
			delete(m, uintptr(unsafe.Pointer(listener)))
		}
		if e.PropagationStopped() {
			me.logger.Debugf(1, "Propagation stopped: type=%s", e.Type())
			return CanceledByEventHandler
		}
	}
	return NotCanceled
}
