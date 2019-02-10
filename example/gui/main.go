package main

import (
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

	b1, err := gui.NewImageButton(DirData() + "/images/tiger1.jpg")
	if err != nil {
		panic(err)
	}
	b1.SetPosition(20, 20)
	app.Gui().Add(b1)

	app.CameraPersp().SetPosition(0, 0, 3)
	app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})

	app.Run()

}
