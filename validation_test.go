package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

// ---------- {url, method, numRequests, duration, warmUp, change, period}
var validParams = &Params{"http://example.com/", "GET", 100, 100, 0, 0, 0, "respFile.txt"}

func TestLtesterWithValidParams(t *testing.T) {
	params := copyParams(validParams)
	testValidLtester(t, params)
}

// url
func TestLtesterWithInvalidUrl(t *testing.T) {
	params := copyParams(validParams)
	params.url = "Invalid URL"

	testInvalidLtester(t, params)
}

// method
func TestLtesterWithInvalidMethod(t *testing.T) {
	params := copyParams(validParams)
	params.method = "get"

	testInvalidLtester(t, params)
}

// numRequests
func TestLtesterWithZeroNumRequests(t *testing.T) {
	params := copyParams(validParams)
	params.numRequests = 0

	testInvalidLtester(t, params)
}

func TestLtesterWithNegativeNumRequests(t *testing.T) {
	params := copyParams(validParams)
	params.numRequests = -1

	testInvalidLtester(t, params)
}

// duration
func TestLtesterWithZeroDuration(t *testing.T) {
	params := copyParams(validParams)
	params.duration = 0

	testInvalidLtester(t, params)
}

func TestLtesterWithNegativeDuration(t *testing.T) {
	params := copyParams(validParams)
	params.duration = -1

	testInvalidLtester(t, params)
}

// warmUp
func TestLtesterWithPositiveWarmUp(t *testing.T) {
	params := copyParams(validParams)
	params.warmUp = 1

	testValidLtester(t, params)
}

func TestLtesterWithNegativeWarmUp(t *testing.T) {
	params := copyParams(validParams)
	params.warmUp = -1

	testInvalidLtester(t, params)
}

// period
func TestLtesterWithPositivePeriodZeroChange(t *testing.T) {
	params := copyParams(validParams)
	params.change = 0
	params.period = -1

	testValidLtester(t, params)
}

func TestLtesterWithZeroPeriodZeroChange(t *testing.T) {
	params := copyParams(validParams)
	params.change = 0
	params.period = 0

	testValidLtester(t, params)
}

func TestLtesterWithNegativePeriodZeroChange(t *testing.T) {
	params := copyParams(validParams)
	params.change = 0
	params.period = -1

	testValidLtester(t, params)
}

func TestLtesterWithZeroPeriodWithChange(t *testing.T) {
	params := copyParams(validParams)
	params.change = 1
	params.period = -1

	testInvalidLtester(t, params)
}

func TestLtesterWithNegativePeriodWithChange(t *testing.T) {
	params := copyParams(validParams)
	params.change = 1
	params.period = -1

	testInvalidLtester(t, params)
}

func TestLtesterWithChangePeriodMoreThanDuration(t *testing.T) {
	params := copyParams(validParams)
	params.change = 1
	params.duration = 10
	params.period = 11

	testInvalidLtester(t, params)
}

func TestLtesterWithChangePeriodEqDuration(t *testing.T) {
	params := copyParams(validParams)
	params.change = 1
	params.duration = 10
	params.period = 10

	testInvalidLtester(t, params)
}

// Helper methods

func copyParams(params *Params) *Params {
	return &Params{params.url, params.method, params.numRequests,
		params.duration, params.warmUp, params.change, params.period,
		params.respFile}
}

func testInvalidLtester(t *testing.T, params *Params) {
	lt, err := ltesterFromParams(params)

	if err == nil {
		out := fmt.Sprintf("Creating Ltester with values:\n") +
			fmt.Sprintf("%v\n", lt.toJSON()) +
			fmt.Sprintf("Expected result: error\n") +
			fmt.Sprintf("Actual result: %T\n", lt)
		t.Fatalf(out)
	}
}

func testValidLtester(t *testing.T, params *Params) {
	lt, err := ltesterFromParams(params)

	if err != nil {
		out := fmt.Sprintf("Creating Ltester with values:\n") +
			fmt.Sprintf("%v\n", newPubLtester(params).toJSON()) +
			fmt.Sprintf("Expected result: %T\n", lt) +
			fmt.Sprintf("Actual result: error(%v)\n", err.Error())
		t.Fatalf(out)
	}
}

// This struct is the same with the Ltester struct but with
// public fields. The purpose behind it to have a nice string
// presentation of Ltester instance.
type pubLtester struct {
	URL        string
	Method     string
	NumRequest int
	Duration   int
	WarmUp     int
	Change     int
	Period     int
	RespFile   string
}

func (plt *pubLtester) toJSON() string {
	ltesterJSON, err := json.Marshal(plt)
	if err != nil {
		return err.Error()
	}
	return string(ltesterJSON)
}

func newPubLtester(params *Params) *pubLtester {
	return &pubLtester{params.url, params.method, params.numRequests,
		params.duration, params.warmUp, params.change, params.period,
		params.respFile}
}

func copyLtester(lt *Ltester) *pubLtester {
	return &pubLtester{lt.request.URL.String(), lt.request.Method,
		lt.numRequests, lt.duration, lt.warmUp, lt.change, lt.period,
		lt.respFile}
}

func (lt *Ltester) toJSON() string {
	pubLtester := copyLtester(lt)
	return pubLtester.toJSON()
}
