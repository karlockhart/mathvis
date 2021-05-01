package main

import (
	"log"

	"github.com/karlockhart/mathvis"
	"go.uber.org/zap"
)

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	lm := mathvis.NewLogMap(mathvis.LogMapConfig{
		RMin:           0.0,
		RMax:           4.0,
		RStep:          0.001,
		N:              0.4,
		MaxDelta:       0.0000001,
		MaxIterations:  1000000,
		MaxConcurrency: 16,
	}, l)
	re := mathvis.NewRenderer(lm, l)
	re.Run()
}
