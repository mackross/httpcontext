package httpcontext

import (
	"io"
	"net/http"
)

// Set sets a context value on req.
// It currently accomplishes this by replacing the http.Request’s Body with
// a ContextReadCloser, which wraps the original io.ReadCloser.
// See “Invasion of the Body Snatchers.”
func Set(req *http.Request, key interface{}, value interface{}) {
	crc := getContextReadCloser(req)
	crc.Context()[key] = value
}

// Get gets a context value from req.
// Returns nil if key not found in the request context.
func Get(req *http.Request, key interface{}) interface{} {
	crc := getContextReadCloser(req)
	if value, ok := crc.Context()[key]; ok {
		return value
	}
	return nil
}

// GetOk gets a context value from req.
// Returns (nil, false) if key not found in the request context.
func GetOk(req *http.Request, key interface{}) (interface{}, bool) {
	crc := getContextReadCloser(req)
	if value, ok := crc.Context()[key]; ok {
		return value, ok
	}
	return nil, false
}

// GetAll returns all stored context values for a request.
// Will always return a valid map. Returns an empty map for
// requests context data previously set.
func GetAll(req *http.Request) map[interface{}]interface{} {
	crc := getContextReadCloser(req)
	return crc.Context()
}

// Delete deletes a stored value from a request’s context.
func Delete(req *http.Request, key interface{}) {
	crc := getContextReadCloser(req)
	delete(crc.Context(), key)
}

// Clear clears all stored values from a request’s context.
func Clear(req *http.Request) {
	crc := getContextReadCloser(req)
	ctx := crc.Context()
	for key, _ := range ctx {
		delete(ctx, key)
	}
}

// ContextReadCloser implements the io.ReadCloser interface
// with two additional methods: Context() and SetContext().
type ContextReadCloser interface {
	io.ReadCloser
	Context() map[interface{}]interface{}
}

type contextReadCloser struct {
	io.ReadCloser
	context map[interface{}]interface{}
}

func (crc *contextReadCloser) Context() map[interface{}]interface{} {
	return crc.context
}

func getContextReadCloser(req *http.Request) ContextReadCloser {
	crc, ok := req.Body.(ContextReadCloser)
	if !ok {
		crc = &contextReadCloser{
			ReadCloser: req.Body,
			context:    make(map[interface{}]interface{}),
		}
		req.Body = crc
	}
	return crc
}
