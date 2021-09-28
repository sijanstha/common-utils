package rest

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	retryCount = 5
	timeOut    = 1000 * time.Millisecond
	isDebug    = true
)

var (
	RestClient *resty.Client
)

func init() {
	RestClient = resty.New()
	RestClient.SetRetryCount(retryCount).SetHeaders(map[string]string{"Content-Type": "application/json"}).SetTimeout(timeOut).SetDebug(isDebug)
}
