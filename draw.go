package main

import (
	"os"

	chart "github.com/wcharczuk/go-chart/v2"
)

func draw(XValues []float64, YValues []float64, outputImage string) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Duration from start (ms)",
		},
		YAxis: chart.YAxis{
			Name: "Response time (ms)",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: XValues,
				YValues: YValues,
			},
		},
	}

	f, _ := os.Create(outputImage)
	defer f.Close()
	graph.Render(chart.PNG, f)
}
