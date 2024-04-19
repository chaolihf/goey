//go:build cocoa || (darwin && !gtk)
// +build cocoa darwin,!gtk

package dialog

import (
	"github.com/chaolihf/goey/internal/cocoa"
)

// Owner holds a pointer to the owning window.
// This type varies between platforms.
type Owner struct {
	Window *cocoa.Window
}
