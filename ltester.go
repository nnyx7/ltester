package main

import (
	"fmt"
	"net/http"
	"net/url"
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

func validateParams(u string,
	method string,
	numRequests int,
	duration int,
	warmUp int,
	change int,
	period int) (bool, error) {

	accumulatedError := ""

	if _, err := url.ParseRequestURI(u); err != nil {
		accumulatedError += (err.Error() + ", ")
	}

	if _, err := validateMethod(method); err != nil {
		accumulatedError += (err.Error() + ". ")
	}

	if numRequests <= 0 {
		accumulatedError += "numRequests must be positive, "
	}
	if duration <= 0 {
		accumulatedError += "duration must be positive, "
	}
	if warmUp < 0 {
		accumulatedError += "duration must be non-negative, "
	}
	if change != 0 {
		if period <= 0 {
			accumulatedError += "period must be positive, "
		}
		if period >= duration {
			accumulatedError += "period must be less than duration"
		}
	}

	if accumulatedError != "" {
		accumulatedError = accumulatedError[0 : len(accumulatedError)-2]
		return false, fmt.Errorf(accumulatedError)
	}
	return true, nil
}

var httpMethods = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace}

func validateMethod(method string) (bool, error) {
	for _, httpMethod := range httpMethods {
		if method == httpMethod {
			return true, nil
		}
	}
	return false, fmt.Errorf("invalid method: %v", method)
}
