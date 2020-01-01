package cocoa

/*
#cgo CFLAGS: -x objective-c -DNTRACE -I/usr/include/GNUstep
#cgo LDFLAGS: -lgnustep-gui -lgnustep-base -lobjc
#include "loop.h"
*/
import "C"
import "sync"
import "bitbucket.org/rj/goey/internal/nopanic"

func Init() error {
	C.init()
	return nil
}

func Run() error {
	// Run the event loop.
	C.run()
	return nil
}

var (
	thunkAction func() error
	thunkErr    error
	thunkMutex  sync.Mutex
)

func PerformOnMainThread(action func() error) error {
	// Lock thunk to avoid overwriting of thunkAction or thunkErr
	thunkMutex.Lock()
	defer thunkMutex.Unlock()
	// Is additional syncronization required to provide memory barriers to
	// coordinate with the GUI thread?

	thunkAction = action
	C.performOnMainThread()
	return nopanic.Unwrap(thunkErr)
}

//export callbackDo
func callbackDo() {
	thunkErr = nopanic.Wrap(thunkAction)
}

func Stop() {
	C.stop()
}

func IsMainThread() bool {
	return C.isMainThread() != 0
}
