package supertest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/franela/goreq"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
)

type Request struct {
	base    string
	method  string
	path    string
	done    reflect.Value
	body    interface{}
	headers []headerTuple
	query   url.Values
}

type headerTuple struct {
	name  string
	value string
}

func NewRequest(base string) *Request {
	r := &Request{}
	r.base = base
	return r
}

func (r *Request) Get(path string) *Request {
	r.method = "GET"
	r.path = path
	return r
}

func (r *Request) Post(path string) *Request {
	r.method = "POST"
	r.path = path
	return r
}

func (r *Request) Put(path string) *Request {
	r.method = "PUT"
	r.path = path
	return r
}

func (r *Request) Delete(path string) *Request {
	r.method = "DELETE"
	r.path = path
	return r
}

func (r *Request) Patch(path string) *Request {
	r.method = "PATCH"
	r.path = path
	return r
}

func (r *Request) Options(path string) *Request {
	r.method = "OPTIONS"
	r.path = path
	return r
}

func (r *Request) Head(path string) *Request {
	r.method = "HEAD"
	r.path = path
	return r
}

func (r *Request) Send(body interface{}) *Request {
	r.body = body
	return r
}

func (r *Request) Set(name, value string) *Request {
	r.headers = append(r.headers, headerTuple{name: name, value: value})
	return r
}

func (r *Request) Query(name, value string) *Request {
	if r.query == nil {
		r.query = url.Values{}
	}
	r.query.Add(name, value)
	return r
}

func (r *Request) Expect(code int, args ...interface{}) error {

	var bodyToCompare interface{}

	if len(args) == 1 {
		if reflect.ValueOf(args[0]).Kind() == reflect.Func {
			r.done = reflect.ValueOf(args[0])
		} else {
			bodyToCompare = args[0]
		}
	}

	if len(args) == 2 {
		bodyToCompare = args[0]
		r.done = reflect.ValueOf(args[1])
	}

	var err error

	var body io.Reader

	if r.body != nil {
		body, err = prepareRequestBody(r.body)
		if err != nil {
			return err
		}
	}

	req := goreq.Request{Method: r.method, Uri: r.base + r.path + "?" + r.query.Encode(), Body: body}
	for _, tuple := range r.headers {
		req.AddHeader(tuple.name, tuple.value)
	}
	res, e := req.Do()

	if e != nil {
		err = e
	} else {
		if res.StatusCode != code {
			err = fmt.Errorf("Expected %d, was %d", code, res.StatusCode)
		} else if bodyToCompare != nil {
			// Read the entire response body
			b, _ := ioutil.ReadAll(res.Body)

			if s, ok := bodyToCompare.(string); ok {
				// It is a string
				str := string(b)

				if s != str {
					err = fmt.Errorf(fmt.Sprintf("%#v", s) + " does not equal " + fmt.Sprintf("%#v", str))
				}
			} else {
				// Try to parse to JSON
				ptrNewValue := reflect.New(reflect.TypeOf(bodyToCompare))
				newValue := reflect.Indirect(ptrNewValue)

				e := json.Unmarshal(b, ptrNewValue.Interface())
				if e != nil {
					err = fmt.Errorf("Expected: %#v, but got %#v. %s", bodyToCompare, string(b), e)
				} else {
					if !objectsAreEqual(bodyToCompare, newValue.Interface()) {
						err = fmt.Errorf(fmt.Sprintf("%#v", bodyToCompare) + " does not equal " + fmt.Sprintf("%#v", newValue.Interface()))
					}
				}
			}
		}
	}

	if r.done.IsValid() {
		if err != nil {
			r.done.Call([]reflect.Value{reflect.ValueOf(err)})
		} else {
			r.done.Call([]reflect.Value{})
		}
	}
	return err
}

func objectsAreEqual(a, b interface{}) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}

	if reflect.ValueOf(a) == reflect.ValueOf(b) {
		return true
	}

	if fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b) {
		return true
	}

	return false
}

func prepareRequestBody(b interface{}) (io.Reader, error) {
	var body io.Reader

	if sb, ok := b.(string); ok {
		// treat is as text
		body = strings.NewReader(sb)
	} else if rb, ok := b.(io.Reader); ok {
		// treat is as text
		body = rb
	} else if bb, ok := b.([]byte); ok {
		//treat as byte array
		body = bytes.NewReader(bb)
	} else {
		// try to jsonify it
		j, err := json.Marshal(b)
		if err == nil {
			body = bytes.NewReader(j)
		} else {
			return nil, err
		}
	}

	return body, nil
}
