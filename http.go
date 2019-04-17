package xobj

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RequestBuilder is only useful for quick prototyping.
// You should not expect that to be capable of everything, and you are probably
// better of using the standard client or resty in production code.
type RequestBuilder struct {
	http                  bool
	host                  string
	port                  int
	path                  string
	retry                 int
	header                http.Header
	query                 url.Values
	responseHeaderTimeout time.Duration
	dialTimeout           time.Duration
	method                string
	body                  io.Reader
	pendingCancelFunc     func()
}

// NewRequest returns a RequestBuilder with some useful defaults.
// You have to create a new request for each query to want to make. Do not recycle it.
func NewRequest() *RequestBuilder {
	return &RequestBuilder{retry: 3, header: http.Header{}, query: url.Values{}, responseHeaderTimeout: time.Second * 30, dialTimeout: time.Second * 30}
}

// Http sets the connection to http
func (r *RequestBuilder) Http() *RequestBuilder {
	r.http = true
	return r
}

// Https sets the connection to use https only, which is the default
func (r *RequestBuilder) Https() *RequestBuilder {
	r.http = false
	return r
}

// Host sets the host part
func (r *RequestBuilder) Host(host string) *RequestBuilder {
	r.host = host
	return r
}

// Ports defines the port to use, if not defined (0) it takes 80 for http and 443 for https
func (r *RequestBuilder) Port(port int) *RequestBuilder {
	r.port = port
	return r
}

// Path concates the given path segments with /. But can also be just a single string already with slashes.
func (r *RequestBuilder) Path(p ...string) *RequestBuilder {
	r.path = strings.Join(p, "/")
	return r
}

// Retry sets the amount of retries before a failure is actually returned. The default is 3.
func (r *RequestBuilder) Retry(retry int) *RequestBuilder {
	r.retry = retry
	return r
}

// Header adds a key/value into the header part of the request
func (r *RequestBuilder) Header(key, value string) *RequestBuilder {
	r.header.Add(key, value)
	return r
}

// Query adds a key/value combination into the query part of the url
func (r *RequestBuilder) Query(key, value string) *RequestBuilder {
	r.header.Add(key, value)
	return r
}

// DialTimeout influences the timeout to establish a connection. The default is 30 seconds.
func (r *RequestBuilder) DialTimeout(duration time.Duration) *RequestBuilder {
	r.dialTimeout = duration
	return r
}

// ResponseHeaderTimeout influences the timeout before the first byte has been received. The default is 30 seconds.
func (r *RequestBuilder) ResponseHeaderTimeout(duration time.Duration) *RequestBuilder {
	r.responseHeaderTimeout = duration
	return r
}

// Body sets a reader to be consumed as a body for the request
func (r *RequestBuilder) Body(reader io.Reader) *RequestBuilder {
	r.body = reader
	return r
}

// JSONBody converts the given obj into a json and sets the content type accordingly
func (r *RequestBuilder) JSONBody(obj Obj) *RequestBuilder {
	r.header.Set("Content-Type", "application/json")
	r.body = bytes.NewReader([]byte(obj.String()))
	return r
}

// Cancel has only an effect while no timeout has run out, the request has ever been started, using an http verb
// and the request has never been cancelled before and the request is still pending.
func (r *RequestBuilder) Cancel() {
	cfunc := r.pendingCancelFunc
	if cfunc != nil {
		cfunc()
	}
}

// Get performs a get request based on the current configuration and tries to parse the result into an obj.
// Otherwise returns an error. Tries to always return the http status code.
func (r *RequestBuilder) Get() (Obj, int, error) {
	return r.genericObjRequest("GET")
}

// Put performs a put request based on the current configuration and tries to parse the result into an obj.
// Otherwise returns an error. Tries to always return the http status code.
func (r *RequestBuilder) Put() (Obj, int, error) {
	return r.genericObjRequest("PUT")
}

// Delete performs a delete request based on the current configuration and tries to parse the result into an obj.
// Otherwise returns an error. Tries to always return the http status code.
func (r *RequestBuilder) Delete() (Obj, int, error) {
	return r.genericObjRequest("DELETE")
}

// Post performs a post request based on the current configuration and tries to parse the result into an obj.
// Otherwise returns an error. Tries to always return the http status code.
func (r *RequestBuilder) Post() (Obj, int, error) {
	return r.genericObjRequest("POST")
}

// genericObjRequest always tries to parse the result as obj
func (r *RequestBuilder) genericObjRequest(method string) (Obj, int, error) {
	r.method = method
	var obj Obj
	var statusCode int
	err := r.doRequest(func(req *http.Request, res *http.Response) error {
		statusCode = res.StatusCode
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		obj, err = Parse(data)
		return err
	})
	return obj, statusCode, err
}

// doRequest performs the actual request based on builder settings and calls the onResult function which
// in turn may return an error which is just delegated.
func (r *RequestBuilder) doRequest(onResult func(req *http.Request, res *http.Response) error) error {
	var transport http.RoundTripper = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: r.responseHeaderTimeout,
		TLSHandshakeTimeout:   r.dialTimeout,
	}

	client := &http.Client{Transport: transport}

	sb := &strings.Builder{}
	if r.http {
		sb.WriteString("http://")
		if r.port == 0 {
			r.port = 80
		}
	} else {
		sb.WriteString("https://")
		if r.port == 0 {
			r.port = 443
		}
	}

	sb.WriteString(r.host)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(r.port))
	if !strings.HasPrefix(r.path, "/") {
		sb.WriteString("/")
	}
	sb.WriteString(r.path)
	if len(r.query) > 0 {
		sb.WriteString("?")
		sb.WriteString(r.query.Encode())
	}

	req, err := http.NewRequest(r.method, sb.String(), r.body)
	if err != nil {
		return err
	}

	for k, v := range r.header {
		for _, s := range v {
			req.Header.Add(k, s)
		}
	}

	backOffTime := time.Millisecond * 200
	var res *http.Response
	for i := 0; i < r.retry; i++ {
		ctx, cancelFunc := context.WithTimeout(context.Background(), r.dialTimeout)
		r.pendingCancelFunc = cancelFunc
		defer cancelFunc()

		req = req.WithContext(ctx)

		res, err = client.Do(req)
		if err != nil {
			logger.Info(Fields{"msg": "http query failed", "backoff": i, "url": sb.String(), "err": err.Error()})
			time.Sleep(backOffTime)
			backOffTime *= 2
			//backoff give up
			if i == r.retry-1 {
				return err
			}

		}
	}

	err = onResult(req, res)
	if err != nil {
		return err
	}

	defer silentClose(res.Body)
	return nil
}

func silentClose(closer io.Closer) {
	_ = closer.Close()
}
