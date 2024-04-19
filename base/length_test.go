package base_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/chaolihf/goey/base"
)

func ExampleLength() {
	// Since there are 96 device-independent pixels per inch, and 6 picas
	// per inch, the following two lengths should be equal.
	length1 := 96 * base.DIP
	length2 := 6 * base.PC

	if length1 == length2 {
		fmt.Printf("All is OK with the world.")
	} else {
		fmt.Printf("This should not happen, unless there is a rounding error.")
	}

	// Output:
	// All is OK with the world.
}

func ExampleLength_Scale() {
	// There are 96 DIP in an inch, and 6 pica in a inch, so the following
	// should work.

	if length := (1 * base.DIP).Scale(96, 6); length == (1 * base.PC) {
		fmt.Printf("The ratio of pica to DIP is 96 to 6.")
	}

	// Output:
	// The ratio of pica to DIP is 96 to 6.
}

func ExampleLength_String() {
	fmt.Printf("Converting:  1pt is equal to %sdip\n", 1*base.PT)
	fmt.Printf("Converting:  1pt is equal to %1.2fdip\n", (1 * base.PT).DIP())
	fmt.Printf("Converting:  1pc is equal to %1.1fdip\n", (1 * base.PC).DIP())

	// Output:
	// Converting:  1pt is equal to 1:21dip
	// Converting:  1pt is equal to 1.33dip
	// Converting:  1pc is equal to 16.0dip
}

func ExampleGuardInf() {
	width1 := base.Inf
	width2 := 1 * base.Inch

	fmt.Printf("Half of Inf is %s.\n", base.GuardInf(width1, width1/2))
	fmt.Printf("Half of %s is %s.\n", width2, base.GuardInf(width2, width2/2))

	// Output:
	// Half of Inf is Inf.
	// Half of 96:00 is 48:00.
}

func ExampleRectangle() {
	r := base.Rectangle{
		Min: base.Point{10 * base.DIP, 20 * base.DIP},
		Max: base.Point{90 * base.DIP, 80 * base.DIP},
	}

	fmt.Printf("Rectangle %s has dimensions %.0fdip by %.0fdip.",
		r, r.Dx().DIP(), r.Dy().DIP(),
	)

	// Output:
	// Rectangle (10:00,20:00)-(90:00,80:00) has dimensions 80dip by 60dip.
}

func ExampleRectangle_Add() {
	r := base.Rectangle{
		Min: base.Point{10 * base.DIP, 20 * base.DIP},
		Max: base.Point{90 * base.DIP, 80 * base.DIP},
	}
	v := base.Point{5 * base.DIP, 5 * base.DIP}

	fmt.Printf("Rectangle %s, moved by %s,\n", r, v)
	fmt.Printf("---- %s", r.Add(v))

	// Output:
	// Rectangle (10:00,20:00)-(90:00,80:00), moved by (5:00,5:00),
	// ---- (15:00,25:00)-(95:00,85:00)
}

func ExampleRectangle_Pixels() {
	// The following line is for the example only, and should not appear in
	// user code, as the platform-specific code should update the DPI based
	// on the system.  However, for the purpose of this example, set a known
	// DPI.
	base.DPI = image.Point{2 * 96, 2 * 96}

	// Construct an example rectangle.
	r := base.Rectangle{
		Min: base.Point{10 * base.DIP, 20 * base.DIP},
		Max: base.Point{90 * base.DIP, 80 * base.DIP},
	}
	rpx := r.Pixels()

	fmt.Printf("Rectangle %s when translated to pixels is %s.", r, rpx)

	// Output:
	// Rectangle (10:00,20:00)-(90:00,80:00) when translated to pixels is (20,40)-(180,160).
}

func TestFromPixels(t *testing.T) {
	cases := []struct {
		dpix, dpiy       int
		pixels           int
		lengthx, lengthy base.Length
	}{
		// Standard DPI tests
		{96, 96, 2, 2 * base.DIP, 2 * base.DIP},
		{96, 96, 3, 3 * base.DIP, 3 * base.DIP},
		{96, 96 * 3 / 2, 2, 2 * base.DIP, 2 * base.DIP * 2 / 3},
		{96, 96 * 3 / 2, 3, 3 * base.DIP, 2 * base.DIP},
		// 300 DPI tests
		{300, 300, 1, base.DIP * 96 / 300, base.DIP * 96 / 300},
		{300, 300, 4096, 4096 * base.DIP * 96 / 300, 4096 * base.DIP * 96 / 300},
		// Very high DPI stress test.
		{1024, 1024, 1, base.DIP * 96 / 1024, base.DIP * 96 / 1024},
		{1024, 1024, 1024 * 16, 16 * base.Inch, 16 * base.Inch},
	}

	for i, v := range cases {
		base.DPI = image.Point{v.dpix, v.dpiy}
		if got := base.FromPixelsX(v.pixels); got != v.lengthx {
			t.Errorf("Unexpected conversion in FromPixelsX on case %d, got %v, want %v", i, got, v.lengthx)
		}
		if got := base.FromPixelsY(v.pixels); got != v.lengthy {
			t.Errorf("Unexpected conversion in FromPixelsY on case %d, got %v, want %v", i, got, v.lengthy)
		}
	}
}

func TestLength(t *testing.T) {
	if rt := (1 * base.DIP).DIP(); rt != 1 {
		t.Errorf("Unexpected round-trip for Length, %v =/= %v", rt, 1)
	}
	if rt := (1 * base.PT).PT(); rt != 1 {
		t.Errorf("Unexpected round-trip for PT,  %v =/= %v", rt, 1)
	}
	if rt := (1 * base.PC).PC(); rt != 1 {
		t.Errorf("Unexpected round-trip for PC,  %v =/= %v", rt, 1)
	}
	if rt := (1 * base.Inch).Inch(); rt != 1 {
		t.Errorf("Unexpected round-trip for inch,  %v =/= %v", rt, 1)
	}
	if rt := (1 * base.PT) * (1 << 6) / (1 * base.DIP); rt != 96*(1<<6)/72 {
		t.Errorf("Unexpected ratio between DIP and PT, %v =/= %v", rt, 96*(1<<6)/72)
	}
	if rt := (1 * base.PC) * (1 << 6) / (1 * base.DIP); rt != 96*(1<<6)/6 {
		t.Errorf("Unexpected ratio between DIP and PC, %v =/= %v", rt, 96*(1<<6)/72)
	}
	if rt := (1 * base.Inch) * (1 << 6) / (1 * base.DIP); rt != 96*(1<<6) {
		t.Errorf("Unexpected ratio between DIP and inch, %v =/= %v", rt, 96*(1<<6))
	}
}

func TestLength_Clamp(t *testing.T) {
	const DIP = base.DIP

	cases := []struct {
		in       base.Length
		min, max base.Length
		out      base.Length
	}{
		{10 * DIP, 0 * DIP, 20 * DIP, 10 * DIP},
		{30 * DIP, 0 * DIP, 20 * DIP, 20 * DIP},
		{-10 * DIP, 0 * DIP, 20 * DIP, 0 * DIP},
		{10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{30 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{-10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{10 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
		{30 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
		{-10 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
	}

	for i, v := range cases {
		if out := v.in.Clamp(v.min, v.max); out != v.out {
			t.Errorf("Error in case %d, want %s, got %s", i, v.out, out)
		}
	}
}

func TestPoint(t *testing.T) {
	const DIP = base.DIP

	cases := []struct {
		a, b base.Point
		add  base.Point
		sub  base.Point
	}{
		{base.Point{}, base.Point{1 * DIP, 2 * DIP}, base.Point{1 * DIP, 2 * DIP}, base.Point{-1 * DIP, -2 * DIP}},
		{base.Point{1 * DIP, 2 * DIP}, base.Point{}, base.Point{1 * DIP, 2 * DIP}, base.Point{1 * DIP, 2 * DIP}},
		{base.Point{3 * DIP, 5 * DIP}, base.Point{7 * DIP, 11 * DIP}, base.Point{10 * DIP, 16 * DIP}, base.Point{-4 * DIP, -6 * DIP}},
	}

	for i, v := range cases {
		if out := v.a.Add(v.b); out != v.add {
			t.Errorf("Error in case %d, want %s, got %s", i, v.add, out)
		}
		if out := v.a.Sub(v.b); out != v.sub {
			t.Errorf("Error in case %d, want %s, got %s", i, v.add, out)
		}
	}
}

func TestLength_Pixels(t *testing.T) {
	const DIP = base.DIP

	cases := []struct {
		in  base.Point
		dpi image.Point
		out image.Point
	}{
		{base.Point{1 * DIP, 2 * DIP}, image.Point{96, 96}, image.Point{1, 2}},
		{base.Point{1 * DIP, 2 * DIP}, image.Point{2 * 96, 3 * 96}, image.Point{2, 6}},
	}

	for i, v := range cases {
		base.DPI = v.dpi
		if out := v.in.Pixels(); out != v.out {
			t.Errorf("Error in case %d, want %s, got %s", i, v.out, out)
		}
	}
}

func TestRectangle(t *testing.T) {
	const DIP = base.DIP

	cases := []struct {
		x0, y0, x1, y1 base.Length
		min            base.Point
		width          base.Length
		height         base.Length
	}{
		{1 * DIP, 2 * DIP, 10 * DIP, 12 * DIP, base.Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
		{1 * DIP, 12 * DIP, 10 * DIP, 2 * DIP, base.Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
		{10 * DIP, 2 * DIP, 1 * DIP, 12 * DIP, base.Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
		{10 * DIP, 12 * DIP, 1 * DIP, 2 * DIP, base.Point{1 * DIP, 2 * DIP}, 9 * DIP, 10 * DIP},
	}

	for i, v := range cases {
		out := base.Rect(v.x0, v.y0, v.x1, v.y1)

		if out.Min != v.min {
			t.Errorf("Error in case %d, want %s, got %s", i, out.Min, v.min)
		}
		if got := out.Dx(); got != v.width {
			t.Errorf("Error in case %d, want %s, got %s", i, got, v.width)
		}
		if got := out.Dy(); got != v.height {
			t.Errorf("Error in case %d, want %s, got %s", i, got, v.height)
		}
		expected := base.Point{v.width, v.height}
		if got := out.Size(); got != expected {
			t.Errorf("Error in case %d, want %s, got %s", i, got, expected)
		}
	}
}
