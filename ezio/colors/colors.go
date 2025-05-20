package colors

import "image/color"

type ColorType uint8

const (
	Red ColorType = iota
	Green
	Blue
	Yellow
	White
	Black
	Gray
	Orange
	Purple
	Pink

	// Pastel colors
	LightRed
	LightGreen
	LightBlue
	LightYellow
	LightWhite
	LightBlack
	LightGray
	LightOrange
	LightPurple
	LightPink

	// Special backgrounds
	DarkBackground
)

// Get Light NRGBA color enum from uint8 value
func GetLightNRGBA(value int) color.NRGBA {
	switch value {
	case 0:
		return LightRed.NRGBA()
	case 1:
		return LightGreen.NRGBA()
	case 2:
		return LightBlue.NRGBA()
	case 3:
		return LightYellow.NRGBA()
	case 4:
		return LightWhite.NRGBA()
	case 5:
		return LightBlack.NRGBA()
	case 6:
		return LightGray.NRGBA()
	case 7:
		return LightOrange.NRGBA()
	case 8:
		return LightPurple.NRGBA()
	case 9:
		return LightPink.NRGBA()
	default:
		return Red.NRGBA()
	}
}

func GetBrightNRGBA(value int) color.NRGBA {
	switch value {
	case 0:
		return Red.NRGBA()
	case 1:
		return Green.NRGBA()
	case 2:
		return Blue.NRGBA()
	case 3:
		return Yellow.NRGBA()
	case 4:
		return White.NRGBA()
	case 5:
		return Black.NRGBA()
	case 6:
		return Gray.NRGBA()
	case 7:
		return Orange.NRGBA()
	case 8:
		return Purple.NRGBA()
	case 9:
		return Pink.NRGBA()
	default:
		return LightRed.NRGBA()
	}
}

// NRGBA returns the color.NRGBA representation for this color
func (c ColorType) NRGBA() color.NRGBA {
	switch c {
	case Red:
		return color.NRGBA{R: uint8(255), G: uint8(0), B: uint8(0), A: uint8(255)}
	case Green:
		return color.NRGBA{R: uint8(0), G: uint8(255), B: uint8(0), A: uint8(255)}
	case Blue:
		return color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(255), A: uint8(255)}
	case Yellow:
		return color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(0), A: uint8(255)}
	case White:
		return color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}
	case Black:
		return color.NRGBA{R: uint8(0), G: uint8(0), B: uint8(0), A: uint8(255)}
	case Gray:
		return color.NRGBA{R: uint8(128), G: uint8(128), B: uint8(128), A: uint8(255)}
	case Pink:
		return color.NRGBA{R: uint8(255), G: uint8(192), B: uint8(203), A: uint8(255)}
	case Orange:
		return color.NRGBA{R: uint8(255), G: uint8(165), B: uint8(0), A: uint8(255)}
	case Purple:
		return color.NRGBA{R: uint8(128), G: uint8(0), B: uint8(128), A: uint8(255)}
	case LightRed:
		return color.NRGBA{R: uint8(255), G: uint8(182), B: uint8(193), A: uint8(255)}
	case LightGreen:
		return color.NRGBA{R: uint8(182), G: uint8(255), B: uint8(193), A: uint8(255)}
	case LightBlue:
		return color.NRGBA{R: uint8(182), G: uint8(193), B: uint8(255), A: uint8(255)}
	case LightYellow:
		return color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(182), A: uint8(255)}
	case LightWhite:
		return color.NRGBA{R: uint8(255), G: uint8(255), B: uint8(255), A: uint8(255)}
	case LightBlack:
		return color.NRGBA{R: uint8(182), G: uint8(182), B: uint8(182), A: uint8(255)}
	case LightGray:
		return color.NRGBA{R: uint8(193), G: uint8(193), B: uint8(193), A: uint8(255)}
	case LightOrange:
		return color.NRGBA{R: uint8(255), G: uint8(204), B: uint8(182), A: uint8(255)}
	case LightPurple:
		return color.NRGBA{R: uint8(193), G: uint8(182), B: uint8(255), A: uint8(255)}
	case LightPink:
		return color.NRGBA{R: uint8(255), G: uint8(182), B: uint8(255), A: uint8(255)}
	case DarkBackground:
		return color.NRGBA{R: uint8(37), G: uint8(35), B: uint8(49), A: uint8(255)}
	default:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	}
}
