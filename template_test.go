package main

import (
	"fmt"
	"os"
	"testing"
)

var templatePath string = "template.html"
var validHTMLParams = &HTMLParams{"http://example.com/", "GET", 100, 100,
	0, 0, 0, 0.0, 0.0, "00:00:00.000", "00:00:00.000", 1000, 1000}

func TestTemplateGeneratorValidParams(t *testing.T) {
	resultPath := "resultHTMLPath.html"

	err := genResultsHTML(templatePath, resultPath, validHTMLParams)
	os.Remove(resultPath)

	if err != nil {
		out := fmt.Sprintf("Expected nil\n") +
			fmt.Sprintf("Received: %v\n", err)
		t.Fatalf(out)
	}
}

func TestTemplateGeneratorNilParams(t *testing.T) {
	resultPath := "resultHTMLPath.html"

	err := genResultsHTML(templatePath, resultPath, nil)
	os.Remove(resultPath)

	if err == nil {
		out := fmt.Sprintf("Expected error\n") +
			fmt.Sprintf("Received: %v\n", err)
		t.Fatalf(out)
	}
}
