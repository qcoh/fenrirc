package mondrian

// Dummy is a simple widget consisting of a solid color. Its purpose is to be used when designing layouts.
type Dummy struct {
	Region
	Visible
}

// Draw draws a solid rectangle with color given by `d.Bg`.
func (d *Dummy) Draw() {
	d.Clear()
}

// Resize resizes `d`.
func (d *Dummy) Resize(r *Region) {
	d.Region = *r
}
