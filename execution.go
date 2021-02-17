package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type result struct {
	duration time.Duration
	response *http.Response
	err      error
}

func (lt *Ltester) makeRequest(request *http.Request, resultChan chan<- *result, wg *sync.WaitGroup) (rs *result) {
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
	response, err := lt.client.Do(request)
	duration := time.Since(start)

	if err != nil {
		return &result{duration, nil, err}
	}
	return &result{duration, response, nil}
}

func (lt *Ltester) execute() int {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	var duration int64

	var wg sync.WaitGroup
	resultChan := make(chan *result, lt.numRequests)

	ctr := 0
	for i := 0; i < lt.numRequests; i++ {
		wg.Add(1)
		go lt.makeRequest(lt.request.Clone(lt.request.Context()), resultChan, &wg)
		ctr++
	}

	for range resultChan {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		duration = now - start
		if duration >= int64(lt.duration) {
			break
		}
		wg.Add(1)
		go lt.makeRequest(lt.request.Clone(lt.request.Context()), resultChan, &wg)
		ctr++
	}
	wg.Wait()

	return ctr
}
