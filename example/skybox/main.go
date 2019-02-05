package main

import (
	"github.com/wangzun/gogame/engine/graphic"
	"github.com/wangzun/gogame/engine/light"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/util/application"
)

func main() {
	app, _ := application.Create(application.Options{Control: true})

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)

	skybox, err := graphic.NewSkybox(graphic.SkyboxData{
		"assets" + "/images/sanfrancisco/", "jpg",
		[6]string{"posx", "negx", "posy", "negy", "posz", "negz"}})
	if err != nil {
		panic(err)
	}
	app.Scene().Add(skybox)

	// Add axis helper
	axis := graphic.NewAxisHelper(2)
	app.Scene().Add(axis)

	app.CameraPersp().SetPosition(0, 0, 3)
	app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})

	app.Run()

}
