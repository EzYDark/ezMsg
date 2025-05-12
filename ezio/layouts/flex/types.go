package flex

import (
	"gioui.org/layout"
)

type FlexOpts struct {
	Alignment layout.Axis
}

type FlexChildOpts struct {
	Weight     float32
	W, H       int
	MinW, MinH int
	Widgets    []layout.Widget
}

type FlexChildStaticOpts struct {
	W, H       int
	MinW, MinH int
}
