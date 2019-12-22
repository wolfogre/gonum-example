package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	xrand "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/spatial/r1"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	points := plotter.XYs{}
	for i := 0; i < 10; i++ {
		points = append(points, plotter.XY{
			X: 100 * rand.Float64(),
			Y: 100 * rand.Float64(),
		})
	}

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	scatter.Shape = draw.CircleGlyph{}

	heat := Heat{Points: points}
	heatmap := plotter.NewHeatMap(heat, moreland.SmoothBlueRed().Palette(100))

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
	Grad := func(grad, x []float64) {
		if len(grad) != len(x) {
			panic("illegal grad or x")
		}
		delta := 1e-9
		for i, v := range x {
			x[i] = v - delta
			f1 := Func(x)
			x[i] = v + delta
			f2 := Func(x)
			x[i] = v
			grad[i] = (f2 - f1) / (2 * delta)
		}
	}

	problem := optimize.Problem{
		Func: Func,
		Grad: Grad,
	}

	methods := []struct {
		Name   string
		Method optimize.Method
	}{
		{"BFGS", &optimize.BFGS{}},
		{"CG", &optimize.CG{}},
		{"CmaEsChol", &optimize.CmaEsChol{}},
		{"GradientDescent", &optimize.GradientDescent{}},
		{"GuessAndCheck", &optimize.GuessAndCheck{
			Rander: distmv.NewUniform([]r1.Interval{{0, 100}, {0, 100}}, xrand.NewSource(0)),
		}},
		{"LBFGS", &optimize.LBFGS{}},
		{"ListSearch", &optimize.ListSearch{
			Locs: mat.NewDense(6, 2, []float64{
				0, 10,
				20, 30,
				40, 50,
				60, 70,
				80, 90,
				90, 100,
			}),
		}},
		{"NelderMead", &optimize.NelderMead{}},
		//{"Newton", &optimize.Newton{}}, // what a pity
	}

	for _, method := range methods {
		recorder := &Recorder{}

		result, err := optimize.Minimize(problem, []float64{1, 1}, &optimize.Settings{
			Recorder: recorder,
		}, method.Method)
		if err != nil {
			panic(err)
		}

		pathLines, pathPoints, err := plotter.NewLinePoints(recorder.XYs)
		if err != nil {
			panic(err)
		}

		aim, err := plotter.NewScatter(plotter.XYs{
			{
				X: result.X[0],
				Y: result.X[1],
			},
		})
		if err != nil {
			panic(err)
		}
		aim.Shape = draw.CircleGlyph{}
		aim.Color = color.White

		plt, err := plot.New()
		if err != nil {
			panic(err)
		}
		plt.Y.Min, plt.X.Min, plt.Y.Max, plt.X.Max = 0, 0, 100, 100
		plt.Title.Text = method.Name

		plt.Add(heatmap, scatter, pathPoints, pathLines, aim)

		if err := plt.Save(5*vg.Inch, 5*vg.Inch, fmt.Sprintf("09-other-methods.%s.png", method.Name)); err != nil {
			panic(err)
		}
	}
}

type Heat struct {
	Points plotter.XYs
}

func (h Heat) Dims() (c, r int) { return 100, 100 }
func (h Heat) X(c int) float64  { return float64(c) }
func (h Heat) Y(r int) float64  { return float64(r) }
func (h Heat) Z(c, r int) float64 {
	var sum float64
	for _, p := range h.Points {
		sum += math.Sqrt(math.Pow(p.X-h.X(c), 2) + math.Pow(p.Y-h.Y(r), 2))
	}
	return -sum
}

type Recorder struct {
	XYs plotter.XYs
}

func (r *Recorder) Init() error {
	return nil
}

func (r *Recorder) Record(location *optimize.Location, op optimize.Operation, _ *optimize.Stats) error {
	if op != optimize.MajorIteration && op != optimize.InitIteration && op != optimize.PostIteration {
		return nil
	}
	r.XYs = append(r.XYs, plotter.XY{
		X: location.X[0],
		Y: location.X[1],
	})
	return nil
}
