package events

import (
	"container/list"
	"reflect"
	"unsafe"
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

// Invoke calls handler with the IEvent.
func (me *EventListener) Invoke(e IEvent) {
	value := reflect.ValueOf(e)
	reflect.ValueOf(me.handler).Call([]reflect.Value{value})
}

// Matches returns whether or not it is equal to the argument.
func (me *EventListener) Matches(listener *EventListener) bool {
	return listener == me
}

// MappableEventListenerCollection is a mappable listener collection.
type MappableEventListenerCollection struct {
	List     list.List
	elements map[uintptr]*list.Element
	removed  map[uintptr]*list.Element
}

// Init this class.
func (me *MappableEventListenerCollection) Init() *MappableEventListenerCollection {
	me.List.Init()
	me.elements = make(map[uintptr]*list.Element)
	me.removed = make(map[uintptr]*list.Element)
	return me
}

// Add adds the listener into collection.
func (me *MappableEventListenerCollection) Add(listener *EventListener) {
	key := uintptr(unsafe.Pointer(listener))
	if _, ok := me.elements[key]; !ok {
		me.elements[key] = me.List.PushBack(listener)
		delete(me.removed, key)
	}
}

// Remove removes the listener from collection.
func (me *MappableEventListenerCollection) Remove(listener *EventListener, immediately bool) {
	key := uintptr(unsafe.Pointer(listener))
	if e, ok := me.elements[key]; ok {
		if !immediately {
			me.removed[key] = e
			return
		}
		me.List.Remove(e)
		delete(me.elements, key)
		delete(me.removed, key)
	}
}

// Next returns the next element of listener to fire an event.
func (me *MappableEventListenerCollection) Next(element *list.Element) *list.Element {
	for element = element.Next(); element != nil; element = element.Next() {
		listener := element.Value.(*EventListener)
		key := uintptr(unsafe.Pointer(listener))
		if _, ok := me.removed[key]; !ok {
			return element
		}
	}
	return nil
}

// RemoveEventually removes the previoursly removed elements, eventually.
func (me *MappableEventListenerCollection) RemoveEventually() {
	for key, e := range me.removed {
		me.List.Remove(e)
		delete(me.elements, key)
		delete(me.removed, key)
	}
}

// Len returns the number of elements of collection.
func (me *MappableEventListenerCollection) Len() int {
	return me.List.Len()
}

// NewEventListener returns new EventListener.
func NewEventListener(handler interface{}, options ...EventListenerOptions) *EventListener {
	return new(EventListener).Init(handler, options...)
}
