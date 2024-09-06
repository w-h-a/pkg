package httputils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpGet(url string) ([]byte, error) {
	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}

	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
