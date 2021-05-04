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
		RMin:           3.3,        // Start Value of R, usually > 0 <RMax
		RMax:           4.0,        // End Value of R, < 4 or infinity will appear
		RStep:          0.0005,     // The step value from RMin -> RMax, smaller = more resolution, more CPU
		N:              0.4,        // Starting Population, 0.4 is my standard test value
		MaxDelta:       0.00000001, // The largest difference in population change considered stable
		MaxIterations:  10000000,   // The maximum number of iterations to try before giving up on stability
		MaxConcurrency: 16,         // Max number of go routines to use
		CirclePlot:     true,       // Plot results on a circle
	}, l)
	re := mathvis.NewRenderer(lm, l)
	re.Run()
}
