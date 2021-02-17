package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type ExecResult struct {
	start           time.Time
	end             time.Time
	totalExecutions int
}

type result struct {
	fromStart int64
	duration  int64
	response  *http.Response
	err       error
}

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

func writeResponse(f *os.File, rs *result, wg *sync.WaitGroup,
	mu *sync.Mutex) (int, error) {
	defer func() {
		wg.Done()
	}()

	mu.Lock()
	line := fmt.Sprintf("%d %d %d\n", rs.fromStart, rs.duration,
		rs.response.StatusCode)
	mu.Unlock()
	return f.WriteString(line)
}

func (lt *Ltester) execute() (*ExecResult, error) {
	startTime := time.Now()
	start := startTime.UnixNano() / int64(time.Millisecond)
	var duration int64

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
	for rs := range resultChan {
		totalExecutions++
		wgFile.Add(1)
		go writeResponse(f, rs, &wgFile, &mu)
		now := time.Now().UnixNano() / int64(time.Millisecond)
		duration = now - start
		if duration >= int64(lt.duration) {
			break
		}
		wgResults.Add(1)
		go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
			resultChan, start, &wgResults)
	}

	wgResults.Wait()
	wgFile.Wait()

	if err := f.Sync(); err != nil {
		return nil, err
	}

	return &ExecResult{startTime, time.Now(), totalExecutions}, nil
}
