package components

import (
	"image"
	"image/color"
	"math"
	"net/http"
	"sync" // Added for image caching

	"gioui.org/f32" // Added for image transformation
	"gioui.org/layout"
	"gioui.org/op" // Added for op.Affine
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/rs/zerolog/log" // Added for logging
)

type RectOpts struct {
	W, H       unit.Dp
	MinW, MinH unit.Dp
	Color      color.NRGBA
	Img        image.Image // Pre-loaded image (local)
	ImgURL     string      // URL for remote image
}

// Image caching mechanisms (similar to Circle component)
var (
	rectImageCache        = make(map[string]image.Image)
	rectImageCacheMutex   = &sync.Mutex{}
	rectUrlFetchOnce      = make(map[string]*sync.Once)
	rectUrlFetchOnceMutex = &sync.Mutex{}
)

// blockingFetchAndCache fetches an image from a URL and caches it.
// This is a helper function adapted from your Circle component.
func blockingFetchAndCacheRect(url string) image.Image {
	rectImageCacheMutex.Lock()
	img, found := rectImageCache[url]
	rectImageCacheMutex.Unlock()
	if found {
		return img
	}

	rectUrlFetchOnceMutex.Lock()
	once, ok := rectUrlFetchOnce[url]
	if !ok {
		once = &sync.Once{}
		rectUrlFetchOnce[url] = once
	}
	rectUrlFetchOnceMutex.Unlock()

	var fetchedImg image.Image
	once.Do(func() {
		log.Debug().Msgf("Fetching image for Rect from URL: %s", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to HTTP GET image for Rect from URL: %s", url)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error().Msgf("Failed to download image for Rect from URL %s: status code %d", url, resp.StatusCode)
			return
		}

		decodedImg, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to decode image for Rect from URL: %s", url)
			return
		}
		fetchedImg = decodedImg
		rectImageCacheMutex.Lock()
		rectImageCache[url] = fetchedImg
		rectImageCacheMutex.Unlock()
	})

	if fetchedImg != nil {
		return fetchedImg
	}

	// Re-check cache in case another goroutine populated it
	rectImageCacheMutex.Lock()
	img, found = rectImageCache[url]
	rectImageCacheMutex.Unlock()
	if found {
		return img
	}
	return nil // Fetch failed
}

func Rect(opts RectOpts) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		w_dp := opts.W
		h_dp := opts.H

		isFixedWidth := w_dp > 0
		isFixedHeight := h_dp > 0

		// Resolve final dimensions for the rectangle
		finalSizePx := gtx.Constraints.Min // Start with minimum constraints

		if isFixedWidth {
			finalSizePx.X = gtx.Dp(w_dp)
		}
		if isFixedHeight {
			finalSizePx.Y = gtx.Dp(h_dp)
		}

		// Ensure minimums from opts are met if size is not fixed
		min_w_px := gtx.Dp(opts.MinW)
		min_h_px := gtx.Dp(opts.MinH)

		if !isFixedWidth && finalSizePx.X < min_w_px {
			finalSizePx.X = min_w_px
		}
		if !isFixedHeight && finalSizePx.Y < min_h_px {
			finalSizePx.Y = min_h_px
		}

		// Constrain by max available if not fixed
		if !isFixedWidth && finalSizePx.X > gtx.Constraints.Max.X {
			finalSizePx.X = gtx.Constraints.Max.X
		}
		if !isFixedHeight && finalSizePx.Y > gtx.Constraints.Max.Y {
			finalSizePx.Y = gtx.Constraints.Max.Y
		}

		var imageToRender image.Image
		if opts.ImgURL != "" {
			imageToRender = blockingFetchAndCacheRect(opts.ImgURL)
		} else if opts.Img != nil {
			imageToRender = opts.Img
		}

		if imageToRender != nil {
			// Image rendering logic
			imgOp := paint.NewImageOp(imageToRender)
			imgOp.Filter = paint.FilterLinear // Or paint.FilterNearest for pixel art

			// Clip to the rectangle's bounds
			rectArea := clip.Rect{Max: finalSizePx}.Push(gtx.Ops)

			imgBounds := imageToRender.Bounds()
			if imgBounds.Dx() > 0 && imgBounds.Dy() > 0 {
				imgWidthF := float32(imgBounds.Dx())
				imgHeightF := float32(imgBounds.Dy())
				targetWidthF := float32(finalSizePx.X)
				targetHeightF := float32(finalSizePx.Y)

				// Scale to cover the rectangle area, maintaining aspect ratio
				// This is similar to CSS 'background-size: cover'
				scaleXVal := targetWidthF / imgWidthF
				scaleYVal := targetHeightF / imgHeightF
				finalScale := float32(math.Max(float64(scaleXVal), float64(scaleYVal)))

				scaledWidth := imgWidthF * finalScale
				scaledHeight := imgHeightF * finalScale

				// Center the image
				dx := (targetWidthF - scaledWidth) / 2
				dy := (targetHeightF - scaledHeight) / 2

				transform := f32.Affine2D{}.
					Scale(f32.Pt(0, 0), f32.Pt(finalScale, finalScale)).
					Offset(f32.Pt(dx, dy))

				stack := op.Affine(transform).Push(gtx.Ops)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops) // Ensure the paint operation is added
				stack.Pop()
			} else {
				// Fallback to color if image is invalid but was specified
				log.Warn().Msg("Image specified for Rect has invalid dimensions. Drawing color fallback.")
				paint.FillShape(gtx.Ops, opts.Color, clip.Rect{Max: finalSizePx}.Op())
			}
			rectArea.Pop()

		} else {
			// Fallback to drawing a colored rectangle
			paint.FillShape(gtx.Ops, opts.Color, clip.Rect{Max: finalSizePx}.Op())
		}

		return layout.Dimensions{Size: finalSizePx}
	}
}
