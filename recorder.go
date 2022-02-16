// SPDX-FileCopyrightText: © 2022 Jade Meskill <iamruinous@ruinous.social>
//
// SPDX-License-Identifier: MIT

package geminitest

import (
	"bytes"
	"io"

	gemini "git.sr.ht/~adnano/go-gemini"
)

// ResponseRecorder is an implementation of gemini.ResponseWriter that
// records its mutations for later inspection in tests.
type ResponseRecorder struct {
	// Status is the response status code set by WriteHeader.
	Status gemini.Status

	// Meta returns the response meta.
	// For successful responses, the meta should contain the media type of the response.
	// For failure responses, the meta should contain a short description of the failure.
	Meta string

	// Body is the buffer to which the Handler's Write calls are sent.
	// If nil, the Writes are silently discarded.
	Body *bytes.Buffer

	// Flushed is whether the Handler called Flush.
	Flushed bool

	// The MediaType set by SetMediaType
	MediaType string

	result *gemini.Response // cache of Result's return value
}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		Body: new(bytes.Buffer),
	}
}

// Flush implements gemini.Flusher. To test whether Flush was
// called, see rw.Flushed.
func (w *ResponseRecorder) Flush() error {
	w.Flushed = true
	return nil
}

// SetMediaType implements gemini.ResponseWriter.
func (w *ResponseRecorder) SetMediaType(mediatype string) {
	w.MediaType = mediatype
}

// Write implements gemini.ResponseWriter. The data in buf is written to
// rw.Body, if not nil.
func (rw *ResponseRecorder) Write(buf []byte) (int, error) {
	rw.WriteHeader(gemini.StatusSuccess, "text/gemini")
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

// WriteHeader implements gemini.ResponseWriter.
func (rw *ResponseRecorder) WriteHeader(status gemini.Status, meta string) {
	rw.Status = status
	rw.Meta = meta
}

// Result returns the response generated by the handler.
//
// The returned Response will have at least its Status, MediaType
// and  Body populated. More fields may be populated in the future,
// so callers should not DeepEqual the result in tests.
//
// The Response.Body is guaranteed to be non-nil and Body.Read call is
// guaranteed to not return any error other than io.EOF.
//
// Result must only be called after the handler has finished running.
func (rw *ResponseRecorder) Result() *gemini.Response {
	if rw.result != nil {
		return rw.result
	}
	res := &gemini.Response{
		Status: rw.Status,
		Meta:   rw.Meta,
	}
	rw.result = res
	if res.Status == 0 {
		res.Status = gemini.StatusSuccess
	}
	if rw.Body != nil {
		res.Body = io.NopCloser(bytes.NewReader(rw.Body.Bytes()))
	}
	return res
}