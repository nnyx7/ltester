package main

import (
	"fmt"
	"os"
	"text/template"
)

type HTMLParams struct {
	URL             string
	Method          string
	NumRequests     int
	Duration        int
	WarmUp          int
	Change          int
	Period          int
	Mean            float64
	Median          float64
	Start           string
	End             string
	TotalExecutions int
}

func genResultsHTML(templateFile string, resultFile string, params *HTMLParams) {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Println("template parsing error: ", err)
	}
	f, _ := os.Create(resultFile)

	err = t.Execute(f, params)
	if err != nil {
		fmt.Println("template executing error: ", err)
	}
}
