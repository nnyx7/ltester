package main

import (
	"net/http"
	"os"
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
	respFile    string
}

// NewLtester is a NewLtester constructor
func NewLtester(url string,
	method string,
	numRequests int,
	duration int,
	warmUp int,
	change int,
	period int,
	respFile string,
) (*Ltester, error) {

	if _, err := validateParams(url, method, numRequests, duration,
		warmUp, change, period); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(respFile); err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Create(respFile); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &Ltester{request, numRequests, duration, warmUp,
		change, period, &http.Client{}, respFile}, nil
}
