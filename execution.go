package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// ExecResult holds execution metrics from the process
// of sending requests from Ltester
type ExecResult struct {
	start                time.Time
	end                  time.Time
	totalExecutions      int
	successfulExecutions int
}

// result holds metrics and response from a specific request
type result struct {
	fromStart int64
	duration  int64
	response  *http.Response
	err       error
}

var te int = 0
var totalExecutions *int = &te

var se int = 0
var successfulExecutions *int = &se

var wgResults sync.WaitGroup
var wgFile sync.WaitGroup
var muResults sync.Mutex
var muFile sync.Mutex

// makeRequest makes http.Request and saves the response and the response time
// in result, as well as the ms passed from the start of making requests at all.
// The result from the function is put in result channel
func makeRequest(c *http.Client, req *http.Request, resCh chan<- *result,
	warmUp int, sTime int64, respFile *os.File, resFile *os.File) (rs *result) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				rs.err = err
			} else {
				rs.err = fmt.Errorf("Panic happened with %v", r)
			}
		}
		resCh <- rs
		// if the warm-up time passed, start recording
		if int64(warmUp) < rs.fromStart {
			wgFile.Add(1)
			go saveResponse(rs, respFile, resFile)
		}
		// Count total Executions
		muResults.Lock()
		*totalExecutions++
		muResults.Unlock()

		wgResults.Done()
	}()

	start := time.Now()
	response, err := c.Do(req)
	current := time.Now()
	duration := current.Sub(start).Milliseconds()
	fromStart := current.UnixNano()/int64(time.Millisecond) - sTime

	if err != nil {
		return &result{fromStart, duration, nil, err}
	}
	return &result{fromStart, duration, response, nil}
}

// saveResponse saves response information in the file f by taking
// the lock of f
func saveResponse(rs *result, respFile *os.File, resFile *os.File) error {
	defer func() {
		wgFile.Done()
	}()

	if rs.err == nil {
		muFile.Lock()
		*successfulExecutions++
		line := fmt.Sprintf("%d %d %d\n", rs.fromStart, rs.duration,
			rs.response.StatusCode)
		_, err := resFile.WriteString(line)
		muFile.Unlock()

		muFile.Lock()
		line = fmt.Sprintf("%v\n", rs.response)
		_, err = respFile.WriteString(line)
		muFile.Unlock()
		return err
	}
	return rs.err
}

// Execute sends concurrent HTTP requests taking into account
// the corresponding Ltester configuration
func (lt *Ltester) Execute(resFile *os.File) (*ExecResult, error) {
	startTime := time.Now()
	start := startTime.UnixNano() / int64(time.Millisecond)
	var duration int64

	respFile, err := os.Create(lt.respFile)
	defer respFile.Close()
	if err != nil {
		return nil, err
	}

	numRequests := lt.numRequests

	for numRequests > 0 && duration < int64(lt.duration) {
		checkpoint := time.Now().UnixNano() / int64(time.Millisecond)

		resultChan := make(chan *result, numRequests)
		// Start numRequest times makeRequest
		for i := 0; i < numRequests; i++ {
			wgResults.Add(1)
			go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
				resultChan, lt.warmUp, start, respFile, resFile)
		}

		for range resultChan {
			duration = time.Now().UnixNano()/int64(time.Millisecond) - start
			if duration > int64(lt.duration) {
				break
			}
			if lt.change != 0 && (checkpoint-int64(lt.period)) > 0 {
				break
			}
			// Executes goroutine on the place of the one that just finished
			wgResults.Add(1)
			go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
				resultChan, lt.warmUp, start, respFile, resFile)
		}
		wgResults.Wait()
		wgFile.Wait()
		close(resultChan)

		if lt.change != 0 {
			numRequests = lt.numRequests + (lt.change * int(duration) / lt.period)
		}
		duration = time.Now().UnixNano()/int64(time.Millisecond) - start
	}

	if err := respFile.Sync(); err != nil {
		return nil, err
	}
	return &ExecResult{startTime, time.Now(), te, se}, nil
}
