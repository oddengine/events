package events

import (
	"reflect"
)

// EventListenerOptions specifies characteristics about the event listener.
type EventListenerOptions struct {
	Once bool
}

// EventListener holds the event handler.
type EventListener struct {
	handler interface{}
	options EventListenerOptions
}

// Init this class.
func (me *EventListener) Init(handler interface{}, options ...EventListenerOptions) *EventListener {
	me.handler = handler
	if len(options) > 0 {
		me.options = options[0]
	}
	return me
}

func (me *EventListener) Invoke(e IEvent) {
	value := reflect.ValueOf(e)
	reflect.ValueOf(me.handler).Call([]reflect.Value{value})
}

func (me *EventListener) Matches(listener *EventListener) bool {
	return listener == me
}

// NewEventListener returns new EventListener.
func NewEventListener(handler interface{}, options ...EventListenerOptions) *EventListener {
	return new(EventListener).Init(handler, options...)
}
