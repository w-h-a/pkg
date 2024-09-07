package httputils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func HttpGet(url string) ([]byte, error) {
	rsp, err := http.Get(SanitizeHttpUrl(url))
	if err != nil {
		return nil, err
	}

	return ExtractBody(rsp.Body)
}

func HttpPost(url string, data []byte) ([]byte, error) {
	rsp, err := http.Post(SanitizeHttpUrl(url), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return ExtractBody(rsp.Body)
}

func SanitizeHttpUrl(url string) string {
	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("http://%s", url)
	}

	return url
}

func ExtractBody(r io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	r.Close()

	return body, nil
}
