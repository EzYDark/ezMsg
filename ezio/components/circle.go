package components

import (
	"fmt"
	"image"
	"image/color"
	"image/draw" // For image manipulation and masking
	"math"
	"net/http"
	"os"
	"sync"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/rs/zerolog/log"
)

// CircleOpts defines the options for rendering a circle.
type CircleOpts struct {
	R      unit.Dp
	Color  color.NRGBA
	Img    image.Image // Pre-loaded image (local)
	ImgURL string      // URL for remote image
}

var (
	imageCache        = make(map[string]image.Image)
	imageCacheMutex   = &sync.Mutex{}
	urlFetchOnce      = make(map[string]*sync.Once)
	urlFetchOnceMutex = &sync.Mutex{}

	rasterizedCircleCache      = make(map[string]image.Image)
	rasterizedCircleCacheMutex = &sync.Mutex{}

	fixedRasterizedDiameterPx = 128
)

// generateCircularImage creates a new square image of fixedRasterizedDiameterPx,
// containing the original image scaled (nearest neighbor), centered, and masked to be circular.
func generateCircularImage(original image.Image) image.Image {
	if original == nil {
		return nil
	}

	// --- Prepare srcForScaling (center-cropped square from original) ---
	srcForScaling := original
	origBounds := original.Bounds()
	origW, origH := origBounds.Dx(), origBounds.Dy()

	if origW <= 0 || origH <= 0 { // Check for invalid original image dimensions
		log.Error().Msg("generateCircularImage: Original image has zero or negative dimensions.")
		return nil
	}

	if origW != origH { // If not square, crop the center square portion
		var cropRect image.Rectangle
		if origW > origH { // Wider
			size := origH
			x0 := origBounds.Min.X + (origW-size)/2
			cropRect = image.Rect(x0, origBounds.Min.Y, x0+size, origBounds.Max.Y)
		} else { // Taller
			size := origW
			y0 := origBounds.Min.Y + (origH-size)/2
			cropRect = image.Rect(origBounds.Min.X, y0, origBounds.Max.X, y0+size)
		}
		type subImager interface {
			SubImage(r image.Rectangle) image.Image
		}
		if si, ok := original.(subImager); ok {
			srcForScaling = si.SubImage(cropRect)
		} else {
			log.Warn().Str("image_type", fmt.Sprintf("%T", original.ColorModel())).Msg("Original image type does not support SubImage for cropping. Using full image for scaling; aspect ratio might be impacted if not square.")
		}
	}
	// srcForScaling is now the square portion to be scaled.
	srcForScalingBounds := srcForScaling.Bounds()
	sfsDx := srcForScalingBounds.Dx()
	sfsDy := srcForScalingBounds.Dy()

	if sfsDx <= 0 || sfsDy <= 0 { // Check for invalid cropped image dimensions
		log.Error().Msg("generateCircularImage: Cropped source for scaling has zero or negative dimensions.")
		return nil
	}

	// --- Scale srcForScaling to fixedRasterizedDiameterPx x fixedRasterizedDiameterPx (Nearest Neighbor) ---
	// Why: Create an intermediate image that holds the correctly scaled version of the cropped source.
	scaledIntermediate := image.NewRGBA(image.Rect(0, 0, fixedRasterizedDiameterPx, fixedRasterizedDiameterPx))

	for y_dst := 0; y_dst < fixedRasterizedDiameterPx; y_dst++ {
		for x_dst := 0; x_dst < fixedRasterizedDiameterPx; x_dst++ {
			// Map (x_dst, y_dst) in scaledIntermediate to (x_src, y_src) in srcForScaling
			x_src := srcForScalingBounds.Min.X + (x_dst * sfsDx / fixedRasterizedDiameterPx)
			y_src := srcForScalingBounds.Min.Y + (y_dst * sfsDy / fixedRasterizedDiameterPx)

			// Clamp source coordinates to be within bounds of srcForScaling.
			// Bounds.Max is exclusive, so valid coordinates are Min to Max-1.
			if x_src >= srcForScalingBounds.Max.X {
				x_src = srcForScalingBounds.Max.X - 1
			}
			if y_src >= srcForScalingBounds.Max.Y {
				y_src = srcForScalingBounds.Max.Y - 1
			}
			// Ensure Min bounds are also respected (though less likely an issue with the formula above if Dx/Dy > 0)
			if x_src < srcForScalingBounds.Min.X {
				x_src = srcForScalingBounds.Min.X
			}
			if y_src < srcForScalingBounds.Min.Y {
				y_src = srcForScalingBounds.Min.Y
			}

			scaledIntermediate.Set(x_dst, y_dst, srcForScaling.At(x_src, y_src))
		}
	}

	// --- Apply circular mask to the scaledIntermediate image ---
	// Why: Now that we have a correctly scaled version of the cat (or other image content),
	// we create the final destination image and apply the circular alpha mask to this scaled content.
	finalDst := image.NewRGBA(image.Rect(0, 0, fixedRasterizedDiameterPx, fixedRasterizedDiameterPx))
	circularMask := image.NewAlpha(image.Rect(0, 0, fixedRasterizedDiameterPx, fixedRasterizedDiameterPx))
	centerX := float32(fixedRasterizedDiameterPx) / 2.0
	// centerY := float32(fixedRasterizedDiameterPx) / 2.0 // Redundant if centerX is used for radiusSq
	radiusSq := centerX * centerX

	for y := 0; y < fixedRasterizedDiameterPx; y++ {
		for x := 0; x < fixedRasterizedDiameterPx; x++ {
			dxPos := float32(x) + 0.5 - centerX
			dyPos := float32(y) + 0.5 - centerX // Should be centerY, but since it's a square, centerX == centerY
			if dxPos*dxPos+dyPos*dyPos <= radiusSq {
				circularMask.SetAlpha(x, y, color.Alpha{A: 0xff})
			} else {
				circularMask.SetAlpha(x, y, color.Alpha{A: 0x00})
			}
		}
	}

	// Draw the already scaled image (scaledIntermediate) onto finalDst using the mask.
	// Source and destination are now the same size (fixedRasterizedDiameterPx), so no unintended scaling by DrawMask.
	draw.DrawMask(finalDst, finalDst.Bounds(), scaledIntermediate, image.Point{}, circularMask, image.Point{}, draw.Over)

	return finalDst
}

// Circle creates a circular widget.
func Circle(opts CircleOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		radiusPx := gtx.Dp(opts.R)
		diameterPx := radiusPx * 2
		if diameterPx <= 0 {
			diameterPx = 1
		}

		targetWidgetSize := image.Pt(diameterPx, diameterPx)
		finalDrawSize := gtx.Constraints.Constrain(targetWidgetSize)

		if finalDrawSize.X <= 0 || finalDrawSize.Y <= 0 {
			return layout.Dimensions{Size: finalDrawSize}
		}

		var imageToRender image.Image
		isPreRasterized := false

		if opts.ImgURL != "" {
			originalImg := blockingFetchAndCache(opts.ImgURL)
			if originalImg != nil {
				rasterizedCircleCacheMutex.Lock()
				cachedRasterizedImg, found := rasterizedCircleCache[opts.ImgURL]
				rasterizedCircleCacheMutex.Unlock()

				if found {
					imageToRender = cachedRasterizedImg
					isPreRasterized = true
				} else {
					log.Debug().Msgf("Generating pre-rasterized circle for URL: %s", opts.ImgURL)
					newRasterizedImg := generateCircularImage(originalImg)
					if newRasterizedImg != nil {
						rasterizedCircleCacheMutex.Lock()
						rasterizedCircleCache[opts.ImgURL] = newRasterizedImg
						rasterizedCircleCacheMutex.Unlock()
						imageToRender = newRasterizedImg
						isPreRasterized = true
					} else {
						imageToRender = originalImg
						log.Warn().Msgf("Failed to generate pre-rasterized circle, falling back to original for URL: %s", opts.ImgURL)
					}
				}
			} else if opts.Img != nil {
				imageToRender = opts.Img
			}
		} else if opts.Img != nil {
			imageToRender = opts.Img
		}

		if imageToRender != nil {
			imgOp := paint.NewImageOp(imageToRender)
			imgOp.Filter = paint.FilterLinear

			if isPreRasterized {
				sourceDiameter := float32(fixedRasterizedDiameterPx)
				targetDiameter := float32(finalDrawSize.X)
				if finalDrawSize.Y < finalDrawSize.X {
					targetDiameter = float32(finalDrawSize.Y)
				}

				scale := targetDiameter / sourceDiameter

				scaledW := sourceDiameter * scale
				scaledH := sourceDiameter * scale

				offsetX := (float32(finalDrawSize.X) - scaledW) / 2.0
				offsetY := (float32(finalDrawSize.Y) - scaledH) / 2.0

				transform := f32.Affine2D{}.
					Scale(f32.Pt(0, 0), f32.Pt(scale, scale)).
					Offset(f32.Pt(offsetX, offsetY))

				stack := op.Affine(transform).Push(gtx.Ops)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				stack.Pop()

			} else { // Fallback to original method: scale, center, and clip original/local image
				clippingArea := clip.Ellipse{Max: finalDrawSize}
				defer clippingArea.Push(gtx.Ops).Pop()

				imgBounds := imageToRender.Bounds()
				if imgBounds.Dx() <= 0 || imgBounds.Dy() <= 0 {
					log.Error().Msg("Fallback image has invalid dimensions. Drawing color.")
					paint.ColorOp{Color: opts.Color}.Add(gtx.Ops)
				} else {
					imgWidthF := float32(imgBounds.Dx())
					imgHeightF := float32(imgBounds.Dy())
					targetWidthF := float32(finalDrawSize.X)
					targetHeightF := float32(finalDrawSize.Y)

					scaleXVal := targetWidthF / imgWidthF
					scaleYVal := targetHeightF / imgHeightF
					finalScale := float32(math.Max(float64(scaleXVal), float64(scaleYVal)))

					scaledWidth := imgWidthF * finalScale
					scaledHeight := imgHeightF * finalScale

					dx := (targetWidthF - scaledWidth) / 2
					dy := (targetHeightF - scaledHeight) / 2

					transform := f32.Affine2D{}.
						Scale(f32.Pt(0, 0), f32.Pt(finalScale, finalScale)).
						Offset(f32.Pt(dx, dy))

					stack := op.Affine(transform).Push(gtx.Ops)
					imgOp.Add(gtx.Ops)
					stack.Pop()
				}
				paint.PaintOp{}.Add(gtx.Ops)
			}
		} else { // No image available, draw fallback color circle
			clippingArea := clip.Ellipse{Max: finalDrawSize}
			defer clippingArea.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: opts.Color}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		return layout.Dimensions{Size: finalDrawSize}
	}
}

// blockingFetchAndCache, LoadImageHelper, minPt remain the same from previous version.
func blockingFetchAndCache(url string) image.Image {
	imageCacheMutex.Lock()
	img, found := imageCache[url]
	imageCacheMutex.Unlock()
	if found {
		return img
	}

	urlFetchOnceMutex.Lock()
	once, ok := urlFetchOnce[url]
	if !ok {
		once = &sync.Once{}
		urlFetchOnce[url] = once
	}
	urlFetchOnceMutex.Unlock()

	var fetchedImg image.Image
	once.Do(func() {
		log.Debug().Msgf("Fetching (original) for URL: %s", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to HTTP GET original image from URL: %s", url)
			return
		}
		defer resp.Body.Close()

		decodedImg, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to decode original image from URL: %s", url)
			return
		}
		fetchedImg = decodedImg
		imageCacheMutex.Lock()
		imageCache[url] = fetchedImg
		imageCacheMutex.Unlock()
	})

	if fetchedImg != nil {
		return fetchedImg
	} else {
		imageCacheMutex.Lock()
		img, found = imageCache[url] // Re-check cache in case another goroutine populated it during Do.
		imageCacheMutex.Unlock()
		if found {
			return img
		}
		return nil // Fetch failed
	}
}

func LoadImageHelper(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		return nil, errDecode
	}
	return img, nil
}

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
