package base_test

import (
	"fmt"

	"bitbucket.org/rj/goey/base"
)

func ExampleKind_String() {
	kind := base.NewKind("bitbucket.org/rj/goey/base.Example")

	fmt.Println("Kind is", kind.String())

	// Output:
	// Kind is bitbucket.org/rj/goey/base.Example
}
