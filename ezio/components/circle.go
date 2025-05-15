package components

import (
	"image"
	"image/color"

	// Standard library image decoders.
	// Ensure you have these blank imports for your app to recognize PNG, JPG, GIF.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os" // We'll keep os for the helper, but it's used outside Circle now

	"gioui.org/f32"      // For 2D float points and affine transformations
	"gioui.org/layout"   // For layout context, dimensions, etc.
	"gioui.org/op"       // For operations like op.Affine
	"gioui.org/op/clip"  // For clipping operations (making it circular)
	"gioui.org/op/paint" // For painting colors and images
	"gioui.org/unit"     // For density-independent pixels (Dp)
	"github.com/rs/zerolog/log"
)

// CircleOpts defines the options for rendering a circle.
type CircleOpts struct {
	R     unit.Dp     // Radius of the circle.
	Color color.NRGBA // Fallback color if Img is nil or has issues.
	Img   image.Image // The pre-loaded image to display.
}

// Circle creates a circular widget.
// If opts.Img is provided, it displays the image clipped to a circle.
// Otherwise, it fills the circle with opts.Color.
func Circle(opts CircleOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		radiusPx := gtx.Dp(opts.R)
		diameterPx := radiusPx * 2
		widgetSize := image.Pt(diameterPx, diameterPx)

		// Constrain the size but don't exceed the desired diameter.
		constrainedSize := gtx.Constraints.Constrain(widgetSize)
		finalSize := minPt(constrainedSize, widgetSize)

		// If the final size is zero or negative, there's nothing to draw.
		if finalSize.X <= 0 || finalSize.Y <= 0 {
			return layout.Dimensions{Size: finalSize}
		}

		// Define the circular clip path. All drawing will be clipped by this.
		circleArea := clip.Ellipse{Max: finalSize}
		defer circleArea.Push(gtx.Ops).Pop()

		if opts.Img != nil {
			// Create an image operation from the pre-loaded image.
			imgOp := paint.NewImageOp(opts.Img)
			imgBounds := imgOp.Size() // Get dimensions of the source image.

			// Check for valid image dimensions.
			if imgBounds.X <= 0 || imgBounds.Y <= 0 {
				log.Error().Msg("Pre-loaded image has invalid (zero or negative) dimensions. Falling back to color.")
				paint.ColorOp{Color: opts.Color}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
			} else {
				// Convert dimensions to float32 for scaling calculations.
				imgWidthF := float32(imgBounds.X)
				imgHeightF := float32(imgBounds.Y)
				targetWidthF := float32(finalSize.X)
				targetHeightF := float32(finalSize.Y)

				// Calculate scale factors.
				// We want the image to cover the circle, so we find the larger scale needed.
				scaleX := targetWidthF / imgWidthF
				scaleY := targetHeightF / imgHeightF
				scale := float32(math.Max(float64(scaleX), float64(scaleY)))

				// Calculate the dimensions of the image after scaling.
				scaledWidth := imgWidthF * scale
				scaledHeight := imgHeightF * scale

				// Calculate offsets to center the scaled image within the circle area.
				dx := (targetWidthF - scaledWidth) / 2
				dy := (targetHeightF - scaledHeight) / 2

				// Define the transformation: scale then offset.
				transform := f32.Affine2D{}.
					Scale(f32.Pt(0, 0), f32.Pt(scale, scale)). // Scale relative to image's origin (0,0)
					Offset(f32.Pt(dx, dy))                     // Then, translate the scaled image

				imgOp.Filter = paint.FilterLinear // Use linear filter for smoother scaling.

				// Apply the transformation, draw the image, then revert the transformation.
				stack := op.Affine(transform).Push(gtx.Ops)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				stack.Pop()
			}
		} else {
			// Fallback to solid color if no image is provided.
			paint.ColorOp{Color: opts.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}

		return layout.Dimensions{Size: finalSize}
	}
}

// LoadImageHelper is a utility function you can use in your application
// to load an image file once.
// You might want to move this to a common utility package in your project.
func LoadImageHelper(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// image.Decode automatically detects the format (PNG, JPEG, GIF)
	// if the respective decoders are imported with blank imports.
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// minPt returns a new image.Point with the minimum X and Y values from two points.
func minPt(p1, p2 image.Point) image.Point {
	x := p1.X
	if p2.X < x {
		x = p2.X
	}
	y := p1.Y
	if p2.Y < y {
		y = p2.Y
	}
	return image.Point{X: x, Y: y}
}
