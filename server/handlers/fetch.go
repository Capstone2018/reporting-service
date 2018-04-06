package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 30,
}

//fetchHTML does an HTTP GET for the pageURL
//and returns an error if the status code is >= 400
//or the content type doesn't start with `text/html`.
//If all goes well, it returns the response body.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, err
	}

	//check status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("response status code was %d", resp.StatusCode)
	}

	//check content-type
	contentType := resp.Header.Get(headerContentType)
	if !strings.HasPrefix(contentType, contentTypeTextHTML) {
		return nil, fmt.Errorf("requested URL is not an HTML page: content type was %s", contentType)
	}

	return resp.Body, nil
}
