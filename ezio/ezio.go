package ezio

import (
	"github.com/ezydark/ezMsg/ezio/components"
	"github.com/ezydark/ezMsg/ezio/layouts"
	"github.com/ezydark/ezMsg/ezio/layouts/flex"
)

// --- Layouts ---

// Re-export: Flex layout
type (
	FlexOpts            = flex.FlexOpts
	FlexChildOpts       = flex.FlexChildOpts
	FlexChildStaticOpts = flex.FlexChildStaticOpts
)

const (
	Horizontal = flex.Horizontal
	Vertical   = flex.Vertical
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
