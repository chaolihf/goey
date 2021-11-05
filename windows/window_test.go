package windows_test

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"testing"
	"testing/quick"
	"time"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
	"bitbucket.org/rj/goey/mock"
	"bitbucket.org/rj/goey/windows"
)

func ExampleNewWindow() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := windows.NewWindow("Test", nil /*empty window*/)
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			fmt.Println("Up")
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Down")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			_ = loop.Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Output:
	// Up
	// Down
}

func ExampleWindow_Message() {
	// All calls that modify GUI objects need to be schedule ont he GUI thread.
	// This callback will be used to create the top-level window.
	createWindow := func() error {
		// Create a top-level window.
		mw, err := windows.NewWindow("Test", nil /*empty window*/)
		if err != nil {
			// This error will be reported back up through the call to
			// Run below.  No need to print or log it here.
			return err
		}

		// We can start a goroutine, but note that we can't modify GUI objects
		// directly.
		go func() {
			// Show the error message.
			_ = loop.Do(func() error {
				return mw.Message("This is an example message.").WithInfo().Show()
			})

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			_ = loop.Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	// Start the GUI thread.
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func TestMain(m *testing.M) {
	// On Cocoa, the GUI even
	loop.TestMain(m)
}

func TestNewWindow(t *testing.T) {
	errSentinel := errors.New("sentinel error")

	cases := []struct {
		widget base.Widget
		out    error
	}{
		{nil, nil},
		{&mock.Widget{}, nil},
		{&mock.Widget{Size: base.Size{base.DIP.Scale(16*1024, 96), base.DIP.Scale(16*1024, 96)}}, nil}, // 16k pixels by 16k pixels
		{&mock.Widget{Err: errSentinel}, errSentinel},
	}

	for i, v := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			createWindow := func() error {
				// Create the window.  Some of the tests here are not expected in
				// production code, but we can be a little paranoid here.
				window, err := windows.NewWindow(t.Name(), v.widget)
				if err != nil {
					return err
				}
				if window == nil {
					t.Errorf("unexpected nil for window")
					return nil
				}

				window.Close()
				return nil
			}

			if err := loop.Run(createWindow); err != v.out {
				t.Errorf("unexpected return: want %s, got %s", v.out, err)
			}
		})
	}

	t.Run("QuickCheck", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping test in short mode")
		}

		f := func(text string, width, height uint16) bool {
			createWindow := func() error {
				size := base.Size{
					Width:  base.DIP.Scale(int(width), 96),
					Height: base.DIP.Scale(int(height), 96),
				}
				window, err := windows.NewWindow(t.Name(), &mock.Widget{Size: size})
				if err != nil {
					return err
				}
				if window == nil {
					return errors.New("unexpected nil for window")
				}

				window.Close()
				return nil
			}

			err := loop.Run(createWindow)
			return err == nil
		}
		if err := quick.Check(f, nil); err != nil {
			t.Errorf("quick: %s", err)
		}
	})
}

func testingWindow(t *testing.T, action func(*testing.T, *windows.Window)) {
	createWindow := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		mw, err := windows.NewWindow(t.Name(), nil)
		if err != nil {
			t.Fatalf("failed to create window: %s", err)
		}
		if mw == nil {
			t.Fatalf("unexpected nil for window")
		}

		go func() {
			// Delegate to test specific actions.
			action(t, mw)
			if testing.Verbose() {
				time.Sleep(250 * time.Millisecond)
			}

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			err := loop.Do(func() error {
				mw.Close()
				return nil
			})
			if err != nil {
				// Would like to report this error using t.Fatalf, but we are
				// not in the same goroutine.  Could send a message using a
				// channel, but if the call to Do failed, it is not certain that
				// we closed the window, and could deadlock.
				panic(err)
			}
		}()

		return nil
	}

	err := loop.Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}

func TestWindow_MinSize(t *testing.T) {
	cases := []struct {
		child            base.Widget
		hscroll, vscroll bool
		minSize          base.Size
	}{
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10 * base.DIP}}, false, false, base.Size{10 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10 * base.DIP}}, false, true, base.Size{10 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10 * base.DIP}}, true, false, base.Size{10 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10 * base.DIP}}, true, true, base.Size{10 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10000 * base.DIP, 10 * base.DIP}}, false, false, base.Size{10000 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10000 * base.DIP, 10 * base.DIP}}, false, true, base.Size{10000 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10000 * base.DIP, 10 * base.DIP}}, true, false, base.Size{120 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10000 * base.DIP, 10 * base.DIP}}, true, true, base.Size{120 * base.DIP, 10 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10000 * base.DIP}}, false, false, base.Size{10 * base.DIP, 10000 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10000 * base.DIP}}, false, true, base.Size{10 * base.DIP, 120 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10000 * base.DIP}}, true, false, base.Size{10 * base.DIP, 10000 * base.DIP}},
		{&mock.Widget{Size: base.Size{10 * base.DIP, 10000 * base.DIP}}, true, true, base.Size{10 * base.DIP, 120 * base.DIP}},
	}

	for i, v := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			testingWindow(t, func(t *testing.T, mw *windows.Window) {
				err := loop.Do(func() error {
					if err := mw.SetChild(v.child); err != nil {
						t.Errorf("failed to set child: %s", err)
						return nil
					}

					mw.SetScroll(v.hscroll, v.vscroll)
					if got := mw.MinSize(); got != v.minSize {
						t.Errorf("incorrect minimum size, want %s, got %s", v.minSize, got)
					}

					return nil
				})
				if err != nil {
					t.Errorf("failed loop.Do: %s", err)
				}
			})
		})
	}
}

func makeImage(t *testing.T, index int) image.Image {
	colors := [3]color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	bounds := image.Rect(0, 0, 32, 32)
	img := image.NewRGBA(bounds)
	draw.Draw(img, image.Rect(0, 0, 11, 32), image.NewUniform(colors[index%3]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(11, 0, 22, 32), image.NewUniform(colors[(index+1)%3]), image.Point{}, draw.Src)
	draw.Draw(img, image.Rect(22, 0, 32, 32), image.NewUniform(colors[(index+2)%3]), image.Point{}, draw.Src)
	return img
}

func TestWindow_SetIcon(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *windows.Window) {
		for i := 0; i < 6; i++ {
			img := makeImage(t, i)

			err := loop.Do(func() error {
				return mw.SetIcon(img)
			})
			if err != nil {
				t.Errorf("Error calling SetIcon, %s", err)
			}
			time.Sleep(50 * time.Millisecond)
		}
	})
}

func TestWindow_SetScroll(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *windows.Window) {
		cases := []struct {
			horizontal, vertical bool
		}{
			{false, false},
			{false, true},
			{true, false},
			{true, true},
		}

		for i, v := range cases {
			err := loop.Do(func() error {
				mw.SetScroll(v.horizontal, v.vertical)
				out1, out2 := mw.Scroll()
				if out1 != v.horizontal {
					t.Errorf("Case %d: Returned flag for horizontal scroll does not match, got %v, want %v", i, out1, v.horizontal)
				}
				if out2 != v.vertical {
					t.Errorf("Case %d: Returned flag for vertical scroll does not match, got %v, want %v", i, out2, v.vertical)
				}
				return nil
			})
			if err != nil {
				t.Errorf("Error calling SetTitle, %s", err)
			}
		}
	})
}

func TestWindow_SetTitle(t *testing.T) {
	testingWindow(t, func(t *testing.T, mw *windows.Window) {
		err := loop.Do(func() error {
			err := mw.SetTitle("Flash!")
			if err != nil {
				return err
			}

			if got := mw.Title(); got != "Flash!" {
				t.Errorf("Failed to set title correctly, got %s", got)
			}

			return nil
		})
		if err != nil {
			t.Errorf("Error calling SetTitle, %s", err)
		}
	})
}
