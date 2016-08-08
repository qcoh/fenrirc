package mondrian

// Dummy is a simple widget consisting of a solid color. Its purpose is to be used when designing layouts.
type Dummy struct {
	Region
	Visible
}

// NewDummy returns a Dummy.
func NewDummy() *Dummy {
	return &Dummy{Region: *defaultRegion}
}

// Draw draws a solid rectangle with color given by `d.Bg`.
func (d *Dummy) Draw() {
	d.Clear()
}

// Resize resizes `d`.
func (d *Dummy) Resize(r *Region) {
	color := d.Region.Bg
	d.Region = *r
	d.Region.Bg = color
}
