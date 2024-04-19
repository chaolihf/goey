package base_test

import (
	"fmt"

	"github.com/chaolihf/goey/base"
)

func ExampleKind_String() {
	kind := base.NewKind("github.com/chaolihf/goey/base.Example")

	fmt.Println("Kind is", kind.String())

	// Output:
	// Kind is github.com/chaolihf/goey/base.Example
}
