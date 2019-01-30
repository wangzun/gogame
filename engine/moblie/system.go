package moblie

import (
	"github.com/wangzun/gogame/engine/core"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

type Moblie struct {
	core.Dispatcher // Embedded event dispatcher
	mApp            app.App
}

const SystemAlive = "moblie.SystemAlive"
const SystemForeground = "moblie.SystemForeground"
const SystemBackground = "moblie.SystemBackground"
const SystemTouch = "moblie.SystemTouch"
const SystemSize = "moblie.SystemSize"

func NewMoblie() *Moblie {
	m := &Moblie{}
	m.Dispatcher.Initialize()
	return m
}

type AliveEvent struct {
	Context gl.Context
}

type ForegroundEvent struct {
	Context gl.Context
}

type BackgroundEvent struct {
}

type SizeEvent struct {
}

type TouchEvent struct {
}

func (m *Moblie) Run() {
	app.Main(func(a app.App) {
		// var sz size.Event
		// var glctx gl.Context
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				// app.log.Info("lifecycle event : %s", e.String())
				// app.log.Info(" dead : %d", e.Crosses(lifecycle.StageDead))
				// app.log.Info(" visible : %d", e.Crosses(lifecycle.StageVisible))
				// app.log.Info(" alive : %d", e.Crosses(lifecycle.StageAlive))
				// app.log.Info(" focused : %d", e.Crosses(lifecycle.StageFocused))
				switch e.Crosses(lifecycle.StageAlive) {
				case lifecycle.CrossOn:
					glctx, _ := e.DrawContext.(gl.Context)
					m.mApp = a
					ae := &AliveEvent{Context: glctx}
					m.Dispatch(SystemAlive, ae)
				case lifecycle.CrossOff:
				}

				switch e.Crosses(lifecycle.StageFocused) {
				case lifecycle.CrossOn:
					glctx, _ := e.DrawContext.(gl.Context)
					ae := &ForegroundEvent{Context: glctx}
					m.Dispatch(SystemForeground, ae)
				case lifecycle.CrossOff:
					ae := &BackgroundEvent{}
					m.Dispatch(SystemBackground, ae)
				}
			case size.Event:
				se := &SizeEvent{}
				m.Dispatch(SystemSize, se)
				//sz = e
				//touchX = float32(sz.WidthPx / 2)
				//touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if e.External {
					continue
				}
			case touch.Event:
				te := &TouchEvent{}
				m.Dispatch(SystemTouch, te)
				//touchX = e.X
				//touchY = e.Y
			}
		}
	})
}

func (m *Moblie) Publish() {
	m.mApp.Publish()

}
