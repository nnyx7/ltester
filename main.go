package main

import (
	"flag"
	"fmt"
)

// Params hold parameters needed to initialize Ltester
type Params struct {
	url    string
	method string
	// headers
	// input
	numRequests int
	duration    int
	warmUp      int
	change      int
	period      int
	respFile    string
}

func main() {
	var params Params
	setFlags(&params)
	flag.Parse()

	lt, err := ltesterFromParams(&params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	execResult, err := lt.execute()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	XValues, YValues := fromFile(lt.respFile)
	draw(XValues, YValues, "result.png")

	htmlParams := &HTMLParams{params.url, params.method, params.numRequests,
		params.duration, params.warmUp, params.change, params.period,
		mean(YValues), median(YValues), execResult.start.Format("15:04:45.000"),
		execResult.end.Format("15:04:45.000"), execResult.totalExecutions,
		execResult.successfulExecutions}

	genResultsHTML("template.html", "result.html", htmlParams)
}
