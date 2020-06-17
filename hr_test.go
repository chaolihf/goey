package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestHRMount(t *testing.T) {
	testMountWidgets(t,
		&HR{},
		&HR{},
		&HR{},
	)
}

func TestHRClose(t *testing.T) {
	testCloseWidgets(t,
		&HR{},
		&HR{},
	)
}

func TestHRUpdate(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&HR{},
		&HR{},
	}, []base.Widget{
		&HR{},
		&HR{},
	})
}

func TestHRLayout(t *testing.T) {
	testLayoutWidget(t, &HR{})
}

func TestHRMinSize(t *testing.T) {
	testMinSizeWidget(t, &HR{})
}
