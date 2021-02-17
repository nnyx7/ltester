package main

import "flag"

func setFlags(params *Params) {
	flag.StringVar(&params.url, "url", "http://example.com/", "URL of the application to test")
	flag.StringVar(&params.method, "method", "GET", "HTTP request method")
	flag.IntVar(&params.numRequests, "numRequest", 100, "Number concurrent request to execute")
	flag.IntVar(&params.duration, "duration", 10000, "Duration time in milliseconds")
	flag.IntVar(&params.warmUp, "warmUp", 0, "Warp-up time in milliseconds")
	flag.IntVar(&params.change, "change", 0, "n + change number requests to send")
	flag.IntVar(&params.period, "period", 0, "Period of time for change in milliseconds")
	flag.StringVar(&params.respFile, "respFile", "respFile.txt", "File to store responses in")
}

func ltesterFromParams(params *Params) (*Ltester, error) {
	lt, err := NewLtester(params.url, params.method, params.numRequests,
		params.duration, params.warmUp, params.change, params.period,
		params.respFile)
	if err != nil {
		return nil, err
	}
	return lt, nil
}
