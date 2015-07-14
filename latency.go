/*
MAF Latency Package

Measure the latency of HTTP requests.
*/
package latency

import (
	"fmt"
	"net/http"
	"time"
)

// LatencyRequest is a wrapper for http.Request.
// You can specify an additonal timeout
type LatencyRequest struct {
	Request *http.Request
	Timeout time.Duration
}

func (l *LatencyRequest) String() string {
	return fmt.Sprintf("%v %v, timout is %v", l.Request.Method, l.Request.URL, l.Timeout)
}

// Performans the LatencyRequest and returns a LatencyResponse
func (l *LatencyRequest) Execute() (resp *LatencyResponse, err error) {

	l.Request.Header.Add("Cache control", "no-cache")
	client := &http.Client{Timeout: l.Timeout}

	timeStart := time.Now()
	response, err := client.Do(l.Request)
	if err != nil {
		return nil, err
	}
	latency := time.Since(timeStart)

	resp = &LatencyResponse{Response: response, Latency: latency}
	return resp, nil
}

// LatencyResponse is a wrapper for http.Response.
// It contains the latency that was measured
type LatencyResponse struct {
	Response *http.Response
	Latency  time.Duration
}

func (l LatencyResponse) String() string {
	return fmt.Sprintf("%v %v %v", l.Response.Request.Method, l.Response.Request.URL, l.Latency)
}
