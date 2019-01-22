// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphic

import (
	"github.com/wangzun/gogame/engine/core"
	"github.com/wangzun/gogame/engine/geometry"
	"github.com/wangzun/gogame/engine/gls"
	"github.com/wangzun/gogame/engine/material"
	"golang.org/x/mobile/gl"
)

// LineStrip is a Graphic which is rendered as a collection of connected lines.
type LineStrip struct {
	Graphic             // Embedded graphic object
	uniMVPm gls.Uniform // Model view projection matrix uniform location cache
}

// NewLineStrip creates and returns a pointer to a new LineStrip graphic
// with the specified geometry and material.
func NewLineStrip(igeom geometry.IGeometry, imat material.IMaterial) *LineStrip {

	l := new(LineStrip)
	l.Graphic.Init(igeom, gls.LINE_STRIP)
	l.AddMaterial(l, imat, 0, 0)
	l.uniMVPm.Init("MVP")
	return l
}

// RenderSetup is called by the engine before drawing this geometry.
func (l *LineStrip) RenderSetup(gs *gls.GLS, rinfo *core.RenderInfo) {

	// Transfer model view projection matrix uniform
	mvpm := l.ModelViewProjectionMatrix()
	location := l.uniMVPm.Location(gs)
	gs.UniformMatrix4fv(gl.Uniform{Value: location}, 1, false, mvpm[:])
}

// Raycast satisfies the INode interface and checks the intersections
// of this geometry with the specified raycaster.
func (l *LineStrip) Raycast(rc *core.Raycaster, intersects *[]core.Intersect) {

	lineRaycast(l, rc, intersects, 1)
}
