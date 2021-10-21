//go:build go1.12
// +build go1.12

package base

import (
	"syscall/js"
)

const (
	// PLATFORM specifies the GUI toolkit being used.
	PLATFORM = "js"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkContainer when targeting GTK.
//
// Unless developing new widgets, users should not need to use this type.
//
// Any methods on this type will be platform specific.
type Control struct {
	Handle js.Value
}

// NativeElement contains platform-specific methods that all widgets
// must support on JS/WASM.
type NativeElement interface{}
