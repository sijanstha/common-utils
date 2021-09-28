package rest

import (
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	retryCount = 5
	timeOut    = 1000 * time.Millisecond
)

var (
	RestClient *resty.Client
)

func init() {
	isDebugStr := os.Getenv("is_debug")

	isDebug := false
	if isDebugStr == "true" {
		isDebug = true
	}

	RestClient = resty.New()
	RestClient.SetRetryCount(retryCount).SetHeaders(map[string]string{"Content-Type": "application/json"}).SetTimeout(timeOut).SetDebug(isDebug)
}
