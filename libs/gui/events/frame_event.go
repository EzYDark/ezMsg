package gui

import (
	"image"
	"image/color"
	"log" // For debugging output

	gio_app "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"github.com/ezydark/ezMsg/libs/gui/pages/overview"
)

// HandleFrameEvent remains the same
func HandleFrameEvent(event gio_app.FrameEvent) {
	var ops op.Ops
	gtx := gio_app.NewContext(&ops, event)
	// TestLayout(gtx)
	overview.Overview(gtx)
	event.Frame(gtx.Ops)
}

var th = material.NewTheme() // Theme variable

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Set this to true to use Max Height for Red Box, false for Min Height
const useMaxHeightForRedBox = false // <<<<------ TOGGLE BEHAVIOR HERE
const desiredMaxHeightRed = 50      // Example Max height for Red box
const desiredMinHeightRed = 100     // Example Min height for Red box

func TestLayout(gtx layout.Context) {
	log.Printf("TestLayout Gtx Constraints: Min=%+v, Max=%+v", gtx.Constraints.Min, gtx.Constraints.Max)

	totalAvailableHeight := gtx.Constraints.Max.Y

	weightBlue := float32(7)
	weightRed := float32(3)
	weightGreen := float32(7)
	totalWeight := weightBlue + weightRed + weightGreen

	// 1. Calculate Red Box's actual height
	var targetRedHeight int
	// Calculate red's initial offer based on its weight
	offerRed := 0
	if totalWeight > 0 {
		offerRed = int(weightRed / totalWeight * float32(totalAvailableHeight))
	} else if weightRed > 0 { // Edge case: only red has weight
		offerRed = totalAvailableHeight
	}

	if useMaxHeightForRedBox {
		targetRedHeight = min(desiredMaxHeightRed, offerRed)
	} else {
		targetRedHeight = max(desiredMinHeightRed, offerRed)
	}
	// Clamp red's height to be non-negative and not exceed total available height
	targetRedHeight = max(0, targetRedHeight)
	targetRedHeight = min(targetRedHeight, totalAvailableHeight)

	// 2. Calculate remaining space for Blue and Green
	heightForBlueAndGreen := totalAvailableHeight - targetRedHeight
	heightForBlueAndGreen = max(0, heightForBlueAndGreen) // Ensure non-negative

	// 3. Distribute remaining space to Blue and Green
	var targetBlueHeight, targetGreenHeight int
	weightBlueAndGreen := weightBlue + weightGreen

	if heightForBlueAndGreen > 0 {
		if weightBlueAndGreen > 0 {
			targetBlueHeight = int(weightBlue / weightBlueAndGreen * float32(heightForBlueAndGreen))
			targetGreenHeight = heightForBlueAndGreen - targetBlueHeight // Green gets the remainder
		} else {
			// If both blue and green have 0 weight, they get 0 height.
			targetBlueHeight = 0
			targetGreenHeight = 0
		}
	} else {
		targetBlueHeight = 0
		targetGreenHeight = 0
	}

	// Ensure non-negative heights for blue and green
	targetBlueHeight = max(0, targetBlueHeight)
	targetGreenHeight = max(0, targetGreenHeight)

	log.Printf("Pre-calculated Heights: TotalAvailable=%d, OfferRed=%d", totalAvailableHeight, offerRed)
	log.Printf("Targets: Blue=%d, Red=%d, Green=%d. Sum=%d",
		targetBlueHeight, targetRedHeight, targetGreenHeight,
		targetBlueHeight+targetRedHeight+targetGreenHeight)

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Blue Box (Top)
		layout.Flexed(weightBlue, func(gtx layout.Context) layout.Dimensions {
			// Blue box uses its calculated target height.
			size := image.Pt(gtx.Constraints.Max.X, targetBlueHeight)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 255, A: 255}, clip.Rect{Max: size}.Op())
			log.Printf("Blue Box: FlexOffer (MaxY In) = %d, ActualPaintedHeight = %d, Constraints In: Min=%+v, Max=%+v",
				gtx.Constraints.Max.Y, targetBlueHeight, gtx.Constraints.Min, gtx.Constraints.Max)
			return layout.Dimensions{Size: size}
		}),

		// Red Box (Middle)
		layout.Flexed(weightRed, func(gtx layout.Context) layout.Dimensions {
			// Red box's internal logic should yield targetRedHeight.
			// gtx.Constraints.Max.Y here is its original weighted offer (offerRed).
			currentActualHeight := 0
			if useMaxHeightForRedBox {
				currentActualHeight = min(desiredMaxHeightRed, gtx.Constraints.Max.Y)
			} else {
				currentActualHeight = max(desiredMinHeightRed, gtx.Constraints.Max.Y)
			}
			currentActualHeight = max(0, currentActualHeight) // Ensure non-negative
			// This currentActualHeight should match the pre-calculated targetRedHeight.

			finalSize := image.Point{X: gtx.Constraints.Max.X, Y: currentActualHeight}
			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 0, B: 0, A: 255}, clip.Rect{Max: finalSize}.Op())
			log.Printf("Red Box: FlexOffer (MaxY In) = %d, ActualPaintedHeight = %d (TargetRed=%d), Constraints In: Min=%+v, Max=%+v",
				gtx.Constraints.Max.Y, currentActualHeight, targetRedHeight, gtx.Constraints.Min, gtx.Constraints.Max)
			return layout.Dimensions{Size: finalSize}
		}),

		// Green Box (Bottom)
		layout.Flexed(weightGreen, func(gtx layout.Context) layout.Dimensions {
			// Green box uses its calculated target height.
			size := image.Pt(gtx.Constraints.Max.X, targetGreenHeight)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 255, B: 0, A: 255}, clip.Rect{Max: size}.Op())
			log.Printf("Green Box: FlexOffer (MaxY In) = %d, ActualPaintedHeight = %d, Constraints In: Min=%+v, Max=%+v",
				gtx.Constraints.Max.Y, targetGreenHeight, gtx.Constraints.Min, gtx.Constraints.Max)
			return layout.Dimensions{Size: size}
		}),
	)
}
