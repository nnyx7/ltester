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

// makeRequest makes http.Request and saves the response and the response time
// in result, as well as the ms passed from the start of making requests at all.
// The result from the function is put in result channel
func makeRequest(client *http.Client, request *http.Request,
	resultChan chan<- *result, startTime int64, wg *sync.WaitGroup) (rs *result) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				rs.err = err
			} else {
				rs.err = fmt.Errorf("Panic happened with %v", r)
			}
		}
		resultChan <- rs
		wg.Done()
	}()

	start := time.Now()
	response, err := client.Do(request)
	current := time.Now()
	duration := current.Sub(start).Milliseconds()
	fromStart := current.UnixNano()/int64(time.Millisecond) - startTime

	if err != nil {
		return &result{fromStart, duration, nil, err}
	}
	return &result{fromStart, duration, response, nil}
}

// saveResponse saves response information in the file f by taking
// the lock of f
func saveResponse(f *os.File, rs *result, successfulExecutions *int, wg *sync.WaitGroup,
	mu *sync.Mutex) (int, error) {
	defer func() {
		wg.Done()
	}()

	if rs.err == nil {
		mu.Lock()
		*successfulExecutions++
		line := fmt.Sprintf("%d %d %d\n", rs.fromStart, rs.duration,
			rs.response.StatusCode)
		n, err := f.WriteString(line)
		mu.Unlock()
		return n, err
	}
	return 0, rs.err
}

// Execute sends concurrent HTTP requests taking into account
// the corresponding Ltester configuration
func (lt *Ltester) Execute() (*ExecResult, error) {
	startTime := time.Now()
	start := startTime.UnixNano() / int64(time.Millisecond)

	var wgFile sync.WaitGroup
	var mu sync.Mutex

	f, err := os.Create(lt.respFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var wgResults sync.WaitGroup
	resultChan := make(chan *result, lt.numRequests)

	for i := 0; i < lt.numRequests; i++ {
		wgResults.Add(1)
		go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
			resultChan, start, &wgResults)
	}

	totalExecutions := 0
	successfulExecutions := 0
	for rs := range resultChan {
		totalExecutions++
		duration := time.Now().UnixNano()/int64(time.Millisecond) - start

		// if the warm-up time passed, start recording
		if int64(lt.warmUp) < duration {
			wgFile.Add(1)
			go saveResponse(f, rs, &successfulExecutions, &wgFile, &mu)
		}

		if duration >= int64(lt.duration) {
			break
		}

		// Executes goroutine on the place of the one that just finished
		wgResults.Add(1)
		go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
			resultChan, start, &wgResults)
	}

	wgResults.Wait()
	wgFile.Wait()

	if err := f.Sync(); err != nil {
		return nil, err
	}

	return &ExecResult{startTime, time.Now(), totalExecutions, successfulExecutions}, nil
}
