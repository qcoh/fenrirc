package mondrian

// Box is a widget and container for other widgets.
type Box struct {
	*Region
	Visible

	Children   []Widget
	ResizeFunc func(*Region) []*Region
}

// NewBox returns a Box.
func NewBox() *Box {
	return &Box{Region: defaultRegion, Children: []Widget{}}
}

// Draw draws all of the box' children.
func (b *Box) Draw() {
	for _, w := range b.Children {
		w.Draw()
	}
}

// Resize resizes the box and calls the Resize method on all children.
func (b *Box) Resize(r *Region) {
	b.Region = r
	for k, v := range b.ResizeFunc(r) {
		b.Children[k].Resize(v)
	}
}

// SetVisibility sets the visibility of the box and all its children.
func (b *Box) SetVisibility(v bool) {
	for _, w := range b.Children {
		w.SetVisibility(v)
	}
	b.Visible.SetVisibility(v)
}
