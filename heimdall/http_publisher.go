package heimdall

import (
	"gopkg.in/inconshreveable/log15.v2"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type httpHookPublisher struct {
	config *HttpHookConfig
}

func NewHttpHookPublisher(config HttpHookConfig) *httpHookPublisher {
	return &httpHookPublisher{
		&config,
	}
}

func (s *httpHookPublisher) Send(publishable Publishable) {
	logger := log15.New("MODULE", "HttpHookPublisher")
	req, err := http.NewRequest(s.config.Method, s.config.Url, nil)

	if err != nil {
		logger.Error("unable to send HTTP notification", "error", err)
	}

	for header, headerValue := range s.config.Headers {
		req.Header.Add(header, headerValue)
	}

	client := &http.Client{
		Timeout: time.Duration(s.config.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("failed executing http request", "error", err)
	}

	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		logger.Error(
			"remote http server returned error",
			"status", resp.Status,
			"code", resp.StatusCode,
			"body", readResponseOrDefault(resp.Body, ""),
		)
	}
}

func readResponseOrDefault(reader io.ReadCloser, defaultValue string) string {
	if body, err := ioutil.ReadAll(reader); err == nil {
		return string(body)
	}

	return defaultValue
}