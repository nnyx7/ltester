package main

import (
	"net/http"
)

// Ltester is a struct that provides load testing functionality
type Ltester struct {
	request     *http.Request
	numRequests int
	duration    int
	warmUp      int
	change      int
	period      int
	client      *http.Client
}

// NewLtester is a NewLtester constructor
func NewLtester(url string,
	method string,
	numRequests int,
	duration int,
	warmUp int,
	change int,
	period int,
) (*Ltester, error) {

	if _, err := validateParams(url, method, numRequests, duration,
		warmUp, change, period); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	return &Ltester{request, numRequests, duration, warmUp,
		change, period, &http.Client{}}, nil
}
