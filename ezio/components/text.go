package components

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
	"github.com/rs/zerolog/log"
)

type TextOpts struct {
	ThemePtr  *material.Theme
	TextState *richtext.InteractiveText
}

func Text(opts TextOpts, spans ...richtext.SpanStyle) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// // Use to detect mouse button PRESS event
		// // and not only button RELEASE event by default
		// richtext.LongPressDuration = 0

		return richtext.Text(opts.TextState, opts.ThemePtr.Shaper, spans...).Layout(gtx)
	}
}

// Helper function to create a text span with custom properties
func TextSpan(span_style richtext.SpanStyle) richtext.SpanStyle {
	return span_style
}

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
