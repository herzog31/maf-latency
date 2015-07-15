package latency

import (
	"net/http"
	"testing"
	"time"
)

// TODO(mjb): implement real tests ;)
func TestGoogle(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://google.com", nil)
	request := LatencyRequest{
		Request: req,
		Timeout: 10 * time.Second,
	}
	_, err := request.Execute()
	if err != nil {
		t.Errorf("Failed at performing the request: %v", err)
	}
}

func TestMarbec(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://marb.ec", nil)
	request := LatencyRequest{
		Request: req,
		Timeout: 10 * time.Second,
	}
	_, err := request.Execute()
	if err != nil {
		t.Errorf("Failed at performing the request: %v", err)
	}
}
