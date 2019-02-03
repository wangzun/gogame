package main

import (
	"math"

	"github.com/wangzun/gogame/engine/geometry"
	"github.com/wangzun/gogame/engine/graphic"
	"github.com/wangzun/gogame/engine/light"
	"github.com/wangzun/gogame/engine/material"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/util/application"
	"github.com/wangzun/gogame/example"
)

type PointLight struct {
	vl    *example.PointLightMesh
	hl    *example.PointLightMesh
	count float64
}

func main() {
	p := &PointLight{}
	p.Init()
}

func (t *PointLight) Init() {

	a, _ := application.Create(application.Options{Control: true})

	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.5)
	a.Scene().Add(ambientLight)

	geom1 := geometry.NewSphere(0.5, 32, 32, 0, math.Pi*2, 0, math.Pi)
	mat1 := material.NewStandard(&math32.Color{0, 0, 0.6})
	sphere1 := graphic.NewMesh(geom1, mat1)
	sphere1.SetPositionX(1)
	a.Scene().Add(sphere1)

	// Creates left sphere
	geom2 := geometry.NewSphere(0.5, 32, 32, 0, math.Pi*2, 0, math.Pi)
	mat2 := material.NewPhong(&math32.Color{0, 0.5, 0.0})
	sphere2 := graphic.NewMesh(geom2, mat2)
	sphere2.SetPositionX(-1)
	a.Scene().Add(sphere2)

	// // Creates left plane
	geom3 := geometry.NewPlane(4, 4, 8, 8)
	mat3 := material.NewStandard(&math32.Color{1, 1, 1})
	pleft := graphic.NewMesh(geom3, mat3)
	pleft.SetPosition(-2, 0, 0)
	pleft.SetRotationY(math.Pi / 2)
	a.Scene().Add(pleft)

	// Creates right plane
	geom4 := geometry.NewPlane(4, 4, 8, 8)
	mat4 := material.NewStandard(&math32.Color{1, 1, 1})
	pright := graphic.NewMesh(geom4, mat4)
	pright.SetPosition(2, 0, 0)
	pright.SetRotationY(-math.Pi / 2)
	a.Scene().Add(pright)

	// // Creates top plane
	geom5 := geometry.NewPlane(4, 4, 8, 8)
	mat5 := material.NewStandard(&math32.Color{1, 1, 1})
	ptop := graphic.NewMesh(geom5, mat5)
	ptop.SetPosition(0, 2, 0)
	ptop.SetRotationX(math.Pi / 2)
	a.Scene().Add(ptop)

	// // Creates bottom plane
	geom6 := geometry.NewPlane(4, 4, 8, 8)
	mat6 := material.NewStandard(&math32.Color{1, 1, 1})
	pbot := graphic.NewMesh(geom6, mat6)
	pbot.SetPosition(0, -2, 0)
	pbot.SetRotationX(-math.Pi / 2)
	a.Scene().Add(pbot)

	// // Creates back plane
	geom7 := geometry.NewPlane(4, 4, 8, 8)
	mat7 := material.NewStandard(&math32.Color{1, 1, 1})
	pback := graphic.NewMesh(geom7, mat7)
	pback.SetPosition(0, 0, -2)
	a.Scene().Add(pback)

	axis := graphic.NewAxisHelper(1)
	a.Scene().Add(axis)

	// Creates vertical point light
	t.vl = example.NewPointLightMesh(&math32.Color{1, 1, 1})
	a.Scene().Add(t.vl.Mesh)

	// Creates horizontal point light
	t.hl = example.NewPointLightMesh(&math32.Color{1, 1, 1})
	a.Scene().Add(t.hl.Mesh)
	a.Scene().SetScale(0.5, 0.5, 0.5)

	a.CameraPersp().SetPosition(0, 0, 3)
	a.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})

	a.Subscribe(application.OnBeforeRender, func(evname string, ev interface{}) {
		t.vl.SetPosition(0, 1.5*float32(math.Sin(t.count)), 0)
		t.hl.SetPosition(1.5*float32(math.Sin(t.count)), 1, 0)
		t.count += 0.02

	})

	a.Run()

}
