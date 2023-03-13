package netstatusevent

import (
	"fmt"

	"github.com/oddengine/events"
	Event "github.com/oddengine/events/event"
)

// NetStatusEvent types.
const (
	NET_STATUS = "netStatus"
)

// NetStatusEvent dispatched when a net status event occurred.
type NetStatusEvent struct {
	Event.Event
	Level       string
	Code        string
	Description string
	Info        map[string]interface{}
}

// Init this class
func (me *NetStatusEvent) Init(event string, level string, code string, description string, info map[string]interface{}) *NetStatusEvent {
	me.Event.Init(event)
	me.Level = level
	me.Code = code
	me.Description = description
	me.Info = info
	return me
}

// Clone an instance of an NetStatusEvent subclass.
func (me *NetStatusEvent) Clone() events.IEvent {
	return New(me.Type(), me.Target(), me.Level, me.Code, me.Description, me.Info)
}

// String returns a string containing all the properties of the NetStatusEvent object.
func (me *NetStatusEvent) String() string {
	return fmt.Sprintf("[NetStatusEvent type=%s level=%s code=%s description=%s]", me.Type(), me.Level, me.Code, me.Description)
}

// New creates a new NetStatusEvent object.
func New(event string, target events.IEventTarget, level string, code string, description string, info map[string]interface{}) *NetStatusEvent {
	e := new(NetStatusEvent).Init(event, level, code, description, info)
	e.SetTarget(target)
	e.SetCurrentTarget(target)
	return e
}
