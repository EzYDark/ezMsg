package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/libs/gui"
	"github.com/rs/zerolog/log"
)

type TextOpts struct {
	ThemePtr  *material.Theme
	TextState *richtext.InteractiveText
}

type ResponsiveOpts struct {
	Enabled   bool
	Alignment layout.Axis
	Scale     float32
}

type Span struct {
	SpanStyle  richtext.SpanStyle
	Responsive ResponsiveOpts
}

func Text(opts TextOpts, spans ...Span) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Use to detect mouse button PRESS event
		// and not only button RELEASE event
		// richtext.LongPressDuration = 0

		for {
			span, event, ok := opts.TextState.Update(gtx)
			if !ok {
				break
			}
			content, _ := span.Content()
			switch event.Type {
			case richtext.Click:
				log.Debug().Msgf("RichText Clicked: %s", event.ClickData.Kind)
				// if event.ClickData.Kind == gesture.KindClick {
				// 	interactColorIndex++
				// 	gtx.Execute(op.InvalidateCmd{})
				// }
			case richtext.Hover:
				log.Debug().Msgf("Hovered: %s", content)
			case richtext.Unhover:
				log.Debug().Msgf("Unhovered: %s", content)
			case richtext.LongPress:
				log.Debug().Msgf("Long-pressed: %s", content)
			}
		}

		var finished_spans []richtext.SpanStyle

		for _, span := range spans {
			var maxConstraint float32
			if span.Responsive.Enabled {
				switch span.Responsive.Alignment {
				case gui.Horizontal:
					maxConstraint = float32(gtx.Constraints.Max.X)
				case gui.Vertical:
					maxConstraint = float32(gtx.Constraints.Max.Y)
				default:
					log.Fatal().Msg("Invalid alignment of span")
				}

				var scaledSize unit.Sp
				if span.Responsive.Scale == 0 {
					// If scale is 0 or not set, use max size of parent
					switch span.Responsive.Alignment {
					case gui.Horizontal:
						// Reasonable default size for horizontal alignment
						scaledSize = unit.Sp(maxConstraint * 0.05)
					case gui.Vertical:
						scaledSize = unit.Sp(maxConstraint * 0.1)
					default:
						log.Fatal().Msg("Invalid alignment of span")
					}
				} else {
					scaledSize = unit.Sp(maxConstraint * span.Responsive.Scale)
				}

				new_span := span.SpanStyle
				new_span.Size = scaledSize

				finished_spans = append(finished_spans, new_span)
			} else {
				finished_spans = append(finished_spans, span.SpanStyle)
			}
		}

		return richtext.Text(opts.TextState, opts.ThemePtr.Shaper, finished_spans...).Layout(gtx)
	}
}

// Helper function to create a text span with custom properties
func TextSpan(span_style richtext.SpanStyle, responsive ResponsiveOpts) Span {
	return Span{
		SpanStyle:  span_style,
		Responsive: responsive,
	}
}

// var interactColorIndex = 0

// func TextCustom(opts TextOpts, str string) layout.Widget {
// 	interactColors := []color.NRGBA{colors.LightBlack, colors.Green, colors.Blue, colors.Red}

// 	return func(gtx layout.Context) layout.Dimensions {
// 		maxConstraint := float32(gtx.Constraints.Max.Y)
// 		scaledSize := unit.Sp(maxConstraint * 1)

// 		var spans = []richtext.SpanStyle{
// 			{
// 				Content: "Hello ",
// 				Color:   colors.Black,
// 				Size:    scaledSize,
// 				Font:    gui.Fonts[0].Font,
// 			},
// 			{
// 				Content: "in ",
// 				Color:   colors.Green,
// 				Size:    unit.Sp(36),
// 				Font:    gui.Fonts[0].Font,
// 			},
// 			{
// 				Content: "rich ",
// 				Color:   colors.Blue,
// 				Size:    unit.Sp(30),
// 				Font:    gui.Fonts[0].Font,
// 			},
// 			{
// 				Content: "text\n",
// 				Color:   colors.Red,
// 				Size:    unit.Sp(40),
// 				Font:    gui.Fonts[0].Font,
// 			},
// 			{
// 				Content:     "Interact with me!",
// 				Color:       interactColors[interactColorIndex%len(interactColors)],
// 				Size:        unit.Sp(40),
// 				Font:        gui.Fonts[0].Font,
// 				Interactive: true,
// 			},
// 		}

// 		for {
// 			span, event, ok := opts.TextState.Update(gtx)
// 			if !ok {
// 				break
// 			}
// 			content, _ := span.Content()
// 			switch event.Type {
// 			case richtext.Click:
// 				log.Debug().Msgf("RichText Clicked: %s", event.ClickData.Kind)
// 				if event.ClickData.Kind == gesture.KindClick {
// 					interactColorIndex++
// 					gtx.Execute(op.InvalidateCmd{})
// 				}
// 			case richtext.Hover:
// 				log.Debug().Msgf("Hovered: %s", content)
// 			case richtext.Unhover:
// 				log.Debug().Msgf("Unhovered: %s", content)
// 			case richtext.LongPress:
// 				log.Debug().Msgf("Long-pressed: %s", content)
// 			}
// 		}

// 		// render the rich text into the operation list
// 		return richtext.Text(opts.InteractiveText, opts.ThemePtr.Shaper, spans...).Layout(gtx)
// 	}
// }

func H1(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.H1(opts.ThemePtr, text).Layout(gtx)
	}
}

func H2(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.H2(opts.ThemePtr, text).Layout(gtx)
	}
}

func H3(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.H3(opts.ThemePtr, text).Layout(gtx)
	}
}

func H4(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.H4(opts.ThemePtr, text).Layout(gtx)
	}
}

func H5(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.H5(opts.ThemePtr, text).Layout(gtx)
	}
}

func H6(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.H6(opts.ThemePtr, text).Layout(gtx)
	}
}

func P(opts TextOpts, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if opts.ThemePtr == nil {
			log.Fatal().Msg("opts.ThemePtr must be initialized!")
		}

		return material.Body1(opts.ThemePtr, text).Layout(gtx)
	}
}
