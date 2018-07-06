package goey

var (
	hboxKind = Kind{"bitbucket.org/rj/goey.HBox"}
)

// HBox describes a layout widget that arranges its child widgets into a row.
// Children are positioned in order from the left towards the right.  The main
// axis for alignment is therefore horizontal, with the cross axis for alignment is vertical.
//
// The size of the box will try to set a width sufficient to contain all of its
// children.  Extra space will be distributed according to the value of
// AlignMain.  Subject to the box constraints during layout, the height should
// match the largest minimum height of the child widgets.
type HBox struct {
	AlignMain  MainAxisAlign  // Control distribution of excess horizontal space when positioning children.
	AlignCross CrossAxisAlign // Control distribution of excess vertical space when positioning children.
	Children   []Widget       // Children.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HBox) Kind() *Kind {
	return &hboxKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HBox) Mount(parent Control) (Element, error) {
	c := make([]Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			CloseElements(c)
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &hboxElement{
		parent:       parent,
		children:     c,
		alignMain:    w.AlignMain,
		alignCross:   w.AlignCross,
		childrenSize: make([]Size, len(c)),
	}, nil
}

func (*hboxElement) Kind() *Kind {
	return &hboxKind
}

type hboxElement struct {
	parent     Control
	children   []Element
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign

	childrenSize []Size
	totalWidth   Length
}

func (w *hboxElement) Close() {
	CloseElements(w.children)
	w.children = nil
}

func (w *hboxElement) Layout(bc Constraint) Size {
	if len(w.children) == 0 {
		w.totalWidth = 0
		return bc.Constrain(Size{})
	}

	// Determine the constraints for layout of child elements.
	cbc := bc.LoosenWidth()
	if w.alignMain == Homogeneous {
		if count := len(w.children); count > 1 {
			gap := calculateHGap(nil, nil)
			cbc.TightenWidth(cbc.Max.Width.Scale(1, count) - gap.Scale(1, count-1))
		} else {
			cbc.TightenWidth(cbc.Max.Width.Scale(1, count))
		}
	}
	if w.alignCross == Stretch {
		if cbc.HasBoundedHeight() {
			cbc = cbc.TightenHeight(cbc.Max.Height)
		} else {
			cbc = cbc.TightenHeight(max(cbc.Min.Height, w.MinIntrinsicHeight(Inf)))
		}
	} else {
		cbc = cbc.LoosenHeight()
	}

	width := Length(0)
	minHeight := Length(0)
	previous := Element(nil)
	for i, v := range w.children {
		if i > 0 {
			if w.alignMain.IsPacked() {
				width += calculateHGap(previous, v)
			} else {
				width += calculateHGap(nil, nil)
			}
			previous = v
		}
		w.childrenSize[i] = v.Layout(cbc)
		minHeight = max(minHeight, w.childrenSize[i].Height)
		width += w.childrenSize[i].Width
	}
	w.totalWidth = width

	if w.alignCross == Stretch {
		return bc.Constrain(Size{width, cbc.Min.Height})
	}
	return bc.Constrain(Size{width, minHeight})
}

func (w *hboxElement) MinIntrinsicHeight(width Length) Length {
	if len(w.children) == 0 {
		return 0
	}

	if w.alignMain == Homogeneous {
		width = guardInf(width, width.Scale(1, len(w.children)))
		size := w.children[0].MinIntrinsicHeight(width)
		for _, v := range w.children[1:] {
			size = max(size, v.MinIntrinsicHeight(width))
		}
		return size
	}

	size := w.children[0].MinIntrinsicHeight(Inf)
	for _, v := range w.children[1:] {
		size = max(size, v.MinIntrinsicHeight(Inf))
	}
	return size
}

func (w *hboxElement) MinIntrinsicWidth(height Length) Length {
	if len(w.children) == 0 {
		return 0
	}

	size := w.children[0].MinIntrinsicWidth(height)
	if w.alignMain.IsPacked() {
		previous := w.children[0]
		for _, v := range w.children[1:] {
			// Add the preferred gap between this pair of widgets
			size += calculateHGap(previous, v)
			// Find minimum size for this widget, and update
			size += v.MinIntrinsicWidth(height)
		}
		return size
	}

	if w.alignMain == Homogeneous {
		for _, v := range w.children[1:] {
			size = max(size, v.MinIntrinsicWidth(height))
		}

		// Add a minimum gap between the controls.
		size = size.Scale(len(w.children), 1) + calculateHGap(nil, nil).Scale(len(w.children)-1, 1)
		return size
	}

	for _, v := range w.children[1:] {
		size += v.MinIntrinsicWidth(height)
	}

	// Add a minimum gap between the controls.
	if w.alignMain == SpaceBetween {
		size += calculateHGap(nil, nil).Scale(len(w.children)-1, 1)
	} else {
		size += calculateHGap(nil, nil).Scale(len(w.children)+1, 1)
	}

	return size
}

func (w *hboxElement) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	if w.alignMain == Homogeneous {
		gap := calculateHGap(nil, nil)
		dx := bounds.Dx() + gap
		count := len(w.children)

		for i, v := range w.children {
			x1 := bounds.Min.X + dx.Scale(i, count)
			x2 := bounds.Min.X + dx.Scale(i+1, count) - gap
			w.setBoundsForChild(i, v, x1, bounds.Min.Y, x2, bounds.Max.Y)
		}
		return
	}

	// Adjust the bounds so that the minimum Y handles vertical alignment
	// of the controls.  We also calculate 'extraGap' which will adjust
	// spacing of the controls for non-packed alignments.
	extraGap := Length(0)
	switch w.alignMain {
	case MainStart:
		// Do nothing
	case MainCenter:
		bounds.Min.X += (bounds.Dx() - w.totalWidth) / 2
	case MainEnd:
		bounds.Min.X = bounds.Max.X - w.totalWidth
	case SpaceAround:
		extraGap = (bounds.Dx() - w.totalWidth).Scale(1, len(w.children)+1)
		bounds.Min.X += extraGap
	case SpaceBetween:
		if len(w.children) > 1 {
			extraGap = (bounds.Dx() - w.totalWidth).Scale(1, len(w.children)-1)
		} else {
			// There are no controls between which to put the extra space.
			// The following essentially convert SpaceBetween to SpaceAround
			bounds.Min.X += (bounds.Dx() - w.totalWidth) / 2
		}
	}

	// Position all of the child controls.
	posX := bounds.Min.X
	previous := Element(nil)
	for i, v := range w.children {
		if w.alignMain.IsPacked() {
			if i > 0 {
				posX += calculateHGap(previous, v)
			}
			previous = v
		}

		dx := w.childrenSize[i].Width
		w.setBoundsForChild(i, v, posX, bounds.Min.Y, posX+dx, bounds.Max.Y)
		posX += dx + extraGap
	}
}

func (w *hboxElement) setBoundsForChild(i int, v Element, posX, posY, posX2, posY2 Length) {
	dy := w.childrenSize[i].Height
	switch w.alignCross {
	case CrossStart:
		v.SetBounds(Rectangle{
			Point{posX, posY},
			Point{posX2, posY + dy},
		})
	case CrossCenter:
		v.SetBounds(Rectangle{
			Point{posX, posY + (posY2-posY-dy)/2},
			Point{posX2, posY + (posY2-posY+dy)/2},
		})
	case CrossEnd:
		v.SetBounds(Rectangle{
			Point{posX, posY2 - dy},
			Point{posX2, posY2},
		})
	case Stretch:
		v.SetBounds(Rectangle{
			Point{posX, posY},
			Point{posX2, posY2},
		})
	}
}

func (w *hboxElement) updateProps(data *HBox) (err error) {
	// Update properties
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	w.children, err = DiffChildren(w.parent, w.children, data.Children)
	// Clear cached values
	w.childrenSize = make([]Size, len(w.children))
	w.totalWidth = 0
	return err
}

func (w *hboxElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*HBox))
}
