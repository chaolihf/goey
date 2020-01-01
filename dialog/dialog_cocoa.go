// +build cocoa darwin,!gtk

package dialog

import (
	"bitbucket.org/rj/goey/internal/cocoa"
)

type dialogImpl struct {
	parent *cocoa.Window
}
