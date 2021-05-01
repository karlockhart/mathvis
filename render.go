package mathvis

import (
	"context"
	"image/color"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"go.uber.org/zap"
)

type Simulation interface {
	GetPointChannel() chan vector2
	Simulate(context.Context)
}

type vector2 struct {
	X float64
	Y float64
}

type Renderer struct {
	logger         *zap.Logger
	pointBuffer    []vector2
	pointBufferMtx sync.Mutex
	simulation     Simulation
}

func NewRenderer(s Simulation, logger *zap.Logger) *Renderer {
	r := Renderer{}
	r.simulation = s
	r.logger = logger
	return &r
}

func (r *Renderer) enqueuePoints(ctx context.Context) {
	ch := r.simulation.GetPointChannel()
	for {
		select {
		case v := <-ch:
			r.pointBufferMtx.Lock()
			r.pointBuffer = append(r.pointBuffer, v)
			r.pointBufferMtx.Unlock()
			r.logger.Info("added point to buffer")
		case <-ctx.Done():
			return
		}
	}

}

func (r *Renderer) drawNewPoints(imd *imdraw.IMDraw) {
	r.pointBufferMtx.Lock()
	ipb := r.pointBuffer
	r.pointBuffer = []vector2{}
	r.pointBufferMtx.Unlock()

	for _, v := range ipb {
		r.logger.Sugar().Infof("%f,%fs", v.X, v.Y)
	}

}

func (r *Renderer) run() {
	cfg := pixelgl.WindowConfig{
		Title:  "MathVis",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		r.logger.Fatal(err.Error())
	}

	win.Clear(color.White)

	imd := imdraw.New(nil)

	ctx, cf := context.WithCancel(context.Background())

	r.logger.Info("starting enque")
	go r.enqueuePoints(ctx)
	r.logger.Info("starting smilulation")
	go r.simulation.Simulate(ctx)
	for !win.Closed() {
		r.drawNewPoints(imd)
		win.Update()
	}
	cf()
}

func (r *Renderer) Run() {
	pixelgl.Run(r.run)
}
