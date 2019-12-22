package main

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/plotter"
)

func main() {
	points := plotter.XYs{}
	for i := 0; i < 10; i++ {
		points = append(points, plotter.XY{
			X: 100 * rand.Float64(),
			Y: 100 * rand.Float64(),
		})
	}

	Func := func(x []float64) float64 {
		if len(x) != 2 {
			panic("illegal x")
		}
		var sum float64
		for _, point := range points {
			sum += math.Sqrt(math.Pow(point.X-x[0], 2) + math.Pow(point.Y-x[1], 2))
		}
		return sum
	}

	problem := optimize.Problem{
		Func:   Func,
		Grad:   nil,
		Hess:   nil,
		Status: nil,
	}

	_ = problem
}
