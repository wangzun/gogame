// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// Consolidate window events plus GUI events
const (
	OnClick       = "gui.OnClick"       // Widget clicked by mouse left button or key
	OnCursorEnter = "gui.OnCursorEnter" // cursor enters the panel area
	OnCursorLeave = "gui.OnCursorLeave" // cursor leaves the panel area
	OnMouseOut    = "gui.OnMouseOut"    // mouse button pressed outside of the panel
	OnResize      = "gui.OnResize"      // panel size changed (no parameters)
	OnEnable      = "gui.OnEnable"      // panel enabled state changed (no parameters)
	OnChange      = "gui.OnChange"      // onChange is emitted by List, DropDownList, CheckBox and Edit
	OnChild       = "gui.OnChild"       // child added to or removed from panel
	OnRadioGroup  = "gui.OnRadioGroup"  // radio button from a group changed state
	OnRightClick  = "gui.OnRightClick"  // Widget clicked by mouse right button
)
