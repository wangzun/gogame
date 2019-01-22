package main

import (
	"log"

	"github.com/g3n/engine/graphic"

	"github.com/wangzun/gogame/engine/geometry"
	"github.com/wangzun/gogame/engine/graphic"
	"github.com/wangzun/gogame/engine/texture"
	"github.com/wangzun/gogame/engine/util/application"
)

// func main() {
// 	start()
// appInfo, _ := application.Create(application.Options{})
// // Create a blue torus and add it to the scene
//  geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
//  mat := material.NewPhong(math32.NewColor("DarkBlue"))
//  torusMesh := graphic.NewMesh(geom, mat)
// app.Scene().Add(torusMesh)

// // Add lights to the scene
// ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
// app.Scene().Add(ambientLight)
// pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
// pointLight.SetPosition(1, 0, 2)
// app.Scene().Add(pointLight)

// // Add an axis helper to the scene
// axis := graphic.NewAxisHelper(5)
// app.Scene().Add(axis)

// app.CameraPersp().SetPosition(0, 0, 3)
// appInfo.Run()
// }

func main() {
	xxx := texture.Animator{}
	log.Println(xxx)
	yyy := geometry.Geometry{}
	zzz := graphic.Graphic{}
	app, _ := application.Create(application.Options{})
	// geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	app.Run()
}
