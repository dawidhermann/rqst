package rqst

import (
	"net/http"
	"sync"
)

type ParalelResult struct {
	Result any
	Error  error
}

type ParalelRqst struct {
	executor *http.Client
	requests []*ParalelRequestConfig
}

type ParalelRequestConfig struct {
	Id            string
	CreateRequest func() *http.Request
	ResultMapper  func(response *http.Response) (any, error)
}

func NewParalel(requestExecutor *http.Client) *ParalelRqst {
	return &ParalelRqst{
		executor: requestExecutor,
		requests: make([]*ParalelRequestConfig, 0),
	}
}

func (rqst *ParalelRqst) AddNextParalelRequest(requestConfig *ParalelRequestConfig) *ParalelRqst {
	rqst.requests = append(rqst.requests, requestConfig)
	return rqst
}

func (rqst *ParalelRqst) AddNextMultipleParalelRequests(requestConfigs ...*ParalelRequestConfig) *ParalelRqst {
	rqst.requests = append(rqst.requests, requestConfigs...)
	return rqst
}

func (rqst *ParalelRqst) Execute() map[string]ParalelResult {
	results := make(map[string]ParalelResult)
	var wg sync.WaitGroup
	for _, requestConfig := range rqst.requests {
		requestConfig := requestConfig
		wg.Add(1)
		go func() {
			req := requestConfig.CreateRequest()
			response, err := rqst.executor.Do(req)
			if err != nil {
				results[requestConfig.Id] = ParalelResult{Error: err}
				defer wg.Done()
				return
			}
			mappedData, err := requestConfig.ResultMapper(response)
			if err != nil {
				results[requestConfig.Id] = ParalelResult{Error: err}
				defer wg.Done()
				return
			}
			results[requestConfig.Id] = ParalelResult{Result: mappedData}
			defer wg.Done()
			return
		}()
	}
	wg.Wait()
	return results
}
