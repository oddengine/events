package events

type EventResult int

const (
	// Event was not canceled by event handler or default event handler.
	NotCanceled EventResult = iota
	// Event was canceled by event handler; i.e. a script handler calling StopPropagation.
	CanceledByEventHandler
	// Event was canceled by the default event handler; i.e. executing the default action.
	// This result should be used sparingly as it deviates from the Event Dispatch model.
	// Default event handlers really shouldn't be invoked inside of dispatch.
	CanceledByDefaultEventHandler
)

// IEvent defines basic event methods.
type IEvent interface {
	SetType(event string)
	Type() string
	SetTarget(target IEventTarget)
	Target() IEventTarget
	SetCurrentTarget(target IEventTarget)
	CurrentTarget() IEventTarget
	StopPropagation()
	PropagationStopped() bool
	Clone() IEvent
	String() string
}

// IEventTarget objects allow us to add and remove an event listeners of a specific event type.
// Each IEventTarget object also represents the target to which an event is dispatched when something has occurred.
type IEventTarget interface {
	AddEventListener(event string, listener *EventListener)
	RemoveEventListener(event string, listener *EventListener)
	DispatchEvent(e IEvent) EventResult
}
