package timerevent

import (
	"fmt"

	"github.com/oddcancer/events"
	Event "github.com/oddcancer/events/event"
)

// TimerEvent types.
const (
	TIMER    = "timer"
	COMPLETE = "timer-complete"
)

// TimerEvent dispatched whenever the Timer object reaches the interval specified by the Timer.delay property.
type TimerEvent struct {
	Event.Event
}

// Init this class.
func (me *TimerEvent) Init(event string) *TimerEvent {
	me.Event.Init(event)
	return me
}

// Clone an instance of an TimerEvent subclass.
func (me *TimerEvent) Clone() events.IEvent {
	return New(me.Type(), me.Target())
}

// String returns a string containing all the properties of the TimerEvent object.
func (me *TimerEvent) String() string {
	return fmt.Sprintf("[TimerEvent type=%s]", me.Type())
}

// New creates a new TimerEvent object.
func New(event string, target events.IEventTarget) *TimerEvent {
	e := new(TimerEvent).Init(event)
	e.SetTarget(target)
	e.SetCurrentTarget(target)
	return e
}
