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
}

func main() {
	var params Params
	setFlags(&params)
	flag.Parse()

	ls, err := ltesterFromParams(&params)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%T\n", ls)
	}
}

func setFlags(params *Params) {
	flag.StringVar(&params.url, "url", "http://example.com/", "URL of the application to test")
	flag.StringVar(&params.method, "method", "GET", "HTTP request method")
	flag.IntVar(&params.numRequests, "numRequest", 10, "Number concurrent request to execute")
	flag.IntVar(&params.duration, "duration", 1000, "Duration time in milliseconds")
	flag.IntVar(&params.warmUp, "warmUp", 0, "Warp-up time in milliseconds")
	flag.IntVar(&params.change, "change", 0, "n + change number requests to send")
	flag.IntVar(&params.period, "period", 0, "Period of time for change in milliseconds")
}

func ltesterFromParams(params *Params) (*Ltester, error) {
	lt, err := NewLtester(params.url, params.method, params.numRequests, params.duration, params.warmUp, params.change, params.period)
	if err != nil {
		return nil, err
	}
	return lt, nil
}
