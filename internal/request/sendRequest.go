package request

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout:   time.Second,
	Transport: http.DefaultTransport.(*http.Transport).Clone(),
}

func Get() {
}
