package rqst_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/dawidhermann/rqst"
)

func TestRequestConfigNew(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	rqstInstance := rqst.New(server.Client())
	if rqstInstance == nil {
		t.Error("Rqst instance not created")
	}
}

func TestAddRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	requestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	rqstInstance := rqst.New(server.Client()).AddNextRequest(&requestConfig)
	result := rqstInstance.Execute()
	if result == nil {
		t.Error("Invalid result")
	}
}

func TestAddMultipleRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	firstRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	secondRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	rqstInstance := rqst.New(server.Client()).AddNextMultipleRequests(&firstRequestConfig, &secondRequestConfig)
	result := rqstInstance.Execute()
	if result == nil {
		t.Error("Invalid result")
	}
}

func TestCombinedMultipleRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	firstRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	secondRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	thirdRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	rqstInstance := rqst.New(server.Client()).AddNextRequest(&firstRequestConfig).AddNextMultipleRequests(&secondRequestConfig, &thirdRequestConfig)
	result := rqstInstance.Execute()
	if result == nil {
		t.Error("Invalid result")
	}
}

func TestMappingRequest(t *testing.T) {
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
	requestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
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
	rqstInstance := rqst.New(server.Client()).AddNextRequest(&requestConfig)
	result := rqstInstance.Execute()
	if result != 4 {
		t.Error("Failed to get mapped response: ", result)
	}
}

func TestRawResultInterceptor(t *testing.T) {
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
	firstRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
			req, err := http.NewRequest(http.MethodGet, server.URL+"/?number=2", nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			buf := &bytes.Buffer{}
			tee := io.TeeReader(response.Body, buf)
			body, err := io.ReadAll(tee)
			if err != nil {
				t.Error("Failed to read", err)
			}
			number, err := strconv.Atoi(string(body))
			if err != nil {
				t.Error("Failed to convert to number", err)
			}
			return number + 1, nil
		},
		MappedResultInterceptor: func(mappedData any) {
			if mappedData != 3 {
				t.Error("Incorrect result", mappedData)
			}
		},
	}
	secondRequestConfig := rqst.RequestConfig{
		CreateRequest: func(data any) *http.Request {
			if data != 3 {
				t.Error("Incorrect input data", data)
			}
			num, ok := data.(int)
			if !ok {
				t.Error("incorrect value", data)
			}
			url := fmt.Sprintf("%v/?number=%d", server.URL, num)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Error("Failed to create request")
			}
			return req
		},
		ResultMapper: func(response *http.Response) (any, error) {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error("Failed to read", err)
			}
			number, err := strconv.Atoi(string(body))
			if err != nil {
				t.Error("Failed to convert to number", err)
			}
			return number + 1, nil
		},
	}
	rqstInstance := rqst.New(server.Client()).AddNextMultipleRequests(&firstRequestConfig, &secondRequestConfig)
	result := rqstInstance.Execute()
	if result != 4 {
		t.Error("Failed to get mapped response: ", result)
	}
}
