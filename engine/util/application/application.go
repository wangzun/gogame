package application

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"golang.org/x/mobile/gl"

	"github.com/wangzun/gogame/engine/camera"
	"github.com/wangzun/gogame/engine/camera/control"
	"github.com/wangzun/gogame/engine/core"
	"github.com/wangzun/gogame/engine/gls"
	"github.com/wangzun/gogame/engine/gui"
	"github.com/wangzun/gogame/engine/math32"
	"github.com/wangzun/gogame/engine/moblie"
	"github.com/wangzun/gogame/engine/renderer"
	"github.com/wangzun/gogame/engine/util/logger"
)

type Application struct {
	show              bool                  // show on screen
	core.Dispatcher                         // Embedded event dispatcher
	core.TimerManager                       // Embedded timer manager
	gl                *gls.GLS              // OpenGL state
	log               *logger.Logger        // Default application logger
	renderer          *renderer.Renderer    // Renderer object
	camPersp          *camera.Perspective   // Perspective camera
	camOrtho          *camera.Orthographic  // Orthographic camera
	camera            camera.ICamera        // Current camera
	orbit             *control.OrbitControl // Camera orbit controller
	guiroot           *gui.Root             // Gui root panel
	frameRater        *FrameRater           // Render loop frame rater
	scene             *core.Node            // Node container for 3D tests
	frameCount        uint64                // Frame counter
	frameTime         time.Time             // Time at the start of the frame
	frameDelta        time.Duration         // Time delta from previous frame
	startTime         time.Time             // Time at the start of the render loop
	swapInterval      *int                  // Swap interval option
	targetFPS         *uint                 // Target FPS option
	noglErrors        *bool                 // No OpenGL check errors options
	cpuProfile        *string               // File to write cpu profile to
	execTrace         *string               // File to write execution trace data to
	moblie            *moblie.Moblie
	control           bool
}

// Options defines initial options passed to the application creation function
type Options struct {
	LogPrefix   string // Log prefix (default = "")
	LogLevel    int    // Initial log level (default = DEBUG)
	EnableFlags bool   // Enable command line flags (default = false)
	TargetFPS   uint   // Desired frames per second rate (default = 60)
	Control     bool   //
}

// OnBeforeRender is the event generated by Application just before rendering the scene/gui
const OnBeforeRender = "util.application.OnBeforeRender"

// OnAfterRender is the event generated by Application just after rendering the scene/gui
const OnAfterRender = "util.application.OnAfterRender"

// OnQuit is the event generated by Application when the user tries to close the window
// or the Quit() method is called.
const OnQuit = "util.application.OnQuit"

// appInstance contains the pointer to the single Application instance
var appInstance *Application

// Create creates and returns the application object using the specified options.
// This function must be called only once.
func Create(ops Options) (*Application, error) {

	if appInstance != nil {
		return nil, fmt.Errorf("Application already created")
	}
	app := new(Application)
	appInstance = app
	app.Dispatcher.Initialize()
	app.TimerManager.Initialize()

	// Initialize options defaults
	app.control = ops.Control

	app.show = false
	app.swapInterval = new(int)
	app.targetFPS = new(uint)
	app.noglErrors = new(bool)
	app.cpuProfile = new(string)
	app.execTrace = new(string)
	*app.swapInterval = -1
	*app.targetFPS = 60

	// Options parameter overrides some options
	if ops.TargetFPS != 0 {
		*app.targetFPS = ops.TargetFPS
	}

	// Creates flags if requested (override options defaults)
	if ops.EnableFlags {
		app.swapInterval = flag.Int("swapinterval", -1, "Sets the swap buffers interval to this value")
		app.targetFPS = flag.Uint("targetfps", 60, "Sets the frame rate in frames per second")
		app.noglErrors = flag.Bool("noglerrors", false, "Do not check OpenGL errors at each call (may increase FPS)")
		app.cpuProfile = flag.String("cpuprofile", "", "Activate cpu profiling writing profile to the specified file")
		app.execTrace = flag.String("exectrace", "", "Activate execution tracer writing data to the specified file")
	}
	flag.Parse()

	// Creates application logger
	app.log = logger.New(ops.LogPrefix, nil)
	app.log.AddWriter(logger.NewConsole(true))
	go func() {
		nnet, err := logger.NewNet("tcp", "192.168.1.4:6666")
		if err == nil {
			app.log.AddWriter(nnet)
			app.log.Info("init net log succ")
		}
	}()
	app.log.SetFormat(logger.FTIME | logger.FMICROS)
	app.log.SetLevel(ops.LogLevel)

	// app.InitGls(mobileApp.GetContext())

	// Get the window manager
	// Create OpenGL state
	// gl, err := gls.New()
	// if err != nil {
	// 	return nil, err
	// }
	// app.gl = gl
	// // Checks OpenGL errors
	// app.gl.SetCheckErrors(!*app.noglErrors)

	// Logs OpenGL version
	// glVersion := app.Gl().GetString(gls.VERSION)
	// app.log.Info("OpenGL version: %s", glVersion)

	// Clears the screen
	// cc := math32.NewColor("gray")
	// app.gl.ClearColor(cc.R, cc.G, cc.B, 1)
	// app.gl.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

	// Creates perspective camera
	// aspect := float32(width) / float32(height)
	// aspect := float32(750) / float32(1334)
	// app.camPersp = camera.NewPerspective(65, aspect, 0.01, 1000)

	// // Creates orthographic camera
	// app.camOrtho = camera.NewOrthographic(-2, 2, 2, -2, 0.01, 100)
	// app.camOrtho.SetPosition(0, 0, 3)
	// app.camOrtho.LookAt(&math32.Vector3{0, 0, 0})
	// app.camOrtho.SetZoom(1.0)

	// // Default camera is perspective
	// app.camera = app.camPersp

	// Creates orbit camera control
	// It is important to do this after the root panel subscription
	// to avoid GUI events being propagated to the orbit control.
	// app.orbit = control.NewOrbitControl(app.camera, app.win)

	app.CreateCamera()
	// Creates scene for 3D objects
	app.scene = core.NewNode()

	// Creates renderer
	// app.renderer = renderer.NewRenderer(gl)
	// err = app.renderer.AddDefaultShaders()
	// if err != nil {
	// 	return nil, fmt.Errorf("Error from AddDefaulShaders:%v", err)
	// }
	// app.renderer.SetScene(app.scene)
	// Create frame rater
	app.frameRater = NewFrameRater(*app.targetFPS)

	app.moblie = moblie.NewMoblie()
	app.guiroot = gui.NewRoot(app.moblie)
	app.guiroot.SetColor(math32.NewColor("silver"))

	// Sets the default window resize event handler
	return app, nil
}

// Get returns the application single instance or nil
// if the application was not created yet
func Get() *Application {

	return appInstance
}

// Log returns the application logger
func (app *Application) Log() *logger.Logger {

	return app.log
}

// Gl returns the OpenGL state
func (app *Application) Gl() *gls.GLS {

	return app.gl
}

// Scene returns the current application 3D scene
func (app *Application) Scene() *core.Node {

	return app.scene
}

// SetScene sets the 3D scene to be rendered
func (app *Application) SetScene(scene *core.Node) {

	app.renderer.SetScene(scene)
}

func (app *Application) Gui() *gui.Root {

	return app.guiroot
}

// SetGui sets the root panel of the gui to be rendered
func (app *Application) SetGui(root *gui.Root) {

	app.guiroot = root
	app.renderer.SetGui(app.guiroot)
}

// // SetPanel3D sets the gui panel inside which the 3D scene is shown.
// func (app *Application) SetPanel3D(panel3D gui.IPanel) {

// 	app.renderer.SetGuiPanel3D(panel3D)
// }

// Panel3D returns the current gui panel where the 3D scene is shown.
// func (app *Application) Panel3D() gui.IPanel {

// 	return app.renderer.Panel3D()
// }

// CameraPersp returns the application perspective camera
func (app *Application) CameraPersp() *camera.Perspective {

	return app.camPersp
}

// CameraOrtho returns the application orthographic camera
func (app *Application) CameraOrtho() *camera.Orthographic {

	return app.camOrtho
}

// Camera returns the current application camera
func (app *Application) Camera() camera.ICamera {

	return app.camera
}

// SetCamera sets the current application camera
func (app *Application) SetCamera(cam camera.ICamera) {

	app.camera = cam
}

// // Orbit returns the current camera orbit control
// func (app *Application) Orbit() *control.OrbitControl {

// 	return app.orbit
// }

// // SetOrbit sets the camera orbit control
// func (app *Application) SetOrbit(oc *control.OrbitControl) {

// 	app.orbit = oc
// }

// FrameRater returns the FrameRater object
func (app *Application) FrameRater() *FrameRater {

	return app.frameRater
}

// FrameCount returns the total number of frames since the call to Run()
func (app *Application) FrameCount() uint64 {

	return app.frameCount
}

// FrameDelta returns the duration of the previous frame
func (app *Application) FrameDelta() time.Duration {

	return app.frameDelta
}

// FrameDeltaSeconds returns the duration of the previous frame
// in float32 seconds
func (app *Application) FrameDeltaSeconds() float32 {

	return float32(app.frameDelta.Seconds())
}

// RunTime returns the duration since the call to Run()
func (app *Application) RunTime() time.Duration {

	return time.Now().Sub(app.startTime)
}

// RunSeconds returns the elapsed time in seconds since the call to Run()
func (app *Application) RunSeconds() float32 {

	return float32(time.Now().Sub(app.startTime).Seconds())
}

// Renderer returns the application renderer
func (app *Application) Renderer() *renderer.Renderer {

	return app.renderer
}

// SetCPUProfile must be called before Run() and sets the file name for cpu profiling.
// If set the cpu profiling starts before running the render loop and continues
// till the end of the application.
func (app *Application) SetCPUProfile(fname string) {

	*app.cpuProfile = fname
}

// Run runs the application render loop
func (app *Application) Run() error {

	// Set swap interval
	if *app.swapInterval >= 0 {
		app.log.Debug("Swap interval set to: %v", *app.swapInterval)
	}

	// Start profiling if requested
	if *app.cpuProfile != "" {
		f, err := os.Create(*app.cpuProfile)
		if err != nil {
			return err
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			return err
		}
		defer pprof.StopCPUProfile()
		app.log.Info("Started writing CPU profile to: %s", *app.cpuProfile)
	}

	// Start execution trace if requested
	if *app.execTrace != "" {
		f, err := os.Create(*app.execTrace)
		if err != nil {
			return err
		}
		defer f.Close()
		err = trace.Start(f)
		if err != nil {
			return err
		}
		defer trace.Stop()
		app.log.Info("Started writing execution trace to: %s", *app.execTrace)
	}

	app.startTime = time.Now()
	app.frameTime = time.Now()
	app.moblie.Subscribe(moblie.SystemAlive, func(evname string, ev interface{}) {
		aliveEvent := ev.(*moblie.AliveEvent)
		glctx := aliveEvent.Context
		app.InitGls(glctx)
		go func() {
			for {
				app.Loop()
			}
		}()
	})
	app.moblie.Subscribe(moblie.SystemForeground, func(evname string, ev interface{}) {
		foreEvent := ev.(*moblie.ForegroundEvent)
		glctx := foreEvent.Context
		app.gl.SetContext(glctx)
		app.show = true

	})
	app.moblie.Subscribe(moblie.SystemBackground, func(evname string, ev interface{}) {
		app.show = false
	})
	app.moblie.Run()
	return nil
}

func (app *Application) Loop() error {
	if app.gl == nil {
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			str := fmt.Sprintln(err)
			app.log.Error("recover err : %s", str)
		}
	}()

	app.frameRater.Start()

	// Updates frame start and time delta in context
	now := time.Now()
	app.frameDelta = now.Sub(app.frameTime)
	app.frameTime = now

	// Process application timers
	app.ProcessTimers()

	// Dispatch before render event
	app.Dispatch(OnBeforeRender, nil)

	// Renders the current scene and/or gui
	if app.show {
		isRender, err := app.renderer.Render(app.camera)
		if err != nil {
			panic(err)
		}

		if isRender {
			app.moblie.Publish()
		}
	}

	// fmt.Println("rendered ", rendered)
	// Poll input events and process them
	// Dispatch after render event
	app.Dispatch(OnAfterRender, nil)

	// Controls the frame rate
	app.frameRater.Wait()
	app.frameCount++

	return nil
}

func (app *Application) ClearUI() {
	cc := math32.NewColor("gray")
	app.gl.ClearColor(cc.R, cc.G, cc.B, 1)
	// app.gl.ClearColor(0, 0, 0, 1)
	app.gl.Clear(gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT | gl.COLOR_BUFFER_BIT)
}

func (app *Application) CreateCamera() {
	aspect := float32(750) / float32(1334)
	app.camPersp = camera.NewPerspective(65, aspect, 0.01, 1000)

	// Creates orthographic camera
	app.camOrtho = camera.NewOrthographic(-2, 2, 2, -2, 0.01, 100)
	app.camOrtho.SetPosition(0, 0, 3)
	app.camOrtho.LookAt(&math32.Vector3{0, 0, 0})
	app.camOrtho.SetZoom(1.0)

	// Default camera is perspective
	app.camera = app.camPersp
}

func (app *Application) InitGls(glctx gl.Context) {
	gs, err := gls.New(glctx, app.log)
	if err != nil {
		app.log.Error("init gls error : %s ", err)
		panic(err)
	}
	app.gl = gs
	// Checks OpenGL errors
	app.gl.SetCheckErrors(!*app.noglErrors)

	// // Logs OpenGL version
	glVersion := app.Gl().GetString(gl.VERSION)
	app.log.Info("OpenGL version: %s", glVersion)

	// // Clears the screen
	app.ClearUI()

	// Creates orbit camera control
	// It is important to do this after the root panel subscription
	// to avoid GUI events being propagated to the orbit control.
	// app.orbit = control.NewOrbitControl(app.camera, app.win)
	if app.control {
		app.orbit = control.NewOrbitControl(app.camera, app.moblie)
	}

	// Creates renderer
	app.renderer = renderer.NewRenderer(gs)
	err = app.renderer.AddDefaultShaders()
	if err != nil {
		app.log.Error("Error from AddDefaulShaders:%v", err)
		panic(err)
	}

	app.renderer.SetScene(app.scene)
	app.renderer.SetGui(app.guiroot)

}
