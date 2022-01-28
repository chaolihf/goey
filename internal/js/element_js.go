package goeyjs

import "syscall/js"

func AppendChildToBody(child js.Value) {
	document := js.Global().Get("document")

	body := document.Call("getElementsByTagName", "body").Index(0)
	body.Call("appendChild", child)
}

func CreateElement(tagName string, className string) js.Value {
	document := js.Global().Get("document")

	element := document.Call("createElement", tagName)
	if className != "" {
		element.Set("className", className)
	}

	return element
}
