package main

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	var a, b float64 = 0.7, 3
	points := plotter.XYs{}
	for i := 0; i <= 10; i++ {
		points = append(points, plotter.XY{
			X: float64(i),
			Y: a*float64(i) + b,
		})
	}

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}
	plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, 10, 10

	if err := plotutil.AddLines(plt,
		"line1", points,
	); err != nil {
		panic(err)
	}

	if err := plt.Save(5*vg.Inch, 5*vg.Inch, "01-draw-line.png"); err != nil {
		panic(err)
	}
}
