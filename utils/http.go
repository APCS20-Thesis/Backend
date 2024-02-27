package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// HTTP Header
	Header_AUTHORIZATION = "Authorization"
	Header_CONTENT_TYPE  = "Content-Type"
)

const (
	Method_POST  = "POST"
	Method_GET   = "GET"
	Method_PUT   = "PUT"
	Method_PATCH = "PATCH"
)

type HttpClient struct {
	httpClient *http.Client
	log        logr.Logger
	host       string
}

type Request struct {
	Endpoint string
	Method   string
	Body     interface{}
	Params   map[string]string
	Headers  map[string]string
}

type BasicAuth struct {
	Username string
	Password string
}

func (s *HttpClient) SendHttpRequest(ctx context.Context, request Request, response interface{}) error {
	queryString := ""
	for k, v := range request.Params {
		if queryString != "" {
			queryString += "&"
		}
		queryString += k + "=" + url.QueryEscape(v)
	}

	urlPath := fmt.Sprintf("%v%v", s.host, request.Endpoint)
	if queryString != "" {
		urlPath += "?" + queryString
	}

	data, err := json.Marshal(request.Body)
	if err != nil {
		return err
	}

	reqBody := bytes.NewReader(data)

	req, err := http.NewRequestWithContext(ctx, request.Method, urlPath, reqBody)
	if err != nil {
		return err
	}

	for h, v := range request.Headers {
		req.Header.Set(h, v)
	}

	s.log.Info("[HTTP-Client] Begin to call http request", "url", urlPath, "method", request.Method, "body", request.Body)
	resp, err := s.httpClient.Do(req)
	s.log.Info("[HTTP-Client] Finish to call http request", "url", urlPath, "method", request.Method, "body", request.Body, "error", err)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error(err, "[HTTP-Client] Fail to read response body")
		return err
	}

	responseBodyStr := string(responseBody)
	switch resp.StatusCode / 100 {
	case 2:
		s.log.Info("[HTTP-Client] Succeed processing http request", "response", responseBodyStr)
		if response == nil {
			return nil
		}
		return json.Unmarshal(responseBody, &response)
	case 4:
		s.log.Error(err, "[HTTP-Client] Fail processing http request", "code", resp.StatusCode, "response", responseBodyStr)
		return status.Errorf(codes.Aborted, "Http Status=%v: %s", resp.StatusCode, responseBodyStr)
	default:
		s.log.Error(err, "[HTTP-Client] Error processing http request", "code", resp.StatusCode, "response", responseBodyStr)
		return status.Errorf(codes.Aborted, "Http Status=%v: %s", resp.StatusCode, responseBodyStr)
	}
}

func (s *HttpClient) SendHttpRequestWithBasicAuth(ctx context.Context, basicAuth BasicAuth, request Request, response interface{}) error {
	queryString := ""
	for k, v := range request.Params {
		if queryString != "" {
			queryString += "&"
		}
		queryString += k + "=" + url.QueryEscape(v)
	}

	urlPath := fmt.Sprintf("%v%v", s.host, request.Endpoint)
	if queryString != "" {
		urlPath += "?" + queryString
	}

	data, err := json.Marshal(request.Body)
	if err != nil {
		return err
	}

	reqBody := bytes.NewReader(data)

	req, err := http.NewRequestWithContext(ctx, request.Method, urlPath, reqBody)
	if err != nil {
		return err
	}

	for h, v := range request.Headers {
		req.Header.Set(h, v)
	}

	req.SetBasicAuth(basicAuth.Username, basicAuth.Password)

	s.log.Info("[HTTP-Client] Begin to call http request", "url", urlPath, "method", request.Method, "body", request.Body)
	resp, err := s.httpClient.Do(req)
	s.log.Info("[HTTP-Client] Finish to call http request", "url", urlPath, "method", request.Method, "body", request.Body, "error", err)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.log.Error(err, "[HTTP-Client] Fail to read response body")
		return err
	}

	responseBodyStr := string(responseBody)
	switch resp.StatusCode / 100 {
	case 2:
		s.log.Info("[HTTP-Client] Succeed processing http request", "response", responseBodyStr)
		if response == nil {
			return nil
		}
		return json.Unmarshal(responseBody, &response)
	case 4:
		s.log.Error(err, "[HTTP-Client] Fail processing http request", "code", resp.StatusCode, "response", responseBodyStr)
		return status.Errorf(codes.Aborted, "Http Status=%v: %s", resp.StatusCode, responseBodyStr)
	default:
		s.log.Error(err, "[HTTP-Client] Error processing http request", "code", resp.StatusCode, "response", responseBodyStr)
		return status.Errorf(codes.Aborted, "Http Status=%v: %s", resp.StatusCode, responseBodyStr)
	}
}

func (s *HttpClient) Init(clientName string, log logr.Logger, host string) {
	s.host = host
	s.httpClient = &http.Client{}
	s.log = log.WithName("http_client/" + clientName)
}

func (s *HttpClient) InitForTest(clientName string, log logr.Logger, host string, client *http.Client) {
	s.host = host
	if client != nil {
		s.httpClient = client
	} else {
		s.httpClient = &http.Client{}
	}
	s.log = log.WithName("http_client/" + clientName)
}
