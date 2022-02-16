// SPDX-FileCopyrightText: © 2022 Jade Meskill <iamruinous@ruinous.social>
//
// SPDX-License-Identifier: MIT

package geminitest

import (
	"io"
	"net/url"
	"reflect"
	"testing"

	gemini "git.sr.ht/~adnano/go-gemini"
)

func TestNewRequest(t *testing.T) {
	for _, tt := range [...]struct {
		name string

		method, uri string
		body        io.Reader

		want     *gemini.Request
		wantBody string
	}{
		{
			name: "full URL",
			uri:  "gemini://foo.com/path/%2f/bar/",
			want: &gemini.Request{
				URL: &url.URL{
					Scheme:  "gemini",
					Path:    "/path///bar/",
					RawPath: "/path/%2f/bar/",
					Host:    "foo.com",
				},
			},
			wantBody: "",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRequest(tt.uri)
			if !reflect.DeepEqual(got.URL, tt.want.URL) {
				t.Errorf("Request.URL mismatch:\n got: %#v\nwant: %#v", got.URL, tt.want.URL)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request mismatch:\n got: %#v\nwant: %#v", got, tt.want)
			}
		})
	}
}
