package components

import (
	"image"
	"image/color"
	"math" // Import math for min

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// CircleOpts defines the options for rendering a simple circle.
type CircleOpts struct {
	R     unit.Dp     // Radius of the circle.
	Color color.NRGBA // Color of the circle.
}

// Circle creates a simple circular widget filled with a solid color.
func Circle(opts CircleOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		radiusPx := gtx.Dp(opts.R)
		diameterPx := radiusPx * 2

		// Define the widget's size.
		widgetSize := image.Pt(diameterPx, diameterPx)

		// Use the layout constraints to determine a size that fits within constraints,
		// but then ensure the final size does not exceed the desired widgetSize.
		constrainedSize := gtx.Constraints.Constrain(widgetSize)

		// The actual size used should be the minimum of the constrained size and the desired widget size.
		// This prevents the circle from becoming larger than intended if the minimum constraint is large.
		size := min(constrainedSize, widgetSize)

		// Define the circular clip path for the widget's area.
		// All subsequent drawing operations within this widget will be clipped to this circle.
		circleArea := clip.Ellipse{Max: size} // Use the clamped size here
		// The Push method applies the clip, and Pop (called by defer) removes it.
		defer circleArea.Push(gtx.Ops).Pop()

		// Set the color for painting.
		paint.ColorOp{Color: opts.Color}.Add(gtx.Ops)

		// Fill the current clip (which is the circleArea) with the set color.
		paint.PaintOp{}.Add(gtx.Ops)

		// Return the dimensions of the widget.
		return layout.Dimensions{Size: size} // Return the clamped size
	}
}

// min returns a new image.Point with the minimum X and Y values of the two input points.
func min(p1, p2 image.Point) image.Point {
	return image.Pt(int(math.Min(float64(p1.X), float64(p2.X))), int(math.Min(float64(p1.Y), float64(p2.Y))))
}
