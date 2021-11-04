package windows

import (
	"bitbucket.org/rj/goey/base"
)

func max(a, b base.Length) base.Length {
	if a > b {
		return a
	}
	return b
}

func min(a, b base.Length) base.Length {
	if a < b {
		return a
	}
	return b
}
