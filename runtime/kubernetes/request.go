package kubernetes

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/telemetry/log"
)

type request struct {
	host      string
	namespace string
	client    *http.Client
	header    http.Header
	params    url.Values
	method    string
	resource  string
	body      io.Reader
}

func (r *request) setHeader(key, value string) *request {
	r.header.Add(key, value)
	return r
}

func (r *request) get() *request {
	return r.setMethod("GET")
}

func (r *request) setMethod(method string) *request {
	r.method = method
	return r
}

func (r *request) setResource(n string) *request {
	r.resource = n
	return r
}

func (r *request) setParams(p *params) *request {
	for k, v := range p.labelSelector {
		// create new key=value
		pair := fmt.Sprintf("%s=%s", k, v)

		// set it
		r.params.Set("labelSelector", pair)
	}

	return r
}

func (r *request) do() *response {
	req, err := r.request()
	if err != nil {
		return newResponse(nil, err)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return newResponse(nil, err)
	}

	return newResponse(res, err)
}

func (r *request) request() (*http.Request, error) {
	// get the base url
	var url string
	switch r.resource {
	case "pod", "service":
		url = fmt.Sprintf("%s/api/v1/namespaces/%s/%ss/", r.host, r.namespace, r.resource)
	case "deployment":
		url = fmt.Sprintf("%s/apis/apps/v1/namespaces/%s/%ss/", r.host, r.namespace, r.resource)
	}

	// append query params
	if len(r.params) > 0 {
		url += "?" + r.params.Encode()
	}

	log.Infof("HERE IS MY URL: %+v", url)

	// build request
	req, err := http.NewRequest(r.method, url, r.body)
	if err != nil {
		return nil, err
	}

	// set header
	req.Header = r.header

	return req, nil
}

func newRequest(namespace string, options *runtime.RuntimeOptions) *request {
	req := &request{
		host:      options.Host,
		namespace: namespace,
		client:    options.Client,
		header:    make(http.Header),
		params:    make(url.Values),
	}

	if options.BearerToken != "" {
		req.setHeader("Authorization", "Bearer "+options.BearerToken)
	}

	return req
}

type params struct {
	labelSelector map[string]string
}
