package main

import (
	"math/rand"

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

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, 10, 10

	if err := plotutil.AddLinePoints(plt,
		"line1", points1,
		"line2", points2,
	); err != nil {
		panic(err)
	}

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "02-draw-linepoint.png"); err != nil {
		panic(err)
	}
}
