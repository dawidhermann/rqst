package rqst_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dawidhermann/rqst"
)

func TestParalelRequestConfigNew(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	rqstInstance := rqst.NewParalel(server.Client())
	if rqstInstance == nil {
		t.Error("Rqst instance not created")
	}
}

func TestAddParalelRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	requestConfig := rqst.ParalelRequestConfig{
		Id: "requestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, "example.com", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			return response, nil
		},
	}
	rqstInstance := rqst.NewParalel(server.Client()).AddNextParalelRequest(requestConfig)
	result := rqstInstance.Execute()
	if result == nil {
		t.Error("Invalid result")
	}
}

func TestAddMultipleParalelRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	firstRequestConfig := rqst.ParalelRequestConfig{
		Id: "firstRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, "example.com", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			return response, nil
		},
	}
	secondRequestConfig := rqst.ParalelRequestConfig{
		Id: "secondRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, "example.com", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			return response, nil
		},
	}
	rqstInstance := rqst.NewParalel(server.Client()).AddNextMultipleParalelRequests(firstRequestConfig, secondRequestConfig)
	result := rqstInstance.Execute()
	if result == nil {
		t.Error("Invalid result")
	}
}

func TestCombinedMultipleParalelRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	firstRequestConfig := rqst.ParalelRequestConfig{
		Id: "firstRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, "example.com", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			return response, nil
		},
	}
	secondRequestConfig := rqst.ParalelRequestConfig{
		Id: "secondRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, "example.com", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			return response, nil
		},
	}
	thirdRequestConfig := rqst.ParalelRequestConfig{
		Id: "thirdRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, "example.com", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			return response, nil
		},
	}
	rqstInstance := rqst.NewParalel(server.Client()).AddNextParalelRequest(firstRequestConfig).AddNextMultipleParalelRequests(secondRequestConfig, thirdRequestConfig)
	result := rqstInstance.Execute()
	if result == nil {
		t.Error("Invalid result")
	}
}

func TestMappingParalelRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		numberStr := req.URL.Query().Get("number")
		if numberStr == "" {
			t.Error("Number param is empty")
		}
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			t.Error("Failed to convert number")
		}
		rw.Write([]byte(strconv.Itoa(number)))
	}))
	defer server.Close()
	requestConfig := rqst.ParalelRequestConfig{
		Id: "requestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, server.URL+"/?number=3", nil)
			if err != nil {
				t.Error("Failed create request", err)
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error("Failed to get body")
			}
			number, err := strconv.Atoi(string(body))
			if err != nil {
				t.Error("Failed to convert to number", err)
			}
			return number + 1, nil
		},
	}
	secondRequestConfig := rqst.ParalelRequestConfig{
		Id: "secondRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, server.URL+"/?number=5", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error("Failed to get body")
			}
			number, err := strconv.Atoi(string(body))
			if err != nil {
				t.Error("Failed to convert to number", err)
			}
			return number + 1, nil
		},
	}
	rqstInstance := rqst.NewParalel(server.Client()).AddNextParalelRequest(requestConfig).AddNextParalelRequest(secondRequestConfig)
	result := rqstInstance.Execute()
	if result["requestConfig"].Result != 4 {
		t.Error("Failed to get mapped response: ", result)
	}
	if result["secondRequestConfig"].Result != 6 {
		t.Error("Failed to get mapped response: ", result)
	}
}

func TestMappingParalelChanRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		numberStr := req.URL.Query().Get("number")
		if numberStr == "" {
			t.Error("Number param is empty")
		}
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			t.Error("Failed to convert number")
		}
		rw.Write([]byte(strconv.Itoa(number)))
	}))
	defer server.Close()
	requestConfig := rqst.ParalelRequestConfig{
		Id: "requestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, server.URL+"/?number=3", nil)
			if err != nil {
				t.Error("Failed create request", err)
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error("Failed to get body")
			}
			number, err := strconv.Atoi(string(body))
			if err != nil {
				t.Error("Failed to convert to number", err)
			}
			return number + 1, nil
		},
	}
	secondRequestConfig := rqst.ParalelRequestConfig{
		Id: "secondRequestConfig",
		CreateRequest: func() *http.Request {
			req, err := http.NewRequest(http.MethodGet, server.URL+"/?number=5", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error("Failed to get body")
			}
			number, err := strconv.Atoi(string(body))
			if err != nil {
				t.Error("Failed to convert to number", err)
			}
			return number + 1, nil
		},
	}
	rqstInstance := rqst.NewParalel(server.Client()).AddNextParalelRequest(requestConfig).AddNextParalelRequest(secondRequestConfig)
	resultChan := rqstInstance.ExecuteWithChan()
	for result := range resultChan {
		if result.Id == "requestConfig" && result.Result != 4 {
			t.Error("Failed to get mapped response: ", result)
		}
		if result.Id == "secondRequestConfig" && result.Result != 6 {
			t.Error("Failed to get mapped response: ", result)
		}
	}
}
