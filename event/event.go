package event

import (
	"fmt"

	"github.com/oddcancer/events"
)

// Event types.
const (
	ACTIVATE   = "activate"
	ADDED      = "added"
	CANCEL     = "cancel"
	CHANGE     = "change"
	CLEAR      = "clear"
	CLOSE      = "close"
	COMPLETE   = "complete"
	CONNECT    = "connect"
	DEACTIVATE = "deactivate"
	IDLE       = "idle"
	INIT       = "init"
	OPEN       = "open"
	RELEASE    = "release"
	REMOVED    = "removed"
)

// Event is used as the base class for the creation of Event objects, which are passed as parameters to event listeners when an event occurs.
type Event struct {
	event              string
	target             events.IEventTarget
	currentTarget      events.IEventTarget
	propagationStopped bool
}

// Init this class.
func (me *Event) Init(event string) *Event {
	me.event = event
	me.propagationStopped = false
	return me
}

// SetType sets the event type.
func (me *Event) SetType(event string) {
	me.event = event
}

// Type gets the event type.
func (me *Event) Type() string {
	return me.event
}

// SetTarget sets the source target.
func (me *Event) SetTarget(target events.IEventTarget) {
	me.target = target
}

// Target gets the source target.
func (me *Event) Target() events.IEventTarget {
	return me.target
}

// SetCurrentTarget sets the current target.
func (me *Event) SetCurrentTarget(target events.IEventTarget) {
	me.currentTarget = target
}

// CurrentTarget gets the current target.
func (me *Event) CurrentTarget() events.IEventTarget {
	return me.currentTarget
}

// StopPropagation stops propagation.
func (me *Event) StopPropagation() {
	me.propagationStopped = true
}

// PropagationStopped returns whether the propagation is stopped.
func (me *Event) PropagationStopped() bool {
	return me.propagationStopped
}

// Clone an instance of an Event subclass.
func (me *Event) Clone() events.IEvent {
	return New(me.Type(), me.Target())
}

// String returns a string containing all the properties of the Event object.
func (me *Event) String() string {
	return fmt.Sprintf("[Event type=%s]", me.Type())
}

// New creates a new Event object.
func New(event string, target events.IEventTarget) *Event {
	e := new(Event).Init(event)
	e.SetTarget(target)
	e.SetCurrentTarget(target)
	return e
}
