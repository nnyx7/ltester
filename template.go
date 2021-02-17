package main

import (
	"fmt"
	"os"
	"text/template"
)

// HTMLParams holds the parametes with which to execute
// a corresponding template
type HTMLParams struct {
	URL                  string
	Method               string
	NumRequests          int
	Duration             int
	WarmUp               int
	Change               int
	Period               int
	Mean                 float64
	Median               float64
	Start                string
	End                  string
	TotalExecutions      int
	SuccessfulExecutions int
}

// genResultsHTML fills the template with provided HTMLParams
// and saves the result in resultFile
func genResultsHTML(templateFile string, resultFile string,
	params *HTMLParams) error {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("template parsing error: %v", err)
	}
	f, err := os.Create(resultFile)
	if err != nil {
		return fmt.Errorf("result file creating error: %v", err)
	}
	err = t.Execute(f, params)
	if err != nil {
		return fmt.Errorf("template executing error: %v", err)
	}
	return nil
}
