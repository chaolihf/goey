package dialog

import (
	"errors"
)

func (m *Message) show() error {
	return errors.New("Message.Show not imiplemented on js")
}

func (m *Message) withError() {
}

func (m *Message) withWarn() {
}

func (m *Message) withInfo() {
}
