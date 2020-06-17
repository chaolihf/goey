package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestEmptyMount(t *testing.T) {
	testMountWidgets(t,
		&Empty{},
		&Empty{},
		&Empty{},
	)
}

func TestEmptyClose(t *testing.T) {
	testCloseWidgets(t,
		&Empty{},
		&Empty{},
	)
}

func TestEmptyUpdate(t *testing.T) {
	testUpdateWidgets(t, []base.Widget{
		&Empty{},
		&Empty{},
	}, []base.Widget{
		&Empty{},
		&Empty{},
	})
}
