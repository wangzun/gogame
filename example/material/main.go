package main

import (
	"github.com/wangzun/gogame/engine/geometry"
	"github.com/wangzun/gogame/engine/gls"
	"github.com/wangzun/gogame/engine/graphic"
	"github.com/wangzun/gogame/engine/light"
	"github.com/wangzun/gogame/engine/loader/obj"
	"github.com/wangzun/gogame/engine/material"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/texture"
	"github.com/wangzun/gogame/engine/util/application"
	"github.com/wangzun/gogame/example"
)

type Blending struct {
	texbg *texture.Texture2D
}

func main() {
	// t := &Blending{}
	// t.Initialize()
	p := &PhysicalHelmet{}
	p.Initialize()

}

func (t *Blending) Initialize() {
	a, _ := application.Create(application.Options{Control: true})

	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 2)
	a.Scene().Add(ambientLight)

	// Creates checker board textures for background
	// c1 := &math32.Color{0.7, 0.7, 0.7}
	// c2 := &math32.Color{0.3, 0.3, 0.3}
	t.texbg, _ = texture.NewTexture2DFromImage("assets/images/wall1.jpg")
	// t.texbg = texture.NewBoard(16, 16, c1, c2, c2, c1, 1)
	t.texbg.SetWrapS(gls.REPEAT)
	t.texbg.SetWrapT(gls.REPEAT)
	t.texbg.SetRepeat(64, 64)

	// Creates background plane
	matbg := material.NewPhong(&math32.Color{1, 1, 1})
	matbg.SetPolygonOffset(1, 1)
	matbg.AddTexture(t.texbg)
	geombg := geometry.NewPlane(4000, 3000, 1, 1)
	// geombg := geometry.NewPlane(40, 30, 1, 1)
	meshbg := graphic.NewMesh(geombg, matbg)
	meshbg.SetPosition(0, 0, -1)
	a.Scene().Add(meshbg)

	// Builds list of textures
	texnames := []string{
		"uvgrid.jpg", "sprite0.jpg",
		"sprite0.png", "lensflare0.png",
		"lensflare0_alpha.png",
	}
	texlist := []*texture.Texture2D{}
	for _, tname := range texnames {
		tex, err := texture.NewTexture2DFromImage("assets/" + "/images/" + tname)
		if err != nil {
			a.Log().Fatal("Error loading texture: %s", err)
		}
		texlist = append(texlist, tex)
	}

	blendings := []struct {
		blending string
		value    material.Blending
	}{
		{"NoBlending", material.BlendingNone},
		{"NormalBlending", material.BlendingNormal},
		{"AdditiveBlending", material.BlendingAdditive},
		{"SubtractiveBlending", material.BlendingSubtractive},
		{"MultiplyBlending", material.BlendingMultiply},
	}

	// This geometry will be shared by several meshes
	// For each mesh which uses this geometry we need to increment its refcount
	geo1 := geometry.NewPlane(100, 100, 1, 1)
	defer geo1.Dispose()

	// Internal function go generate a row of images
	var addImageRow = func(tex *texture.Texture2D, y int) {
		for i := 0; i < len(blendings); i++ {
			material := material.NewPhong(&math32.Color{1, 1, 1})
			material.SetOpacity(1)
			material.SetTransparent(true)
			material.AddTexture(tex)
			material.SetBlending(blendings[i].value)
			x := (float32(i) - float32(len(blendings))/2) * 110
			mesh := graphic.NewMesh(geo1.Incref(), material)
			mesh.SetPosition(x, float32(y), 0)
			a.Scene().Add(mesh)
		}
	}

	addImageRow(texlist[0], 300)
	addImageRow(texlist[1], 150)
	addImageRow(texlist[2], 0)
	addImageRow(texlist[3], -150)
	addImageRow(texlist[4], -300)

	a.Scene().SetScale(0.5, 0.5, 0.5)

	a.CameraPersp().SetPositionZ(600)
	// a.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})
	// a.CameraPersp().SetPosition(0, 0, 3)
	// a.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})

	a.Run()
}

type PhysicalHelmet struct {
	p1    *example.PointLightMesh
	p2    *example.PointLightMesh
	s1    *example.SpotLightMesh
	s2    *example.SpotLightMesh
	count float32
}

func DirData() string {
	return "assets"
}

func (t *PhysicalHelmet) Initialize() {

	a, _ := application.Create(application.Options{Control: true})
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 2)
	a.Scene().Add(ambientLight)

	// Add directional red light from right
	l1 := light.NewDirectional(&math32.Color{1, 0, 0}, 1.0)
	l1.SetPosition(0.1, 0, 0)
	a.Scene().Add(l1)

	// Add directional green light from top
	l2 := light.NewDirectional(&math32.Color{0, 1, 0}, 1.0)
	l2.SetPosition(0, 0.1, 0)
	a.Scene().Add(l2)

	// Add directional blue light from front
	l3 := light.NewDirectional(&math32.Color{0, 0, 1}, 1.0)
	l3.SetPosition(0, 0, 0.1)
	a.Scene().Add(l3)

	t.p1 = example.NewPointLightMesh(&math32.Color{1, 1, 1})
	a.Scene().Add(t.p1)

	t.p2 = example.NewPointLightMesh(&math32.Color{1, 0, 0})
	a.Scene().Add(t.p2)

	t.s1 = example.NewSpotLightMesh(&math32.Color{0, 0, 1})
	//t.s1.SetPosition(0.5, 1, 0.5)
	//a.Scene().Add(t.s1)

	t.s2 = example.NewSpotLightMesh(&math32.Color{1, 0, 0})
	//t.s2.SetPosition(-1, 0, 0)
	//t.s2.SetRotationZ(math.Pi/2)
	//a.Scene().Add(t.s2)

	// DamagedHelmet

	// Decodes obj file and associated mtl file
	dec, err := obj.Decode(DirData()+"/obj/DamagedHelmet.obj", "")
	if err != nil {
		panic(err)
	}

	// Creates a new node with all the objects in the decoded file and adds it to the scene
	geom, err := dec.NewGeometry(&dec.Objects[0])
	if err != nil {
		panic(err)
	}

	// Helper function to load texture and handle errors
	newTexture := func(path string) *texture.Texture2D {
		tex, err := texture.NewTexture2DFromImage(path)
		if err != nil {
			a.Log().Fatal("Error loading texture: %s", err)
		}
		tex.SetWrapS(gls.MIRRORED_REPEAT)
		tex.SetWrapT(gls.MIRRORED_REPEAT)
		return tex
	}

	pbrMat := material.NewPhysical()
	pbrMat.SetEmissiveFactor(math32.NewColor("white"))
	pbrMat.SetBaseColorMap(newTexture(DirData() + "/obj/DamagedHelmet_albedo.jpg"))
	pbrMat.SetMetallicRoughnessMap(newTexture(DirData() + "/obj/DamagedHelmet_metalRoughness.jpg"))
	pbrMat.SetNormalMap(newTexture(DirData() + "/obj/DamagedHelmet_normal.jpg"))
	pbrMat.SetEmissiveMap(newTexture(DirData() + "/obj/DamagedHelmet_emissive.jpg"))
	pbrMat.SetOcclusionMap(newTexture(DirData() + "/obj/DamagedHelmet_AO.jpg"))

	helmet := graphic.NewMesh(geom, pbrMat)
	a.Scene().Add(helmet)
	a.CameraPersp().SetPosition(0, 0, 3)
	a.CameraPersp().LookAt(&math32.Vector3{0, 0, 0})
	a.Subscribe(application.OnBeforeRender, func(evname string, ev interface{}) {
		t.Render()
	})
	a.Run()

}

func (t *PhysicalHelmet) Render() {

	t.p1.SetPosition(float32(math32.Cos(t.count)), float32(math32.Sin(t.count)), 0)
	t.p2.SetPosition(0, 1.5*float32(math32.Sin(t.count)), 1.5*float32(math32.Cos(t.count)))

	t.s1.SetPosition(0, 1.5*float32(math32.Sin(t.count*1.5)), 1.5*float32(math32.Cos(t.count*1.5)))
	t.s2.SetPosition(1.5*float32(math32.Cos(t.count*1.5)), 1.5*float32(math32.Sin(t.count*1.5)), 0)

	t.count += 0.01
}
