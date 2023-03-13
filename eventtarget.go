package events

import (
	"fmt"
	"runtime/debug"

	"github.com/oddengine/events/reentrant"
	"github.com/oddengine/log"
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
	mtx       reentrant.Mutex
	logger    log.ILogger
	listeners map[string]*MappableEventListenerCollection
	recursion int32
}

// Init this class.
func (me *EventTarget) Init(logger log.ILogger) *EventTarget {
	me.logger = logger
	me.listeners = make(map[string]*MappableEventListenerCollection)
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

	m := me.listeners[event]
	if m == nil {
		m = new(MappableEventListenerCollection).Init()
		me.listeners[event] = m
	}

	me.logger.Debugf(1, "Adding event listener: type=%s, listener=%p", event, listener)
	m.Add(listener)
}

// RemoveEventListener removes an event listener from the EventTarget object.
func (me *EventTarget) RemoveEventListener(event string, listener *EventListener) {
	if event == "" || listener == nil {
		me.logger.Debugf(1, "Event type or listener not present: type=%s, listener=%p", event, listener)
		return
	}

	me.mtx.Lock()
	defer me.mtx.Unlock()

	m := me.listeners[event]
	if m == nil {
		me.logger.Debugf(0, "No listener[s] found: type=%s", event)
		return
	}

	me.logger.Debugf(1, "Removing event listener: type=%s, listener=%p", event, listener)
	m.Remove(listener, me.recursion == 0)
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

	// Get the typed listener collection.
	m := me.listeners[e.Type()]
	if m == nil {
		me.logger.Debugf(0, "No listener[s] found: type=%s", e.Type())
		return NotCanceled
	}

	if me.recursion == 1 {
		defer m.RemoveEventually()
	}

	// Loop to invoke the handlers.
	for element := m.List.Front(); element != nil; element = m.Next(element) {
		listener := element.Value.(*EventListener)
		listener.Invoke(e)

		if listener.options.Once {
			me.logger.Debugf(1, "Removing event listener: type=%s, listener=%p", e.Type(), listener)
			m.Remove(listener, me.recursion == 0)
		}
		if e.PropagationStopped() {
			me.logger.Debugf(1, "Propagation stopped: type=%s", e.Type())
			return CanceledByEventHandler
		}
	}
	return NotCanceled
}
