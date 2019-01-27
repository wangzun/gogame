package main

import (
	"github.com/wangzun/gogame/engine/animation"
	"github.com/wangzun/gogame/engine/graphic"
	"github.com/wangzun/gogame/engine/light"
	"github.com/wangzun/gogame/engine/loader/gltf"
	"github.com/wangzun/gogame/engine/loader/obj"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/util/application"
)

var anims []*animation.Animation

func main() {

	app, _ := application.Create(application.Options{})

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	app.Scene().Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(5)
	app.Scene().Add(axis)

	gltfjson, err := gltf.ParseJSON("assets/griffin/scene.gltf")
	if err != nil {
		app.Log().Error("gltf json ", err)
	}
	node, err := gltfjson.LoadScene(0)
	if err != nil {
		app.Log().Error("load scene ", err)
	}

	node.GetNode().SetPosition(0.5, 0, -2)
	node.GetNode().SetScale(0.5, 0.5, 0.5)

	app.Scene().Add(node)

	dec, err := obj.Decode("assets/gopher/gopher.obj", "assets/gopher/gopher.mtl")
	if err != nil {
		panic(err.Error())
	}

	// Creates a new node with all the objects in the decoded file and adds it to the scene
	group, err := dec.NewGroup()
	if err != nil {
		panic(err.Error())
	}
	group.GetNode().SetRotationY(-0.5 * 3.14)
	group.GetNode().SetPosition(-0.5, 0, -1)
	group.GetNode().SetScale(0.3, 0.3, 0.3)
	app.Scene().Add(group)

	for i := range gltfjson.Animations {
		anim, _ := gltfjson.LoadAnimation(i)
		anim.SetLoop(true)
		anims = append(anims, anim)
	}

	app.Subscribe(application.OnBeforeRender, func(evname string, ev interface{}) {
		for _, anim := range anims {
			anim.Update(app.FrameDeltaSeconds())
		}
	})

	app.CameraPersp().SetPosition(0, 0, 3)
	app.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})
	app.Run()

}
