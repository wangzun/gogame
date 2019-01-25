package main

import (

	// "github.com/g3n/engine/graphic"

	"github.com/wangzun/gogame/engine/geometry"
	"github.com/wangzun/gogame/engine/graphic"
	"github.com/wangzun/gogame/engine/light"
	"github.com/wangzun/gogame/engine/material"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/util/application"
)

func main() {
	app, _ := application.Create(application.Options{})
	// Create a blue torus and add it to the scene
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewPhong(math32.NewColor("darkblue"))
	torusMesh := graphic.NewMesh(geom, mat)
	torusMesh.SetScale(0.5, 0.5, 0.5)
	torusMesh.SetPosition(0, 0, 0)
	// torusMesh.SetRotation(0, 1, 0)
	app.Scene().Add(torusMesh)

	// // Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(0, 0, 3)
	app.Scene().Add(pointLight)

	// // // Add an axis helper to the scene
	axis := graphic.NewAxisHelper(5)
	app.Scene().Add(axis)
	app.Scene().RotateX(1)

	app.CameraPersp().SetPosition(0, 0.5, 3)
	app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})
	app.Run()
}
