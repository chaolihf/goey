package main

import (
	"time"

	"bitbucket.org/rj/goey/loop"
)

// Debounce will return a new function that debounces the callback.
func Debounce(cb func(string)) func(string) {
	d := &debouncer{
		cb:       cb,
		duration: 500 * time.Millisecond,
	}

	return d.OnEvent
}

// debouncer managers the internal state to debounce a callback.
type debouncer struct {
	duration time.Duration
	timer    *time.Timer
	cb       func(string)
	value    string
}

// OnEvent should be called instead of the original callback to debounce an
// event.  OnEvent will call the original callback once the events settles.
func (d *debouncer) OnEvent(s string) {
	// This function should only be called from the GUI thread.
	// No locking should be necessary.
	d.value = s

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(
		d.duration,
		d.emitEvent,
	)
}

// emitEvent will call the original callback with the last value received for
// the event.
func (d *debouncer) emitEvent() {
	loop.Do(func() error {
		d.cb(d.value)
		return nil
	})
}
