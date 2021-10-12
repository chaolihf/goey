package base_test

import (
	"fmt"
	"image"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func ExampleFromPixels() {
	// Most code should not need to worry about setting the DPI.  Windows will
	// ensure that the DPI is set.
	base.DPI = image.Point{96, 96}

	size := base.FromPixels(48, 96+96)
	fmt.Printf("The size is %s.\n", size)

	// Output:
	// The size is (48:00x192:00).
}

func TestSize(t *testing.T) {
	cases := []struct {
		in     base.Size
		isZero bool
		out    string
	}{
		{base.Size{}, true, "(0:00x0:00)"},
		{base.Size{1, 2}, false, "(0:01x0:02)"},
		{base.Size{1 * base.DIP, 2 * base.DIP}, false, "(1:00x2:00)"},
	}

	for i, v := range cases {
		if out := v.in.IsZero(); out != v.isZero {
			t.Errorf("Case %d:  Failed predicate IsZero, got %v, want %v", i, out, v.isZero)
		}
		if out := v.in.String(); out != v.out {
			t.Errorf("Case %d:  Failed method String, got %v, want %v", i, out, v.out)
		}
	}
}
