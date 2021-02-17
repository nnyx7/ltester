package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type result struct {
	duration time.Duration
	response *http.Response
	err      error
}

func makeRequest(client *http.Client, request *http.Request,
	resultChan chan<- *result, wg *sync.WaitGroup) (rs *result) {
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
	duration := time.Since(start)

	if err != nil {
		return &result{duration, nil, err}
	}
	return &result{duration, response, nil}
}

func writeResponse(f *os.File, rs *result, wg *sync.WaitGroup,
	mu *sync.Mutex) (int, error) {
	defer func() {
		wg.Done()
		mu.Unlock()
	}()

	mu.Lock()
	f.Sync()
	n, err := f.WriteString(rs.response.Status + "\n")
	return n, err
}

func (lt *Ltester) execute() (int, error) {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	var duration int64

	var wgFile sync.WaitGroup
	var mu sync.Mutex

	f, err := os.Create(lt.respFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var wgResults sync.WaitGroup
	resultChan := make(chan *result, lt.numRequests)

	ctr := 0
	for i := 0; i < lt.numRequests; i++ {
		wgResults.Add(1)
		go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
			resultChan, &wgResults)
	}

	for rs := range resultChan {
		ctr++
		wgFile.Add(1)
		go writeResponse(f, rs, &wgFile, &mu)
		now := time.Now().UnixNano() / int64(time.Millisecond)
		duration = now - start
		if duration >= int64(lt.duration) {
			break
		}
		wgResults.Add(1)
		go makeRequest(lt.client, lt.request.Clone(lt.request.Context()),
			resultChan, &wgResults)
	}

	wgResults.Wait()
	wgFile.Wait()
	f.Sync()

	return ctr, nil
}
