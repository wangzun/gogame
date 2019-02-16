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
	core.Dispatcher   // Embedded event dispatcher
	mApp              app.App
	WidthPx, HeightPx int
}

const SystemAlive = "moblie.SystemAlive"
const SystemForeground = "moblie.SystemForeground"
const SystemBackground = "moblie.SystemBackground"
const SystemTouch = "moblie.SystemTouch"
const SystemSize = "moblie.SystemSize"
const SystemFrame = "moblie.SystemFrame"

func NewMoblie() *Moblie {
	m := &Moblie{}
	m.WidthPx = 750
	m.HeightPx = 1334
	// m.WidthPx = 1200
	// m.HeightPx = 1550
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
	WidthPx, HeightPx int
}

type TouchEvent struct {
	X, Y     float32
	Sequence int64
	Type     Type
}

type Type byte

const (
	TypeBegin Type = iota
	TypeMove
	TypeEnd
)

func (m *Moblie) Run() {
	app.Main(func(a app.App) {
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
				se := &SizeEvent{WidthPx: e.WidthPx, HeightPx: e.HeightPx}
				m.WidthPx = e.WidthPx
				m.HeightPx = e.HeightPx
				m.Dispatch(SystemSize, se)
			case paint.Event:
				if e.External {
					continue
				}
			case touch.Event:
				te := &TouchEvent{X: e.X, Y: e.Y, Sequence: int64(e.Sequence), Type: Type(e.Type)}
				m.Dispatch(SystemTouch, te)
			}
		}
	})
}

func (m *Moblie) Publish() {
	m.mApp.Publish()
}

func (m *Moblie) Frame() {
	m.Dispatch(SystemFrame, nil)
}
