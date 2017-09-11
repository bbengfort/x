/*
Package events provides an event dispatcher and registration framework.
*/
package events

import (
	"reflect"
	"sync"
)

// Some standard event types
const (
	UnknownEvent Type = iota
	TimeoutEvent
)

// Names of event types
var eventTypeStrings = [...]string{
	"unknown", "timeout",
}

//===========================================================================
// Event Types
//===========================================================================

// Type is an enumeration of the kind of events that can occur.
type Type uint16

// String returns the name of event types
func (t Type) String() string {
	if int(t) < len(eventTypeStrings) {
		return eventTypeStrings[t]
	}
	return "custom"
}

// Callback is a function that can receive events.
type Callback func(Event) error

//===========================================================================
// Event Dispatcher Definition and Methods
//===========================================================================

// Dispatcher objects can register callbacks for specific events, then when
// those events occur, dispatch them to all callback functions.
type Dispatcher struct {
	sync.RWMutex
	source    interface{}
	callbacks map[Type][]Callback
}

// Init a dispatcher with the source, creating the callbacks map.
func (d *Dispatcher) Init(source interface{}) {
	d.source = source
	d.callbacks = make(map[Type][]Callback)
}

// Register a callback function for the specified event type.
func (d *Dispatcher) Register(etype Type, callback Callback) {
	if callback == nil {
		return
	}

	d.Lock()
	defer d.Unlock()
	d.callbacks[etype] = append(d.callbacks[etype], callback)
}

// Remove a callback function for the specified event type.
func (d *Dispatcher) Remove(etype Type, callback Callback) {
	d.Lock()
	defer d.Unlock()

	// Grab a reference to the function pointer
	ptr := reflect.ValueOf(callback).Pointer()

	// Find callback by pointer and remove it
	callbacks := d.callbacks[etype]
	for idx, cb := range callbacks {
		if reflect.ValueOf(cb).Pointer() == ptr {
			d.callbacks[etype] = append(callbacks[:idx], callbacks[idx+1:]...)
		}
	}
}

// Dispatch an event, ensuring that the event is properly formatted.
// Currently this method simply warns if there is an error.
// TODO: return list of errors or do better error handling.
func (d *Dispatcher) Dispatch(etype Type, value interface{}) error {
	d.RLock()
	defer d.RUnlock()
	return d.dispatch(etype, value)
}

// Internal dispatch event that is not thread-safe (surrounded by locks).
func (d *Dispatcher) dispatch(etype Type, value interface{}) error {
	// Create the event
	e := &event{
		etype:  etype,
		source: d.source,
		value:  value,
	}

	// Dispatch the event to all callbacks
	for _, cb := range d.callbacks[etype] {
		if err := cb(e); err != nil {
			return err
		}
	}

	return nil
}

//===========================================================================
// Event Definition and Methods
//===========================================================================

// Event represents actions that occur during consensus. Listeners can
// register callbacks with event handlers for specific event types.
type Event interface {
	Type() Type
	Source() interface{}
	Value() interface{}
}

// event is an internal implementation of the Event interface.
type event struct {
	etype  Type
	source interface{}
	value  interface{}
}

// Type returns the event type.
func (e *event) Type() Type {
	return e.etype
}

// Source returns the entity that dispatched the event.
func (e *event) Source() interface{} {
	return e.source
}

// Value returns the current value associated with teh event.
func (e *event) Value() interface{} {
	return e.value
}
