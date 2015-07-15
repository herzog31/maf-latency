/*
MAF Latency Package

Measure the latency of HTTP requests.

Usage

First specify a new http.Request like:

	req, _ := http.NewRequest("GET", "http://google.com", nil)

To use this request with the latency library you have to embed it in a
LatencyRequest and specify a timeout interval. Example:

	request := LatencyRequest{
		Request: req,
		Timeout: 10 * time.Second,
	}

To execute the latency measurement, simply call the Execute() method:

	response, err := request.Execute()

Note that the latency measurement follows redirects, but only measures the time
of the last non-redirecting request. The call will return a LatencyResponse
object, containing the latency and an array of traversed redirections as
strings.
*/
package latency

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrLatencyRedirect = "Redirect in request"
)

// LatencyRequest is a wrapper for http.Request.
// You can specify an additonal timeout
type LatencyRequest struct {
	*http.Request
	Timeout   time.Duration
	Redirects []string
}

func (l *LatencyRequest) String() string {
	return fmt.Sprintf("%v %v, timout is %v", l.Method, l.URL, l.Timeout)
}

// LatencyRedirectError is an error type used, if redirection occured during
// latency measurement
type LatencyRedirectError string

func (l LatencyRedirectError) Error() string {
	return ErrLatencyRedirect
}

// Redirect policy for http.Client
// It will always return an error to detect redirects and to restart the
// latency measurement
func NoRedirectsPolicy(req *http.Request, via []*http.Request) error {
	return LatencyRedirectError(fmt.Sprintf("%v", req.URL))
}

// Performans the LatencyRequest and returns a LatencyResponse
func (l *LatencyRequest) Execute() (resp *LatencyResponse, err error) {

	client := &http.Client{
		Timeout:       l.Timeout,
		CheckRedirect: NoRedirectsPolicy,
	}

	if len(l.Redirects) > 0 {
		// Parse redirect target
		newUrl, err := url.Parse(l.Redirects[len(l.Redirects)-1])
		if err != nil {
			return nil, err
		}
		// Adapt request
		l.Request.URL = newUrl
		l.Request.Host = newUrl.Host
	} else {
		l.Redirects = append(l.Redirects, l.Request.URL.String())
	}
	l.Request.Header.Set("Cache control", "no-cache")

	timeStart := time.Now()
	response, err := client.Do(l.Request)
	if err != nil {
		// Redirect error
		if strings.Contains(err.Error(), ErrLatencyRedirect) {
			if loc := response.Header.Get("Location"); loc != "" {
				l.Redirects = append(l.Redirects, loc)
				return l.Execute()
			} else {
				return nil, errors.New("No redirect location given!")
			}
		} else {
			return nil, err
		}
	}
	latency := time.Since(timeStart)

	response.Body.Close()

	resp = &LatencyResponse{response, latency, l.Redirects}
	return resp, nil
}

// LatencyResponse is a wrapper for http.Response.
// It contains the latency that was measured
type LatencyResponse struct {
	*http.Response
	Latency   time.Duration
	Redirects []string
}

func (l *LatencyResponse) String() string {
	if len(l.Redirects) > 1 {
		return fmt.Sprintf("%v %v %v", l.Request.Method,
			strings.Join(l.Redirects, " -> "), l.Latency)
	}
	return fmt.Sprintf("%v %v %v", l.Request.Method, l.Request.URL, l.Latency)
}
