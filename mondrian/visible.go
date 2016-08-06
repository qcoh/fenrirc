package mondrian

// Visible abstracts the common to all widgets `IsVisible` and `SetVisibility` methods.
type Visible struct {
	visible bool
}

// IsVisible returns `v.visible`.
func (v *Visible) IsVisible() bool {
	return v.visible
}

// SetVisibility sets `v.visible`.
func (v *Visible) SetVisibility(visible bool) {
	v.visible = visible
}
