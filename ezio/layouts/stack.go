package layouts

import "gioui.org/layout"

// StackOpts defines the options for the Stack layout.
type StackOpts struct {
	// Alignment specifies how children are aligned within the stack.
	Alignment layout.Direction
}

// StackBox creates a layout widget that stacks its children on top of each other.
// The size of the stack is determined by the largest child.
func StackBox(opts StackOpts, children ...layout.StackChild) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Use the standard Gio layout.Stack with the provided options and children.
		return layout.Stack{Alignment: opts.Alignment}.Layout(gtx, children...)
	}
}

// StackedChild creates a layout.StackChild that is laid out at its preferred size.
func StackedChild(child layout.Widget) layout.StackChild {
	return layout.Stacked(child)
}

// ExpandedChild creates a layout.StackChild that is expanded to fill the stack's bounds.
func ExpandedChild(child layout.Widget) layout.StackChild {
	return layout.Expanded(child)
}
