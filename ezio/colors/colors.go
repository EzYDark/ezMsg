package colors

import "image/color"

var (
	Red    = color.NRGBA{R: uint8(255), G: uint8(0), B: uint8(0), A: uint8(255)}
	Green  = color.NRGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)}
	Blue   = color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(255), A: uint8(255)}
	Yellow = color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(0), A: uint8(255)}
	White  = color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}
	Black  = color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}
	Gray   = color.NRGBA{R: uint8(128), G: uint8(128), B: uint8(128), A: uint8(255)}
	Orange = color.NRGBA{R: uint8(255), G: uint8(165), B: uint8(0), A: uint8(255)}
	Purple = color.NRGBA{R: uint8(128), G: uint8(0), B: uint8(128), A: uint8(255)}

	// Pastel colors
	LightRed    = color.NRGBA{R: uint8(255), G: uint8(182), B: uint8(193), A: uint8(255)}
	LightGreen  = color.NRGBA{R: uint8(182), G: uint8(255), B: uint8(193), A: uint8(255)}
	LightBlue   = color.NRGBA{R: uint8(182), G: uint8(193), B: uint8(255), A: uint8(255)}
	LightYellow = color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(182), A: uint8(255)}
	LightWhite  = color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}
	LightBlack  = color.NRGBA{R: uint8(182), G: uint8(182), B: uint8(182), A: uint8(255)}
	LightGray   = color.NRGBA{R: uint8(193), G: uint8(193), B: uint8(193), A: uint8(255)}
	LightOrange = color.NRGBA{R: uint8(255), G: uint8(204), B: uint8(182), A: uint8(255)}
	LightPurple = color.NRGBA{R: uint8(193), G: uint8(182), B: uint8(255), A: uint8(255)}
	LightPink   = color.NRGBA{R: uint8(255), G: uint8(182), B: uint8(255), A: uint8(255)}
)

var (
	DarkBackground = color.NRGBA{R: uint8(37), G: uint8(35), B: uint8(49), A: uint8(255)}
)
