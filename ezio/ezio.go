package ezio

import (
	"gioui.org/x/richtext"
	"github.com/ezydark/ezMsg/ezio/colors"
	"github.com/ezydark/ezMsg/ezio/components"
	"github.com/ezydark/ezMsg/ezio/layouts"
	"github.com/ezydark/ezMsg/ezio/layouts/flex"
	"github.com/ezydark/ezMsg/libs/gui"
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

// --- Layouts ---

// Re-export: Flex layout
type (
	FlexOpts            = flex.FlexOpts
	FlexChildOpts       = flex.FlexChildOpts
	FlexChildStaticOpts = flex.FlexChildStaticOpts
)

const (
	Horizontal = gui.Horizontal
	Vertical   = gui.Vertical
)

var (
	Flex            = flex.Flex
	FlexChild       = flex.FlexChild
	FlexChildWeight = flex.FlexChildWeight
	FlexChildStatic = flex.FlexChildStatic
)

// Re-export: List layout
type ListOpts = layouts.ListOpts

var (
	List      = layouts.List
	ListChild = layouts.ListChild
)

// Re-export: Direction layout
type DirectionOpts = layouts.DirectionOpts

var Direction = layouts.Direction

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

// Re-export: Text component
type TextOpts = components.TextOpts

type ResponsiveOpts = components.ResponsiveOpts
type Span = components.Span
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
