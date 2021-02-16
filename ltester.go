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
}

// NewLtester is a NewLtester constructor
func NewLtester(request *http.Request,
	numRequests int,
	duration int,
	warmUp int,
	change int,
	period int,
) (*Ltester, error) {

	if _, err := validateParams(request, numRequests, duration,
		warmUp, change, period); err != nil {
		return nil, err
	}

	return &Ltester{request, numRequests, duration, warmUp,
		change, period}, nil
}

func validateParams(request *http.Request,
	numRequests int,
	duration int,
	warmUp int,
	change int,
	period int) (bool, error) {

	accumulatedError := ""

	if _, err := validateRequest(request); err != nil {
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

func validateRequest(request *http.Request) (bool, error) {
	accumulatedError := ""
	// Check if URL is valid
	if _, err := url.ParseRequestURI(request.URL.String()); err != nil {
		accumulatedError += (err.Error() + ", ")
	}

	// Check if Method is valid
	for _, httpMethod := range httpMethods {
		if request.Method == httpMethod {
			if accumulatedError == "" {
				return true, nil
			}
			return false, fmt.Errorf(accumulatedError)
		}
	}
	accumulatedError += fmt.Sprintf("invalid method: %v", request.Method)
	return false, fmt.Errorf(accumulatedError)
}
