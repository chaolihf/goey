package dialog

// Dialog contains some functionality used by all dialog types.
type Dialog struct {
	owner Owner
	err   error
}

// Err returns the first error that was encountered while building the dialog.
func (d *Dialog) Err() error {
	return d.err
}

// WithOwner sets the owner of the dialog box.
func (d *Dialog) WithOwner(owner Owner) *Dialog {
	d.owner = owner
	return d
}
