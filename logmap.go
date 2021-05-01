package mathvis

import (
	"context"
	"math"
	"sync"

	"go.uber.org/zap"
)

type LogMap struct {
	buffer chan vector2
	logger *zap.Logger
	tokens chan interface{}
	config LogMapConfig
	wg     *sync.WaitGroup
}

type LogMapConfig struct {
	RMin           float64
	RMax           float64
	RStep          float64
	N              float64
	MaxDelta       float64
	MaxIterations  int
	MaxConcurrency int
}

func NewLogMap(config LogMapConfig, logger *zap.Logger) *LogMap {
	l := LogMap{}
	l.logger = logger
	l.config = config
	l.buffer = make(chan vector2, 1024)
	l.tokens = make(chan interface{}, config.MaxConcurrency)
	var wg sync.WaitGroup
	l.wg = &wg

	l.logger.Info("adding tokens")
	for i := 0; i < l.config.MaxConcurrency; i++ {
		l.logger.Info("adding token")
		l.tokens <- nil
	}
	l.logger.Info("done adding tokens")

	return &l
}

func (l *LogMap) GetPointChannel() chan vector2 {
	return l.buffer
}

func (l *LogMap) calcStablePop(ctx context.Context, r float64) {
	i := 0
	n := l.config.N

	for delta := l.config.MaxDelta; delta >= l.config.MaxDelta && i < l.config.MaxIterations; i++ {
		nplusone := r * n * (1 - n)
		delta = math.Abs(nplusone - n)
		n = nplusone
	}

	l.buffer <- vector2{X: r, Y: n}

	l.wg.Done()
	l.tokens <- nil
}

func (l *LogMap) Simulate(ctx context.Context) {
outer:
	for curr := l.config.RMin; curr < l.config.RMax; curr += l.config.RStep {
		select {
		case <-l.tokens:
			l.wg.Add(1)
			go l.calcStablePop(ctx, curr)
		case <-ctx.Done():
			break outer
		}
	}
	l.wg.Wait()
	l.logger.Info("calc done")
}
