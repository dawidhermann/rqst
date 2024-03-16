// Request helper
package rqst

import (
	"net/http"
)

// Rqst is basic struct for helping transform chain of requests
type Rqst struct {
	executor *http.Client
	requests []RequestConfig
}

// A single request configuration
type RequestConfig struct {
	CreateRequest func(data any) *http.Request
	ResultMapper  func(response *http.Response) (any, error)
	// RawResultInterceptor    func(rawData *http.Response)
	MappedResultInterceptor func(mappedData any)
}

// Create new Rqst
func New(requestExecutor *http.Client) *Rqst {
	return &Rqst{
		executor: requestExecutor,
		requests: make([]RequestConfig, 0),
	}
}

// Add next request to execution chain
func (rqst *Rqst) AddNextRequest(requestConfig RequestConfig) *Rqst {
	rqst.requests = append(rqst.requests, requestConfig)
	return rqst
}

// Add multiple requests to execution chain at once
func (rqst *Rqst) AddNextMultipleRequests(requestConfigs ...RequestConfig) *Rqst {
	rqst.requests = append(rqst.requests, requestConfigs...)
	return rqst
}

// Execute request chain
func (rqst *Rqst) Execute() any {
	var lastResult any
	for _, requestConfig := range rqst.requests {
		req := requestConfig.CreateRequest(lastResult)
		response, err := rqst.executor.Do(req)
		if err != nil {
			return err
		}
		mappedData, err := requestConfig.ResultMapper(response)
		if err != nil {
			return err
		}
		if requestConfig.MappedResultInterceptor != nil {
			requestConfig.MappedResultInterceptor(mappedData)
		}
		lastResult = mappedData
	}
	return lastResult
}
