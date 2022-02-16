// SPDX-FileCopyrightText: © 2022 Jade Meskill <iamruinous@ruinous.social>
//
// SPDX-License-Identifier: MIT

package geminitest

import (
	"context"
	"fmt"
	"io"
	"testing"

	gemini "git.sr.ht/~adnano/go-gemini"
)

func TestRecorder(t *testing.T) {
	type checkFunc func(*ResponseRecorder) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasStatus := func(wantCode gemini.Status) checkFunc {
		return func(rec *ResponseRecorder) error {
			if rec.Status != wantCode {
				return fmt.Errorf("Status = %d; want %d", rec.Status, wantCode)
			}
			return nil
		}
	}
	hasResultStatus := func(want gemini.Status) checkFunc {
		return func(rec *ResponseRecorder) error {
			if rec.Result().Status != want {
				return fmt.Errorf("Result().Status = %q; want %q", rec.Result().Status, want)
			}
			return nil
		}
	}
	hasResultContents := func(want string) checkFunc {
		return func(rec *ResponseRecorder) error {
			contentBytes, err := io.ReadAll(rec.Result().Body)
			if err != nil {
				return err
			}
			contents := string(contentBytes)
			if contents != want {
				return fmt.Errorf("Result().Body = %s; want %s", contents, want)
			}
			return nil
		}
	}
	hasContents := func(want string) checkFunc {
		return func(rec *ResponseRecorder) error {
			if rec.Body.String() != want {
				return fmt.Errorf("wrote = %q; want %q", rec.Body.String(), want)
			}
			return nil
		}
	}
	hasFlush := func(want bool) checkFunc {
		return func(rec *ResponseRecorder) error {
			if rec.Flushed != want {
				return fmt.Errorf("Flushed = %v; want %v", rec.Flushed, want)
			}
			return nil
		}
	}

	for _, tt := range [...]struct {
		name   string
		h      func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request)
		checks []checkFunc
	}{
		{
			"20 default",
			func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {},
			check(hasStatus(gemini.StatusSuccess), hasContents(""), hasResultStatus(gemini.StatusSuccess)),
		},
		{
			"write sends 20",
			func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
				w.Write([]byte("hi first"))
				w.WriteHeader(gemini.StatusSuccess, "")
				w.WriteHeader(gemini.StatusPermanentFailure, "")
			},
			check(hasStatus(gemini.StatusSuccess), hasContents("hi first"), hasFlush(false)),
		},
		{
			"write string",
			func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
				io.WriteString(w, "hi first")
			},
			check(
				hasStatus(gemini.StatusSuccess),
				hasContents("hi first"),
				hasFlush(false),
			),
		},
		// {
		// 	"flush",
		// 	func(w http.ResponseWriter, r *http.Request) {
		// 		w.(http.Flusher).Flush() // also sends a 200
		// 		w.WriteHeader(gemini.StatusPermanentFailure)
		// 	},
		// 	check(hasStatus(gemini.StatusSuccess), hasFlush(true), hasContentLength(-1)),
		// },
		{
			"nil ResponseRecorder.Body",
			func(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
				w.(*ResponseRecorder).Body = nil
				io.WriteString(w, "hi")
			},
			check(hasResultContents("")), // check we don't crash reading the body

		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := gemini.NewRequest("gemini://foo.com/")
			h := gemini.HandlerFunc(tt.h)
			rec := NewRecorder()
			ctx := context.Background()
			h.ServeGemini(ctx, rec, r)
			for _, check := range tt.checks {
				if err := check(rec); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
