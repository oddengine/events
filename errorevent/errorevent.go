package errorevent

import (
	"fmt"

	"github.com/oddcancer/events"
	Event "github.com/oddcancer/events/event"
)

// ErrorEvent types.
const (
	ERROR = "error"
)

// ErrorEvent dispatched when an error causes an asynchronous operation to fail.
type ErrorEvent struct {
	Event.Event
	Name    string
	Message error
}

// Init this class.
func (me *ErrorEvent) Init(event string, name string, message error) *ErrorEvent {
	me.Event.Init(event)
	me.Name = name
	me.Message = message
	return me
}

// Clone an instance of an ErrorEvent subclass.
func (me *ErrorEvent) Clone() events.IEvent {
	return New(me.Type(), me.Target(), me.Name, me.Message)
}

// String returns a string containing all the properties of the ErrorEvent object.
func (me *ErrorEvent) String() string {
	return fmt.Sprintf("[ErrorEvent type=%s name=%s message=%v]", me.Type(), me.Name, me.Message)
}

// New creates a new ErrorEvent object.
func New(typ string, target events.IEventTarget, name string, message error) *ErrorEvent {
	e := new(ErrorEvent).Init(typ, name, message)
	e.SetTarget(target)
	e.SetCurrentTarget(target)
	return e
}
