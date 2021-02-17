package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var params Params
	setFlags(&params)
	flag.Parse()
	resPath := "resFile.txt"

	lt, err := ltesterFromParams(&params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resFile, err := os.Create(resPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	execResult, err := lt.Execute(resFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	XValues, YValues := fromFile(resPath)
	draw(XValues, YValues, "result.png")

	htmlParams := &HTMLParams{params.url, params.method, params.numRequests,
		params.duration, params.warmUp, params.change, params.period,
		mean(YValues), median(YValues), execResult.start.Format("15:04:45.000"),
		execResult.end.Format("15:04:45.000"), execResult.totalExecutions,
		execResult.successfulExecutions}

	err = genResultsHTML("template.html", "result.html", htmlParams)
	if err != nil {
		fmt.Println(err)
	}
}
