package main

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	var a, b float64 = 0.7, 3
	points1 := plotter.XYs{}
	points2 := plotter.XYs{}

	for i := 0; i <= 10; i++ {
		points1 = append(points1, plotter.XY{
			X: float64(i),
			Y: a*float64(i) + b,
		})
		points2 = append(points2, plotter.XY{
			X: float64(i),
			Y: a*float64(i) + b + (2*rand.Float64() - 1),
		})
	}

	result, err := optimize.Minimize(optimize.Problem{
		Func: func(x []float64) float64 {
			if len(x) != 2 {
				panic("illegal x")
			}
			a := x[0]
			b := x[1]
			var sum float64
			for _, point := range points2 {
				y := a*point.X + b
				sum += math.Abs(y - point.Y)
			}
			return sum
		},
	}, []float64{1, 1}, &optimize.Settings{}, &optimize.NelderMead{})
	if err != nil {
		panic(err)
	}

	fa, fb := result.X[0], result.X[1]
	points3 := plotter.XYs{}
	for i := 0; i <= 10; i++ {
		points3 = append(points3, plotter.XY{
			X: float64(i),
			Y: fa*float64(i) + fb,
		})
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, 10, 10

	if err := plotutil.AddLinePoints(plt,
		"line1", points1,
		"line2", points2,
		"line3", points3,
	); err != nil {
		panic(err)
	}

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "10-optimize-fit.png"); err != nil {
		panic(err)
	}
}
