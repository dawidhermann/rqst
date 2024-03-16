package rqst

import (
	"net/http"
	"sync"
)

// Result of a single request
type ParalelResult struct {
	Id     string
	Result any
	Error  error
}

// Helper for executing multiple requests at once
type ParalelRqst struct {
	executor *http.Client
	requests []ParalelRequestConfig
}

// Configuration for a single paralel request
type ParalelRequestConfig struct {
	Id            string
	CreateRequest func() *http.Request
	ResultMapper  func(response *http.Response) (any, error)
}

// Creates new ParalelRqst
func NewParalel(requestExecutor *http.Client) *ParalelRqst {
	return &ParalelRqst{
		executor: requestExecutor,
		requests: make([]ParalelRequestConfig, 0),
	}
}

// Add new paralel request to execution list
func (rqst *ParalelRqst) AddNextParalelRequest(requestConfig ParalelRequestConfig) *ParalelRqst {
	rqst.requests = append(rqst.requests, requestConfig)
	return rqst
}

// Add multiple new requests to execution list
func (rqst *ParalelRqst) AddNextMultipleParalelRequests(requestConfigs ...ParalelRequestConfig) *ParalelRqst {
	rqst.requests = append(rqst.requests, requestConfigs...)
	return rqst
}

// Execute all requests at once
func (rqst *ParalelRqst) Execute() map[string]ParalelResult {
	results := make(map[string]ParalelResult)
	var wg sync.WaitGroup
	wg.Add(len(rqst.requests))
	for _, requestConfig := range rqst.requests {
		requestConfig := requestConfig
		requestId := requestConfig.Id
		req := requestConfig.CreateRequest()
		go func(request *http.Request) {
			defer wg.Done()
			response, err := rqst.executor.Do(request)
			if err != nil {
				results[requestId] = ParalelResult{Id: requestId, Error: err}
				return
			}
			mappedData, err := requestConfig.ResultMapper(response)
			if err != nil {
				results[requestId] = ParalelResult{Id: requestId, Error: err}
				return
			}
			results[requestId] = ParalelResult{Id: requestId, Result: mappedData}
			return
		}(req)
	}
	wg.Wait()
	return results
}

// Execute all requests at once and emit all results to channel
func (rqst *ParalelRqst) ExecuteWithChan() <-chan ParalelResult {
	resChan := make(chan ParalelResult, len(rqst.requests))
	var wg sync.WaitGroup
	wg.Add(len(rqst.requests))
	for _, requestConfig := range rqst.requests {
		requestConfig := requestConfig
		requestId := requestConfig.Id
		req := requestConfig.CreateRequest()
		go func(request *http.Request, resultChan chan<- ParalelResult) {
			defer wg.Done()
			response, err := rqst.executor.Do(request)
			if err != nil {
				resultChan <- ParalelResult{Id: requestId, Error: err}
				return
			}
			mappedData, err := requestConfig.ResultMapper(response)
			if err != nil {
				resultChan <- ParalelResult{Id: requestId, Error: err}
				return
			}
			resultChan <- ParalelResult{Id: requestId, Result: mappedData}
			return
		}(req, resChan)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()
	return resChan
}
