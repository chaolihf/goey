package gtk

import "C"
import "unsafe"

type Window interface {
	Widget
	OnDeleteEvent() bool
	OnSizeAllocate(width, height int)
}

//export onDeleteEvent
func onDeleteEvent(handle unsafe.Pointer) bool {
	return widgets[uintptr(handle)].(Window).OnDeleteEvent()
}

//export onSizeAllocate
func onSizeAllocate(handle unsafe.Pointer, width, height int) {
	widgets[uintptr(handle)].(Window).OnSizeAllocate(width, height)
}
