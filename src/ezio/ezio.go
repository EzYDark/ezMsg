package ezio

import (
	"gioui.org/layout"
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/src/ezio/colors"
	"github.com/ezydark/ezMsg/src/ezio/components"
	"github.com/ezydark/ezMsg/src/ezio/layouts"
	"github.com/ezydark/ezMsg/src/ezio/layouts/flex"
)

// --- Colors ---
var (
	Red    = colors.Red
	Green  = colors.Green
	Blue   = colors.Blue
	Yellow = colors.Yellow
	White  = colors.White
	Black  = colors.Black
	Gray   = colors.Gray
	Orange = colors.Orange
	Purple = colors.Purple

	LightRed    = colors.LightRed
	LightGreen  = colors.LightGreen
	LightBlue   = colors.LightBlue
	LightYellow = colors.LightYellow
	LightWhite  = colors.LightWhite
	LightBlack  = colors.LightBlack
	LightGray   = colors.LightGray
	LightOrange = colors.LightOrange
	LightPurple = colors.LightPurple
	LightPink   = colors.LightPink
)

var (
	DarkBackground = colors.DarkBackground
)

// --- Layouts ---

// Re-export: Flex layout
type (
	FlexBoxOpts   = flex.FlexBoxOpts
	FlexChildOpts = flex.FlexChildOpts
)

// Re-defined from layout.Axis lib
const (
	Horizontal layout.Axis = iota
	Vertical
)

// Re-defined from layout.Alignment lib
const (
	Start layout.Alignment = iota
	End
	Middle
	Baseline
)

// Re-defined from layout.Spacing lib
const (
	// SpaceEnd leaves space at the end.
	SpaceEnd layout.Spacing = iota
	// SpaceStart leaves space at the start.
	SpaceStart
	// SpaceSides shares space between the start and end.
	SpaceSides
	// SpaceAround distributes space evenly between children,
	// with half as much space at the start and end.
	SpaceAround
	// SpaceBetween distributes space evenly between children,
	// leaving no space at the start and end.
	SpaceBetween
	// SpaceEvenly distributes space evenly between children and
	// at the start and end.
	SpaceEvenly
)

var (
	NW     = layout.NW
	N      = layout.N
	NE     = layout.NE
	E      = layout.E
	SE     = layout.SE
	S      = layout.S
	SW     = layout.SW
	W      = layout.W
	Center = layout.Center
)

var (
	FlexBox   = flex.FlexBox
	FlexChild = flex.FlexChild
)

// Re-export Stack layout
type StackBoxOpts = layouts.StackBoxOpts

var (
	StackBox      = layouts.StackBox
	StackedChild  = layouts.StackedChild
	ExpandedChild = layouts.ExpandedChild
)

// Re-export Background layout
var BackgroundBox = layouts.BackgroundBox
var BackgroundStackBox = layouts.BackgroundStackBox

// Re-export: List layout
type ListOpts = layouts.ListOpts

var (
	ListBox   = layouts.ListBox
	ListChild = layouts.ListChild
)

// Re-export: Direction layout
type AlignOpts = layouts.AlignOpts

var (
	Align = layouts.Align
)

// Re-export: Margin layout
type MarginOpts = layouts.MarginOpts

var Margin = layouts.Margin

// --- Components ---

// Re-export: Button component
type ButtonOpts = components.ButtonOpts

var Button = components.Button

// Re-export: ContextMenu component
type ContextMenuOpts = components.ContextMenuOpts

var ContextMenu = components.ContextMenu

// Re-export: Input component
type InputOpts = components.InputOpts

var Input = components.Input

// Re-export: Rect component
type RectOpts = components.RectOpts

var Rect = components.Rect

// Re-export: Circle component
type CircleOpts = components.CircleOpts

var Circle = components.Circle
var LoadImageHelper = components.LoadImageHelper

// Re-export: Text component
type TextOpts = components.TextOpts
type SpanStyle = richtext.SpanStyle

var Text = components.Text
var TextSpan = components.TextSpan
var H1 = components.H1
var H2 = components.H2
var H3 = components.H3
var H4 = components.H4
var H5 = components.H5
var H6 = components.H6
var P = components.P
