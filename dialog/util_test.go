package dialog

import (
	"flag"
	"time"
)

var (
	asyncWait = 500 * time.Millisecond
)

type durationValue time.Duration

func (tv durationValue) String() string {
	return time.Duration(tv).String()
}

func (tv *durationValue) Set(s string) error {
	println("set")
	value, err := time.ParseDuration(s)
	if err == nil {
		*(*time.Duration)(tv) = value
	}
	return err
}

func init() {
	flag.Var((*durationValue)(&asyncWait), "async-wait", "Set delay before typing keys")
}

func asyncKeyEnter() {
	go func() {
		time.Sleep(asyncWait)
		typeKeys("\n")
	}()
}

func asyncKeyEscape() {
	go func() {
		time.Sleep(asyncWait)
		typeKeys("\x1b")
	}()
}

func asyncType(s string) {
	go func() {
		time.Sleep(asyncWait)
		typeKeys(s)
	}()
}
