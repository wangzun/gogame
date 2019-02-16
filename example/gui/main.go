package main

import (
	"fmt"

	"github.com/wangzun/gogame/engine/gui"
	"github.com/wangzun/gogame/engine/light"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/util/application"
)

func DirData() string {
	return "assets"
}

func main() {
	app, _ := application.Create(application.Options{Control: true})

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)

	// geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	// mat := material.NewPhong(math32.NewColor("darkblue"))
	// torusMesh := graphic.NewMesh(geom, mat)
	// torusMesh.SetScale(0.5, 0.5, 0.5)
	// torusMesh.SetPosition(0, 0, 0)
	// // torusMesh.SetRotation(0, 1, 0)
	// app.Scene().Add(torusMesh)

	// b1, err := gui.NewImage(DirData() + "/images/tiger1.jpg")

	// b1 := gui.NewImageLabel("jjjjjjjjjjjjjjj")

	b1, err := gui.NewImageButton(DirData() + "/images/tiger1.jpg")
	if err != nil {
		panic(err)
	}
	b1.SetPosition(2, 2)
	// b1.SetScaleX(0.00003)
	// b1.SetScaleY(0.00003)
	fmt.Println("w h : ", b1.TotalWidth(), b1.TotalHeight())

	b1.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fmt.Println("click image button !!!!!!!!!!!!!!")
	})

	// app.Scene().Add(b1)
	app.Gui().Add(b1)

	// app.CameraPersp().SetPosition(0, 0, 3)
	// app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})

	// app.Subscribe(application.OnBeforeRender, func(evname string, ev interface{}) {
	// 	b1.SetChanged(true)
	// })

	app.Run()

}
